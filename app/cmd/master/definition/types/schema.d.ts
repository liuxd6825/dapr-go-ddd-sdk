export type SchemaType = "object" | "string" | "integer" | "number" | "array" | "boolean";

export type SchemaFormat =
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


export interface Schema {
    $schema?: string;
    $id?: string;
    title?: string;
    type: string[];
    properties?: SchemaProperties;
    definitions?: SchemaProperties;
    items?: Schema;
    description?: string;
    required?: string[];
    validate?: (val: any) => void;
}

export type SchemaProperties = { [key: string]: SchemaProperty };

export interface SchemaProperty {
    title?: string;
    type: string | string[];
    pattern?: string; // 正则表达式
    format?: SchemaFormat;  // 表格
    properties?: SchemaProperties;
    required?: string[];
    description?: string;
    ref?: string;
    items?: Schema;
    minItems?: number;
    uniqueItems?: boolean;
    exclusiveMinimum?: number;
    minimum?: number;
    maximum?: number;
}


export interface SchemaPkg {
    loadFile(fileName: string, basePath?: string): Schema;
}
