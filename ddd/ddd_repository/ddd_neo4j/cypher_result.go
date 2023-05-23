package ddd_neo4j

type CypherResult interface {
	Cypher() string
	GetCountCypher() string
	Params() map[string]any
	ResultKeys() []string
	ResultOneKey() string
}

type cypherBuilderResult struct {
	cypher    string
	countCypher string
	params    map[string]any
	resultKey []string
}

type CypherResultOptions  struct{
	CountCypher *string
}


func NewCypherBuilderResult(cypher string,  params map[string]any, resultKey []string, opts ...*CypherResultOptions ) CypherResult {
	res := &cypherBuilderResult{
		cypher:    cypher,
		countCypher: "",
		params:    params,
		resultKey: resultKey,
	}

	for _, o := range opts {
		if o==nil{
			continue
		}
		if o.CountCypher!=nil{
			res.countCypher = *o.CountCypher
		}
	}
	return res;
}

func (c *cypherBuilderResult) Cypher() string {
	return c.cypher
}

func (c *cypherBuilderResult) GetCountCypher() string {
	return c.countCypher
}

func (c *cypherBuilderResult) SetCountCypher(value string)   {
	 c.countCypher = value
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


func NewCypherResultOptions() *CypherResultOptions {
	return  &CypherResultOptions{}
}

func (c *CypherResultOptions) SetCountCypher(v string) *CypherResultOptions{
	c.CountCypher = &v;
	return c
}

func (c *CypherResultOptions) GetCountCypher() string{
	if c.CountCypher==nil{
		return ""
	}
	return *c.CountCypher
}

