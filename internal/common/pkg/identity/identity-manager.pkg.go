package identity_manager

import "github.com/google/uuid"

func NewUUIDV7() string {
	id := uuid.Must(uuid.NewV7())
	return id.String()
}

func IsValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}

func IsNotValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err != nil
}
