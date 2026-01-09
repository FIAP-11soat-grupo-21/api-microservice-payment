package repositories

import "testing"

func TestPaymentModelTableName(t *testing.T) {
	model := PaymentModel{}
	table := model.TableName()

	if table != "payments" {
		t.Fatalf("expected table name 'payments', got '%s'", table)
	}
}
