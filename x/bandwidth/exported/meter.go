package exported

//type AccStakeProvider interface {
//	GetAccStakePercentage(ctx sdk.Context, address sdk.AccAddress) float64
//}
//
//// General bw handler tx flow
//// 1. Update bw to current for signer
//// 2. Check if signer have enough bw
//// 3. Run Tx, if tx succeed, than:
//// 4. Consume bw and save with old max bw value
//// 5. Load bw with new max value and save it
//type BandwidthMeter interface {
//	// load current bandwidth state after restart
//	Load(ctx sdk.Context)
//	// add value to consumed bandwidth for current block
//	AddToBlockBandwidth(value int64)
//	// add value to overall linked bandwidth
//	AddToBlockKarma(value int64)
//	// adjust price based on 24h loading
//	AdjustPrice(ctx sdk.Context)
//	// get current bandwidth price
//	GetCurrentCreditPrice() float64
//	// get used in window band proportional to desirable
//	GetCurrentNetworkLoad(ctx sdk.Context) float64
//	// commit bandwidth value spent for current block
//	CommitBlockBandwidth(ctx sdk.Context)
//	// commit bandwidth value spent for current block
//	CommitTotalKarma(ctx sdk.Context)
//	// Update acc max bandwidth for current stake. Also, performs recover.
//	UpdateAccMaxBandwidth(ctx sdk.Context, address sdk.AccAddress)
//	// Returns recovered to current block height acc bandwidth
//	GetCurrentAccBandwidth(ctx sdk.Context, address sdk.AccAddress) types.AcсountBandwidth
//	// Returns acc max bandwidth
//	GetAccMaxBandwidth(ctx sdk.Context, address sdk.AccAddress) int64
//	// Returns tx bandwidth cost
//	GetTxCost(ctx sdk.Context, tx sdk.Tx) int64
//	// Return max block bandwidth param value
//	GetMaxBlockBandwidth(ctx sdk.Context) uint64
//	// Return tx bandwidth cost considering the price
//	GetPricedTxCost(ctx sdk.Context, tx sdk.Tx) int64
//	// Return current block spend bandwidth
//	GetCurBlockSpentBandwidth(ctx sdk.Context) uint64
//	// Return links bandwidth cost considering the price
//	GetPricedLinksCost(ctx sdk.Context, tx sdk.Tx) int64
//	//
//	// Performs bw consumption for given acc
//	// To get right number, should be called after tx delivery with bw state obtained prior delivery
//	//
//	// Pseudo code:
//	// bw := getCurrentBw(addr)
//	// bwCost := deliverTx(tx)
//	// consumeBw(bw, bwCost)
//	ConsumeAccBandwidth(ctx sdk.Context, bw types.AcсountBandwidth, amt int64)
//	// Performs updating of total bandwidth used by account for linking
//	UpdateLinkedBandwidth(ctx sdk.Context, bw types.AcсountBandwidth, amt int64)
//
//}
