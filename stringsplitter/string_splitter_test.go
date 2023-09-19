package stringssplitter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	commandSplitter := StringSplitter("\r\n", "'", "\"")
	commands := commandSplitter("Hello\r\n'World \"this\" is a test'\r\nsome\r\nmore\r\narguments")
	require.Equal(t, 5, len(commands))
	require.Equal(t, "Hello", commands[0])
	require.Equal(t, "'World \"this\" is a test'", commands[1])
	require.Equal(t, "some", commands[2])
}

func TestSingleCommand(t *testing.T) {
	commandSplitter := StringSplitter(" ", "'", "\"")
	commands := commandSplitter("Test ")
	require.Equal(t, 1, len(commands))
	require.Equal(t, "Test", commands[0])
}

func TestNoCommands(t *testing.T) {
	commandSplitter := StringSplitter(" ", "'", "\"")
	commands := commandSplitter("")
	require.Equal(t, 0, len(commands))
}

func TestSplittingABasicCommand(t *testing.T) {
	commandSplitter := StringSplitter(" ", "'", "\"")
	commands := commandSplitter("print (a+b)")
	require.Equal(t, 2, len(commands))
}
