module github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service

go 1.18

replace github.com/tamararankovic/microservices_demo/common => ../common

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service => ../auth_service

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	go.mongodb.org/mongo-driver v1.8.4
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.9.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

require (
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.14.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tamararankovic/microservices_demo/common v1.0.0
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9 // indirect
	golang.org/x/text v0.3.7 // indirect
)
