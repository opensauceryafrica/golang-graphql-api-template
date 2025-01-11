module cendit.io/signal

go 1.21.0

replace cendit.io/garage => ../garage

replace cendit.io/signal => ../signal

replace cendit.io/auth => ../auth

require (
	cendit.io/garage v0.0.0-00010101000000-000000000000
	github.com/aymerick/raymond v2.0.2+incompatible
	github.com/opensaucerer/goaxios v0.0.6
)

require gopkg.in/yaml.v2 v2.4.0 // indirect
