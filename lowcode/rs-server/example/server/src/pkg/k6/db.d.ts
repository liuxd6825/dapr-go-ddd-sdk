import {Context, GoError, Result} from "./common";
import {Schema} from "./schema";

export interface Entity {
    id:string;
    tenantId: string;
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
    maxPoolSize?: number;
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



type DataType =
    "string"
    | "int"
    | "float"
    | "money"
    | "date"
    | "dateTime"
    | "bool"
    | "array"
    | "object"
    | "year"
    | "month"
    | "day";

type AggFunc = "sum" | "count" | "avg" | "first" | "last" | "max" | "min" | "zero";

export interface GroupCol {
    field: string;
    dataType: DataType;
}


export interface ValueCol {
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


export interface FindPagingQuery {
    tenantId: string;
    fields: string;
    filter: string;
    mustFilter: string;
    sort: string;
    pageNum: number;
    pageSize: number;
    isTotalRows: boolean;
    groupCols: GroupCol[];
    groupKeys: any[];
    valueCols: ValueCol[];
}

export interface Model<T> {
    create(ctx: Context, entity: any, opts?: Options): GoError;

    update(ctx: Context, entity: any, opts?: Options): GoError;

    deleteById(ctx: Context, tenantId: string, id: string, opts?: Options): GoError;

    findById(ctx: Context, tenantId: string, id: string, opts?: Options): Result<T>;

    findPaging(ctx: Context, query: FindPagingQuery, opts?: Options): FindPagingQueryResult<T>;
}

export interface DB {
    open(cfg: DBConfig): GoError | undefined;

    model<T>(tableName: string): Model<T>;
}

export interface Table {
    name(): string;
    create(ctx: Context): GoError;
    drop(ctx: Context): GoError;
}

export interface Model<T> {
    create(ctx: Context, entity: any, opts?: Options): GoError;

    update(ctx: Context, entity: any, opts?: Options): GoError;

    deleteById(ctx: Context, tenantId: string, id: string, opts?: Options): GoError;

    findById(ctx: Context, tenantId: string, id: string, opts?: Options): Result<T>;

    findPaging(ctx: Context, query: FindPagingQuery, opts?: Options): FindPagingQueryResult<T>;

    table(ctx: Context, schema: Schema, opts?: Options) :Table

}


export const db: DB;