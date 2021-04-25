package parser

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/gihyeonsung/pott/internal/model"
)

func Parse(r io.Reader) (*model.TextDocument, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New("could not read the text document: " + err.Error())
	}
	return &model.TextDocument{Content: string(content)}, nil
}
