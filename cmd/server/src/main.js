// @ts-ignore
import { TestService } from "./test_service";
// @ts-ignore
import { db } from "../pkg/k6/db";
export default function () {
    let err = db.open({
        appName: "jsServer",
        host: "192.168.65.5:27018,192.168.65.5:27019,192.168.65.5:27020",
        replicaSet: "mongors",
        dbName: "dev_master_cmd",
        user: "master",
        pwd: "123456",
        maxPoolSize: 20,
        operationTimeout: "300s",
        socketTimeout: "60s",
        maxConnIdleTime: "60s",
        serverSelectionTimeout: "20s"
    });
    if (err) {
        console.log(err.Error());
        return;
    }
    new TestService();
}
