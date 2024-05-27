export const schemas = {
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
                type: "object",
                properties: {
                    length: {
                        type: "number"
                    },
                    width: {
                        type: "number"
                    },
                    height: {
                        type: "number"
                    }
                },
                required: ["length", "width", "height"]
            },
            password: {
                description: "Coordinates of the warehouse where the product is located.",
                $ref: "https://example.com/geographical-location.schema.json"
            },
            telephone: {
                description: "Coordinates of the warehouse where the product is located.",
                $ref: "https://example.com/geographical-location.schema.json"
            },
            gender: {
                description: "Coordinates of the warehouse where the product is located.",
                $ref: "https://example.com/geographical-location.schema.json"
            }
        },
        required: ["firstName", "lastName", "password"]
    }

};

