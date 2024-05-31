import {GoError, Context, Result, Object} from "./common";
import {Entity, FindPagingQuery} from "./db";
import {Schema} from "./schema";


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



export interface Time {
    now(): any;
}

type CmdFunc = (ctx: Context) => GoError;
type QueryFunc = (ctx: Context) => Result<any>;

export interface DoOptions {
    checkAuth?: boolean // 是否检查 Header Auth
}

export interface Command<T extends any> {
    commandId:string;
    typeName:string;
    isValidOnly:boolean;
    data:T;
}

// CmdAndQueryOptions
// @Description: 命令执行参数
export interface CmdAndQueryOptions extends DoOptions {
    waitSecond: number; // 超时时间，单位秒
}

export interface HandleOption  {
    method?:"GET"|"POST"|"PUT"|"DELETE"|"PATCH"|"HEAD"|"OPTIONS";
    path:string;
    handle:Handle;
    schema?:Schema;
    desc?: string; // 方法说明
}

export interface HandleDesc {
    desc:string;
    params: {name:string; dataType: string; desc:string; }[]
}

export type Handle = (ctx: WebContext, data?:any)=>void;

export interface Server {
    get(opts: HandleOption):void;
    post(opts: HandleOption):void;
    put(opts: HandleOption):void;
    delete(opts: HandleOption):void;
    patch(opts: HandleOption):void;
    handle(opt:HandleOption): void;
    handles(opts:HandleOption[]): void;

    doRequest:(wctx: WebContext, fun: (ctx: Context) => void, opts?: DoOptions)=> Result<any>;
    doQuery:(wctx: WebContext,  fun: QueryFunc, opt?: DoOptions)=> Result<any>;
    doQueryOne:(wctx: WebContext, fun: QueryFunc, opt?: DoOptions)=> Result<any>;
    doCmd:(wctx: WebContext,cmdFun: CmdFunc, opts?: CmdAndQueryOptions)=> Result<any>;
    doCmdAndQueryOne:(wctx: WebContext, queryAppId: string, cmd: Command<any>, cmdFun: CmdFunc, queryFun: QueryFunc, opts?: CmdAndQueryOptions)=> Result<any>;
    doCmdAndQueryList:(wctx: WebContext, queryAppId: string, cmd: Command<any>, cmdFun: CmdFunc, queryFun: QueryFunc, opts?: CmdAndQueryOptions)=> Result<any>;
}
// iris context
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
//  request context
export interface WebContext {
    params():Params;
    readJson(data?:any):Result<any>;
    readObject(schema:Schema):Result<Object>;
    writeJson(data:any):GoError;
    executor<T>():Executor<T>;
    getTenantId():string;
    getCaseId():string;
    getId():string;
    getFindPaging():FindPagingQuery;
    printf(a1?:any, a2?:any, a3?:any, a4?:any, a5?:any, a6?:any, a7?:any, a8?:any, a9?:any, a10?:any, a11?:any, a12?:any, a13?:any, a14?:any, a15?:any, a16?:any, a17?:any, a18?:any, a19?:any, a20?:any, a21?:any, a22?:any, a23?:any, a24?:any, a25?:any, a26?:any, a27?:any, a28?:any, a29?:any, a30?:any);

}

export interface Params {
    getId():string;
    getTenantId(): string;
    getCaseId(): string;
    getFindPaging(): FindPagingQuery;
    string(key:string): Result<string>;
    strings(key:string):Result<string[]>;
    bool(key:string): Result<boolean>;
    float64(key:string): Result<number>;
    int64(key:string): Result<bigint>;
}

export interface Executor<T> {
    doQuery(fun:(ctx:Context)=>Result<any>):Executor<T>;
    doQueryOne(fun:(ctx:Context)=>Result<any>):Executor<T>;
    doCommand(fun:(ctx:Context)=>GoError):Executor<T>;
    doError(fun:(err:Error)=>void):Executor<T>;
    getData():T;
    getError():GoError;
    getResult():Result<T>;
    setResponse():void;
}


export function newObject(v: any): any;
export const server: Server;
export const time: Time;



