export var uiSchemas = {
    default: {
        layout: {
            titleWidth: "130px",
            rows: [
                {
                    cols: [
                        {name: "firstName", span: 12},
                        {name: "lastName", span: 12}
                    ]
                },
                {
                    cols: [
                        {name: "province", span: 12},
                        {name: "city", span: 12}
                    ]
                },
                {
                    height: "200px",
                    cols: [{name: "bio", span: 24}]
                },
                {
                    cols: [
                        {name: "password", span: 12},
                        {name: "telephone", span: 12}
                    ]
                },
                {
                    cols: [{name: "gender", span: 24}]
                }
            ]
        },
        title: "xxxxxx",
        properties: {
            firstName: {
                widget: "input",
                autofocus: "true",
                emptyValue: "",
                placeholder: "first Name",
                autocomplete: "family-name",
                enableMarkdownInDescription: "true",
                description: "Make text **bold** or *italic*. Take a look at other options [here](https://markdown-to-jsx.quantizor.dev/).",
                options: {}
            },
            lastName: {
                widget: "input",
                autocomplete: "given-name",
                enableMarkdownInDescription: "true",
                description: "Make things **bold** or *italic*. Embed snippets of `code`. <small>And this is a small texts.</small> "
            },
            province: {
                widget: "select"
            },
            city: {
                widget: "select"
            },
            bio: {
                widget: "textarea",
                readonly: "(param) => { console.log('============bio ui:readonly============='); \n  console.log(param); \n return true; \n}"
            },
            password: {
                widget: "input",
                help: "Hint: Make it strong!"
            },
            telephone: {
                widget: "input",
                options: {
                    inputType: "tel"
                }
            },
            gender: {
                widget: "select",
                readonly: "false"
            }
        }
    }
}
