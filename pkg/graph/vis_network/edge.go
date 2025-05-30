package vis_network

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"

type Edge struct {
	graph.Edge
	Title              *string `json:"title,omitempty"`
	SubEdges           []any   `json:"subEdges,omitempty"`
	Arrows             *string `json:"arrows,omitempty"`
	ArrowStrikethrough *bool   `json:"arrowStrikethrough,omitempty"`
	Chosen             *bool   `json:"chosen,omitempty"`
	Color              *string `json:"color,omitempty"`
	Dashes             *bool   `json:"dashes,omitempty"`
	Font               *string `json:"font,omitempty"`
	Hidden             *bool   `json:"hidden,omitempty"`
	HoverWidth         *int64  `json:"hoverWidth,omitempty"` // please note, hoverWidth could be also a function. This case is not represented here
	LabelHighlightBold *bool   `json:"labelHighlightBold,omitempty"`
	Length             *int64  `json:"length,omitempty"`
	Physics            *bool   `json:"physics,omitempty"`
	//Scaling OptionsScaling `json:"scaling,omitempty"`
	SelectionWidth    *int64 `json:"selectionWidth,omitempty"` // please note, selectionWidth could be also a function. This case is not represented here
	SelfReferenceSize *int64 `json:"selfReferenceSize,omitempty"`
	SelfReference     *struct {
		Size                *int64 `json:"size,omitempty"`
		Angle               *int64 `json:"angle,omitempty"`
		RenderBehindTheNode *bool  `json:"renderBehindTheNode,omitempty"`
	} `json:"selfReference,omitempty"`
	//edges.smooth
	Shadow          *bool  `json:"shadow,omitempty"`
	Smooth          *bool  `json:"smooth,omitempty"`
	Value           *int64 `json:"value,omitempty"`
	Width           *int64 `json:"width,omitempty"`
	WidthConstraint *int64 `json:"widthConstraint,omitempty"`
}
