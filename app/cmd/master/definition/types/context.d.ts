/**
 * Go中context.Context接口
 */
export interface Context {
    err(): GoError;

    value(key: any): any;

    done(): void;

    deadline(): Date;
}

export interface ContextPkg {
    background(): Context;

    withValue(parent: Context, key: string, val: any): Context;

    withTimeout(parent: Context, timeout: any): { ctx: Context; func: any };

    withCancel(parent: Context): { ctx: Context; cancel: any };
}
