import {Context, doQuery, doRequest, error, IContext, readJson, ResultQuery, server, time} from "../k6/server";

interface Data {
    name: string;
    product: string;
    data: string;
    time: any;
    id: string;
    tenantId: string;
}

export class TestService {

    constructor() {
        server.get("/test/{tenantId}/get/{id}", this.get )
        server.post("/test/{tenantId}/post",  this.post )
    }

    get(ictx:IContext){
        let tenantId = ictx.params().get("tenantId")
        let id = ictx.params().get("id")
        doQuery(ictx, tenantId, (ctx:Context):ResultQuery=>{
            let data : Data = {
                tenantId: tenantId,
                id: id,
                name:'000sss',
                product:"0001",
                data:"data1",
                time: time.now(),
            }
            return {data:data, isFound:true}
        } )
    }

    post(ictx:IContext):void {
        let tenantId = ictx.params().get("tenantId")
        doRequest(ictx, tenantId,(ctx:Context):error=>{
            let data= readJson(ictx)
            console.log(data);
            ictx.json(data);
            return null
        })
    }

}
