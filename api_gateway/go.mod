module github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway

go 1.18

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/common => ../common

require (
	github.com/XWS-BSEP-TIM2/dislinkt-backend/common v0.0.0-20220426203826-7d3aeb3f989f
	github.com/gorilla/handlers v1.5.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.9.0
	google.golang.org/grpc v1.45.0
)

require (
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	google.golang.org/genproto v0.0.0-20220414192740-2d67ff6cf2b4 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
