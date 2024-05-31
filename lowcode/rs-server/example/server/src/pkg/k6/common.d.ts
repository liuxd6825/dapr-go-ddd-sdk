export interface GoError {
    error(): string;
}


export type SuccessHandler= (data:any)=>GoError;
export type ErrorHandler = (err:GoError)=>void;

export interface Result<T> {
    data?: T;
    error?: GoError;
    doSuccess?: (func:SuccessHandler)=>Result<T>;
    doError?: (error:ErrorHandler)=>Result<T>;
    results?:()=>{data:T, error?:GoError}
}

export type Object = { [key:string]:any };

export function newGoError (error:string):GoError;

export interface Context {
    err(): GoError;
    value(key:any):any;
    done():void;
    deadline():Date;
}

export interface Logs {
    info(a1?:any):void;
}

export const fmt:Logs;