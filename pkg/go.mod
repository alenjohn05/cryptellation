module cryptellation/pkg

go 1.22.4

replace (
	cryptellation/client => ../clients/go
	cryptellation/internal => ../internal

	cryptellation/svc/backtests => ../svc/backtests
	cryptellation/svc/candlesticks => ../svc/candlesticks
	cryptellation/svc/exchanges => ../svc/exchanges
	cryptellation/svc/forwardtests => ../svc/forwardtests
	cryptellation/svc/indicators => ../svc/indicators
	cryptellation/svc/ticks => ../svc/ticks
)

require (
	cryptellation/client v0.0.0-00010101000000-000000000000
	cryptellation/svc/candlesticks v0.0.0-00010101000000-000000000000
	cryptellation/svc/indicators v0.0.0-00010101000000-000000000000
	cryptellation/svc/ticks v0.0.0-00010101000000-000000000000
	github.com/go-echarts/go-echarts/v2 v2.4.1
	github.com/go-git/go-git/v5 v5.12.0
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/mod v0.20.0
)

require (
	cryptellation/internal v0.0.0-00010101000000-000000000000 // indirect
	cryptellation/svc/backtests v0.0.0-00010101000000-000000000000 // indirect
	cryptellation/svc/exchanges v0.0.0-00010101000000-000000000000 // indirect
	cryptellation/svc/forwardtests v0.0.0-00010101000000-000000000000 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.0.0 // indirect
	github.com/bluele/gcache v0.0.2 // indirect
	github.com/cloudflare/circl v1.3.9 // indirect
	github.com/cyphar/filepath-securejoin v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lerenn/asyncapi-codegen v0.43.0 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/skeema/knownhosts v1.3.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
