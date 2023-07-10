package market

import (
	"trade-robot-bd/app/grid-strategy-svc/util/huobi/pkg/response/base"
)

type SubscribeDepthResponse struct {
	base.WebSocketResponseBase
	Data *Depth
	Tick *Depth
}
