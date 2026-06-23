export interface JsonPkg {
    /**
     * 根据json创建map对象
     * @param jsonVal
     */
    newMap(jsonVal: string): any;

    /**
     * 对Json字符串进行格式化
     * @param rawJSON
     */
    format(rawJSON: Uint8Array | string): Uint8Array;
}