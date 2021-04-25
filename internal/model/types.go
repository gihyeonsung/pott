package model

type TextDocument struct {
	content string
}

func NewTextDocument(content string) *TextDocument {
	return &TextDocument{content: content}
}
