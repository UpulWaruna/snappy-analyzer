module webSocket-server

go 1.25.1

require (
	common/logger v0.0.0
	github.com/gorilla/websocket v1.5.3
)

replace common/logger => ../common/logger
