module cryptellation/svc/candlesticks/dagger

go 1.22.4

replace cryptellation/internal => ../../../../internal

require (
	cryptellation/internal v0.0.0-00010101000000-000000000000
	github.com/99designs/gqlgen v0.17.49
	github.com/Khan/genqlient v0.7.0
	github.com/vektah/gqlparser/v2 v2.5.16
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.0.0-20240812163816-c609b12d9815
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.4.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.28.0
	go.opentelemetry.io/otel/log v0.4.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/sdk/log v0.4.0
	go.opentelemetry.io/otel/trace v1.28.0
	go.opentelemetry.io/proto/otlp v1.3.1
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa
	golang.org/x/sync v0.8.0
	google.golang.org/grpc v1.65.0
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240812133136-8ffd90a71988 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240812133136-8ffd90a71988 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
