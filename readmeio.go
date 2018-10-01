package main

import (
	"bytes"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	bf "gopkg.in/russross/blackfriday.v2"
)

// Renderer is the rendering interface for output
type Renderer struct {
	Base bf.Renderer

	w bytes.Buffer

	Flags Flag

	lastOutputLen int
}

// Flag control optional behavior of this renderer
type Flag int

const (
	// FlagsNone does not allow customizing this renderer's behavior
	FlagsNone Flag = 0

	// InformationMacros allow using info, tip, note, and warning macros
	InformationMacros Flag = 1 << iota
)

var (
	linkTitleTag       = []byte("[")
	linkTitleCloseTag  = []byte("]")
	linkDataTag        = []byte("(")
	linkDataCloseTag   = []byte(")")
	linkInternalSymbol = []byte("doc:")
	h1Tag              = []byte("#")
	h2Tag              = []byte("#")
	h3Tag              = []byte("##")
	h4Tag              = []byte("###")
	h5Tag              = []byte("####")
	h6Tag              = []byte("#####")
)

var (
	nlBytes    = []byte{'\n'}
	spaceBytes = []byte(" ")
)

func (r *Renderer) cr(w io.Writer) {
	if r.lastOutputLen > 0 {
		r.out(w, nlBytes)
		r.out(w, nlBytes)
	}
}

func (r *Renderer) out(w io.Writer, text []byte) {
	w.Write(text)
	r.lastOutputLen = len(text)
}

func headingTagFromLevel(level int) []byte {
	switch level {
	case 1:
		return h1Tag
	case 2:
		return h1Tag
	case 3:
		return h2Tag
	case 4:
		return h3Tag
	case 5:
		return h4Tag
	default:
		return h5Tag
	}
}

// RenderNode is a renderer of a single node of a syntax tree.
func (r *Renderer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	// case bf.Heading:
	// 	headingTag := headingTagFromLevel(node.Level)
	// 	if entering {
	// 		r.out(w, headingTag)
	// 		w.Write(spaceBytes)
	// 	} else {
	// 		r.cr(w)
	// 	}
	// 	r.Base.RenderNode(w, node, entering)
	case bf.Link:
		if entering {
			r.out(w, linkTitleTag)
		} else {
			r.out(w, node.LinkData.Title)
			r.out(w, linkTitleCloseTag)
			r.out(w, linkDataTag)
			u, _ := url.Parse(string(node.LinkData.Destination))
			if u.Scheme == "" {
				r.out(w, linkInternalSymbol)
				dest := string(node.LinkData.Destination)
				r.out(w, []byte(strings.TrimSuffix(dest, filepath.Ext(dest))))
			} else {
				r.out(w, node.LinkData.Destination)
			}
			r.out(w, linkDataCloseTag)
		}
	// case bf.CodeBlock:
	default:
		return r.Base.RenderNode(w, node, entering)
	}
	return bf.GoToNext
}

// Render prints out the whole document from the ast
func (r *Renderer) Render(ast *bf.Node) []byte {
	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		return r.RenderNode(&r.w, node, entering)
	})

	return r.w.Bytes()
}

// RenderHeader writes document header
func (r *Renderer) RenderHeader(w io.Writer, ast *bf.Node) {
	r.Base.RenderHeader(w, ast)
}

// RenderFooter writes document footer
func (r *Renderer) RenderFooter(w io.Writer, ast *bf.Node) {
	r.Base.RenderFooter(w, ast)
}

// Render renders the README.io flavored markdown
func Render(input []byte, opts ...bf.Option) []byte {
	r := &Renderer{
		Flags: InformationMacros,
		Base: bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.SkipHTML,
		}),
	}
	optList := []bf.Option{bf.WithRenderer(r), bf.WithExtensions(bf.CommonExtensions)}
	optList = append(optList, opts...)
	parser := bf.New(optList...)
	ast := parser.Parse([]byte(input))
	return r.Render(ast)
}
