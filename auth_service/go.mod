module github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service

go 1.18

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/common => ../common

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service => ../auth_service

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service => ../profile_service

require (
	github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway v0.0.0-20220627220935-0cfa5948c274
	github.com/XWS-BSEP-TIM2/dislinkt-backend/common v0.0.0-20220617075150-48a21e06c92f
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/joho/godotenv v1.4.0
	github.com/opentracing/opentracing-go v1.2.0
	go.mongodb.org/mongo-driver v1.9.1
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	rsc.io/qr v0.2.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.2.0 // indirect
	github.com/nats-io/nats.go v1.14.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)

require (
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/pkg/errors v0.9.1
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/text v0.3.7 // indirect
)
