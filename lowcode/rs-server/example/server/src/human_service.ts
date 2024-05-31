// @ts-ignore
import {db, Model} from "./pkg/k6/db";
// @ts-ignore
import {Context, GoError, Result, context, logs} from "./pkg/k6/common";
// @ts-ignore
import {newSchema, Property, Schema} from "./pkg/k6/schema";
// @ts-ignore
import {server, WebContext} from "./pkg/k6/server";

import {get, put, post, del} from "./pkg/k6server/decorator";
import {schemas as humanSchemas, Human} from "./schema/human.schema";

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

export class HumanService {
    // 数据表
    humans:Model<Human>;
    schema:Schema;
    createSchema:Schema;
    constructor() {
        this.schema = newSchema(humanSchemas.default);
        this.createSchema = newSchema(humanSchemas.createCommand);

        let ctx = context.background();
        this.humans = db.model("human");
        this.humans.table(ctx, this.schema).create(ctx)
    }

    // @ts-ignore
    @post("/tenants/{tenantId}/human", humanSchemas.createCommand)
    create_handler=(rctx:WebContext, cmd:CreateCommand)=>{
        let data:any= cmd.data
        rctx.executor().doCommand((ctx:Context):GoError=>{
            return this.humans.create(ctx, data);
        }).setResponse();
    }

    // @ts-ignore
    @put("/tenants/{tenantId}/human")
    update_handler=(rctx:WebContext, cmd:CreateCommand)=> {
        rctx.executor().doCommand((ctx:Context):GoError=>{
            let err = this.schema.validate(cmd.data)
            return err
        }).setResponse();
    }

    // 按id查找数据
    // @ts-ignore
    @get("/tenants/{tenantId}/human/{id}")
    findById_handler=(rctx:WebContext)=>{
        rctx.executor().doQueryOne((ctx:Context):Result<any>=>{
            return this.humans.findById(ctx, rctx.getTenantId(), rctx.getId());
        }).setResponse()
    }


    // 分页查询数据
    // @ts-ignore
    @get("/tenants/{tenantId}/human")
    findPaging_handler=(rctx:WebContext)=>{
        rctx.executor().doQuery((ctx:Context):Result<any>=>{
            let data = this.humans.findPaging(ctx, rctx.getFindPaging());
            return {data:data, error:null};
        }).setResponse()
    }

}

