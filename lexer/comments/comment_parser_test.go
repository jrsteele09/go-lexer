package comments_test

import (
	"testing"

	"github.com/jrsteele09/go-lexer/lexer/comments"
	"github.com/stretchr/testify/require"
)

var testComments = map[string]string{
	"//":  "\n",
	"/*":  "*/",
	"rem": "\n",
}

func TestCommentParser(t *testing.T) {
	cp := comments.NewCommentParser(testComments)

	require.False(t, cp.InComment())
	require.False(t, cp.IsStartOfComment("abc"))
	require.True(t, cp.IsStartOfComment("/*"))
	require.True(t, cp.InComment())

	testEndOfCommentString := "abc dd eee fff*/some more text"
	var commentEnd bool
	for _, r := range testEndOfCommentString {
		commentEnd = cp.ParseEndOfComment(r)
	}

	require.True(t, commentEnd)
	require.False(t, cp.InComment())
}
