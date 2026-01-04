package repositories

import (
	"context"
	"database/sql"
	"errors"
	"payment_microservice/internal/common/mocks"
	"payment_microservice/internal/core/domain/entities"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock, sqlDB
}

func TestNewGormPaymentDataSource(t *testing.T) {
	t.Run("should create GormPaymentDataSource successfully", func(t *testing.T) {
		gormDB, _, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		assert.NotNil(t, dataSource)
		assert.NotNil(t, dataSource.db)
	})
}

func TestGormPaymentDataSource_Insert(t *testing.T) {
	t.Run("should insert payment successfully", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		payment := mocks.GetValidPaymentEntity()
		ctx := context.Background()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "payments"`)).
			WithArgs(
				payment.ID,
				payment.OrderID,
				payment.Amount.Value(),
				payment.Status.Value(),
				payment.Method.Value(),
				sqlmock.AnyArg(),
				payment.PaidAt,
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dataSource.Insert(ctx, payment)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when insert fails", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		payment := mocks.GetValidPaymentEntity()
		ctx := context.Background()

		expectedError := errors.New("database error")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "payments"`)).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		err := dataSource.Insert(ctx, payment)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should insert payment without optional fields", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		qrCode := "00020101021243650016COM.MERCADOLIBRE"
		createdAt := time.Now()
		payment, _ := entities.NewPayment(
			"payment-id",
			"order-id",
			100.00,
			"pending",
			"pix",
			&qrCode,
			nil,
			createdAt,
		)
		ctx := context.Background()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "payments"`)).
			WithArgs(
				payment.ID,
				payment.OrderID,
				payment.Amount.Value(),
				payment.Status.Value(),
				payment.Method.Value(),
				sqlmock.AnyArg(), // qr_code_url
				nil,              // paid_at
				sqlmock.AnyArg(), // created_at
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := dataSource.Insert(ctx, *payment)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormPaymentDataSource_Update(t *testing.T) {
	t.Run("should update payment successfully", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		payment := mocks.GetValidPaymentEntity()
		ctx := context.Background()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "payments"`)).
			WithArgs(
				payment.OrderID,
				payment.Amount.Value(),
				payment.Status.Value(),
				payment.Method.Value(),
				sqlmock.AnyArg(),
				payment.PaidAt,
				sqlmock.AnyArg(),
				payment.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := dataSource.Update(ctx, payment)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		payment := mocks.GetValidPaymentEntity()
		ctx := context.Background()

		expectedError := errors.New("update failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "payments"`)).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		err := dataSource.Update(ctx, payment)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should update payment marking as paid", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		payment := mocks.GetPendingPaymentEntity()

		payment.MarkAsPaid()

		ctx := context.Background()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "payments"`)).
			WithArgs(
				payment.OrderID,
				payment.Amount.Value(),
				"paid",
				payment.Method.Value(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				payment.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := dataSource.Update(ctx, payment)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormPaymentDataSource_FindByOrderID(t *testing.T) {
	t.Run("should find payment by order ID successfully", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		modelData := mocks.GetPaymentModelData()
		ctx := context.Background()

		rows := sqlmock.NewRows([]string{
			"id", "order_id", "amount", "status", "payment_method",
			"qr_code_url", "paid_at", "created_at",
		}).AddRow(
			modelData.ID,
			modelData.OrderID,
			modelData.Amount,
			modelData.Status,
			modelData.PaymentMethod,
			modelData.QRCodeURL,
			modelData.PaidAt,
			modelData.CreatedAt,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE order_id = $1`)).
			WithArgs(modelData.OrderID, 1).
			WillReturnRows(rows)

		payment, err := dataSource.FindByOrderID(ctx, modelData.OrderID)
		assert.NoError(t, err)
		assert.Equal(t, modelData.ID, payment.ID)
		assert.Equal(t, modelData.OrderID, payment.OrderID)
		assert.Equal(t, modelData.Amount, payment.Amount.Value())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when payment not found", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		orderID := "non-existent-order"
		ctx := context.Background()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE order_id = $1`)).
			WithArgs(orderID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		payment, err := dataSource.FindByOrderID(ctx, orderID)
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Empty(t, payment.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when database fails", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		orderID := "order-123"
		ctx := context.Background()

		expectedError := errors.New("database connection error")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "payments" WHERE order_id = $1`)).
			WithArgs(orderID, 1).
			WillReturnError(expectedError)

		payment, err := dataSource.FindByOrderID(ctx, orderID)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, payment.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormPaymentDataSource_Delete(t *testing.T) {
	t.Run("should delete payment successfully", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		paymentID := "payment-to-delete"
		ctx := context.Background()

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "payments" WHERE id = $1`)).
			WithArgs(paymentID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := dataSource.Delete(ctx, paymentID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		paymentID := "payment-to-delete"
		ctx := context.Background()

		expectedError := errors.New("delete failed")

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "payments" WHERE id = $1`)).
			WithArgs(paymentID).
			WillReturnError(expectedError)
		mock.ExpectRollback()

		err := dataSource.Delete(ctx, paymentID)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should handle delete with context cancellation", func(t *testing.T) {
		gormDB, mock, sqlDB := setupMockDB(t)
		defer sqlDB.Close()

		dataSource := NewGormPaymentDataSource(gormDB)
		paymentID := "payment-to-delete"
		ctx, cancel := context.WithCancel(context.Background())

		cancel()
		err := dataSource.Delete(ctx, paymentID)

		assert.NotNil(t, err)
		_ = mock.ExpectationsWereMet()
	})
}
