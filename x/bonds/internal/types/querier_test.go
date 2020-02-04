package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupQuerierTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	return func(t *testing.T) {
		t.Log("teardown test case")
	}
}

func TestStringMessageForQueryBondIsAsExpected(t *testing.T) {
	teardown := setupQuerierTestCase(t)
	defer teardown(t)

	b := QueryBonds{"Education", "Government", "Market"}
	expectedResult := fmt.Sprintf("%s\n%s\n%s", b[0], b[1], b[2])

	require.Equal(t, expectedResult, b.String())
}
