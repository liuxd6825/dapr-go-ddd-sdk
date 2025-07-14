package mongo

const PromptFormat = `
你是一个SQL到MongoDB映射器，根据用户要求，将自然语言转换为mongo统计查询语句

### 处理步骤
1. 根据以下schema将自然语言转换为sql
2. 将sql转换为mongo语句

### 用户要求
{query}

### 输出结果
仅以json格式返回mongo统计查询语句。
{schema}
`
