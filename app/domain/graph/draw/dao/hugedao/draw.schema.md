# schema

## 删除所有
```
g.V().hasLabel("draw_node").drop();
graph.schema().vertexLabel("draw_node").remove();

g.E().hasLabel("draw_rel").drop();
graph.schema().edgeLabel("draw_rel").remove();

```

## 创建属性类型
```gremlin
def schema = graph.schema();
schema.propertyKey("id").asText().ifNotExist().create();
schema.propertyKey("name").asText().ifNotExist().create();
schema.propertyKey("node_type").asText().ifNotExist().create();
schema.propertyKey("case_id").asText().ifNotExist().create();
schema.propertyKey("tenant_id").asText().ifNotExist().create();
schema.propertyKey("draw_id").asText().ifNotExist().create();
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
## 创建顶点类型
```
def schema = graph.schema();
schema.vertexLabel("draw_node")
.properties(
  "id",           "name",        "type",       "case_id",   "tenant_id",  "draw_id", 
  "description",  "table",       "keywords",   "tags",
  "src_url",      "src_name",    "src_id",     "src_type"
).useCustomizeStringId().ifNotExist().create();
schema.indexLabel("draw_node_name").onV("draw_node").by("name").secondary().ifNotExist().create();
schema.indexLabel("draw_node_type").onV("draw_node").by("type").secondary().ifNotExist().create();
schema.indexLabel("draw_node_case_id").onV("draw_node").by("case_id").secondary().ifNotExist().create();
schema.indexLabel("draw_node_draw_id").onV("draw_node").by("draw_id").secondary().ifNotExist().create();
schema.indexLabel("draw_node_table").onV("draw_node").by("table").secondary().ifNotExist().create();
schema.indexLabel("draw_node_keywords").onV("draw_node").by("keywords").search().ifNotExist().create();  
schema.indexLabel("draw_node_tags").onV("draw_node").by("tags").search().ifNotExist().create();  

//schema.indexLabel("draw_node_tenant_id").onV("draw_node").by("tenant_id").secondary().ifNotExist().create();
//schema.indexLabel("draw_node_src_id").onV("draw_node").by("src_id").secondary().ifNotExist().create();
//schema.indexLabel("draw_node_src_type").onV("draw_node").by("src_type").secondary().ifNotExist().create();
//schema.indexLabel("draw_node_src_url").onV("draw_node").by("src_url").secondary().ifNotExist().create();
//schema.indexLabel("draw_node_src_name").onV("draw_node").by("src_name").secondary().ifNotExist().create();
```
## 创建边类型
```
def schema = graph.schema();
schema.edgeLabel("draw_rel").sourceLabel("draw_node").targetLabel("draw_node")
.properties(
"id",             "name",            "case_id",      "draw_id",    "keywords",  "description",   
"rel_target_id",  "rel_source_id",   "rel_type",      "tags", 
"src_id",         "src_name",        "src_type",     "src_url"
).ifNotExist().create();
schema.indexLabel("draw_rel_name").onE("draw_rel").by("name").secondary().ifNotExist().create();
schema.indexLabel("draw_rel_rel_type").onE("draw_rel").by("rel_type").secondary().ifNotExist().create();
schema.indexLabel("draw_rel_rel_target_id").onE("draw_rel").by("rel_target_id").secondary().ifNotExist().create();
schema.indexLabel("draw_rel_rel_source_id").onE("draw_rel").by("rel_source_id").secondary().ifNotExist().create();
schema.indexLabel("draw_rel_case_id").onE("draw_rel").by("case_id").secondary().ifNotExist().create();
schema.indexLabel("draw_rel_draw_id").onE("draw_rel").by("draw_id").secondary().ifNotExist().create();
schema.indexLabel("draw_rel_src_id").onE("draw_rel").by("src_id").secondary().ifNotExist().create();
schema.indexLabel("draw_rel_keywords").onE("draw_rel").by("keywords").search().ifNotExist().create();
schema.indexLabel("draw_rel_tags").onE("draw_rel").by("tags").search().ifNotExist().create();  

//schema.indexLabel("draw_rel_src_type").onE("draw_rel").by("src_type").secondary().ifNotExist().create();
//schema.indexLabel("draw_rel_tenant_id").onE("draw_rel").by("tenant_id").secondary().ifNotExist().create();
//schema.indexLabel("draw_rel_table").onE("draw_rel").by("table").secondary().ifNotExist().create();
//schema.indexLabel("draw_rel_src_url").onE("draw_rel").by("src_url").secondary().ifNotExist().create();
//schema.indexLabel("draw_rel_src_name").onE("draw_rel").by("src_name").secondary().ifNotExist().create();

```