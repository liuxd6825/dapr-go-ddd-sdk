package class

import "context"

type ParseContext struct {
	umlBlocks   []*UMLBlock
	nodes       map[string]Element
	connections []*Connection
}

type key struct {
}

var parseCtxKey = key{}

func NewContext(parent context.Context, umlBlock *UMLBlock) context.Context {
	blocks := make([]*UMLBlock, 0)
	blocks = append(blocks, umlBlock)
	nodes := make(map[string]Element)
	return context.WithValue(parent, parseCtxKey, &ParseContext{umlBlocks: blocks, nodes: nodes})
}

func GetParseContext(ctx context.Context) *ParseContext {
	val := ctx.Value(parseCtxKey)
	if val == nil {
		return nil
	}
	pCtx := val.(*ParseContext)
	return pCtx
}

func (p *ParseContext) AddUMLBlock(umlBlock *UMLBlock) {

}

func (p *ParseContext) AddNode(node Element) {
	p.nodes[node.GetNamespaceName()+"."+node.GetName()] = node
}

func (p *ParseContext) GetNode(namespaceName string, nodeName string) (Element, bool) {
	node := p.nodes[namespaceName+"."+nodeName]
	if node != nil {
		return node, true
	}
	return nil, false
}

func (p *ParseContext) GetNodes() map[string]Element {
	return p.nodes
}

func (p *ParseContext) AddConnection(conn ...*Connection) {
	p.connections = append(p.connections, conn...)
}

func (p *ParseContext) GetConnections() []*Connection {
	return p.connections
}
