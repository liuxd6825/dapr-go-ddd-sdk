package template

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHtmlParse(t *testing.T) {
	fileService := file.NewLoader("file", "/Users/lxd/Projects/gitea/draw-web")
	if err := fileService.Connect(); err != nil {
		t.Error(err)
	}

	parse := NewParse(fileService, "src/master/schema/human/form.html", []byte(content))
	res, err := parse.Parse()
	assert.NoError(t, err)
	t.Log(res.HTML)
}

var content = `
<html>
    <head>
        <link type="application/json" id="schema" src="/schema/master/human/human.schema.json"></link>
        <link type="application/json" id="uischema" src="/schema/master/human/human.uischema.json"></link>
        <link type="application/html" id="formtpl" src="/schema/system/tpl/form/form.html"></link>
        <script type="text/javascript">
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            function init() {
                const vdata = { id: "001", name: "lxd", account: "6800001" }
                return vdata
            }
        </script>
    </head>
    <body>
        <script>
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            const data = {
                $schema: "https://json-schema.org/draft/2020-12/schema"
            }
        </script>
        <div>
            {{repo}}<br>
            {{ref}}<br>
        </div>
        <div>id:<ui5-input value="{{vdata.id}}"></ui5-input></div>
        <div>name:<ui5-input value="{{vdata.name}}"></ui5-input></div>
        <div>account:<ui5-input value="{{vdata.account}}"></ui5-input></div>
		<hyk-slot></hyk-slot>
		<ui5-button>按钮5</ui5-button>
    </body>    
</html>
`
