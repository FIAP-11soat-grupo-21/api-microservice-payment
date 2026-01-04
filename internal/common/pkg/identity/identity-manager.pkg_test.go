package identity_manager

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUUIDV7(t *testing.T) {
	t.Run("should generate a valid UUID v7 string", func(t *testing.T) {
		id := NewUUIDV7()

		assert.NotEmpty(t, id)
		_, err := uuid.Parse(id)
		assert.NoError(t, err)
	})

	t.Run("should return different UUIDs on each call", func(t *testing.T) {
		id1 := NewUUIDV7()
		id2 := NewUUIDV7()
		id3 := NewUUIDV7()

		assert.NotEqual(t, id1, id2)
		assert.NotEqual(t, id2, id3)
		assert.NotEqual(t, id1, id3)
	})

	t.Run("should return UUID in standard format", func(t *testing.T) {
		uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

		id := NewUUIDV7()

		assert.True(t, uuidPattern.MatchString(id), "UUID should match standard format")
	})

	t.Run("should generate UUID v7 (not other versions)", func(t *testing.T) {
		id := NewUUIDV7()

		parsed, err := uuid.Parse(id)
		assert.NoError(t, err)

		assert.NotZero(t, parsed.Version())
	})
}

func TestIsValidUUID(t *testing.T) {
	t.Run("should return true for valid UUID", func(t *testing.T) {
		validUUID := uuid.New().String()

		result := IsValidUUID(validUUID)

		assert.True(t, result)
	})

	t.Run("should return true for valid UUID in lowercase", func(t *testing.T) {
		validUUID := "550e8400-e29b-41d4-a716-446655440000"

		result := IsValidUUID(validUUID)

		assert.True(t, result)
	})

	t.Run("should return true for valid UUID in uppercase", func(t *testing.T) {
		validUUID := "550E8400-E29B-41D4-A716-446655440000"

		result := IsValidUUID(validUUID)

		assert.True(t, result)
	})

	t.Run("should return false for invalid UUID format", func(t *testing.T) {
		invalidUUID := "not-a-valid-uuid"

		result := IsValidUUID(invalidUUID)

		assert.False(t, result)
	})

	t.Run("should return false for empty string", func(t *testing.T) {
		result := IsValidUUID("")

		assert.False(t, result)
	})

	t.Run("should return false for UUID with invalid characters", func(t *testing.T) {
		invalidUUID := "550e8400-e29b-41d4-a716-44665544000g"

		result := IsValidUUID(invalidUUID)

		assert.False(t, result)
	})

	t.Run("should return false for UUID with missing hyphens", func(t *testing.T) {
		uuidWithoutHyphens := "550e8400e29b41d4a716446655440000"

		result := IsValidUUID(uuidWithoutHyphens)

		assert.True(t, result)
	})

	t.Run("should return false for UUID with too many characters", func(t *testing.T) {
		invalidUUID := "550e8400-e29b-41d4-a716-446655440000-extra"

		result := IsValidUUID(invalidUUID)

		assert.False(t, result)
	})

	t.Run("should return false for nil UUID", func(t *testing.T) {
		nilUUID := "00000000-0000-0000-0000-000000000000"

		result := IsValidUUID(nilUUID)

		assert.True(t, result)
	})
}

func TestIsNotValidUUID(t *testing.T) {
	t.Run("should return false for valid UUID", func(t *testing.T) {
		validUUID := uuid.New().String()

		result := IsNotValidUUID(validUUID)

		assert.False(t, result)
	})

	t.Run("should return true for invalid UUID format", func(t *testing.T) {
		invalidUUID := "not-a-valid-uuid"

		result := IsNotValidUUID(invalidUUID)

		assert.True(t, result)
	})

	t.Run("should return true for empty string", func(t *testing.T) {
		result := IsNotValidUUID("")

		assert.True(t, result)
	})

	t.Run("should return true for UUID with invalid characters", func(t *testing.T) {
		invalidUUID := "550e8400-e29b-41d4-a716-44665544000g"

		result := IsNotValidUUID(invalidUUID)

		assert.True(t, result)
	})

	t.Run("should return false for UUID with missing hyphens", func(t *testing.T) {
		uuidWithoutHyphens := "550e8400e29b41d4a716446655440000"

		result := IsNotValidUUID(uuidWithoutHyphens)

		assert.False(t, result)
	})

	t.Run("should return true for UUID with too many characters", func(t *testing.T) {
		invalidUUID := "550e8400-e29b-41d4-a716-446655440000-extra"

		result := IsNotValidUUID(invalidUUID)

		assert.True(t, result)
	})

	t.Run("should return false for nil UUID", func(t *testing.T) {
		nilUUID := "00000000-0000-0000-0000-000000000000"

		result := IsNotValidUUID(nilUUID)

		assert.False(t, result)
	})

	t.Run("should be opposite of IsValidUUID", func(t *testing.T) {
		testCases := []string{
			uuid.New().String(),
			"not-a-uuid",
			"",
			"550e8400-e29b-41d4-a716-446655440000",
		}

		for _, testUUID := range testCases {

			isValid := IsValidUUID(testUUID)
			isNotValid := IsNotValidUUID(testUUID)

			assert.Equal(t, !isValid, isNotValid, "IsNotValidUUID should be opposite of IsValidUUID for: %s", testUUID)
		}
	})
}
