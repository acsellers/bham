package bham

import (
	"bytes"
	"testing"
	"text/template"
)

func TestParse(t *testing.T) {
	within(t, func(test *aTest) {
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

func TestParse2(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%head\n\t\t%title\n\t\t\twat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, nil)
		test.AreEqual("<html><head><title>wat</title></head></html>", b.String())
	})
}

func TestParseIf(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%head\n\t\t= if .ShowWat\n\t\t\t%title wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"ShowWat": true})
		test.AreEqual("<html><head><title>wat</title></head></html>", b.String())

		b.Reset()
		t.Execute(b, map[string]interface{}{"ShowWat": false})
		test.AreEqual("<html><head></head></html>", b.String())
	})
}

func TestParseIfElse(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%head\n\t\t= if .ShowWat\n\t\t\t%title wat\n\t\t= else\n\t\t\t%title taw")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"ShowWat": true})
		test.AreEqual("<html><head><title>wat</title></head></html>", b.String())

		b.Reset()
		t.Execute(b, map[string]interface{}{"ShowWat": false})
		test.AreEqual("<html><head><title>taw</title></head></html>", b.String())
	})
}
