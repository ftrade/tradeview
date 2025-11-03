package tscene

import (
	"github.com/ftrade/tradeview/config"
	"github.com/ftrade/tradeview/market"
	"github.com/ftrade/tradeview/scene"
	"github.com/nullboundary/glfont"
)

const Indent = 0.5

type TradeStats struct {
	MinPrice  float32
	MaxPrice  float32
	MaxVolome int32
}

type TradeScene struct {
	Scene      *scene.Scene
	CandleAxis *CandleAxis
	ViewPort   *scene.ViewPort
	TradeStats TradeStats
}

func BuildScene(report market.Report, font *glfont.Font, height, width int32) *TradeScene {
	xAxis := NewXAxis(report.Candles.Items)
	candleRect, candleLine := BuildCandles(report.Candles.Items)
	trades, tradesAxis := BuildTrades(report.Trades.Items, xAxis)
	scene2D := scene.NewScene()

	fullWindowStage := NewFullWindowStage(width, height)
	fullWindowStage.AddDrawItem(NewCrosslines())
	scene2D.AddStage(fullWindowStage)

	candlesStage := scene.NewStage(scene.Rect[float32]{
		Left:   0,
		Right:  1,
		Bottom: 0,
		Top:    config.CandelsComponentHeight,
	})
	candlesStage.AddDrawItem(candleRect)
	candlesStage.AddDrawItem(candleLine)
	candlesStage.AddDrawItem(trades)
	scene2D.AddStage(candlesStage)

	volumeStage := NewFullWindowStage(width, height)
	volumeStage.AddDrawItem(BuildVolumes(report.Candles.Items))
	scene2D.AddStage(volumeStage)

	candlesStage.SceneView = scene.Rect[float32]{
		Left:  -Indent,
		Right: float32(xAxis.WidthX()) + Indent,
	}
	volumeStage.SceneView = scene.Rect[float32]{
		Left:   -Indent,
		Right:  float32(xAxis.WidthX()) + Indent,
		Bottom: 0,
	}

	vp := &scene.ViewPort{}
	tradeScene := &TradeScene{
		Scene:      scene2D,
		CandleAxis: xAxis,
		ViewPort:   vp,
	}
	stats := &tradeScene.TradeStats

	recomputeBottomTop := func(ws scene.WindowState) {
		minPrice, maxPrice, maxVol := xAxis.MinMaxPriceAndMaxVolume(
			int(candlesStage.SceneView.Left),
			int(candlesStage.SceneView.Right),
		)
		priceHeight := (maxPrice - minPrice) / config.CandelsComponentHeight
		candlesStage.SceneView.Bottom = maxPrice - priceHeight
		candlesStage.SceneView.Top = maxPrice
		volHeigh := float32(maxVol) / (1 - config.CandelsComponentHeight)
		volumeStage.SceneView.Top = volHeigh

		vp.Window.X = ws.MousePosition.X
		vp.Window.Y = ws.MousePosition.Y
		vp.Window.Width = ws.WindowSize.Width
		vp.Window.Height = ws.WindowSize.Height
		vp.Scene.X = candlesStage.SceneView.Left
		vp.Scene.Width = candlesStage.SceneView.Right - candlesStage.SceneView.Left
		vp.Scene.Y = candlesStage.SceneView.Bottom
		vp.Scene.Height = priceHeight

		stats.MinPrice = minPrice
		stats.MaxPrice = maxPrice
		stats.MaxVolome = maxVol
	}
	recomputeBottomTop(scene.WindowState{
		WindowSize: scene.Size[int32]{
			Width:  width,
			Height: height,
		},
	})

	events := &scene2D.WindowEvents
	onDragChain := events.OnDrag
	events.OnDrag = func(ws scene.WindowState) {
		onDragChain(ws)

		dxPercentOfWindow := float32(float32(ws.DragXDelta) / float32(ws.WindowSize.Width))
		view := candlesStage.SceneView
		view = view.Move(dxPercentOfWindow*view.Width(), 0)
		candlesStage.SceneView = view
		volumeStage.SceneView = view
		recomputeBottomTop(ws)
		candlesStage.SceneViewUpdated = true
		volumeStage.SceneViewUpdated = true
	}

	onZoomChain := events.OnZoom
	events.OnZoom = func(ws scene.WindowState) {
		onZoomChain(ws)

		var scaleFactor float32 = 0.8
		if ws.ZoomIn {
			scaleFactor = 1.2
		}
		whereScale := float32(ws.MousePosition.X) / float32(ws.WindowSize.Width)
		view := &candlesStage.SceneView
		scaleAroundX := view.Left + view.Width()*whereScale
		newWidth := view.Width() * scaleFactor
		view.Left = scaleAroundX - (newWidth * whereScale)
		view.Right = scaleAroundX + (newWidth * (1 - whereScale))
		volumeStage.SceneView = *view
		recomputeBottomTop(ws)
		candlesStage.SceneViewUpdated = true
		volumeStage.SceneViewUpdated = true
	}

	candleInfoDrawer := newCandleInfoDrawer(font, vp, xAxis)
	scene2D.AddDrawer(candleInfoDrawer)
	tradeInfoDrawer := newTradeInfoDrawer(font, vp, tradesAxis)
	scene2D.AddDrawer(tradeInfoDrawer)
	minMaxLabelsDrawer := newMinMaxDrawer(font, stats)
	scene2D.AddDrawer(minMaxLabelsDrawer)

	return tradeScene
}
