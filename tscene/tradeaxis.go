package tscene

import (
	"math"

	"github.com/ftrade/tradeview/market"
)

type XTrade struct {
	ViewX float32
	Trade market.Trade
}

type TradeAxis struct {
	trades []XTrade
}

func newTradeAxis(size int) *TradeAxis {
	return &TradeAxis{
		trades: make([]XTrade, size),
	}
}

func (ta *TradeAxis) FindTrade(viewX float32, distance float32) []XTrade {
	//TODO: hot function on mouse moving that can be improve. Matter if trades are a lot.
	// Binary search can be used.
	left, right := 0, 0
	for i := range len(ta.trades) {
		if math.Abs(float64(ta.trades[i].ViewX)-float64(viewX)) < float64(distance) {
			if right == 0 {
				left = i
			}
			right = i + 1
		} else if right != 0 {
			break
		}
	}
	return ta.trades[left:right]
}
