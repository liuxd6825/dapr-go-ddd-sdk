# schema

## 删除所有
```
g.V().hasLabel("master_node").drop();
graph.schema().vertexLabel("master_node").remove();

g.V().hasLabel("master_same").drop();
graph.schema().vertexLabel("master_same").remove();

g.E().hasLabel("master_rel").drop();
graph.schema().edgeLabel("master_rel").remove();
```

## 创建属性
```gremlin
def schema = graph.schema();
schema.propertyKey("id").asText().ifNotExist().create();
schema.propertyKey("name").asText().ifNotExist().create();
schema.propertyKey("type").asText().ifNotExist().create();
schema.propertyKey("case_id").asText().ifNotExist().create();
schema.propertyKey("tenant_id").asText().ifNotExist().create();
schema.propertyKey("description").asText().ifNotExist().create();
schema.propertyKey("table").asText().ifNotExist().create();
schema.propertyKey("keywords").asText().valueSet().ifNotExist().create();  
schema.propertyKey("tags").asText().valueSet().ifNotExist().create();
schema.propertyKey("rel_type").asText().ifNotExist().create();
schema.propertyKey("rel_target_id").asText().ifNotExist().create();
schema.propertyKey("rel_source_id").asText().ifNotExist().create();
schema.propertyKey("src_id").asText().ifNotExist().create();
schema.propertyKey("src_type").asText().ifNotExist().create();
schema.propertyKey("src_url").asText().ifNotExist().create();
schema.propertyKey("src_name").asText().ifNotExist().create();

```

## 创建顶点
```
def schema = graph.schema();
schema.vertexLabel("master_node")
.properties(
  "name",       "case_id",    "tenant_id",  "description",
  "type",       "table",      "keywords",   "tags",
  "src_id",     "src_type",   "src_name",   "src_url"
).useCustomizeStringId().ifNotExist().create();
schema.indexLabel("master_node_name").onV("master_node").by("name").secondary().ifNotExist().create();
schema.indexLabel("master_node_case_id").onV("master_node").by("case_id").secondary().ifNotExist().create();
schema.indexLabel("master_node_type").onV("master_node").by("type").secondary().ifNotExist().create();
schema.indexLabel("master_node_src_id").onV("master_node").by("src_id").secondary().ifNotExist().create();
schema.indexLabel("master_node_keywords").onV("master_node").by("keywords").search().ifNotExist().create();
schema.indexLabel("master_node_tags").onV("master_node").by("tags").search().ifNotExist().create();  

//schema.indexLabel("master_node_tenant_id").onV("master_node").by("tenant_id").secondary().ifNotExist().create();
//schema.indexLabel("master_node_table").onV("master_node").by("table").secondary().ifNotExist().create();
//schema.indexLabel("master_node_src_name").onV("master_node").by("src_name").secondary().ifNotExist().create();
//schema.indexLabel("master_node_src_type").onV("master_node").by("src_type").secondary().ifNotExist().create();
//schema.indexLabel("master_node_src_url").onV("master_node").by("src_url").secondary().ifNotExist().create();
schema.vertexLabel("master_same")
.properties(
  "name"  
)
.useCustomizeStringId().ifNotExist().create();
schema.indexLabel("master_same_name").onV("master_same").by("name").secondary().ifNotExist().create();
```
## 创建边
```
def schema = graph.schema();
schema.edgeLabel("master_rel")
.sourceLabel("master_node").targetLabel("master_node")
.properties(
  "rel_target_id", "rel_source_id",   "rel_type",     "case_id",   "tenant_id", 
  "keywords",      "tags",            "table",        "description", 
  "src_id",        "src_type",        "src_name",     "src_url"
).ifNotExist().create();
schema.indexLabel("master_rel_target").onE("master_rel").by("rel_target_id").secondary().ifNotExist().create();
schema.indexLabel("master_rel_source").onE("master_rel").by("rel_source_id").secondary().ifNotExist().create();
schema.indexLabel("master_rel_case_id").onE("master_rel").by("case_id").secondary().ifNotExist().create();
schema.indexLabel("master_rel_type").onE("master_rel").by("rel_type").secondary().ifNotExist().create();
schema.indexLabel("master_rel_src_type").onE("master_rel").by("src_type").secondary().ifNotExist().create();
schema.indexLabel("master_rel_src_id").onE("master_rel").by("src_id").secondary().ifNotExist().create();
schema.indexLabel("master_rel_keywords").onE("master_rel").by("keywords").search().ifNotExist().create();
schema.indexLabel("master_rel_tags").onE("master_rel").by("tags").search().ifNotExist().create();  

//schema.indexLabel("master_rel_tenant_id").v("master_rel").by("tenant_id").secondary().ifNotExist().create();
//schema.indexLabel("master_rel_src_name").onE("master_rel").by("src_name").secondary().ifNotExist().create();
//schema.indexLabel("master_rel_src_url").onE("master_rel").by("src_url").secondary().ifNotExist().create();

```