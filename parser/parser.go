package parser

import (
	"strings"
	"unicode"
)

type Parser struct {
	input []rune
	pos   int
}

func New(input string) *Parser {
	return &Parser{input: []rune(input)}
}

func (p *Parser) Parse() []any {
	var nodes []any
	for p.pos < len(p.input) {
		if p.peek() == '[' {
			node, raw := p.parseTag()
			if raw != "" {
				nodes = append(nodes, raw)
			} else {
				nodes = append(nodes, node)
			}
		} else {
			nodes = append(nodes, p.parseTextTokens()...)
		}
	}
	return nodes
}

func (p *Parser) peek() rune {
	return p.input[p.pos]
}

func (p *Parser) parseTextTokens() []any {
	var tokens []any
	for p.pos < len(p.input) && p.peek() != '[' {
		start := p.pos
		if p.peek() == '\n' {
			p.pos++
			tokens = append(tokens, string(p.input[start:p.pos]))
			continue
		}
		isSpace := unicode.IsSpace(p.peek())
		for p.pos < len(p.input) && p.peek() != '[' && p.peek() != '\n' && unicode.IsSpace(p.peek()) == isSpace {
			p.pos++
		}
		tokens = append(tokens, string(p.input[start:p.pos]))
	}
	return tokens
}

func (p *Parser) parseTag() (Node, string) {
	startPos := p.pos
	p.pos++ // skip '['

	if p.pos < len(p.input) && p.peek() == '/' {
		from := p.pos - 1
		p.pos++
		p.readUntil(']')
		p.pos++
		return Node{}, string(p.input[from:p.pos])
	}

	openFrom := p.pos - 1
	raw := p.readUntil(']')
	p.pos++ // skip ']'
	openTo := p.pos

	tag, attrs := parseAttrs(raw)
	content, end := p.parseChildren(tag)

	if end == (Span{}) {
		p.pos = openTo
		return Node{}, string(p.input[startPos:openTo])
	}

	return Node{
		Tag:     tag,
		Attrs:   attrs,
		Start:   Span{From: openFrom, To: openTo},
		End:     end,
		Content: content,
	}, ""
}

func (p *Parser) parseChildren(parentTag string) ([]any, Span) {
	var children []any
	for p.pos < len(p.input) {
		if p.isClosingTag(parentTag) {
			from := p.pos
			p.skipClosingTag()
			return children, Span{From: from, To: p.pos}
		}
		if p.peek() == '[' {
			node, raw := p.parseTag()
			if raw != "" {
				children = append(children, raw)
			} else {
				children = append(children, node)
			}
		} else {
			children = append(children, p.parseTextTokens()...)
		}
	}
	return children, Span{}
}

func (p *Parser) isClosingTag(tag string) bool {
	closing := "[/" + tag + "]"
	end := p.pos + len([]rune(closing))
	if end > len(p.input) {
		return false
	}
	return strings.EqualFold(string(p.input[p.pos:end]), closing)
}

func (p *Parser) skipClosingTag() {
	p.readUntil(']')
	p.pos++
}

func (p *Parser) readUntil(ch rune) string {
	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] != ch {
		p.pos++
	}
	return string(p.input[start:p.pos])
}
