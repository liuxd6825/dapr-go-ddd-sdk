export interface TemplatePkg {
    renderFile(fileName: string, data?: { [key: string]: any }): string;

    renderString(str: string, data?: { [key: string]: any }): string;
}
