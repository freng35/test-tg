module ext_service

go 1.24.2

replace tg_service => ../telegram_service

require (
	github.com/go-chi/chi/v5 v5.2.1
	golang.org/x/sync v0.13.0
	google.golang.org/grpc v1.71.1
	gopkg.in/yaml.v3 v3.0.1
	tg_service v0.0.0
)

require (
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250414145226-207652e42e2e // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
