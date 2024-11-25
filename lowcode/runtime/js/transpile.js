// transpile.js
const ts = require("typescript");

// 从标准输入读取 TypeScript 代码
process.stdin.on("data", (data) => {
    const tsCode = data.toString();

    try {
        // 使用 TypeScript API 编译 TypeScript 代码
        const result = ts.transpileModule(tsCode, {
            compilerOptions: { module: ts.ModuleKind.CommonJS }
        });

        // 输出编译后的 JavaScript
        console.log(result.outputText);
    } catch (error) {
        console.error("Error during transpilation:", error.message);
        process.exit(1);
    }
});