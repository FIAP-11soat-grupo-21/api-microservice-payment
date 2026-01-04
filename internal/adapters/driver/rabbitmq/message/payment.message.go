package message

import "encoding/json"

type CreatePaymentMessage struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func NewCreatePaymentMessageFromJSON(messageJSON []byte) (*CreatePaymentMessage, error) {
	var msg CreatePaymentMessage

	err := json.Unmarshal(messageJSON, &msg)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

type RollbackPaymentMessage struct {
	OrderID string `json:"order_id"`
}

func NewRollbackPaymentMessageFromJSON(messageJSON []byte) (*RollbackPaymentMessage, error) {
	var msg RollbackPaymentMessage

	err := json.Unmarshal(messageJSON, &msg)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}
