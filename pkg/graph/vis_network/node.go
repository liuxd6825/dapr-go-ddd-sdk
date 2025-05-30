package vis_network

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/graph"

type Node struct {
	graph.Node
	SubNodes            []string `json:"subNodes,omitempty"`
	Group               *string  `json:"group,omitempty"`
	BorderWidth         *int64   `json:"borderWidth,omitempty"`
	BorderWidthSelected *int64   `json:"borderWidthSelected,omitempty"`
	BrokenImage         *string  `json:"brokenImage,omitempty"`
	Color               *string  `json:"color,omitempty"`
	Chosen              *bool    `json:"chosen,omitempty"`
	Opacity             *int64   `json:"opacity,omitempty"`
	Fixed               *bool    `json:"fixed,omitempty"`
	Font                *Font    `json:"font,omitempty"`
	Hidden              *bool    `json:"hidden,omitempty"`
	Icon                *struct {
		Face   *string `json:"face,omitempty"`
		Code   *string `json:"code,omitempty"`
		Size   *int64  `json:"size,omitempty"`
		Color  *string `json:"color,omitempty"`
		Weight *int64  `json:"weight,omitempty"`
	} `json:"icon,omitempty"`
	Image              *string `json:"image,omitempty"`
	ImagePadding       *int64  `json:"imagePadding,omitempty"`
	LabelHighlightBold *bool   `json:"labelHighlightBold,omitempty"`
	Level              *int64  `json:"level,omitempty"`
	Margin             *struct {
		Top    *int64 `json:"top,omitempty"`
		Right  *int64 `json:"right,omitempty"`
		Bottom *int64 `json:"bottom,omitempty"`
		Left   *int64 `json:"left,omitempty"`
	} `json:"margin,omitempty"`
	Mass    *int64 `json:"mass,omitempty"`
	Physics *bool  `json:"physics,omitempty"`
	//scaling?: OptionsScaling
	Shadow          *bool   `json:"shadow,omitempty"`
	Shape           *string `json:"shape,omitempty"`
	ShapeProperties *struct {
		BorderDashes       *bool   `json:"borderDashes,omitempty"`
		BorderRadius       *int64  `json:"borderRadius,omitempty"`       // only for box shape
		Interpolation      *bool   `json:"interpolation,omitempty"`      // only for image and circularImage shapes
		UseImageSize       *bool   `json:"useImageSize,omitempty"`       // only for image and circularImage shapes
		UseBorderWithImage *bool   `json:"useBorderWithImage,omitempty"` // only for image shape
		CoordinateOrigin   *string `json:"coordinateOrigin,omitempty"`   // only for image and circularImage shapes
	} `json:"shapeProperties,omitempty"`
	Size  *int64  `json:"size,omitempty"`
	Title *string `json:"title,omitempty"`
	Value *int64  `json:"value,omitempty"`
	/**
	 * If false, no widthConstraint is applied. If a number is specified, the minimum and maximum widths of the node are set to the value.
	 * The node's label's lines will be broken on spaces to stay below the maximum and the node's width
	 * will be set to the minimum if less than the value.
	 */
	WidthConstraint *int64 `json:"widthConstraint,omitempty"`
	X               *int64 `json:"x,omitempty,omitempty"`
	Y               *int64 `json:"y,omitempty,omitempty"`
}
