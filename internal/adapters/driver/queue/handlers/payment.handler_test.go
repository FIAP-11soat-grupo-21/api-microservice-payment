package handlers

import (
	"context"
	"errors"
	"testing"

	"payment_microservice/internal/common/config/env"
	"payment_microservice/internal/core/domain/entities"
	"payment_microservice/internal/core/domain/ports"
	"payment_microservice/internal/core/dto"

	"github.com/stretchr/testify/assert"
)

type mockQueuePublisher struct {
	publishedTopic string
	publishedBody  []byte
	publishErr     error
	calledTopic    bool
}

func (m *mockQueuePublisher) PublishOnQueue(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (m *mockQueuePublisher) PublishOnTopic(_ context.Context, topic string, message []byte) error {
	m.calledTopic = true
	m.publishedTopic = topic
	m.publishedBody = message
	return m.publishErr
}

type mockPaymentGateway struct {
	result dto.PIXBillingResultDTO
	err    error
}

func (m *mockPaymentGateway) CreatePIXBilling(_ dto.CreatePIXBillingDTO) (dto.PIXBillingResultDTO, error) {
	return m.result, m.err
}

type mockPaymentRepository struct {
	inserted   *entities.Payment
	insertErr  error
	findResult entities.Payment
	findErr    error
	deletedID  string
	deleteErr  error
}

func (m *mockPaymentRepository) Insert(_ context.Context, payment entities.Payment) error {
	m.inserted = &payment
	return m.insertErr
}

func (m *mockPaymentRepository) FindByOrderID(_ context.Context, _ string) (entities.Payment, error) {
	return m.findResult, m.findErr
}

func (m *mockPaymentRepository) Update(_ context.Context, _ entities.Payment) error {
	return nil
}

func (m *mockPaymentRepository) Delete(_ context.Context, paymentID string) error {
	m.deletedID = paymentID
	return m.deleteErr
}

func TestCreatePayment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		repo := &mockPaymentRepository{}
		gateway := &mockPaymentGateway{result: dto.PIXBillingResultDTO{QRData: "qr"}}

		originalRepo := newPaymentRepository
		originalGateway := newPaymentGateway
		newPaymentRepository = func() ports.IPaymentRepository { return repo }
		newPaymentGateway = func() ports.IPaymentGateway { return gateway }
		t.Cleanup(func() {
			newPaymentRepository = originalRepo
			newPaymentGateway = originalGateway
		})

		messageJSON := []byte(`{"order_id":"order-1","amount":10}`)

		err := CreatePayment(messageJSON)

		assert.NoError(t, err)
		if assert.NotNil(t, repo.inserted) {
			assert.Equal(t, "order-1", repo.inserted.OrderID)
		}
	})

	t.Run("parse error", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		err := CreatePayment([]byte("invalid-json"))

		assert.Error(t, err)
	})

	t.Run("use case error triggers rollback", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		repo := &mockPaymentRepository{insertErr: errors.New("insert fail")}
		queuePub := &mockQueuePublisher{}
		originalRepo := newPaymentRepository
		originalGateway := newPaymentGateway
		originalQueue := newQueuePublisher
		newPaymentRepository = func() ports.IPaymentRepository { return repo }
		newPaymentGateway = func() ports.IPaymentGateway {
			return &mockPaymentGateway{result: dto.PIXBillingResultDTO{QRData: "qr"}}
		}
		newQueuePublisher = func() ports.IQueuePublisher { return queuePub }
		t.Cleanup(func() {
			newPaymentRepository = originalRepo
			newPaymentGateway = originalGateway
			newQueuePublisher = originalQueue
		})

		messageJSON := []byte(`{"order_id":"order-err","amount":5}`)

		err := CreatePayment(messageJSON)

		assert.Error(t, err)
		assert.True(t, queuePub.calledTopic)
		assert.Contains(t, string(queuePub.publishedBody), "order-err")
	})
}

func TestRollbackPayment(t *testing.T) {
	t.Run("parse error", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		repo := &mockPaymentRepository{}
		originalRepo := newPaymentRepository
		newPaymentRepository = func() ports.IPaymentRepository { return repo }
		t.Cleanup(func() { newPaymentRepository = originalRepo })

		err := RollbackPayment([]byte("bad-json"))

		assert.Error(t, err)
		assert.Equal(t, "", repo.deletedID)
	})

	t.Run("ignores self triggered", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		repo := &mockPaymentRepository{}
		originalRepo := newPaymentRepository
		newPaymentRepository = func() ports.IPaymentRepository { return repo }
		t.Cleanup(func() { newPaymentRepository = originalRepo })

		msg := []byte(`{"order_id":"order-ignored","system_triggered":"payment_service"}`)

		err := RollbackPayment(msg)

		assert.NoError(t, err)
		assert.Equal(t, "", repo.deletedID)
	})

	t.Run("successful delete", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		repo := &mockPaymentRepository{findResult: entities.Payment{ID: "p-1", OrderID: "order-123"}}
		originalRepo := newPaymentRepository
		newPaymentRepository = func() ports.IPaymentRepository { return repo }
		t.Cleanup(func() { newPaymentRepository = originalRepo })

		msg := []byte(`{"order_id":"order-123","system_triggered":"other"}`)

		err := RollbackPayment(msg)

		assert.NoError(t, err)
		assert.Equal(t, "p-1", repo.deletedID)
	})

	t.Run("find error", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		repo := &mockPaymentRepository{findErr: errors.New("find err")}
		originalRepo := newPaymentRepository
		newPaymentRepository = func() ports.IPaymentRepository { return repo }
		t.Cleanup(func() { newPaymentRepository = originalRepo })

		msg := []byte(`{"order_id":"order-500","system_triggered":"other"}`)

		err := RollbackPayment(msg)

		assert.Error(t, err)
	})

	t.Run("no payment found", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		repo := &mockPaymentRepository{findResult: entities.Payment{}}
		originalRepo := newPaymentRepository
		newPaymentRepository = func() ports.IPaymentRepository { return repo }
		t.Cleanup(func() { newPaymentRepository = originalRepo })

		msg := []byte(`{"order_id":"order-nil","system_triggered":"other"}`)

		err := RollbackPayment(msg)

		assert.NoError(t, err)
		assert.Equal(t, "", repo.deletedID)
	})
}

func TestSendRollbackEvent(t *testing.T) {
	t.Run("success publishes topic", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		pub := &mockQueuePublisher{}
		originalQueue := newQueuePublisher
		newQueuePublisher = func() ports.IQueuePublisher { return pub }
		t.Cleanup(func() { newQueuePublisher = originalQueue })

		sendRollbackEvent(context.Background(), "order-xyz")

		assert.True(t, pub.calledTopic)
		assert.Contains(t, string(pub.publishedBody), "order-xyz")
	})

	t.Run("publish error still handled", func(t *testing.T) {
		cleanup := env.SetupTestEnv()
		defer cleanup()

		pub := &mockQueuePublisher{publishErr: errors.New("sns error")}
		originalQueue := newQueuePublisher
		newQueuePublisher = func() ports.IQueuePublisher { return pub }
		t.Cleanup(func() { newQueuePublisher = originalQueue })

		sendRollbackEvent(context.Background(), "order-fail")

		assert.True(t, pub.calledTopic)
		assert.Contains(t, string(pub.publishedBody), "order-fail")
	})
}
