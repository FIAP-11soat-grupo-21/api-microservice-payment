package message

import "testing"

func TestNewCreatePaymentMessageFromJSON(t *testing.T) {
	validPayload := []byte(`{"order_id":"123","amount":42.5}`)
	msg, err := NewCreatePaymentMessageFromJSON(validPayload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msg == nil {
		t.Fatalf("expected message instance, got nil")
	}
	if msg.OrderID != "123" {
		t.Errorf("expected order_id 123, got %s", msg.OrderID)
	}
	if msg.Amount != 42.5 {
		t.Errorf("expected amount 42.5, got %f", msg.Amount)
	}

	invalidPayload := []byte(`{"order_id"`)
	msg, err = NewCreatePaymentMessageFromJSON(invalidPayload)
	if err == nil {
		t.Fatalf("expected error for invalid payload")
	}
	if msg != nil {
		t.Fatalf("expected nil message for invalid payload")
	}
}

func TestNewRollbackPaymentMessageFromJSON(t *testing.T) {
	validPayload := []byte(`{"order_id":"XYZ"}`)
	msg, err := NewRollbackPaymentMessageFromJSON(validPayload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msg == nil {
		t.Fatalf("expected message instance, got nil")
	}
	if msg.OrderID != "XYZ" {
		t.Errorf("expected order_id XYZ, got %s", msg.OrderID)
	}

	invalidPayload := []byte(`{"order_id"`)
	msg, err = NewRollbackPaymentMessageFromJSON(invalidPayload)
	if err == nil {
		t.Fatalf("expected error for invalid payload")
	}
	if msg != nil {
		t.Fatalf("expected nil message for invalid payload")
	}
}
