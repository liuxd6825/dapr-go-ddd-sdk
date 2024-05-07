export interface Request {
    method: string;
    contentLength: bigint;
    host: string;
    remoteAddr: string;
    requestURI: string;
}

export interface IContext {
    writeString(data: string): void;
    setErr(err: error): void;
    statusCode(val: Number): void;
    json(data: any): void;
    readJSON(data: any): void;
    readForm(data: any): void;
    header(name: string, value: string): void;
    getHeader(name: string): string;
    params(): RequestParams;
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


export interface DB {
    Model(schema: Schema): Model;
}

export interface error {
    error(): string;
}

export type Record = any;

export interface Model {
    create(ctx: IContext, record: Record): error;
    update(ctx: IContext, field: string, val: any): error;
}

export interface Schema {

}

export interface Context {

}

export interface ResultQuery {
    data: any
    isFound?: boolean
    error?: error
}


type CmdFunc = (ctx: Context) => error;
type QueryFunc = (ctx: Context) => ResultQuery;

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
}

export function readJson(ictx: IContext): { data:any; error:error };
export function doRequest(ictx: IContext, tenantId: string, fun: (ctx: Context) => void, opts?: DoOptions): error;
export function doQuery(ictx: IContext, tenantId: string, fun: QueryFunc, opt?: DoOptions): ResultQuery;
export function doQueryOne(ictx: IContext, tenantId: string, fun: QueryFunc, opt?: DoOptions): ResultQuery;
export function doCmdAndQueryOne(ictx: IContext, tenantId: string, queryAppId: string, cmd: Command, cmdFun: CmdFunc, queryFun: QueryFunc, opts?: CmdAndQueryOptions):ResultQuery;
export function doCmdAndQueryList(ictx: IContext, tenantId: string, queryAppId: string, cmd: Command, cmdFun: CmdFunc, queryFun: QueryFunc, opts?: CmdAndQueryOptions):ResultQuery;

export const server: Server;
export const time: Time;
export const db: DB;

