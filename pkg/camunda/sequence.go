package camunda

import "github.com/beevik/etree"

type Sequence struct {
	*Element
	SourceRef string `json:"sourceRef"`
	TargetRef string `json:"targetRef"`
}

func NewSequence(el *etree.Element) *Sequence {
	element := NewElement(el)
	s := &Sequence{Element: element}
	s.SourceRef = s.GetSourceRef()
	s.TargetRef = s.GetTargetRef()
	return s
}

func (s *Sequence) GetSourceRef() string {
	return s.el.SelectAttrValue("sourceRef", "")
}

func (s *Sequence) GetTargetRef() string {
	return s.el.SelectAttrValue("targetRef", "")
}
