import {GoError} from "./common";
import {Entity} from "./db";

export type Type = "object" | "string" | "integer" | "number" | "array" | "boolean";

export type Format =
    "date-time"  // 2018-11-13T20:20:39+00:00
    | "time"     // 20:20:39+00:00
    | "date"     // 2018-11-13
    | "duration"
    | "email"
    | "idn-email"
    | "hostname"
    | "idn-hostname"
    | "ipv4"
    | "ipv6"
    | "uuid"
    | "uri"
    | "uri-reference"
    | "iri"
    | "iri-reference"
    | "uri-template"
    | "json-pointer"
    | "relative-json-pointer"
    | "regex";
export interface Items {
    type: string;
}

export type Properties = { [key: string]: Property };

export interface Property {
    title?: string;
    type: string;
    pattern?: string; // 正则表达式
    format?: Format;  // 表格
    properties?: Properties;
    required?: string[];
    description?: string;
    ref?: string;
    items?: Items;
    minItems?: number;
    uniqueItems?: boolean;
    exclusiveMinimum?: number;
    minimum?: number;
    maximum?: number;
}


export interface Options {
    $schema?: string;
    $id: string;
    title?: string;
    description?: string;
    type: string;
    properties?: Properties;
    definitions?: Properties;
    required?: string[];
}

export interface Schema {
    $schema?: string;
    $id: string;
    title?: string;
    description?: string;
    type: string;
    properties?: Properties;
    definitions?: Properties;
    required?: string[];
    validate?: (val: any) => GoError;
}


export function newSchema(opts: Options): Schema;