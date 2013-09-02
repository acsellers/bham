package bham

import (
	"bytes"
	"html/template"
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
