package tscene

import (
	"fmt"
	"math"
	"time"

	"github.com/ftrade/tradeview/config"
	"github.com/ftrade/tradeview/scene"
	"github.com/nullboundary/glfont"
)

type candleInfoDrawer struct {
	font *glfont.Font
	vp   *scene.ViewPort
	axis *CandleAxis
}

func newCandleInfoDrawer(font *glfont.Font, vp *scene.ViewPort, axis *CandleAxis) *candleInfoDrawer {
	return &candleInfoDrawer{
		font: font,
		vp:   vp,
		axis: axis,
	}
}

func (cid *candleInfoDrawer) Draw(dctx scene.DrawContext) {
	barsBottom := int(float32(dctx.WindowSize.Height) * config.CandelsComponentHeight)
	x, y := dctx.MousePosition.X, dctx.MousePosition.Y
	barsHeight := float32(dctx.WindowSize.Height) * config.CandelsComponentHeight
	onHover := float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= dctx.WindowSize.Width
	if !onHover {
		return
	}
	const pad = 20
	viewX := cid.vp.WindowXToSceneX(x)
	bar, ok := cid.axis.FindBar(viewX)
	if ok {
		time := time.UnixMilli(bar.Timestampt)
		timeStr := time.Format("Mon Jan _2 15:04:05 2006")
		_ = cid.font.Printf(float32(x)+pad, float32(barsBottom), 1, timeStr)

		_ = cid.font.Printf(float32(x)+pad, float32(y)-pad*4-pad, 1, fmt.Sprintf("Open:  %0.2f", bar.Open))
		_ = cid.font.Printf(float32(x)+pad, float32(y)-pad*3-pad, 1, fmt.Sprintf("High:  %0.2f", bar.High))
		_ = cid.font.Printf(float32(x)+pad, float32(y)-pad*2-pad, 1, fmt.Sprintf("Low:   %0.2f", bar.Low))
		_ = cid.font.Printf(float32(x)+pad, float32(y)-pad*1-pad, 1, fmt.Sprintf("Close: %0.2f", bar.Close))
		_ = cid.font.Printf(float32(x)+pad, float32(y)-pad*0-pad, 1, fmt.Sprintf("Vol: %d", bar.Volume))
	}
}

type tradeInfoDrawer struct {
	font *glfont.Font
	vp   *scene.ViewPort
	axis *TradeAxis
}

func newTradeInfoDrawer(font *glfont.Font, vp *scene.ViewPort, axis *TradeAxis) *tradeInfoDrawer {
	return &tradeInfoDrawer{
		font: font,
		vp:   vp,
		axis: axis,
	}
}

func (tid *tradeInfoDrawer) Draw(dctx scene.DrawContext) {
	x, y := dctx.MousePosition.X, dctx.MousePosition.Y
	barsHeight := float32(dctx.WindowSize.Height) * config.CandelsComponentHeight
	onHover := float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= dctx.WindowSize.Width
	if !onHover {
		return
	}
	const pad = 20
	viewX := tid.vp.WindowXToSceneX(x)

	// try make convient search distance depending scale and window width
	wdx := math.Pow(float64(tid.vp.Scene.Width*0.1), 0.5) //nolint:mnd // just empirical numbers
	searchDistance := float32(wdx) * float32(math.Log10(wdx))
	const minDistance = 0.25
	if searchDistance < minDistance {
		searchDistance = minDistance
	}
	trades := tid.axis.FindTrade(viewX, searchDistance)
	switch len(trades) {
	case 0:
		return
	case 1:
		trade := trades[0].Trade
		_ = tid.font.Printf(float32(x)+pad, float32(y)+20*1+pad, 1, fmt.Sprintf("Price: %0.2f", trade.Price))
		_ = tid.font.Printf(float32(x)+pad, float32(y)+20*2+pad, 1, fmt.Sprintf("Vol: %d", trade.Volume))
		time := time.UnixMilli(trade.Timestampt)
		timeStr := time.Format("15:04:05")
		_ = tid.font.Printf(float32(x)+pad, float32(y)+20*3+pad, 1, fmt.Sprintf("Time: %s", timeStr))
		if trade.Profit != 0 {
			_ = tid.font.Printf(float32(x)+pad, float32(y)+20*4+pad, 1, fmt.Sprintf("Profit: %0.2f", trade.Profit))
		}
	default:
		_ = tid.font.Printf(float32(x)+pad, float32(y)+20*1+pad, 1, fmt.Sprintf("%d trades", len(trades)))
	}
}

type minMaxLabelsDrawer struct {
	font  *glfont.Font
	stats *TradeStats
}

func newMinMaxDrawer(font *glfont.Font, stats *TradeStats) *minMaxLabelsDrawer {
	return &minMaxLabelsDrawer{
		font:  font,
		stats: stats,
	}
}

func (mmd *minMaxLabelsDrawer) Draw(dctx scene.DrawContext) {
	barsBottom := float32(dctx.WindowSize.Height) * config.CandelsComponentHeight
	x, y := dctx.MousePosition.X, dctx.MousePosition.Y
	barsHeight := float32(dctx.WindowSize.Height) * config.CandelsComponentHeight
	if float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= dctx.WindowSize.Width {
		dPrice := mmd.stats.MaxPrice - mmd.stats.MinPrice
		price := mmd.stats.MaxPrice - (float32(y)/barsHeight)*dPrice

		_ = mmd.font.Printf(0, float32(y), 1, fmt.Sprintf("%0.2f", price))
	}

	_ = mmd.font.Printf(0, barsBottom, 1, fmt.Sprintf("%0.2f", mmd.stats.MinPrice))
	_ = mmd.font.Printf(0, config.FontSize, 1, fmt.Sprintf("%0.2f", mmd.stats.MaxPrice))
}

// func (w *Window) drawPricesLabels() {
// 	// warning. Code is fragile, don't draw text before crosslines drawing because it change program
// 	barsBottom := int(float32(w.Height) * config.CandelsComponentHeight)
// 	x, y := w.window.GetCursorPos()
// 	barsHeight := float32(w.Height) * config.CandelsComponentHeight
// 	if float32(y) <= barsHeight && y >= 0 && x >= 0 && x <= float64(w.Width) {
// 		dPrice := w.ViewInfo.MaxPrice - w.ViewInfo.MinPrice
// 		price := w.ViewInfo.MaxPrice - (float32(y)/barsHeight)*dPrice

// 		_ = w.Font.Printf(0, float32(y), 1, fmt.Sprintf("%0.2f", price))
// 	}

// 	_ = w.Font.Printf(0, float32(barsBottom), 1, fmt.Sprintf("%0.2f", w.ViewInfo.MinPrice))
// 	_ = w.Font.Printf(0, config.FontSize, 1, fmt.Sprintf("%0.2f", w.ViewInfo.MaxPrice))
// }
