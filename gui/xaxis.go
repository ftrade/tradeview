package gui

import (
	"math"
	"tradeview/market"
)

// segment represent time period from leftMillis to rightMillis.
// Segment are splitted by equal parts where first starts with leftMillis time and leftIndex index at a XAxis
// and last ends with rightMillis time and rightIndex index at a XAsis.
type segment struct {
	leftMilllis        int64
	leftIndex          int
	rightMillis        int64
	rightIndex         int
	minPrice, maxPrice float32
	maxVolume          int32
}

// XAxis represents x axis that can map bar timestamp to index and vice versa.
type XAxis struct {
	segments  []segment
	stepWidth int
	bars      []market.Candle
}

func NewXAxis(bars []market.Candle) XAxis {
	stepWidth := bars[1].Timestampt - bars[0].Timestampt

	var segments []segment
	seg := segment{
		leftMilllis: bars[0].Timestampt,
		leftIndex:   0,
		minPrice:    float32(math.MaxFloat32),
		maxPrice:    0,
		maxVolume:   0,
	}
	prevMilli := bars[0].Timestampt
	for i, b := range bars {
		if b.Timestampt-prevMilli > stepWidth {
			seg.rightIndex = i - 1
			seg.rightMillis = prevMilli
			segments = append(segments, seg)
			seg = segment{
				leftMilllis: b.Timestampt,
				leftIndex:   i,
				minPrice:    float32(math.MaxFloat32),
				maxPrice:    0,
				maxVolume:   0,
			}
		}
		if b.High > seg.maxPrice {
			seg.maxPrice = b.High
		}
		if b.Low < seg.minPrice {
			seg.minPrice = b.Low
		}
		if b.Volume > seg.maxVolume {
			seg.maxVolume = b.Volume
		}

		prevMilli = b.Timestampt
	}
	seg.rightIndex = len(bars) - 1
	seg.rightMillis = prevMilli
	segments = append(segments, seg)

	return XAxis{
		segments:  segments,
		stepWidth: int(stepWidth),
		bars:      bars,
	}
}

func (xa *XAxis) XWidth() int {
	return len(xa.bars) - 1
}

func (xa *XAxis) MinMaxPriceAndMaxVolume(from int, upTo int) (float32, float32, int32) {
	if from < 0 {
		from = 0
	}
	if upTo >= len(xa.bars) {
		upTo = len(xa.bars) - 1
	}

	fromSegIndex := xa.searchSegment(from)
	minPrice, maxPrice, maxVol := float32(math.MaxFloat32), float32(0), int32(0)
	i := from

	for segI := fromSegIndex; segI < len(xa.segments) && i <= upTo; segI++ {
		curSeg := xa.segments[segI]
		if i == curSeg.leftIndex && curSeg.rightIndex <= upTo {
			if maxPrice < curSeg.maxPrice {
				maxPrice = curSeg.maxPrice
			}
			if minPrice > curSeg.minPrice {
				minPrice = curSeg.minPrice
			}
			if maxVol < curSeg.maxVolume {
				maxVol = curSeg.maxVolume
			}
			i = curSeg.rightIndex + 1
		} else {
			for ; i <= curSeg.rightIndex && i <= upTo; i++ {
				b := xa.bars[i]
				if b.High > maxPrice {
					maxPrice = b.High
				}
				if b.Low < minPrice {
					minPrice = b.Low
				}
				if b.Volume > maxVol {
					maxVol = b.Volume
				}
			}
		}
	}
	return minPrice, maxPrice, maxVol
}

func (xa *XAxis) searchSegment(index int) (segIndex int) {
	avgSegmentSize := len(xa.bars) / len(xa.segments)
	si := index / avgSegmentSize
	if si >= len(xa.segments) {
		si = len(xa.segments) - 1
	}
	for {
		curSeg := xa.segments[si]
		if curSeg.leftIndex <= index && index <= curSeg.rightIndex {
			return si
		}
		if index < curSeg.leftIndex {
			si--
		} else if index > curSeg.rightIndex {
			si++
		}
	}
}

func (xa *XAxis) TimeToX(millis int64) uint32 {
	return 0
}

func (xa *XAxis) XToTime(x float32) int64 {
	return 0
}
