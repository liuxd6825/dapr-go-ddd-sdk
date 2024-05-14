// @ts-ignore
import { server } from "../pkg/k6/server";
// @ts-ignore
import { db } from "../pkg/k6/db";
// @ts-ignore
import { recordSchemas } from "./schema/record.schema";
// @ts-ignore
import { newSchema } from "../pkg/k6/schema";
export class RecordService {
    constructor() {
        // 插入数据
        this.create_handler = (rctx, cmd) => {
            console.log(cmd.toString());
        };
        // 更新数据
        this.update_handler = (rctx) => {
            rctx.executor().doCommand((ctx, cmd) => {
                let err = this.schema.validate(cmd.data);
                return { error: err };
            }).setResponse();
        };
        // 按id查找数据
        this.findById_handler = (rctx) => {
            rctx.executor().doQueryOne((ctx) => {
                return this.records.findById(ctx, rctx.getTenantId(), rctx.getId());
            }).setResponse();
        };
        // 分页查询数据
        this.findPaging_handler = (rctx) => {
            rctx.executor().doQuery((ctx) => {
                let data = this.records.findPaging(ctx, rctx.getFindPaging());
                return { data: data, ok: true };
            }).setResponse();
        };
        this.records = db.model("record");
        this.schema = newSchema(recordSchemas.domain);
        server.get({ path: "/tenants/{tenantId}/records/{id}", handle: this.findById_handler });
        server.get({ path: "/tenants/{tenantId}/records", handle: this.findPaging_handler });
        server.post({ path: "/tenants/{tenantId}/records", handle: this.create_handler, schema: recordSchemas.createCommand });
        server.put({ path: "/tenants/{tenantId}/records", handle: this.update_handler });
    }
}
