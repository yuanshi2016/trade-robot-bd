package market

import (
	"trade-robot-bd/app/grid-strategy-svc/util/huobi/pkg/response/base"
)

type SubscribeLast24hCandlestickResponse struct {
	base.WebSocketResponseBase
	Data *Candlestick
	Tick *Candlestick
}
