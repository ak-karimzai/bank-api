package validators

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidFullName(t *testing.T) {
	fullName := "Ahmad Khalid Karimzai"
	err := ValidateFullName(fullName)
	require.NoError(t, err)
}
