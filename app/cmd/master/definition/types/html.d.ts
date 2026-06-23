import {Context} from "./context";
import {FsPkg} from "./pkg";

export interface HtmlPKg extends FsPkg {
    build(ctx: Context, schemaDoc: string, tplType: string, tplOpts: { [key: string]: any }): string;
}
