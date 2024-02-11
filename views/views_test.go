package views

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	tmpl := getTemplates()

	for _, subTmpl := range tmpl.Templates() {
		fmt.Println(subTmpl.Name())
	}

	g := goldie.New(t,
		goldie.WithFixtureDir("testdata"),
		goldie.WithNameSuffix(".html"),
	)

	bodyHTML, err := ExecuteHTML("body.html", nil)
	assert.Equal(t, nil, err)

	var buf bytes.Buffer
	err = View(&buf, bodyHTML)
	assert.Equal(t, nil, err)
	g.Assert(t, "full", buf.Bytes())
}
