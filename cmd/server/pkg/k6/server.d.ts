import {Context, Result} from "./common";


export interface IContext {
    writeString(data: string): void;
    setErr(err: Error): void;
    statusCode(val: Number): void;
    json(data: any): void;
    readJSON(data: any): void;
    readForm(data: any): void;
    header(name: string, value: string): void;
    getHeader(name: string): string;
    params(): RequestParams;
}

export interface Request {
    method: string;
    contentLength: bigint;
    host: string;
    remoteAddr: string;
    requestURI: string;
}

export interface RequestParams {
    get(key: string): string;
}

export interface Handler {
    (ctx: IContext): void;
}


export interface Time {
    now(): any;
}

type CmdFunc = (ctx: Context) => Error;
type QueryFunc = (ctx: Context) => Result;

export interface DoOptions {
    checkAuth?: boolean // 是否检查 Header Auth
}

export interface Command {

}

// CmdAndQueryOptions
// @Description: 命令执行参数
export interface CmdAndQueryOptions extends DoOptions {
    waitSecond: number; // 超时时间，单位秒
}

export interface Server {
    get(url: string, fun: Handler): void
    post(url: string, fun: Handler): void
    delete(url: string, fun: Handler): void
    put(url: string, fun: Handler): void;

    readJson: (ictx: IContext, data?: any)=> Result;
    doRequest:(ictx: IContext, tenantId: string, fun: (ctx: Context) => void, opts?: DoOptions)=>  Result;
    doQuery:(ictx: IContext, tenantId: string, fun: QueryFunc, opt?: DoOptions)=>  Result;
    doQueryOne:(ictx: IContext, tenantId: string, fun: QueryFunc, opt?: DoOptions)=>  Result;
    doCmdAndQueryOne:(ictx: IContext, tenantId: string, queryAppId: string, cmd: Command, cmdFun: CmdFunc, queryFun: QueryFunc, opts?: CmdAndQueryOptions)=> Result;
    doCmdAndQueryList:(ictx: IContext, tenantId: string, queryAppId: string, cmd: Command, cmdFun: CmdFunc, queryFun: QueryFunc, opts?: CmdAndQueryOptions)=> Result;
}

export const server: Server;
export const time: Time;

