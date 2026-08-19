module user-service

go 1.26.3

require (
	common v0.0.0
	github.com/go-playground/validator/v10 v10.30.3
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/gorilla/mux v1.8.1
	github.com/jackc/pgx/v5 v5.10.0
	github.com/lib/pq v1.12.3
	google.golang.org/grpc v1.83.1
	google.golang.org/protobuf v1.36.12
)

replace common => ../common

require (
	github.com/gabriel-vasile/mimetype v1.4.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260819154853-08b0e4226688 // indirect
)
