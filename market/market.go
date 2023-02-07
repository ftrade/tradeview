package market

import (
	"encoding/xml"
	"io"
	"os"
)

type Report struct {
	XMLName xml.Name `xml:"report"`
	Candles Candles  `xml:"candles"`
	Trades  Trades   `xml:"trades"`
}

type Candles struct {
	XMLName xml.Name `xml:"candles"`
	Items   []Candle `xml:"candle"`
}
type Candle struct {
	XMLName    xml.Name `xml:"candle"`
	Timestampt int64    `xml:"timestampt,attr"`
	Open       float32  `xml:"open,attr"`
	High       float32  `xml:"high,attr"`
	Low        float32  `xml:"low,attr"`
	Close      float32  `xml:"close,attr"`
	Volume     int32    `xml:"volume,attr"`
}

type Trades struct {
	XMLName xml.Name `xml:"trades"`
	Items   []Trade  `xml:"trade"`
}

type Trade struct {
	XMLName    xml.Name `xml:"trade"`
	Timestampt int64    `xml:"timestampt,attr"`
	Price      float32  `xml:"price,attr"`
	Volume     int32    `xml:"volume,attr"`
	Profit     float32  `xml:"profit,attr"`
}

func LoadReport(path string) Report {
	xmlFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer xmlFile.Close()

	var report Report
	bytes, _ := io.ReadAll(xmlFile)
	err = xml.Unmarshal(bytes, &report)
	if err != nil {
		panic(err)
	}
	return report
}

// func (r *Report) Stats() Stats {
// 	var minTime, maxTime int64 = math.MaxInt64, math.MinInt64
// 	var maxVolume int32 = 0
// 	var minPrice, maxPrice float32 = math.MaxFloat32, math.SmallestNonzeroFloat32
// 	var minTimeStep int64 = math.MaxInt64
// 	var prevTime int64 = 0
// 	for _, bar := range r.Candles.Items {
// 		ts := bar.Timestampt
// 		if ts > maxTime {
// 			maxTime = ts
// 		}
// 		if ts < minTime {
// 			minTime = ts
// 		}
// 		if bar.Volume > maxVolume {
// 			maxVolume = bar.Volume
// 		}
// 		if bar.High > maxPrice {
// 			maxPrice = bar.High
// 		}
// 		if bar.Low < minPrice {
// 			minPrice = bar.Low
// 		}
// 		dTime := bar.Timestampt - prevTime
// 		if dTime < minTimeStep {
// 			minTimeStep = dTime
// 		}
// 		prevTime = bar.Timestampt
// 	}
// 	fmt.Printf("Report. Max time: %v. Min time: %v. Max volume: %v. \n", maxTime, minTime, maxVolume)

// 	return Stats{
// 		MaxTime:     maxTime,
// 		MinTime:     minTime,
// 		MaxVolume:   maxVolume,
// 		MinPrice:    minPrice,
// 		MaxPrice:    maxPrice,
// 		MinTimeStep: minTimeStep,
// 		BarsCount:   uint32(len(r.Candles.Items)),
// 	}
// }
