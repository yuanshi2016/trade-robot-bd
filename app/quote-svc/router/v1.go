package router

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"time"
	"trade-robot-bd/api/constant"
	pb "trade-robot-bd/api/quote/v1"
	"trade-robot-bd/api/response"
	"trade-robot-bd/app/quote-svc/cron"
	"trade-robot-bd/app/quote-svc/internal/service"
	"trade-robot-bd/libs/logger"
)

var (
	quoteService = service.NewQuoteService()
	upGrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WsMessage struct {
	Sub      any    `json:"sub"`
	SubType  string `json:"subType"`
	Exchange string `json:"exchange"`
}

func v1api(group *gin.RouterGroup) {
	group.GET("/ws", SubRealTimeTickers)
	group.GET("/ticks", GetTicks)
	group.GET("/ticks/realtime", SubRealTimeTickers)
}

// SubRealTimeTickers 实时获取行情数据
func SubRealTimeTickers(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		logger.Warnf("ws upgrader conn err: %v", err)
		response.NewErrorCreate(c, "网络出错", nil)
		return
	}
	go StreamHandler(conn)
}

// SubTicker 即便 我们不再 期望 来自websocket 更多的请求,我们仍需要 去 websocket 读取 内容，为了能获取到 close 信号
func SubTicker(ws *websocket.Conn, exchange string, ctxOut context.Context, cancelClose context.CancelFunc) {
	for {
		select {
		case <-ctxOut.Done():
			return
		default:
			// 1. 不断获取行情数据
			resp := quoteService.GetTicks(ctxOut, &pb.GetTicksReq{Exchange: exchange}) //
			if resp == nil {
				time.Sleep(1 * time.Second)
				continue
			}
			// 2. 转发到ws中
			var ticks []cron.Ticker
			if err := jsoniter.Unmarshal(resp.Ticks, &ticks); err != nil {
				logger.Warnf("StreamHandler:Unmarshal数据失败")
				errMsg := response.NewResultInternalErr("StreamHandler:Unmarshal数据失败")
				_ = ws.WriteJSON(errMsg)
				time.Sleep(5 * time.Second)
				continue
			}

			err := ws.WriteJSON(response.NewResultSuccess(ticks))
			if err != nil {
				if isExpectedClose(err) {
					logger.Warnf("expected close on socket")
				}
				logger.Warnf("ws1 writeErr: %v", err)
				cancelClose()
				return
			}
			time.Sleep(2 * time.Second)
		}

	}
}
func SubKline(ws *websocket.Conn, msg WsMessage, ctxOut context.Context, cancelClose context.CancelFunc) {
	for {
		select {
		case <-ctxOut.Done():
			return
		default:
			//aa, _ := cron.BinanceKlineAll.MarshalBinary()
			err := ws.WriteJSON(cron.BinanceKlineAll)
			if err != nil {
				if isExpectedClose(err) {
					logger.Warnf("expected close on socket")
				}
				logger.Warnf("ws1 writeErr: %v", err)
				cancelClose()
				return
			}
			time.Sleep(2 * time.Second)
		}

	}
}
func StreamHandler(ws1 *websocket.Conn) {

	//todo 待优化代码
	go func(ws1 *websocket.Conn) {
		ctx := context.Background()
		ctxOut, cancelOut := context.WithCancel(ctx)
		ctxRun, cancelRun := context.WithCancel(ctx)
		msg := []byte(constant.BINANCE)
		var err error
		//go run(ws1, string(msg), ctxOut, cancelRun)
		//open := true
		time.Sleep(2 * time.Second)
		defer ws1.Close()
		for {
			select {
			case <-ctxRun.Done():
				return
			default:
				_, msg, err = ws1.ReadMessage()
				if err != nil {
					cancelOut()
					return
				}
				var message WsMessage
				if err := jsoniter.Unmarshal(msg, &message); err != nil {
					logger.Warnf("StreamHandler:Unmarshal Message消息失败")
					errMsg := response.NewResultInternalErr("StreamHandler:Unmarshal Message消息失败")
					_ = ws1.WriteJSON(errMsg)
					continue
				}
				switch message.SubType {
				case "ticker": //订阅市场
					go SubTicker(ws1, string(msg), ctxOut, cancelRun)
					break
				case "kline": //订阅K线
					go SubKline(ws1, message, ctxOut, cancelRun)
					break
				}
				// 解析发送消息
				// 这块应该是针对交易所订阅的  目前先按照币安
				//if string(msg) == constant.OKEX || string(msg) == constant.BINANCE || string(msg) == constant.HUOBI {
				//	cancelOut()
				//	time.Sleep(2 * time.Second)
				//	ctxOut, cancelOut = context.WithCancel(ctx)
				//	ctxRun, cancelRun = context.WithCancel(ctx)
				//	go run(ws1, string(msg), ctxOut, cancelRun)
				//}
			}
		}
	}(ws1)
}

func isExpectedClose(err error) bool {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
		logger.Warnf("Unexpected websocket close: %v", err)
		return false
	}
	return true
}

// GetTicks 获取行情数据
func GetTicks(c *gin.Context) {
	resp, err := quoteService.GetTicksWithExchange(context.Background(), &pb.GetTicksReq{All: false})
	if err != nil {
		fromError := errors.FromError(err)
		response.NewErrWithCodeAndMsg(c, fromError.Code, fromError.Message)
		return
	}
	var ticks = make(map[string]map[string]interface{})
	if err := json.Unmarshal(resp.Ticks, &ticks); err != nil {
		response.NewInternalServerErr(c, nil)
		return
	}
	response.NewSuccess(c, ticks)
}
