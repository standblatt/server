package domain

import "slices"

// Layout is a print form: values placed onto a pre-printed sheet. Sizes and positions are in millimetres;
// the offset moves every element to make up for a printer that does not print exactly where it should.
type Layout struct {
	Width    float64   `json:"width" minimum:"20" maximum:"500"`
	Height   float64   `json:"height" minimum:"20" maximum:"500"`
	OffsetX  float64   `json:"offsetX" required:"false" minimum:"-50" maximum:"50"`
	OffsetY  float64   `json:"offsetY" required:"false" minimum:"-50" maximum:"50"`
	Elements []Element `json:"elements"`
	Source   string    `json:"source,omitempty" required:"false" enum:"results,registration,rank," doc:"What the form prints: the results of a programme (default), a registration or a rank"`
}

// Sources of the values a print form prints.
const (
	SourceResults      = "results"
	SourceRegistration = "registration"
	SourceRank         = "rank"
)

// DataSource returns the source of the layout, results when none is set.
func (l Layout) DataSource() string {
	if l.Source == "" {
		return SourceResults
	}
	return l.Source
}

// Element is one thing printed on a form. Value is text with placeholders like {{shooter.lastname}}; with Repeat
// the element is printed that often, moved by DX/DY each time, and # inside a placeholder becomes the count from 1,
// so {{s1.shot#}} with Repeat 10 prints the ten shots of series 1. An image shows Image of the club's library in W×H.
type Element struct {
	Type   string  `json:"type" enum:"text,qr,barcode,line,image"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	W      float64 `json:"w" required:"false" doc:"Width; 0 lets text run to the right edge"`
	H      float64 `json:"h" required:"false" doc:"Height; QR codes are W×W"`
	Size   float64 `json:"size" required:"false" doc:"Font size in points, 10 when 0"`
	Bold   bool    `json:"bold" required:"false"`
	Italic bool    `json:"italic,omitempty" required:"false"`
	Wrap   bool    `json:"wrap,omitempty" required:"false" doc:"Text breaks into lines within W, each H high"`
	Align  string  `json:"align" required:"false" enum:"L,C,R,"`
	Value  string  `json:"value" required:"false"`
	Repeat int     `json:"repeat" required:"false" minimum:"0" maximum:"100"`
	DX     float64 `json:"dx" required:"false"`
	DY     float64 `json:"dy" required:"false"`
	Image  int64   `json:"image,omitempty" required:"false" doc:"The image of the club's library an image element shows"`
}

// ImageIDs returns the images the layout shows, each once.
func (l Layout) ImageIDs() []int64 {
	ids := []int64{}
	for _, e := range l.Elements {
		if e.Type == "image" && e.Image != 0 && !slices.Contains(ids, e.Image) {
			ids = append(ids, e.Image)
		}
	}
	return ids
}
