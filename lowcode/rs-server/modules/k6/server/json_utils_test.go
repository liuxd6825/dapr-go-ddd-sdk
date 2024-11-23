package server

import "testing"

const JSON_SCHEMA = `
{
    default: {
        title: 'Product',
        description: "A product from Acmes catalog",
        type: 'object',
        properties: {
            firstName: {
                description: 'The unique identifier for a product',
                type: 'string',
                title: '名字',
            },
            lastName: {
                description: 'Name of the product',
                type: 'string',
                title: '姓氏',
            },
            province: {
                description: 'The price of the product',
                type: 'number',
                exclusiveMinimum: 0,
            },
            city: {
                description: 'Tags for the product',
                type: 'array',
                items: {
                    type: 'string',
                },
                minItems: 1,
                uniqueItems: true,
            },
            bio: {
                type: 'string',
            },
            password: {
                type: 'string',
            },
            telephone: {
                type: 'string',
            },
            gender: {
                type: 'string',
            },
        },
        required: ['firstName', 'lastName', 'password'],
    },
    createCommand: {
        type: 'object',
        properties: {
            commandId: { type: 'string' },
            isValidOnly: { type: 'boolean' },
            typeName: { type: 'string' },
            data: {
                type: 'object',
                properties: {
                    id: { type: 'string' },
                    tenantId: { type: 'string' },
                    firstName: { type: 'string' },
                    lastName: { type: 'string' },
                    province: { type: 'number' },
                    citry: { type: 'string' },
                    password: { type: 'string' },
                    telephone: { type: 'string' },
                    gender: { type: 'string' }
                },
                required: ['id', 'tenantId', 'firstName', 'lastName'],
            },
        },
        required: ['commandId', 'typeName', 'data'],
    },
    updateCommand: {
        type: 'object',
        properties: {
            commandId: { type: 'string' },
            isValidOnly: { type: 'boolean' },
            typeName: { type: 'string' },
            data: {
                type: 'object',
                properties: {
                    id: { type: 'string' },
                    tenantId: { type: 'string' },
                    firstName: { type: 'string' },
                    lastName: { type: 'string' },
                    province: { type: 'number' },
                    citry: { type: 'string' },
                    password: { type: 'string' },
                    telephone: { type: 'string' },
                    gender: { type: 'string' },
                },
                required: ['id', 'tenantId', 'firstName', 'lastName'],
            },
        },
        required: ['commandId', 'typeName', 'data'],
    }
}`

func Test_NewMap(t *testing.T) {
	jsonUtils := NewJsonUtils()
	mapData := jsonUtils.NewMap(JSON_SCHEMA)
	t.Log(mapData)
}

func Test_preprocessJSON(t *testing.T) {
	jsonStr := preprocessJSON(JSON_SCHEMA)
	t.Log(jsonStr)
}
