package market

type Stats struct {
	BarsCount          uint32
	MinTime, MaxTime   int64
	MinTimeStep        int64
	MaxVolume          int32
	MinPrice, MaxPrice float32
}
