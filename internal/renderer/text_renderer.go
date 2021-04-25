package renderer

import (
	"errors"
	"io"

	"github.com/gihyeonsung/pott/internal/model"
)

func render(w io.Writer, doc *model.TextDocument) error {
	w.Write([]byte("<pre>"))
	if _, err := w.Write([]byte(doc.Content)); err != nil {
		return errors.New("could not write the text document: " + err.Error())
	}
	w.Write([]byte("</pre>"))
	return nil
}
