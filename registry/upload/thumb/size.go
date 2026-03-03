package thumb

import (
	"fmt"

	"github.com/admpub/imageproxy"
	"github.com/coscms/webcore/library/fileupdater"
)

// DefaultSize 缩略图默认尺寸
var DefaultSize = Size{
	Width:   200,
	Height:  200,
	Quality: 100,
}

func AsSizes(ts ...Size) Sizes {
	return Sizes(ts)
}

// Sizes 尺寸列表
type Sizes []Size

func (s Sizes) String() string {
	var r string
	for i, t := range s {
		if i > 0 {
			r += `,`
		}
		r += t.String()
	}
	return r
}

func (s Sizes) AutoCrop() Sizes {
	r := Sizes{}
	for _, t := range s {
		if t.AutoCrop {
			r = append(r, t)
		}
	}
	return r
}

func (s *Sizes) Add(size Size) {
	*s = append(*s, size)
}

func (s Sizes) Has(width, height float64) bool {
	for _, v := range s {
		if v.Width == width && v.Height == height {
			return true
		}
	}
	return false
}

func (s Sizes) Get(width, height float64) *Size {
	for _, v := range s {
		if v.Width == width && v.Height == height {
			return &v
		}
	}
	return nil
}

// Size 缩略图尺寸信息
type Size struct {
	AutoCrop  bool
	Width     float64
	Height    float64
	Quality   int
	SmartCrop *bool
	ScaleUp   *bool
	Fit       *bool
}

func (t Size) String() string {
	return fmt.Sprintf("%vx%v", t.Width, t.Height)
}

func (t *Size) SetAutoCrop(autoCrop bool) {
	t.AutoCrop = autoCrop
}

func (t *Size) SetWidth(width float64) {
	t.Width = width
}

func (t *Size) SetHeight(height float64) {
	t.Height = height
}

func (t *Size) SetQuality(quality int) {
	t.Quality = quality
}

func (t *Size) SetSmartCrop(smartCrop bool) {
	t.SmartCrop = &smartCrop
}

func (t *Size) SetScaleUp(scaleUp bool) {
	t.ScaleUp = &scaleUp
}

func (t *Size) SetFit(fit bool) {
	t.Fit = &fit
}

// Suffix 文件名称尺寸后缀
func (t Size) Suffix() string {
	return fmt.Sprintf("_%v_%v", t.Width, t.Height)
}

func (t Size) ThumbValue() fileupdater.ValueFunc {
	return fileupdater.ThumbValue(int(t.Width), int(t.Height))
}

func (t Size) CropOptions(o *imageproxy.Options) {
	if t.SmartCrop != nil {
		o.SmartCrop = *t.SmartCrop
	}
	if t.ScaleUp != nil {
		o.ScaleUp = *t.ScaleUp
	}
	if t.Fit != nil {
		o.Fit = *t.Fit
	}
	if t.Quality > 0 {
		o.Quality = t.Quality
	}
}
