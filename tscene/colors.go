package tscene

type Color = uint32

const (
	Grey      Color = 0 // default is not setted
	Red       Color = 1
	Grean     Color = 2
	DarkGrey  Color = 3
	TradeOpen Color = 4
	BadTrade  Color = 5
	GoodTrade Color = 6
)

func Colors() []float32 {
	return []float32{
		0.5, 0.5, 0.5,
		0.898, 0.274, 0.282,
		0.258, 0.650, 0.513,
		0.3, 0.3, 0.3,
		0.298, 0.368, 0.956,
		0.439, 0.039, 0.039,
		0.709, 0.682, 0.211,
	}
}
