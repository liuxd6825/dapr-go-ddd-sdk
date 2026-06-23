import {Schema} from "./schema";
import {GoError} from "./types";
import {FindPagingQuery} from "./web";

export interface Entity {
    id: string;
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
    localThreshold?: string; // 时间长度
    connectTimeout?: string; // 时间长度
    heartbeatInterval?: string; // 时间长度
    operationTimeout?: string; // 时间长度
    maxConnIdleTime?: string; // 时间长度
    serverSelectionTimeout?: string; // 时间长度
    socketTimeout?: string; // 时间长度
}

export interface Options {
    aggId?: string;
    eventType?: string;
    sort?: string;
    timeout?: bigint;
    updateFields?: string[];
    updateCancel?: string[];
    upsert?: boolean;
}

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


export interface FindListQuery {
    [key: string]: any;
}

export const findListQuery: FindListQuery;

export interface FindListQueryResult<T> {
    getError?: () => GoError;
    getData: () => T[];
}

export interface Dao<T> {
    create(ctx: Context, entity: any, opts?: Options): void;

    createMany(ctx: Context, entity: any[], opts?: Options): void;

    update(ctx: Context, entity: any, opts?: Options): void;

    updateMany(ctx: Context, entities: any[], opts?: Options): void;

    deleteById(ctx: Context, id: string, opts?: Options): void;

    deleteByIds(ctx: Context, ids: string[], opts?: Options): any;

    deleteByRSQL(ctx: Context, rsql: string, opts?: Options): any;

    findById(ctx: Context, id: string, opts?: Options): T;

    findPaging(ctx: Context, query: FindPagingQuery, opts?: Options): FindPagingQueryResult<T>;

    findByIds(ctx: Context, ids: string[], opts?: Options): T[];

    findByRSQL(ctx: Context, rsql: string, opts?: Options): T[];

    countByRSQL(ctx: Context, rsql: string, opts?: Options): number;

    setAggField(field: string): Model<T>;

    table(): Table;

    getSchema():Schema;
}

export interface Table {
    getName(): string;
    getSchema():Schema;
    drop(ctx?: Context): void;
    autoMigrate(ctx?:Context):void;
    exist(ctx?:Context):bool
}

export interface Builder {
    and(...conditions: Condition): Condition;
    or(...conditions: Condition): Condition;
    like(field: string, value: any): Condition;
    eq(field: string, value: any): Condition;
    start(field: string, value: any): Condition;
    end(field: string, value: any): Condition;
    contains(field: string, value: any): Condition;
    notContains(field: string, value: any): Condition;
    neq(field: string, value: any): Condition;
    gt(field: string, value: any): Condition;
    ge(field: string, value: any): Condition;
    lt(field: string, value: any): Condition;
    le(field: string, value: any): Condition;
    in(field: string, value: any): Condition;
    out(field: string, value: any): Condition;
    null(field: string): Condition;
    notNull(field: string): Condition;
}

export interface Condition {
    build(): string;
}

export  interface NewDaoConfig {
    schema: Schema
}
export interface DBPkg {
    /**
     * 创建数据访问对象
     * @param opts
     */
    newDao<T>(schFile:string): Dao<T>;

    /**
     * getDao
     * @param dbKey
     * @param tableName
     */
    getDao<T>(dbKey: string, tableName: string): Dao<T>;

    /**
     *
     */
    newRSQLBuilder(): Builder;

    /**
     * 创建数据表
     * @param opts
     */
    newTable(opts:{dbKey:string, schema:Schema }):Table;

    /**
     * 开启事务
     * @param ctx
     * @param dbKeys
     * @param txFunc
     */
    startTx(ctx: context.Context, dbKeys:string[], txFunc: (ctx:context.Context)=>GoError):GoError;
}