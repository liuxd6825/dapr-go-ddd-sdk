module.exports = {
    rules:{
        'no-alert': 'off',
        'lines-between-class-members':'off',
        "no-unused-vars": "off",
        "@typescript-eslint/no-unused-vars": "error",
        "import/extensions": ["error", "always", {
            "js": "never",
            "jsx": "never"
        }]
    }
}