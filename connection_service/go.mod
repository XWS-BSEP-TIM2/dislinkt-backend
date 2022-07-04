module github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service

go 1.18

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/common => ../common

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service => ../connection_service

require (
	github.com/XWS-BSEP-TIM2/dislinkt-backend/common v0.0.0-20220617075150-48a21e06c92f
	github.com/joho/godotenv v1.4.0
	github.com/neo4j/neo4j-go-driver/v4 v4.4.3
	google.golang.org/grpc v1.47.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.2.0 // indirect
	github.com/nats-io/nats.go v1.14.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
