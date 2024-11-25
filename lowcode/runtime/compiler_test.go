package runtime

import "testing"

var tsCode = `
            interface Params {
                tenantId:string;
                id:string;
                data:{
                    id:string;
                }
            }
            const params:Parmas;
            let data = service.dao.findById(ctx, params.tenantId, params.id)
            wctx.setData(data)
`

func Test_Compiler(t *testing.T) {
	compiler := NewCompiler()

	jsCode, err := compiler.compileTypeScript("/Users/lxd/Projects/duxm/draw-web", tsCode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsCode))
}
