// @ts-ignore
import {Schema} from "../../pkg/k6/schema";

export interface IRecord {
    id:string;
    caseId:string;
    tenantId:string;
    date:string;
    income:number;
}
export interface RecordSchemas {
    domain: Schema,
    createCommand: Schema,
    createCommand2: Schema,
}

export const recordSchemas : RecordSchemas = {
    domain: {
        id: "https://example.com/schemas/record/domain.json",
        title: "交易记录",
        description: "交易记录",
        type: "object",
        properties: {
            id: {type: "string"},
            tenantId: {type: "string"},
            caseId: {type: "string"},
            docId: {type: "string"},
            date: {type: "string", format: "date-time"},
            name: {type: "string"},
            account: {type: "string"},
            bankName: {type: "string"},
            balance: {type: "number"},
            payout: {type: "number"},
            income: {type: "number"},
            amount: {type: "number"},
            oppName: {type: "string"},
            oppAccount: {type: "string"},
        },
        required: ["id", "tenantId", "caseId", "account", "oppAccount", "date"],
    },
    createCommand2: {
        id: "https://example.com/schemas/record/createCommand.json",
        type: "object",
        properties: {
            commandId: {type: "string"},
            isValidOnly: {type: "boolean"},
            typeName: {type: "string"},
            data: {
                type: "object",
                properties: {},
                required: [],
            }
        },
        required: ["commandId", "isValidOnly", "typeName", "data"],
    },
    createCommand: {
        id: "https://example.com/schemas/record/createCommand.json",
        type: "object",
        properties: {
            commandId: {type: "string"},
            isValidOnly: {type: "boolean"},
            typeName: {type: "string"},
            data: {
                type: "object",
                properties: {
                    id: {type: "string"},
                    tenantId: {type: "string"},
                    caseId: {type: "string"},
                    docId: {type: "string"},
                    date: {type: "string", format: "date-time"},
                    name: {type: "string"},
                    account: {type: "string"},
                    bankName: {type: "string"},
                    balance: {type: "number"},
                    payout: {type: "number"},
                    income: {type: "number"},
                    amount: {type: "number"},
                    oppName: {type: "string"},
                    oppAccount: {type: "string"},
                },
                required: ["id", "tenantId", "caseId", "docId", "date", "name", "account", "bankName", "balance", "payout", "income", "amount", "oppName", "oppAccount"],
            }
        },
        required: ["commandId", "typeName", "data"],
    }
}