// @ts-ignore
import {Command, Executor, IContext, Params, RContext, server, asObject} from "../pkg/k6/server";
// @ts-ignore
import {db, Model} from "../pkg/k6/db";
// @ts-ignore
import {Context, GoError, Result} from "../pkg/k6/common";
// @ts-ignore
import {schemas as recordSchemas, IRecord, schemas} from "./schema/record.schema";
// @ts-ignore
import {newSchema, Property, Schema} from "../pkg/k6/schema";

interface CreateFields {
    tenantId:string;
    id:string;
    name:string;
}
interface CreateCommand {
    commandId:string;
    type:string;
    isValidOnly:boolean;
    data:CreateFields;
}

export class RecordService {
    // 数据表
    records:Model<IRecord>;
    schema:Schema;
    createSchema:Schema;
    constructor() {
        this.records = db.model("record");
        this.schema = newSchema(recordSchemas.domain);
        this.createSchema = newSchema(recordSchemas.createCommand)
        server.get({path: "/tenants/{tenantId}/records/{id}", handle: this.findById_handler});
        server.get({path:"/tenants/{tenantId}/records", handle: this.findPaging_handler});
        server.post({path:"/tenants/{tenantId}/records", handle: this.create_handler, schema: recordSchemas.createCommand});
        server.put({path:"/tenants/{tenantId}/records", handle: this.update_handler});
    }

    // 插入数据
    create_handler=(rctx:RContext, cmd:CreateCommand)=>{
        let data:any= cmd.data
        rctx.printf("cmd.commandId=%s;", cmd.commandId);
        rctx.printf("data.tenantId=%s;", data.tenantId, data.id, data.name);
        rctx.executor().doCommand((ctx:Context):GoError=>{
            return this.records.create(ctx, data);
        }).setResponse();
    }

    // 更新数据
    update_handler=(rctx:RContext, cmd:CreateCommand)=> {
        rctx.executor().doCommand((ctx:Context):Result<any>=>{
            let err = this.schema.validate(cmd.data)
            return {error:err}
        }).setResponse();
    }
    // 按id查找数据1

    findById_handler=(rctx:RContext)=>{
        rctx.executor().doQueryOne((ctx:Context):Result<any>=>{
            return this.records.findById(ctx, rctx.getTenantId(), rctx.getId());
        }).setResponse()
    }

    // 分页查询数据
    findPaging_handler=(rctx:RContext)=>{
        rctx.executor().doQuery((ctx:Context):Result=>{
            let data = this.records.findPaging(ctx, rctx.getFindPaging());
            return {data:data, ok:true};
        }).setResponse()
    }

}
