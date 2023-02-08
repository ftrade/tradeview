package gui

import (
	"tradeview/config"
	"tradeview/market"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Viewport struct {
	XAxis               XAxis
	viewRight, viewLeft float32
	report              market.Report
}

func NewViewport(report market.Report) *Viewport {
	xAxis := NewXAxis(report.Candles.Items)
	v := &Viewport{
		report:    report,
		XAxis:     xAxis,
		viewLeft:  -0.5,
		viewRight: float32(xAxis.XWidth()) + 0.5,
	}
	return v
}

func (v *Viewport) ViewWidth() float32 {
	return v.viewRight - v.viewLeft
}

func (v *Viewport) Scale(scaleFactor float32, whereScale float32) {
	scaleAroundX := v.viewLeft + v.ViewWidth()*whereScale
	newWidth := v.ViewWidth() * scaleFactor
	newLeft := scaleAroundX - (newWidth * whereScale)
	newRight := scaleAroundX + (newWidth * (1 - whereScale))
	v.viewLeft = newLeft
	v.viewRight = newRight
}

func (v *Viewport) Move(dxFactor float32) {
	newLeft := v.viewLeft + v.ViewWidth()*dxFactor
	newRight := v.viewRight + v.ViewWidth()*dxFactor
	v.viewLeft = newLeft
	v.viewRight = newRight
}

func (v *Viewport) CalcView() ViewInfo {
	minPrice, maxPrice, maxVol := v.XAxis.MinMaxPriceAndMaxVolume(int(v.viewLeft), int(v.viewRight))
	priceHeight := (maxPrice - minPrice) / config.BarsComponentHeight
	yWindowBottom := maxPrice - priceHeight
	barsMat := mgl.Ortho2D(v.viewLeft, v.viewRight, yWindowBottom, maxPrice)
	volHeigh := float32(maxVol) / (1 - config.BarsComponentHeight)
	volumeMat := mgl.Ortho2D(v.viewLeft, v.viewRight, 0, volHeigh)
	return ViewInfo{
		BarsMat:    barsMat,
		VolumesMat: volumeMat,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		MaxVolume:  maxVol,
	}
}

func (v *Viewport) WindowXToBar(x, width float32) (bar market.Candle, ok bool) {
	index := int(v.viewLeft + v.ViewWidth()*(x/width) + 0.5)
	if index >= 0 && index < len(v.XAxis.bars) {
		return v.XAxis.bars[index], true
	}
	return bar, false
}
