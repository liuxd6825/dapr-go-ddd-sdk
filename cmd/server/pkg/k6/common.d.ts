export interface Error {
    error(): string;
}
export type Context = any;
export type SuccessHandler= (data:any)=>Error;
export type ErrorHandler = (err:Error)=>void;

export interface Result {
    data: any;
    isFound?: boolean;
    error?: Error;
    doSuccess?: (func:SuccessHandler)=>Result;
    doError?: (error:ErrorHandler)=>Result;
}

