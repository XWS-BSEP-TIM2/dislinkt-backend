module github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service

go 1.18

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/common => ../common

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service => ../connection_service

require (
	github.com/XWS-BSEP-TIM2/dislinkt-backend/common v0.0.0-20220426203826-7d3aeb3f989f
	github.com/joho/godotenv v1.4.0
	github.com/neo4j/neo4j-go-driver/v4 v4.3.3
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.9.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
