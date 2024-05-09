// @ts-ignore
import { doQuery, doRequest, readJson, server } from "../pkg/k6/server";
// @ts-ignore
import { db } from "../pkg/k6/db";
export class TestService {
    constructor() {
        this.findById = (ictx) => {
            let tenantId = ictx.params().get("tenantId");
            let id = ictx.params().get("id");
            tenantId = "test";
            id = "bea191d9-546e-41bd-8335-b052b3fd136c";
            doQuery(ictx, tenantId, (ctx) => {
                let data = this.records.findById(ctx, tenantId, id);
                return data;
            });
        };
        this.findPaging = (ictx) => {
            let tenantId = "test";
            let records = db.model("record");
            doQuery(ictx, tenantId, (ctx) => {
                return records.findPaging(ctx, { tenantId: tenantId, pageSize: 10 });
            });
        };
        this.update = (ictx) => {
            let tenantId = ictx.params().get("tenantId");
            doRequest(ictx, tenantId, (ctx) => {
                let data = readJson(ictx);
                console.log(data);
                ictx.json(data);
                return null;
            });
        };
        server.get("/tenants/{tenantId}/tests/{id}", this.findById);
        server.get("/tenants/{tenantId}/tests/list", this.findPaging);
        server.post("/tenants/{tenantId}/tests", this.update);
        this.records = db.model("record");
        console.log("this.records:", this.records);
    }
}
