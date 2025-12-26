package message

import "encoding/json"

type PaymentMessage struct {
	OrderID string `json:"order_id"`
}

func NewPaymentMessageFromJSON(messageJSON []byte) (*PaymentMessage, error) {
	var msg PaymentMessage

	err := json.Unmarshal(messageJSON, &msg)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}
