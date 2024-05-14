// @ts-ignore
import {Command, Executor, IContext, Params, RContext, server} from "../pkg/k6/server";
// @ts-ignore
import {db, Model} from "../pkg/k6/db";
// @ts-ignore
import {Context, GoError, Result} from "../pkg/k6/common";
// @ts-ignore
import {recordSchemas, IRecord} from "./schema/record.schema";
// @ts-ignore
import {newSchema, Property, Schema} from "../pkg/k6/schema";

interface CreateFields {
    tenantId:string;
    id:string;
    name:string;
}
interface CreateCommand {
    commandId:string;
    typeName:string;
    isValidOnly:boolean;
    data: CreateFields;
    toString():string;
}
export class RecordService {
    // 数据表
    records:Model<IRecord>;
    schema:Schema;

    constructor() {
        this.records = db.model("record");
        this.schema = newSchema(recordSchemas.domain);
        server.get({path: "/tenants/{tenantId}/records/{id}", handle: this.findById_handler});
        server.get({path:"/tenants/{tenantId}/records", handle: this.findPaging_handler});
        server.post({path:"/tenants/{tenantId}/records", handle: this.create_handler, schema: recordSchemas.createCommand});
        server.put({path:"/tenants/{tenantId}/records", handle: this.update_handler});
    }
    // 插入数据
    create_handler=(rctx:RContext, cmd:CreateCommand)=>{
        console.log(cmd.toString())
    }
    // 更新数据
    update_handler=(rctx:RContext)=> {
        rctx.executor().doCommand((ctx:Context, cmd:Command):Result<any>=>{
            let err = this.schema.validate(cmd.data)
            return {error:err}
        }).setResponse();
    }
    // 按id查找数据
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
