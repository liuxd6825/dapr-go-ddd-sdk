import {Context, Result} from "./common";

export interface DB {
    open(cfg: DBConfig): Error | undefined;

    model(tableName: string): Model;
}

export interface DBConfig {
    appName?: string;
    host?: string;
    dbName?: string;
    user?: string;
    pwd?: string;
    replicaSet?: string;
    writeConcern?: string;
    readConcern?: string;
    maxPoolSize?: bigint;
    direct?: boolean;
    localThreshold?: string;             // 时间长度
    connectTimeout?: string;           // 时间长度
    heartbeatInterval?: string;          // 时间长度
    operationTimeout?: string;            // 时间长度
    maxConnIdleTime?: string;             // 时间长度
    serverSelectionTimeout?: string; // 时间长度
    socketTimeout?: string;                  // 时间长度
}

export interface Options {
    sort?: string;
    timeout?: bigint;
    updateFields?: string[];
    updateCancel?: string[];
    upsert?: boolean;
}


export interface FindPagingQueryRequest {
    tenantId  :  string
    fields    :  string
    filter    :  string
    mustFilter:  string
    sort     :   string
    pageNum  :   number
    pageSize  :  number
    isTotalRows: boolean
    groupCols  : GroupCol[]
    groupKeys  : any[]
    valueCols  : ValueCol[]
}


export  interface GroupCol {
    field    : string       ;
    dataType : DataType;
}

type DataType = "string"|"int"|"float"| "money"|"date"| "dateTime"| "bool"| "array"|"object"|"year"|"month"| "day"
type AggFunc = "sum"|"count"|"avg"|"first"|"last"|"max"| "min"|"zero"

export interface  ValueCol {
   aggFunc : AggFunc
   field  :  string
}

export interface FindPagingQueryResult {

}

export interface Model {
    insert(ctx: Context, entity: any, opts?: Options): Error;

    update(ctx: Context, entity: any, opts?: Options): Error;

    deleteById(ctx: Context, tenantId: string, id: string, opts?: Options): Error;

    findById(ctx: Context, tenantId: string, id: string, opts?: Options): Result;

    findPaging(ctx: Context, query: FindPagingQueryRequest, opts?: Options) : Result;
}

export const db: DB;