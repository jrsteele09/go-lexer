package comments

type CommentParser struct {
	commentStart string
	commentEnd   string
	parsedEnd    string
	comments     map[string]string // Map of comment delimiters.
}

func NewCommentParser(comments map[string]string) *CommentParser {
	return &CommentParser{
		comments: comments,
	}
}

func (p *CommentParser) InComment() bool {
	return p.commentStart != "" && p.commentEnd != ""
}

func (p *CommentParser) IsStartOfComment(str string) bool {
	if p.InComment() {
		return false
	}
	commentEnd, found := p.comments[str]
	if found {
		p.commentStart = str
		p.commentEnd = commentEnd
		p.parsedEnd = ""
	}
	return found
}

func (p *CommentParser) ParseEndOfComment(r rune) bool {
	if p.commentEnd == "" {
		return true
	}

	ce := []rune(p.commentEnd)

	if ce[len(p.parsedEnd)] != r {
		p.parsedEnd = ""
		return false
	}

	p.parsedEnd += string(r)
	if p.parsedEnd != p.commentEnd {
		return false
	}

	p.Reset()
	return true
}

func (p *CommentParser) IsNewLineComment() bool {
	return p.commentEnd == "\n"
}

func (p *CommentParser) Reset() {
	p.commentStart = ""
	p.commentEnd = ""
	p.parsedEnd = ""
}
