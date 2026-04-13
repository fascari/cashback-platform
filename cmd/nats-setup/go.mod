module github.com/cashback-platform/cmd/nats-setup

go 1.26.1

require (
	github.com/cashback-platform/kit v0.0.0
	github.com/nats-io/nats.go v1.31.0
)

require (
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
)

replace github.com/cashback-platform/kit => ../../kit
