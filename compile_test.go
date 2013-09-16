package bham

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/acsellers/assert"
)

func TestCompile1(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `!!!
:javascript
  $('#test').hide();`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.AreEqual(2, len(pt.outputTree.Root.Nodes))

		t := template.New("wat").Funcs(map[string]interface{}{})
		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", nil)
		test.AreEqual("<!DOCTYPE html><script type=\"text/javascript\">$('#test').hide();</script>", b.String())
	})
}

func TestCompile2(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `!!!
= .Name`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.AreEqual(2, len(pt.outputTree.Root.Nodes))

		t := template.New("wat").Funcs(map[string]interface{}{})
		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", map[string]interface{}{"Name": "Hello"})
		test.AreEqual("<!DOCTYPE html>Hello", b.String())
	})
}

func TestCompile3(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `= stylesheet "first" "second"`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.IsNil(pt.err)
		t := template.New("wat").Funcs(map[string]interface{}{
			"stylesheet": func(sheets ...string) template.HTML {
				var output []string
				for _, sheet := range sheets {
					output = append(output, fmt.Sprintf(
						`<link href="%s.css" rel="stylesheet">`,
						sheet,
					))
				}
				return template.HTML(strings.Join(output, "\n"))
			},
		})

		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", map[string]interface{}{"Name": "Hello"})
		expected := `<link href="first.css" rel="stylesheet">
<link href="second.css" rel="stylesheet">`
		test.AreEqual(expected, b.String())
	})
}

func TestCompile4(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `%head`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.IsNil(pt.err)
		test.AreEqual(2, len(pt.outputTree.Root.Nodes))
		t := template.New("wat").Funcs(map[string]interface{}{})
		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", map[string]interface{}{})
		test.AreEqual("<head></head>", b.String())

	})
}

func TestCompile5(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `.next holla`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.IsNil(pt.err)
		test.AreEqual(3, len(pt.outputTree.Root.Nodes))
		t := template.New("wat").Funcs(map[string]interface{}{})
		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", map[string]interface{}{})
		test.AreEqual(`<div class="next"> holla </div>`, b.String())

	})
}

func TestCompile6(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `#welcome Hello {{.Name}}`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.IsNil(pt.err)
		t := template.New("wat").Funcs(map[string]interface{}{})
		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", map[string]interface{}{
			"Name": "Human",
		})
		test.AreEqual(`<div id="welcome"> Hello Human</div>`, b.String())
	})
}

func TestCompile7(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `#welcome
  Hello
  = .Name`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.IsNil(pt.err)
		t := template.New("wat").Funcs(map[string]interface{}{})
		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", map[string]interface{}{
			"Name": "Human",
		})
		test.AreEqual(`<div id="welcome">Hello Human</div>`, b.String())
	})
}

func TestCompile8(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `= hello .Name`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.IsNil(pt.err)
		t := template.New("wat").Funcs(map[string]interface{}{
			"hello": func(s string) string {
				return "Hello " + s
			},
		})
		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", map[string]interface{}{
			"Name": "Computer",
		})
		test.AreEqual(`Hello Computer`, b.String())
	})
}

func TestCompile9(t *testing.T) {
	assert.Within(t, func(test *assert.Test) {
		tmpl := `.name#welcome
  Name
  Rank
  Serial Number`
		pt := &protoTree{name: "compile", source: tmpl}
		pt.lex()
		pt.analyze()
		pt.compile()
		test.IsNotNil(pt.outputTree)
		test.IsNil(pt.err)
		t := template.New("wat").Funcs(map[string]interface{}{})
		t.AddParseTree("compile", pt.outputTree)
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, "compile", map[string]interface{}{})
		test.AreEqual(`<div class="name" id="welcome">Name Rank Serial Number </div>`, b.String())
	})
}
