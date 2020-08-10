package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckNoOfReserveTokensInvalidFunctionType(t *testing.T) {
	err := CheckNoOfReserveTokens(nil, "invalid_function_type")
	require.NotNil(t, err)
}

func TestGetExceptionsForFunctionTypeInvalidFunctionType(t *testing.T) {
	_, err := GetExceptionsForFunctionType("invalid_function_type")
	require.NotNil(t, err)
}
