import {GoError, Context, Result} from "./common";
import {Entity} from "./db";
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

export interface Handle {
    (ctx: IContext): void;
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

export interface HandleOptions  {
    method?:"GET"|"POST"|"PUT"|"DELETE"|"PATCH"|"HEAD"|"OPTIONS";
    path:string;
    handle:Handle;
    schema?:Schema;
}

export interface Server {
    handle(opts: HandleOptions): void;
    get(opts: HandleOptions):void;
    post(opts: HandleOptions):void;
    put(opts: HandleOptions):void;
    delete(opts: HandleOptions):void;
    patch(opts: HandleOptions):void;

    doRequest:(ictx: RContext, fun: (ctx: Context) => void, opts?: DoOptions)=>  Result<any>;
    doQuery:(ictx: RContext,  fun: QueryFunc, opt?: DoOptions)=>  Result<any>;
    doQueryOne:(ictx: RContext, fun: QueryFunc, opt?: DoOptions)=>  Result<any>;
    doCmd:(ictx: RContext,cmdFun: CmdFunc, opts?: CmdAndQueryOptions)=> Result<any>;
    doCmdAndQueryOne:(ictx: RContext, queryAppId: string, cmd: Command<any>, cmdFun: CmdFunc, queryFun: QueryFunc, opts?: CmdAndQueryOptions)=> Result<any>;
    doCmdAndQueryList:(ictx: RContext, queryAppId: string, cmd: Command<any>, cmdFun: CmdFunc, queryFun: QueryFunc, opts?: CmdAndQueryOptions)=> Result<any>;
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
export interface RContext {
    getId():string;
    getTenantId(): string;
    getCaseId(): string;
    getFindPaging(): FindPagingRequest;
    paramString(key:string): Result<string>;
    paramStrings(key:string):Result<string[]>;
    paramBool(key:string): Result<boolean>;
    paramFloat64(key:string): Result<number>;
    paramInt(key:string): Result<bigint>;
    readJson(data?:any):Result<any>;
    writeJson(data:any):GoError;
    executor<T>():Executor<T>;
}

export interface Executor<T> {
    doQuery(fun:(ctx:Context)=>Result<any>):Executor<T>;
    doQueryOne(fun:(ctx:Context)=>Result<any>):Executor<T>;
    doCommand(fun:(ctx:Context, cmd:Command<any>)=>GoError):Executor<T>;
    doError(fun:(err:Error)=>void):Executor<T>;
    getData():T;
    getError():GoError;
    getResult():Result<T>;
    setResponse():void;
}

export interface FindPagingRequest {
    tenantId    :string;
    fields      :string;      // 以逗号分隔多个字段
    filter      :string;
    sort        :string;
    pageNum     :bigint;
    pageSize    :bigint;
    isTotalRows :boolean;
    groupCols   :{field:string,dataType:string}[];
    groupKeys   :any[];
    valueCols   :{aggFunc:string;field:string;}[];
}

export const server: Server;
export const time: Time;



