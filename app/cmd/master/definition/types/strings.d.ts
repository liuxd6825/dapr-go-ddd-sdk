export interface StringsPKg {
    /**
     * 字符串转成Byte
     * @param str
     */
    toByte(str: string): Uint8Array;

    /**
     * 将值转成string
     * @param val
     */
    toString(val: any): string;

    /**
     * 格式化
     * @param format
     * @param arg
     */
    fmt(format: string, ...arg: any[]): string;

    /**
     * 连接字符串
     * @param str
     */
    join(...str: string[]): string;

    /**
     * 转小写
     * @param str
     */
    toLower(str: string): string;

    /**
     * 转大写
     * @param str
     */
    toUpper(str: string): string;

    /**
     * 查找子字符串
     * @param str
     * @param substr
     */
    index(str: string, substr: string): bigint;

    /**
     * 从后方查找子字符串
     * @param str
     * @param substr
     */
    lastIndex(str: string, substr: string): bigint;

    /**
     * 分割字符串
     * @param str
     * @param substr
     */
    split(str: string, substr: string): string[];

    /**
     * 字符串长度
     * @param str
     */
    len(str: string): bigint;

    trim(str: string, cutset: string): string;

    trimLeft(str: string, cutset: string): string;

    trimRight(str: string, cutset: string): string;

    trimLeftAll(str: string, cutset: string): string;

    trimRightAll(str: string, cutset: string): string;

    count(str: string, sub: string): bigint;

    fields(str: string): string[];
}
