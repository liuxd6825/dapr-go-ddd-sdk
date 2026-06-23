import {WebParamsType} from "./web";
import {Context} from "./types";

export interface FeignOptions {
    method?: string;
    url: string;
    body?: Schema;
    params?: WebParamsType;
    responses?: WebParamsType
}

export interface FeignPkg {
    get(ctx: Context, options: FeignOptions, params?: { [key: string]: any }): any;

    post(ctx: Context, options: FeignOptions, params?: { [key: string]: any }): any;

    put(ctx: Context, options: FeignOptions, params?: { [key: string]: any }): any;

    delete(ctx: Context, options: FeignOptions, params?: { [key: string]: any }): any;

    patch(ctx: Context, options: FeignOptions, params?: { [key: string]: any }): any;
}


