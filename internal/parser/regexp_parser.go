package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/gihyeonsung/pott/internal/model"
)

func Parse(r io.Reader) (*model.TextDocument, error) {
	var content []byte
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New("could not read the text document: " + err.Error())
	}

	h3, err := regexp.Compile(`(=== [^=]* ===)`)
	h2, err := regexp.Compile(`(== [^=]* ==)`)
	h1, err := regexp.Compile(`(= [^=]* =)`)
	ref, err := regexp.Compile(`(\[\[[^]]*\]\])`)
	b, err := regexp.Compile(`(\*\*[^\*]*\*\*)`)
	i, err := regexp.Compile(`(//[^/]*//)`)
	u, err := regexp.Compile("(__[^_]*__)")
	s, err := regexp.Compile("(--[^-]*--)")
	pre, err := regexp.Compile("(``[^`]*``)")

	content = h3.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<h3>%s</h3>", t[3:len(t)-3])
		return []byte(t)
	})
	content = h2.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<h2>%s</h2>", t[2:len(t)-3])
		return []byte(t)
	})
	content = h1.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<h1>%s</h1>", t[1:len(t)-1])
		return []byte(t)
	})
	content = ref.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<a href=\"%s\">%s</a>", t[2:len(t)-2], t[2:len(t)-2])
		return []byte(t)
	})
	content = b.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<b>%s</b>", t[2:len(t)-2])
		return []byte(t)
	})
	content = i.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<i>%s</i>", t[2:len(t)-2])
		return []byte(t)
	})
	content = u.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<u>%s</u>", t[2:len(t)-2])
		return []byte(t)
	})
	content = s.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<s>%s</s>", t[2:len(t)-2])
		return []byte(t)
	})
	content = pre.ReplaceAllFunc(content, func(m []byte) []byte {
		t := string(m)
		t = fmt.Sprintf("<pre>%s</pre>", t[2:len(t)-2])
		return []byte(t)
	})
	return &model.TextDocument{Content: string(content)}, nil
}
