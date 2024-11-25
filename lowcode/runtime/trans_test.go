package runtime

import "testing"

func Test_TransformTSCodeToJS(t *testing.T) {
	tsCode := `
			interface Params {
                tenantId:string;
                id:string;
                data:{
                    id:string;
                }
            }
            const params:Parmas;
            let data = service.dao.findById(ctx, params.tenantId, params.id);
            wctx.setData(data)
	`

	jsCode := TransformTSCodeToJS(tsCode)
	t.Log(jsCode)
}
