module cendit.io/garage

go 1.21.0

replace cendit.io/auth => ../auth

replace cendit.io/signal => ../signal

require (
	cendit.io/auth v0.0.0-00010101000000-000000000000
	cendit.io/signal v0.0.0-00010101000000-000000000000
	github.com/biter777/countries v1.7.5
	github.com/caarlos0/env/v6 v6.10.1
	github.com/getsentry/sentry-go v0.28.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.4.0
	github.com/uptrace/bun v1.2.1
	github.com/uptrace/bun/dialect/pgdialect v1.2.1
	github.com/uptrace/bun/driver/pgdriver v1.2.1
	github.com/uptrace/bun/extra/bundebug v1.2.1
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.24.0
)

require (
	github.com/aymerick/raymond v2.0.2+incompatible // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/opensaucerer/goaxios v0.0.6 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	mellium.im/sasl v0.3.1 // indirect
)
