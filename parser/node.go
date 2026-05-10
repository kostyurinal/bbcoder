package parser

type Node struct {
	Tag     string
	Attrs   map[string]string
	Content []any
	Start   Span
	End     Span
}

type Span struct {
	From int
	To   int
}
