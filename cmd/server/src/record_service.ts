// @ts-ignore
import {Error,IContext, server} from "../pkg/k6/server";
// @ts-ignore
import {db, Model} from "../pkg/k6/db";
// @ts-ignore
import {Result, Context} from "../pkg/k6/common";
//import schemaData from "./schema/schema.json";

export class RecordService {
    // 数据表
    records:Model;

    constructor() {
        server.get( "/tenants/{tenantId}/records/{id}", this.findById_handler);
        server.get( "/tenants/{tenantId}/records", this.findPaging_handler);
        server.post("/tenants/{tenantId}/records", this.insert_handler);
        server.put("/tenants/{tenantId}/records", this.update_handler);
        this.records = db.model("record");
    }

    insert_handler=(ictx:IContext)=>{
        let tenantId = ictx.params().get("tenantId");
        server.doQuery(ictx, tenantId, (ctx:Context):Result=>{
            return server.readJson(ictx, {})
        })
    }

    update_handler=(ictx:IContext)=> {
        let tenantId = ictx.params().get("tenantId");
        server.doRequest(ictx, tenantId,(ctx:Context):Error=>{
            let data= server.readJson(ictx);
            console.log(data);
            ictx.json(data);
            return null;
        })
    }
    // 按id查找数据
    findById_handler=(ictx:IContext)=>{
        let tenantId = ictx.params().get("tenantId");
        let id = ictx.params().get("id");
        server.doQueryOne(ictx, tenantId, (ctx:Context):Result=>{
            return this.records.findById(ctx, tenantId, id)
        })
    }

    findPaging_handler=(ictx:IContext)=>{
        let tenantId = ictx.params().get("tenantId");
        server.doQuery(ictx, tenantId, (ctx:Context):Result=>{
            return this.records.findPaging(ctx, {tenantId:tenantId, pageSize: 10})
        })
    }

}
