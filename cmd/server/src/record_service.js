// @ts-ignore
import { server } from "../pkg/k6/server";
// @ts-ignore
import { db } from "../pkg/k6/db";
//import schemaData from "./schema/schema.json";
export class RecordService {
    constructor() {
        this.insert_handler = (ictx) => {
            let tenantId = ictx.params().get("tenantId");
            server.doQuery(ictx, tenantId, (ctx) => {
                return server.readJson(ictx, {});
            });
        };
        this.update_handler = (ictx) => {
            let tenantId = ictx.params().get("tenantId");
            server.doRequest(ictx, tenantId, (ctx) => {
                let data = server.readJson(ictx);
                console.log(data);
                ictx.json(data);
                return null;
            });
        };
        // 按id查找数据
        this.findById_handler = (ictx) => {
            let tenantId = ictx.params().get("tenantId");
            let id = ictx.params().get("id");
            server.doQueryOne(ictx, tenantId, (ctx) => {
                return this.records.findById(ctx, tenantId, id);
            });
        };
        this.findPaging_handler = (ictx) => {
            let tenantId = ictx.params().get("tenantId");
            server.doQuery(ictx, tenantId, (ctx) => {
                return this.records.findPaging(ctx, { tenantId: tenantId, pageSize: 10 });
            });
        };
        server.get("/tenants/{tenantId}/records/{id}", this.findById_handler);
        server.get("/tenants/{tenantId}/records", this.findPaging_handler);
        server.post("/tenants/{tenantId}/records", this.insert_handler);
        server.put("/tenants/{tenantId}/records", this.update_handler);
        this.records = db.model("record");
    }
}
