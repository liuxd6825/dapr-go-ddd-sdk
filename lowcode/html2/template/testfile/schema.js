export var schemas = {
    domain: {
        id: "https://example.com/schemas/record/domain.json",
        title: "交易记录",
        description: "交易记录",
        type: "object",
        properties: {
            firstName: { type: "string" },
            lastName: { type: "string" },
            province: { type: "string" },
            city: { type: "string" },
            bio: { type: "string" },
            password: { type: "string" },
            telephone: { type: "string" },
            gender: { type: "string" },
        },
        required: ["id", "tenantId", "caseId", "account", "oppAccount", "date"],
    },
};
