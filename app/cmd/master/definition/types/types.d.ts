/**
 * Go Error接口
 */
export interface GoError {
    error(): string;
}


/**
 * 第三方包加载
 */
export const require: (fileName: string) => any;


