export var schemas = {
    default: {
        $schema: "https://json-schema.org/draft/2020-12/schema",
        $id: "https://example.com/product.schema.json",
        title: "Product",
        description: "A product from Acme's catalog",
        type: "object",
        properties: {
            firstName: {
                description: "The unique identifier for a product",
                type: "string",
                title: "名字"
            },
            lastName: {
                description: "Name of the product",
                type: "string",
                title: "姓氏"
            },
            province: {
                description: "The price of the product",
                type: "number",
                exclusiveMinimum: 0
            },
            city: {
                description: "Tags for the product",
                type: "array",
                items: {
                    type: "string"
                },
                minItems: 1,
                uniqueItems: true
            },
            bio: {
                type: "string",
            },
            password: {
                type: "string",
            },
            telephone: {
                type: "string",
            },
            gender: {
                type: "string",
            }
        },
        required: ["firstName", "lastName", "password"]
    },
    createCommand: {
        $schema: "https://json-schema.org/draft/2020-12/schema",
        $id: "https://example.com/schemas/record/createCommand.json",
        type: "object",
        properties: {
            commandId: { type: "string" },
            isValidOnly: { type: "boolean" },
            typeName: { type: "string" },
            data: {
                type: "object",
                properties: {
                    id: { type: "string" },
                    tenantId: { type: "string" },
                    firstName: { type: "string" },
                    lastName: { type: "string" },
                    province: { type: "number" },
                    citry: { type: "string" },
                    password: { type: "string" },
                    telephone: { type: "string" },
                    gender: { type: "string" },
                },
                required: ["id", "tenantId", "firstName", "lastName"],
            }
        },
        required: ["commandId", "typeName", "data"],
    }
};
