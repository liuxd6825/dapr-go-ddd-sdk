package ddd_neo4j

type CypherResult interface {
	Cypher() string
	Params() map[string]any
	ResultKeys() []string
	ResultOneKey() string
}

type cypherBuilderResult struct {
	cypher    string
	params    map[string]any
	resultKey []string
}

func NewCypherBuilderResult(cypher string, params map[string]any, resultKey []string) CypherResult {
	return &cypherBuilderResult{
		cypher:    cypher,
		params:    params,
		resultKey: resultKey,
	}
}

func (c *cypherBuilderResult) Cypher() string {
	return c.cypher
}

func (c *cypherBuilderResult) Params() map[string]any {
	return c.params
}

func (c *cypherBuilderResult) ResultKeys() []string {
	return c.resultKey
}

func (c *cypherBuilderResult) ResultOneKey() string {
	if len(c.resultKey) > 0 {
		return c.resultKey[0]
	}
	return ""
}
