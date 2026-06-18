module github.com/snisid/platform/services/ws-gateway

go 1.26.0

require (
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674
	github.com/segmentio/kafka-go v0.4.47
	github.com/snisid/platform/backend v0.0.0
)

replace github.com/snisid/platform/backend => ../../
