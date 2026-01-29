package backend

import "github.com/admpub/log"

var WebSocketLogger = log.GetLogger(`websocket`)

func init() {
	WebSocketLogger.SetLevel(`Info`)
}
