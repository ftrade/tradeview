package scene

import (
	"fmt"
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

// BarAxis represents x axis that can map bar timestamp to index and vice versa.
type BarAxis struct {
	Bars      []market.Candle
	segments  []segment
	stepWidth int
}

func NewXAxis(bars []market.Candle) BarAxis {
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

	return BarAxis{
		segments:  segments,
		stepWidth: int(stepWidth),
		Bars:      bars,
	}
}

func (xa *BarAxis) WidthX() int {
	return len(xa.Bars) - 1
}

func (xa *BarAxis) WidthTime() int64 {
	leftMilli := xa.Bars[0].Timestampt
	rightMilli := xa.Bars[len(xa.Bars)-1].Timestampt
	return rightMilli - leftMilli
}

func (xa *BarAxis) MinMaxPriceAndMaxVolume(from int, upTo int) (float32, float32, int32) {
	if from < 0 {
		from = 0
	}
	if upTo >= len(xa.Bars) {
		upTo = len(xa.Bars) - 1
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
				b := xa.Bars[i]
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

func (xa *BarAxis) searchSegment(index int) (segIndex int) {
	avgSegmentSize := len(xa.Bars) / len(xa.segments)
	si := index / avgSegmentSize
	if si >= len(xa.segments) {
		si = len(xa.segments) - 1
	}
	for {
		if si >= len(xa.segments) {
			// try catch error
			fmt.Printf("Search index: %d. Bars count: %d. Last segment %+v", index, len(xa.Bars), xa.segments[len(xa.segments)])
		}
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

func (xa *BarAxis) TimeToX(millis int64) float32 {
	timeCoeff := float32(millis-xa.Bars[0].Timestampt) / float32(xa.WidthTime())
	indexGuess := timeCoeff * float32(len(xa.Bars))
	i := int(indexGuess)
	isFirst, searchLeft := true, false
	leftIndex, rightIndex := -1, -1
	for i >= 0 && i < len(xa.Bars) {
		bar := xa.Bars[i]
		if bar.Timestampt == millis {
			return float32(i)
		}
		if isFirst {
			isFirst = false
			if bar.Timestampt > millis {
				searchLeft = true
			}
		} else {
			if searchLeft {
				if bar.Timestampt < millis {
					leftIndex = i
					rightIndex = i + 1
					break
				}
			} else {
				if bar.Timestampt < millis {
					leftIndex = i
					rightIndex = i + 1
					break
				}
			}
		}
		if searchLeft {
			i--
		} else {
			i++
		}
	}
	if leftIndex < 0 {
		if searchLeft {
			return 0
		} else {
			return float32(xa.WidthX())
		}
	}
	widthBetweenBars := float32(xa.Bars[rightIndex].Timestampt - xa.Bars[leftIndex].Timestampt)
	return float32(leftIndex) + float32(millis-xa.Bars[leftIndex].Timestampt)/widthBetweenBars
}
