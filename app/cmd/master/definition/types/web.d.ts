// iris context
import {Schema} from "./schema";
import {GoError} from "./types";
import {Context} from "./context";

type DataType =
    | 'string'
    | 'int'
    | 'float'
    | 'money'
    | 'date'
    | 'dateTime'
    | 'bool'
    | 'array'
    | 'object'
    | 'year'
    | 'month'
    | 'day';

type AggFunc =
    | 'sum'
    | 'count'
    | 'avg'
    | 'first'
    | 'last'
    | 'max'
    | 'min'
    | 'zero';

export interface PagingGroupCol {
    field: string;
    dataType: DataType;
}

export interface PagingValueCol {
    aggFunc: AggFunc;
    field: string;
}

export interface FindPagingQueryResult<T> {
    data: T[];
    sumData: T[];
    totalRows?: bigint;
    totalPages?: bigint;
    pageNum?: bigint;
    pageSize?: bigint;
    filter?: string;
    fields?: string;
    sort?: string;
    isFound?: boolean;
    isTotalRows?: boolean;
    isSum?: boolean;
    error?: GoError;
}

export interface FindByIdQuery {
    id: string;
}

export interface FindPagingQuery {
    fields: string;
    filter: string;
    mustFilter: string;
    sort: string;
    pageNum: number;
    pageSize: number;
    isTotalRows: boolean;
    groupCols: PagingGroupCol[];
    groupKeys: any[];
    valueCols: PagingValueCol[];
}

export interface FindListQuery {}

export interface WebParams {
    getId(): string;

    getTenantId(): string;

    getCaseId(): string;

    getFindPaging(): FindPagingQuery;

    string(key: string): string;

    strings(key: string): string[];

    bool(key: string): boolean;

    float64(key: string): number;

    int64(key: string): bigint;
}


//  request context
export interface WebContext   {
    ctx(): Context;

    iCtx(): any;

    params(): WebParams;

    readString(): string;

    readJson(data?: any): any;

    readObject(schema: Schema): Object;

    setHeader(key:string, value:string)

    /**
     * 读二进制
     * @param data
     */
    readBytes(): Uint8Array;

    /**
     * 写入Json
     * @param data
     */
    writeJson(data: any): number;

    /**
     * 写入String
     * @param body
     */
    writeString(body: string): number;

    /**
     * 写入HTML
     * @param body
     */
    writeHTML(body: string): number;

    /**
     * 写入二进制内容
     * @param bytes
     */
    writeBytes(bytes: any): void;

    /**
     * 设置返回HTML内容类型
     * @param cType
     */
    setContentType(cType: string): void;

    /**
     * 获取返回的HTML内容类型
     */
    getContentType(): string;

    /**
     * 设置返回的HTML状态类型
     * @param status
     */
    setStatus(status: number): void;

    /**
     * 获取返回的HTML状态类型
     */
    getStatus(): number;

    getTenantId(): string;

    getCaseId(): string;

    getId(): string;

    getFindPaging(): FindPagingQuery;

    setError(error: any, status?: number): void;

    printf(...a1: any): void;

}

/**
 * 参数项
 */
export interface WebParamItem {
    in: "url" | "path" | "body" | "formValue" | "formObject" | "formFile";
    required: boolean;
    type?: string;
    description?: string;
    example?: any;
    schema?: Schema
}

/**
 * 参数类型
 */
export type WebParamsType = { [key: string]: WebParamItem };

export interface IEnvConfig {
    getAppId() : string
    getAppName(): string
    getAppHttpHost() :string
    getAppHttpPort() :int
    getDaprHost(): string
    getDaprHttpPort() :int64
    getDaprGrpcPort() :int64
    getFsManager(): fsm.Manager
    getRsServerSrcPath() :string
    getRsServerEnable() :bool
}

export interface WebServer {
    init(options: {isPubEvent?: boolean, eventPrefix?:string }): void;
    loadPkg(...pkgName: string): void;
    getEnvCfg() : IEnvConfig;
}

export interface WebService {
    name: string;
    description?: string;
}

