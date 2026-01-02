module webSocket-server

go 1.25.1

require (
    github.com/gorilla/websocket v1.5.3
    common/logger v0.0.0
    )
replace common/logger => ../common/logger

require github.com/gorilla/websocket v1.5.3 // indirect