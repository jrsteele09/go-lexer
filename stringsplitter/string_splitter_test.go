package stringssplitter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommaSeparatedSplitWithQuotes(t *testing.T) {
	splitter := StringSplitter(",", "'", "\"")
	list := splitter("Orange, 'Apple   , Banana', \"Cherry\"")
	require.Equal(t, 3, len(list))
	require.ElementsMatch(t, []string{"Orange", "'Apple   , Banana'", "\"Cherry\""}, list)
}

func TestSplitter(t *testing.T) {
	splitter := StringSplitter("\r\n", "'", "\"")
	list := splitter("Hello\r\n'World \"this\" is a test'\r\nsome\r\nmore\r\narguments")
	require.Equal(t, 5, len(list))
	require.Equal(t, "Hello", list[0])
	require.Equal(t, "'World \"this\" is a test'", list[1])
	require.Equal(t, "some", list[2])
}

func TestSingleWordSplit(t *testing.T) {
	splitter := StringSplitter(" ", "'", "\"")
	list := splitter("Test ")
	require.Equal(t, 1, len(list))
	require.Equal(t, "Test", list[0])
}

func TestNoWordsToSplit(t *testing.T) {
	splitter := StringSplitter(" ", "'", "\"")
	list := splitter("")
	require.Equal(t, 0, len(list))
}

func TestSplittingTwoWords(t *testing.T) {
	splitter := StringSplitter(" ", "'", "\"")
	list := splitter("print (a+b)")
	require.Equal(t, 2, len(list))
}
