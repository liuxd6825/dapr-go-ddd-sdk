package vis_network

type Font struct {
	Color       *string     `json:"color,omitempty"`
	Size        *int64      `json:"size,omitempty"`
	Face        *string     `json:"face,omitempty"`
	Background  *string     `json:"background,omitempty"`
	StrokeWidth *int64      `json:"strokeWidth,omitempty"`
	StrokeColor *string     `json:"strokeColor,omitempty"`
	Align       *string     `json:"align,omitempty"`
	Vadjust     *int64      `json:"vadjust,omitempty"`
	Multi       *bool       `json:"multi,omitempty"`
	Bold        *FontStyles `json:"bold,omitempty"`
	Ital        *FontStyles `json:"ital,omitempty"`
	Boldital    *FontStyles `json:"boldital,omitempty"`
	Mono        *FontStyles `json:"mono,omitempty"`
}
