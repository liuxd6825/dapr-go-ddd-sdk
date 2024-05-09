// @ts-ignore
import {Context, doQuery, doRequest, error, IContext, readJson, ResultQuery, server, time} from "../pkg/k6/server";
// @ts-ignore
import {db, Model, FindPagingQueryRequest} from "../pkg/k6/db";
//import schemaData from "./schema/schema.json";

interface Data {
    name: string;
    product: string;
    data: string;
    time: any;
    id: string;
    tenantId: string;
}

export class TestService {
    // 数据表
    records:Model;

    constructor() {
        server.get( "/tenants/{tenantId}/tests/{id}", this.findById);
        server.get( "/tenants/{tenantId}/tests/list", this.findPaging);
        server.post("/tenants/{tenantId}/tests",      this.update);
        this.records = db.model("record");
        console.log("this.records:", this.records);
    }

    findById=(ictx:IContext)=>{
        let tenantId = ictx.params().get("tenantId");
        let id = ictx.params().get("id");
        tenantId="test";
        id="bea191d9-546e-41bd-8335-b052b3fd136c";
        doQuery(ictx, tenantId, (ctx:Context):ResultQuery=>{
            let data = this.records.findById(ctx, tenantId, id)
            return data
        })
    }

    findPaging=(ictx:IContext)=>{
        let tenantId="test";
        let records = db.model("record");
        doQuery(ictx, tenantId, (ctx:Context):ResultQuery=>{
           return records.findPaging(ctx, {tenantId:tenantId, pageSize: 10})
        })
    }

    update=(ictx:IContext)=> {
        let tenantId = ictx.params().get("tenantId");
        doRequest(ictx, tenantId,(ctx:Context):error=>{
            let data= readJson(ictx);
            console.log(data);
            ictx.json(data);
            return null
        })
    }

}
