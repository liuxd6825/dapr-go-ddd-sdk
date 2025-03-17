package template

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/localfs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestHtmlParse(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fs, err := localfs.NewFs(localfs.Config{Path: path})
	assert.NoError(t, err)

	parse := NewParse(fs, "/testfile/form.tpl", []byte(content))
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
