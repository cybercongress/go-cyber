package rpc

type ResultBandwidthPrice struct {
	Price float64 `amino:"unsafe" json:"price"`
}

func CurrentBandwidthPrice() (*ResultBandwidthPrice, error) {
	return &ResultBandwidthPrice{cyberdApp.CurrentBandwidthPrice()}, nil
}
