package renderer

import (
	"errors"
	"io"

	"github.com/gihyeonsung/pott/internal/model"
)

func Render(w io.Writer, doc *model.TextDocument) error {
	if _, err := w.Write([]byte(doc.Content)); err != nil {
		return errors.New("could not write the text document: " + err.Error())
	}
	return nil
}
