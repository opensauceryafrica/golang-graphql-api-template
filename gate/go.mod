module cendit.io/gate

go 1.21.0

replace cendit.io/garage => ../garage

replace cendit.io/signal => ../signal

replace cendit.io/auth => ../auth

require (
	cendit.io/auth v0.0.0-00010101000000-000000000000
	cendit.io/garage v0.0.0-00010101000000-000000000000
	github.com/99designs/gqlgen v0.17.49
	github.com/getsentry/sentry-go v0.28.1
	github.com/go-chi/chi/v5 v5.0.14
	github.com/go-chi/cors v1.2.1
	github.com/vektah/gqlparser/v2 v2.5.16
)

require (
	cendit.io/signal v0.0.0-00010101000000-000000000000 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/aymerick/raymond v2.0.2+incompatible // indirect
	github.com/biter777/countries v1.7.5 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/caarlos0/env/v6 v6.10.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/opensaucerer/goaxios v0.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pquerna/otp v1.4.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun v1.2.1 // indirect
	github.com/uptrace/bun/dialect/pgdialect v1.2.1 // indirect
	github.com/uptrace/bun/driver/pgdriver v1.2.1 // indirect
	github.com/uptrace/bun/extra/bundebug v1.2.1 // indirect
	github.com/urfave/cli/v2 v2.27.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240312152122-5f08fbb34913 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mellium.im/sasl v0.3.1 // indirect
)
