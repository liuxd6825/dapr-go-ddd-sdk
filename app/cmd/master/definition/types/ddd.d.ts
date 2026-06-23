export interface DomainEvent<T> {
    eventId: string; // 事件ID
    eventType: string; // 事件类型
    version: string;  // 事件版本号
    createdTime: Date; // 创建时间
    data: T; // 数据内容
}

export type Metadata = { [key: string]: any };