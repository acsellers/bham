package bham

import (
	"bytes"
	"testing"
	"text/template"
)

func TestSteps(t *testing.T) {
	Within(t, func(test *Test) {
		proto := &protoTree{
			source: `%html
  %head
    %title wat`}
		proto.tokenize()
		test.IsNil(proto.err)

		test.AreEqual(
			[]string{"<html>", "<head>", "<title>", "wat", "</title>", "</head>", "</html>"},
			proto.tokenList,
		)
		proto.classify()
		test.AreEqual(
			"<html><head><title>wat</title></head></html>",
			proto.classified[0].String(),
		)

		tree := proto.treeify()
		t, _ := template.New("test").Parse("{{define \"blank\"}}blank{{end}}")
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "blank", nil)
		test.AreEqual("blank", b.String())

		b.Reset()
		t, _ = t.AddParseTree("tree", tree)
		t.ExecuteTemplate(b, "blank", nil)
		test.AreEqual("blank", b.String())

		b.Reset()
		t.ExecuteTemplate(b, "tree", nil)
		test.AreEqual(
			"<html><head><title>wat</title></head></html>",
			b.String(),
		)

	})
}

func TestParse(t *testing.T) {
	Within(t, func(test *Test) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%head\n\t\t%title wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, nil)
		test.AreEqual("<html><head><title>wat</title></head></html>", b.String())
	})
}
