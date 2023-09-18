package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	a := Parse("10*(2*5)")
	require.Equal(t, 100, a)
}
