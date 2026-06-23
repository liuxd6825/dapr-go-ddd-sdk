export interface FsOpts {
    basePath?: string;
}

export type FsWriteModel = number;

export interface FileInfo {
    /**
     文件名称
     */
    name: string;
    /**
     是否为目录
     */
    isDir: boolean;
    /**
     * 文件大小
     */
    size: bigint;
    /**
     * 文件大小
     */
    sizeTitle: string;
    /**
     * 子目录与子文件
     */
    subFiles: FileInfo[];
}

export interface FsPkg {
    /**
     * 读取文件
     * @param fileName  文件名称
     * @param opts
     */
    readFile(fileName: string, opts?: FsOpts): Uint8Array;

    /**
     * 写文件
     * @param fileName 文件名称
     * @param data    文件内容
     * @param opts
     */
    writeFile(fileName: string, data: Uint8Array | string, opts?: FsOpts): void;

    /**
     * 写Json文件
     * @param fileName 文件名称
     * @param data    文件内容
     * @param opts
     */
    writeJson(filename: string, data: Uint8Array | string, opts?: FsOpts): void;

    /**
     * 删除文件
     * @param fileName 文件名称
     * @param opts
     */
    removeFile(fileName: string, opts?: FsOpts): void;

    /**
     * 删除目录中的所有文件
     * @param path 文件名称
     * @param opts
     */
    removeAll(path: string, opts?: FsOpts): void;

    /**
     * 创建目录
     * @param path 目录
     * @param model 写模式
     */
    mkdir(path: string, model: FsWriteModel): void;

    /**
     * 读取目录中的子目录与文件
     * @param path 目录
     * @param opts
     */
    readPath(path: string, opts?: FsOpts): FileInfo[];

    /**
     * 读取目录中的所有子目录与所有文件
     * @param path
     * @param opts
     */
    readAllPath(path: string, opts?: FsOpts): FileInfo[];

    /**
     * 根据FsName创建的文件接口
     * @param fsName 在config.yaml的fs节配置的名称
     */
    newFs(fsName: string): FsPkg;

    /**
     * 文件改名
     * @param oldName
     * @param newName
     * @param opts
     */
    rename(oldName: string, newName: string, opts?: FsOpts): any;

    /**
     * 创建文件
     * @param name
     * @param opts
     */
    create(name: string, opts?: FsOpts): any;
}
