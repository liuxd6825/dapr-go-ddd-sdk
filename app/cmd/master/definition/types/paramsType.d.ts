import {WebParamsType} from "./web";

export interface ParamsTypePkg {
    loadFile(fileName: string, workPath: string): WebParamsType;
}
