import {LogPkg} from "./logs";
import {DBPkg} from "./db";
import {TemplatePkg} from "./template";
import {SchemaPkg} from "./schema";
import {Context, ContextPkg} from "./context";
import {FeignPkg} from "./feign";
import {JsonPkg} from "./json";
import {FsPkg} from "./fs";
import {ParamsTypePkg} from "./paramsType";
import {StringsPKg} from "./strings";
import {HtmlPKg} from "./html";
import {GoError} from "./types";

export interface PKG {
    logs: LogPkg; // 日志处理
    db: DBPkg; // mongo数据库
    template: TemplatePkg; // 文件模版引擎
    schema: SchemaPkg; // schema数据格式
    context: ContextPkg; // go context上下文
    feign: FeignPkg; //rest ful api 调用包
    json: JsonPkg; // json处理包
    fs: FsPkg; // 文件系统
    paramsType: ParamsTypePkg; // Web参数处理包
    strings: StringsPKg; //字符串处理包
    html: HtmlPKg;
    events: Events;
    sys: SysPkg;
}

export interface SysPkg {
    codeService:CodeService;
}

export interface CodeService {

}

export interface Events {
    publishEvent(ctx: Context, dbKey:string, appId:string, data:any, meta:any):GoError;
    publishEvents(ctx: Context, dbKey:string, appId:string, data:any[], meta:any):GoError;
}




