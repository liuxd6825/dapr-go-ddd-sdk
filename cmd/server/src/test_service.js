import { doQuery, doRequest, readJson, server, time } from "../k6/server";
export class TestService {
    constructor() {
        server.get("/test/{tenantId}/get/{id}", this.get);
        server.post("/test/{tenantId}/post", this.post);
    }
    get(ictx) {
        let tenantId = ictx.params().get("tenantId");
        let id = ictx.params().get("id");
        doQuery(ictx, tenantId, (ctx) => {
            let data = {
                tenantId: tenantId,
                id: id,
                name: '000sss',
                product: "0001",
                data: "data1",
                time: time.now(),
            };
            return { data: data, isFound: true };
        });
    }
    post(ictx) {
        let tenantId = ictx.params().get("tenantId");
        doRequest(ictx, tenantId, (ctx) => {
            let data = readJson(ictx);
            console.log(data);
            ictx.json(data);
            return null;
        });
    }
}
