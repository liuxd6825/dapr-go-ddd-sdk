export interface GoError {
    error(): string;
}
export type Context = any;
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