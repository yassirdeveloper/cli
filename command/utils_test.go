package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValue(t *testing.T) {
	t.Run("Valid Int Parsing", func(t *testing.T) {
		value, err := ParseValue(TypeInt, "42")
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("Invalid Int Parsing", func(t *testing.T) {
		_, err := ParseValue(TypeInt, "not-an-int")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})

	t.Run("Valid Float Parsing", func(t *testing.T) {
		value, err := ParseValue(TypeFloat, "3.14")
		assert.NoError(t, err)
		assert.Equal(t, 3.14, value)
	})

	t.Run("Invalid Float Parsing", func(t *testing.T) {
		_, err := ParseValue(TypeFloat, "not-a-float")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})

	t.Run("Valid Bool Parsing", func(t *testing.T) {
		value, err := ParseValue(TypeBool, "true")
		assert.NoError(t, err)
		assert.Equal(t, true, value)

		value, err = ParseValue(TypeBool, "false")
		assert.NoError(t, err)
		assert.Equal(t, false, value)
	})

	t.Run("Invalid Bool Parsing", func(t *testing.T) {
		_, err := ParseValue(TypeBool, "not-a-bool")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})

	t.Run("Valid String Parsing", func(t *testing.T) {
		value, err := ParseValue(TypeString, "hello")
		assert.NoError(t, err)
		assert.Equal(t, "hello", value)
	})

	t.Run("Unsupported Type", func(t *testing.T) {
		_, err := ParseValue(NoType, "value")
		assert.Error(t, err)
		assert.Equal(t, "unsupported type", err.Error())
	})

	t.Run("Non-String Input", func(t *testing.T) {
		_, err := ParseValue(TypeInt, 42) // Non-string input
		assert.Error(t, err)
		assert.Equal(t, "cant cast to string", err.Error())
	})
}
