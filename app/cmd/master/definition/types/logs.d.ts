import {Context, GoError} from "./types.d.ts";

export type LogFields = { [key: string]: any };

export interface LogPkg {
    trace(ctx: Context, fields: LogFields): void;

    print(ctx: Context, fields: LogFields): void;

    printf(ctx: Context, fields: LogFields, fmt: string, ...args: any): void;

    println(ctx: Context, fields: LogFields, fmt: string, ...args: any): void;

    debug(ctx: Context, fields: LogFields): void;

    debugMsg(ctx: Context, ...arg: any): void;

    debugf(ctx: Context, fields: LogFields, fmt: string, ...args: any): void;

    debugfmt(ctx: Context, fmt: string, ...args: any): void;

    info(ctx: Context, fields: LogFields, ...args: any): void;

    infof(ctx: Context, fields: LogFields, fmt: string, ...args: any): void;

    infoMsg(ctx: Context, ...args: any): void;

    warn(ctx: Context, fields: LogFields): void;

    warnf(ctx: Context, fields: LogFields, fmt: string, ...args: any): void;

    warning(ctx: Context, fields: LogFields, ...args: any): void;

    error(ctx: Context, fields: LogFields, ...args: any): void;

    errorErr(ctx: Context, err: GoError): void;

    errorf(ctx: Context, fields: LogFields, fmt: string, ...args: any): void;

    errorfmt(ctx: Context, fields: LogFields, ...args: any): void;

    errorMsg(ctx: Context, ...args: any): void;

    panic(ctx: Context, fields: LogFields, ...args: any): void;

    panicf(ctx: Context, fields: LogFields, fmt: string, ...args: any): void;

    panicError(ctx: Context, err: GoError): void;

    fatal(ctx: Context, fields: LogFields, ...args: any): void;

    fatalMsg(ctx: Context, ...args: any): void;

    debugStart(ctx: Context, fields: LogFields, fun: () => any): void;
}
