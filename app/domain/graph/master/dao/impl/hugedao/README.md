# HugeGraph Master Graph DAO

## 顶点 ID 数据类型规范

### 1. 类型定义

| 维度 | 规范 | 当前实现 |
|---|---|---|
| Go struct 字段 | `string` (强制) | `Id string` (`master_node.go:14` / `master_relation.go:15`) |
| JSON tag | `json:"id"` | `json:"id"` |
| GORM column | `gorm:"column:id"` | `gorm:"column:id"` |
| HugeGraph PropertyKey `data_type` | `TEXT` | `PropDataTypeText = "TEXT"` |
| HugeGraph Cardinality | `SINGLE` | `PropCardinalitySingle` |
| VertexLabel.PrimaryKeys | `["id"]` | 已配置 |
| Gremlin `addV` 语法 | `property('id','<value>')` | 全部脚本使用此模式 |
| Gremlin `has` 语法 | `has('id','<value>')` | 全部脚本使用此模式 |

### 2. 为什么必须是 String

| 场景 | String 优势 |
|---|---|
| 业务主键 (人员/公司 ID) | UUID / 业务编码可统一表达 |
| 主子表双向关联 | `MasterNode.Id` 与 `MasterRelation.Source/Target` 必须同类型 |
| HugeGraph TEXT 上限 256 字节 | 足够容纳 UUID v4(36) 与业务编码(典型 <64) |
| 中文支持 | UTF-8 字符串原生支持 (测试中 `张三/李四/大豆科技` 已验证) |

### 3. Gremlin 脚本约定

```gremlin
// 写入
g.addV('master').property('id','<string_value>').property('name','<value>').next()

// 查询
g.V().hasLabel('master').has('id','<string_value>')

// 转义: escapeGremlin 函数自动用单引号包裹, 并转义 ' 和 \
```

### 4. ID 命名规范 (推荐)

| 类型 | 前缀 | 示例 |
|---|---|---|
| person (人员主表) | `p_` | `p_zhangsan` |
| company (公司主表) | `c_` | `c_dadou` |
| edge (主表关系边) | `e_` | `e_zs_to_ds` |
| relNode (子表节点) | `rn_` 或 `rn<edge_id>` | `rn1`, `rn_e_zs_to_ds` |

### 5. 字符集与长度

| 项 | 规范 |
|---|---|
| 字符集 | ASCII 字母 / 数字 / 下划线 / 中文 (UTF-8) |
| 最大长度 | 256 字节 (HugeGraph TEXT 上限) |
| 特殊字符 | `'` `\` 由 `escapeGremlin` 自动转义 |
| 空格 | 不推荐 (降低可读性) |

### 6. 集成测试覆盖

| 测试 | ID 类型 | 验证 |
|---|---|---|
| `TestMasterNodeDao_CreateMain` | `string` | CreateMain 顶点按 id 幂等 |
| `TestMasterNodeDao_DeleteMain` | `string` | DeleteMain 按 id 删除 |
| `TestMasterNodeDao_CreateRelNode` | `string` | master_rel 边按 id 查询 |
| `TestMasterNodeDao_Create` | `string` | 主子表双向边 (人员/公司) |
| `TestMasterNodeDao_UpdateRelNode` | `string` | rel_type 属性变更 |
| `TestMasterNodeDao_DeleteRelNode` | `string` | relNode 删除 |
| `TestBusRelationDao_FindByName` | `string` | 按 name 查找 master 顶点 |

### 7. 后续可选增强

- 自动生成 UUID v4: 引入 `github.com/google/uuid`, DAO 层在 Id 为空时填充
- ID 长度校验: 写入前 `len(id) <= 256`
- 字符集白名单: 仅允许 `[A-Za-z0-9_中文]`

## Schema 同步 (add-only)

`SyncMasterSchema` 自动同步以下内容(只增不删):

1. **PropertyKey**: `id, name, case_id, tenant_id, source_ids, source_type, description, type, node_type, table, source, target, keywords, rel_type`
2. **VertexLabel**: `master`(主键=id), `same`(主键=name)
3. **EdgeLabel**: `master_rel`(master→master, rel_type 区分关系), `same`(master→same, name 区分节点)
4. **IndexLabel**: `masterById`(on id), `masterByName`(on name), `masterByNodeType`(on type), `masterRelById`(on id), `masterRelByRelType`(on rel_type), `sameByName`(on name)

## Gremlin 约束

HugeGraph 1.7 Gremlin 不支持的写法:
- `coalesce(__.unfold(), g.addV(...))` 嵌套
- `fold().coalesce()` 创建逻辑

替代写法:
- `if (g.V().has(...).hasNext()) { 0 } else { g.addV(...).next(); 0 }` (幂等创建)
- `g.addV(...).next()` (直接创建,失败抛错)

### 测试服务器配置

| 环境变量 | 默认值 | 说明 |
|---|---|---|
| `HUGEGRAPH_HOST` | `192.168.120.200` | HugeGraph 服务器地址 |
| `HUGEGRAPH_PORT` | `18080` | 端口 |
| `HUGEGRAPH_SKIP_INIT` | (空) | 设为非空跳过 TestMain 中 Client 初始化 |

固定配置:
- Graph: `hugegraph`
- Username/Password: `admin/admin`
