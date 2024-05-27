package class

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"strings"
)

type Connection struct {
	BaseElement
	Arrow DirectionArrow `json:"arrow"`

	LeftName  string `json:"leftName"`
	LeftTitle string `json:"leftTitle"`

	RightName  string `json:"rightName"`
	RightTitle string `json:"rightTitle"`

	leftObject  Element
	rightObject Element

	Type      ConnectionType `json:"type"`
	Connector string         `json:"connector,omitempty"`
	Namespace string         `json:"namespace,omitempty"`
}

type ConnectionType string
type DirectionArrow string

/*
1．--   关联关系(Association): 单向关联\双向关联\自身关联\多维关联
2、<|-- 泛化(Generalization)
3、*--  组合（Composition）
4、o--  聚合(Aggregation)
5、     依赖(Dependency)
6、     实现(Realization)

@startuml
Class11 <|.. Class12
Class13 --> Class14
Class15 ..> Class16
Class17 ..|> Class18
Class19 <--* Class20
@enduml
*/
const (
	ConnectionNone           ConnectionType = "none"
	ConnectionAssociation    ConnectionType = "association"
	ConnectionGeneralization ConnectionType = "generalization"
	ConnectionDependency     ConnectionType = "dependency"
	ConnectionAggregation    ConnectionType = "aggregation"
	ConnectionComposition    ConnectionType = "composition"
	ConnectionRealization    ConnectionType = "realization"
	ConnectionRelation       ConnectionType = "relation" // 关系表
)

const (
	DirectionArrowNone  DirectionArrow = ""
	DirectionArrowLeft  DirectionArrow = "left"
	DirectionArrowRight DirectionArrow = "right"
)

func NewConnection(ctx context.Context, line string, namespace string, notes []*Comment) *Connection {
	conn := &Connection{}
	conn.init(line, namespace)
	conn.Comments = notes
	return conn
}

func (c *Connection) init(line string, namespace string) {
	c.leftObject = nil
	c.Connector = ""
	c.rightObject = nil
	c.Namespace = namespace
	c.Type = c.parseConnectionType(line)
}

func (c *Connection) SetNodes(line string) {
	// AbstractList <|-- ArrayList
	// 类1 "1" *-- "*" 类2 : 包含(contains)
	// foo -left-> dummyLeft
	// Together2 -[hidden]--> Bar1

	str := strings.Trim(line, " ")
	if idx := strings.Index(str, ":"); idx > -1 {
		str = str[:idx]
		str = strings.Trim(str, " ")
	}
	idx := strings.Index(str, "-")
	leftStr := ""
	rightStr := ""
	if idx == -1 {
		return
	}
	leftStr = str[:idx]
	leftList := strings.Split(leftStr, " ")
	count := len(leftList)
	for i := 0; i < count; i++ {
		s := leftList[i]
		if has, res, _, _, other := utils.StringBetween(s, `"`, `"`); has {
			c.LeftTitle = res
			if strings.Contains(other, "<") {
				c.Arrow = DirectionArrowLeft
			} else if strings.Contains(other, ">") {
				c.Arrow = DirectionArrowRight
			}
		} else if strings.Contains(s, "<") {
			c.Arrow = DirectionArrowLeft
		} else if strings.Contains(s, ">") {
			c.Arrow = DirectionArrowRight
		} else {
			c.LeftName = s
		}
	}

	rightStr = str[idx:]
	rightList := strings.Split(rightStr, " ")
	count = len(rightList)
	for i := 0; i < count; i++ {
		s := rightList[i]
		if has, res, _, _, other := utils.StringBetween(s, `"`, `"`); has {
			c.RightTitle = res
			if strings.Contains(other, "<") {
				c.Arrow = DirectionArrowLeft
			} else if strings.Contains(other, ">") {
				c.Arrow = DirectionArrowRight
			}
		} else if strings.Contains(s, "<") {
			c.Arrow = DirectionArrowLeft
		} else if strings.Contains(s, ">") {
			c.Arrow = DirectionArrowRight
		} else {
			c.RightName = s
		}
	}
}
func (c *Connection) GetConnector() string {
	return c.Connector
}

func (c *Connection) GetName() string {
	return ""
}

func (c *Connection) parseConnectionType(line string) ConnectionType {
	if strings.Contains(line, "<|-") {
		return ConnectionGeneralization
	} else if strings.Contains(line, "-|>") {
		return ConnectionGeneralization
	} else if strings.Contains(line, "-*") {
		return ConnectionComposition
	} else if strings.Contains(line, "*-") {
		return ConnectionComposition
	} else if strings.Contains(line, "-o") {
		return ConnectionDependency
	} else if strings.Contains(line, "o-") {
		return ConnectionDependency
	} else if strings.Contains(line, "<-") {
		return ConnectionAssociation
	} else if strings.Contains(line, "->") {
		return ConnectionAssociation
	} else if strings.Contains(line, "-") {
		return ConnectionAssociation
	} else if strings.Contains(line, "..") && strings.HasPrefix(line, "(") && strings.Contains(line, "(") {
		return ConnectionRelation
	}
	return ConnectionNone
}

func IsConnection(line string) bool {
	if strings.Contains(line, "..") && strings.HasPrefix(line, "(") && strings.Contains(line, "(") {
		return true
	} else if strings.Contains(line, "-") {
		return true
	}
	return false
}
