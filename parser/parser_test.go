package parser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []any
	}{
		{
			name:  "plain text",
			input: "Hello World",
			expected: []any{
				"Hello",
				" ",
				"World",
			},
		},
		{
			name:  "simple tag",
			input: "[b]Hello[/b]",
			expected: []any{
				Node{
					Tag:     "b",
					Attrs:   map[string]string{},
					Start:   Span{0, 3},
					End:     Span{8, 12},
					Content: []any{"Hello"},
				},
			},
		},
		{
			name:  "uppercase tag",
			input: "[B]Hello[/B]",
			expected: []any{
				Node{
					Tag:     "b",
					Attrs:   map[string]string{},
					Start:   Span{0, 3},
					End:     Span{8, 12},
					Content: []any{"Hello"},
				},
			},
		},
		{
			name:  "nested tags",
			input: "[b]Hello [i]world[/i][/b]",
			expected: []any{
				Node{
					Tag:   "b",
					Attrs: map[string]string{},
					Start: Span{0, 3},
					End:   Span{21, 25},
					Content: []any{
						"Hello",
						" ",
						Node{
							Tag:     "i",
							Attrs:   map[string]string{},
							Start:   Span{9, 12},
							End:     Span{17, 21},
							Content: []any{"world"},
						},
					},
				},
			},
		},
		{
			name:  "tag with default attr",
			input: "[url=https://example.com]link[/url]",
			expected: []any{
				Node{
					Tag:     "url",
					Attrs:   map[string]string{"_default": "https://example.com"},
					Start:   Span{0, 25},
					End:     Span{29, 35},
					Content: []any{"link"},
				},
			},
		},
		{
			name:  "tag with named attr",
			input: "[best name=value]Foo Bar[/best]",
			expected: []any{
				Node{
					Tag:   "best",
					Attrs: map[string]string{"name": "value"},
					Start: Span{0, 17},
					End:   Span{24, 31},
					Content: []any{
						"Foo",
						" ",
						"Bar",
					},
				},
			},
		},
		{
			name:  "url with special chars",
			input: "[url=https://example.com?q=1&page=2]link[/url]",
			expected: []any{
				Node{
					Tag:     "url",
					Attrs:   map[string]string{"_default": "https://example.com?q=1&page=2"},
					Start:   Span{0, 36},
					End:     Span{40, 46},
					Content: []any{"link"},
				},
			},
		},
		{
			name:  "multiple tags without spaces",
			input: "[b]Tag1[/b][i]Tag2[/i]",
			expected: []any{
				Node{
					Tag:     "b",
					Attrs:   map[string]string{},
					Start:   Span{0, 3},
					End:     Span{7, 11},
					Content: []any{"Tag1"},
				},
				Node{
					Tag:     "i",
					Attrs:   map[string]string{},
					Start:   Span{11, 14},
					End:     Span{18, 22},
					Content: []any{"Tag2"},
				},
			},
		},
		{
			name:     "only closing tag",
			input:    "[/b]",
			expected: []any{"[/b]"},
		},
		{
			name:     "unclosed tag",
			input:    "[b]Hello",
			expected: []any{"[b]", "Hello"},
		},
		{
			name:  "text before and after tag",
			input: "Hello [b]world[/b] foo",
			expected: []any{
				"Hello",
				" ",
				Node{
					Tag:     "b",
					Attrs:   map[string]string{},
					Start:   Span{6, 9},
					End:     Span{14, 18},
					Content: []any{"world"},
				},
				" ",
				"foo",
			},
		},
		{
			name:  "empty tag",
			input: "[b][/b]",
			expected: []any{
				Node{
					Tag:   "b",
					Attrs: map[string]string{},
					Start: Span{0, 3},
					End:   Span{3, 7},
				},
			},
		},
		{
			name:     "mismatched closing tag",
			input:    "[b]Hello[/i]",
			expected: []any{"[b]", "Hello", "[/i]"},
		},
		{
			name:  "newline in content",
			input: "[b]Hello\nworld[/b]",
			expected: []any{
				Node{
					Tag:   "b",
					Attrs: map[string]string{},
					Start: Span{0, 3},
					End:   Span{14, 18},
					Content: []any{
						"Hello",
						"\n",
						"world",
					},
				},
			},
		},
		{
			name:  "deep nesting",
			input: "[b][i][u]deep[/u][/i][/b]",
			expected: []any{
				Node{
					Tag:   "b",
					Attrs: map[string]string{},
					Start: Span{0, 3},
					End:   Span{21, 25},
					Content: []any{
						Node{
							Tag:   "i",
							Attrs: map[string]string{},
							Start: Span{3, 6},
							End:   Span{17, 21},
							Content: []any{
								Node{
									Tag:     "u",
									Attrs:   map[string]string{},
									Start:   Span{6, 9},
									End:     Span{13, 17},
									Content: []any{"deep"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "nested unclosed tags",
			input:    "[b][i]Hello",
			expected: []any{"[b]", "[i]", "Hello"},
		},
		{
			name:     "only whitespace",
			input:    "   ",
			expected: []any{"   "},
		},
		{
			name:  "unclosed tag at end",
			input: "[quote]some[/quote][color=red]test[/color]sdfasdfasdf[url=xxx]xxx[/url][quote]xxxsdfasdf",
			expected: []any{
				Node{
					Tag:     "quote",
					Attrs:   map[string]string{},
					Start:   Span{0, 7},
					End:     Span{11, 19},
					Content: []any{"some"},
				},
				Node{
					Tag:     "color",
					Attrs:   map[string]string{"_default": "red"},
					Start:   Span{19, 30},
					End:     Span{34, 42},
					Content: []any{"test"},
				},
				"sdfasdfasdf",
				Node{
					Tag:     "url",
					Attrs:   map[string]string{"_default": "xxx"},
					Start:   Span{53, 62},
					End:     Span{65, 71},
					Content: []any{"xxx"},
				},
				"[quote]",
				"xxxsdfasdf",
			},
		},
		{
			name:  "outer unclosed tag does not swallow content",
			input: "[quote]xxxsdfasdf[quote]some[/quote][color=red]test[/color]sdfasdfasdf[url=xxx]xxx[/url]",
			expected: []any{
				"[quote]",
				"xxxsdfasdf",
				Node{
					Tag:     "quote",
					Attrs:   map[string]string{},
					Start:   Span{17, 24},
					End:     Span{28, 36},
					Content: []any{"some"},
				},
				Node{
					Tag:     "color",
					Attrs:   map[string]string{"_default": "red"},
					Start:   Span{36, 47},
					End:     Span{51, 59},
					Content: []any{"test"},
				},
				"sdfasdfasdf",
				Node{
					Tag:     "url",
					Attrs:   map[string]string{"_default": "xxx"},
					Start:   Span{70, 79},
					End:     Span{82, 88},
					Content: []any{"xxx"},
				},
			},
		},
		{
			name: "unclosed tag in middle with multiline content",
			input: `[quote]some[/quote][color=red]test[/color]
[quote]xxxsdfasdf
sdfasdfasdf

[url=xxx]xxx[/url]`,
			expected: []any{
				Node{
					Tag:     "quote",
					Attrs:   map[string]string{},
					Start:   Span{0, 7},
					End:     Span{11, 19},
					Content: []any{"some"},
				},
				Node{
					Tag:     "color",
					Attrs:   map[string]string{"_default": "red"},
					Start:   Span{19, 30},
					End:     Span{34, 42},
					Content: []any{"test"},
				},
				"\n",
				"[quote]",
				"xxxsdfasdf",
				"\n",
				"sdfasdfasdf",
				"\n",
				"\n",
				Node{
					Tag:     "url",
					Attrs:   map[string]string{"_default": "xxx"},
					Start:   Span{74, 83},
					End:     Span{86, 92},
					Content: []any{"xxx"},
				},
			},
		},
		{
			name:  "repeated unclosed same tag",
			input: `Hello World[u]Wrong underline[u] This is another text [u]and this, too[/u]`,
			expected: []any{
				"Hello",
				" ",
				"World",
				"[u]",
				"Wrong",
				" ",
				"underline",
				"[u]",
				" ",
				"This",
				" ",
				"is",
				" ",
				"another",
				" ",
				"text",
				" ",
				Node{
					Tag:     "u",
					Attrs:   map[string]string{},
					Start:   Span{54, 57},
					End:     Span{70, 74},
					Content: []any{"and", " ", "this,", " ", "too"},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes := New(tc.input).Parse()
			if !reflect.DeepEqual(nodes, tc.expected) {
				t.Fatalf("\nExpected: %#v\nGot:      %#v", tc.expected, nodes)
			}
		})
	}
}
