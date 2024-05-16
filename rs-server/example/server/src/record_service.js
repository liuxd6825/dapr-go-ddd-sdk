// @ts-ignore
import { server } from "../pkg/k6/server";
// @ts-ignore
import { db } from "../pkg/k6/db";
// @ts-ignore
import { recordSchemas } from "./schema/record.schema";
// @ts-ignore
import { newSchema } from "../pkg/k6/schema";
var RecordService = /** @class */ (function () {
    function RecordService() {
        var _this = this;
        // 插入数据
        this.create_handler = function (rctx, cmd) {
            var data = cmd.data;
            rctx.printf("cmd.commandId=%s;", cmd.commandId);
            rctx.printf("data.tenantId=%s;", data.tenantId, data.id, data.name);
            rctx.executor().doCommand(function (ctx) {
                return _this.records.create(ctx, data);
            }).setResponse();
        };
        // 更新数据
        this.update_handler = function (rctx, cmd) {
            rctx.executor().doCommand(function (ctx) {
                var err = _this.schema.validate(cmd.data);
                return { error: err };
            }).setResponse();
        };
        // 按id查找数据1
        this.findById_handler = function (rctx) {
            rctx.executor().doQueryOne(function (ctx) {
                return _this.records.findById(ctx, rctx.getTenantId(), rctx.getId());
            }).setResponse();
        };
        // 分页查询数据
        this.findPaging_handler = function (rctx) {
            rctx.executor().doQuery(function (ctx) {
                var data = _this.records.findPaging(ctx, rctx.getFindPaging());
                return { data: data, ok: true };
            }).setResponse();
        };
        this.records = db.model("record");
        this.schema = newSchema(recordSchemas.domain);
        this.createSchema = newSchema(recordSchemas.createCommand);
        server.get({ path: "/tenants/{tenantId}/records/{id}", handle: this.findById_handler });
        server.get({ path: "/tenants/{tenantId}/records", handle: this.findPaging_handler });
        server.post({ path: "/tenants/{tenantId}/records", handle: this.create_handler, schema: recordSchemas.createCommand });
        server.put({ path: "/tenants/{tenantId}/records", handle: this.update_handler });
    }
    return RecordService;
}());
export { RecordService };
