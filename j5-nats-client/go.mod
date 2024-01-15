module github.com/jacksondr5/go-monorepo/j5-nats-client

go 1.21.5

require (
	github.com/nats-io/nats.go v1.31.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/kardianos/service v1.2.2

require (
	github.com/jacksondr5/go-monorepo/nats-common v0.0.0-00010101000000-000000000000
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/jacksondr5/go-monorepo/nats-common => ../nats-common
