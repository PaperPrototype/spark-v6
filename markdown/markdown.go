package markdown

import (
	"bytes"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	xurls "mvdan.cc/xurls/v2"

	embed "github.com/13rac1/goldmark-embed"
)

var markdowner goldmark.Markdown = goldmark.New(
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithExtensions(
		highlighting.NewHighlighting(
			highlighting.WithStyle("native"),
			highlighting.WithFormatOptions(
				html.WithLineNumbers(false),
			),
		),
		extension.NewLinkify(
			extension.WithLinkifyAllowedProtocols([][]byte{
				[]byte("http:"),
				[]byte("https:"),
			}),
			extension.WithLinkifyURLRegexp(
				xurls.Strict(),
			),
		),
		embed.New(),
	),
)

func Convert(source []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	return &buf, markdowner.Convert(source, &buf)
}
