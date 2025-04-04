// @ts-ignore
import {Context} from "../../../rsrc/pkg/k6/common.d.ts";
// @ts-ignore
import {feign, FeignOptions} from "../../definition/types/feign";
// @ts-ignore
import paramsUtils from "../../definition/types/params";

let baseUrl:string = "dapr://duxm-master-cmd-service/api/v1.0/human";

export class HumanFeign{
    findPaging (ctx: Context, params: { tenantId: string; pageNum?: number; pageSize?: number; filter?: string; fields?: string; isTotal?: boolean; }): any {
        console.log("HumanFeign.findPaging()");
    }

    create (ctx:Context, params:{}):any {
        console.log("HumanFeign.create()");
    }

}
export default  new HumanFeign();


