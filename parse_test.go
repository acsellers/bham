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

func TestParseRange(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%body\n\t\t= range .Wats\n\t\t\t%p wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Wats": []int{1, 2}})
		test.AreEqual("<html><body><p>wat</p><p>wat</p></body></html>", b.String())

		b.Reset()
		t.Execute(b, map[string]interface{}{"Wats": []int{}})
		test.AreEqual("<html><body></body></html>", b.String())
	})
}

func TestParseRangeElse(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "%html\n\t%body\n\t\t= range .Wats\n\t\t\t%p wat\n\t\t= else\n\t\t\t%p no wat")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Wats": []int{1, 2}})
		test.AreEqual("<html><body><p>wat</p><p>wat</p></body></html>", b.String())

		b.Reset()
		t.Execute(b, map[string]interface{}{"Wats": []int{}})
		test.AreEqual("<html><body><p>no wat</p></body></html>", b.String())
	})
}

func TestTextPassthrough(t *testing.T) {
	within(t, func(test *aTest) {
		t := template.New("test").Funcs(map[string]interface{}{})
		tree, err := Parse("test.bham", "<!DOCTYPE html>\n%html\n\t%body")
		test.IsNil(err)
		t, err = t.AddParseTree("tree", tree["test"])
		test.IsNil(err)

		b := new(bytes.Buffer)
		t.Execute(b, map[string]interface{}{"Wats": []int{1, 2}})
		test.AreEqual("<!DOCTYPE html><html><body></body></html>", b.String())
	})
}
