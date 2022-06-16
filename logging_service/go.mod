module github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service

go 1.18

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/common => ../common

require (
	github.com/XWS-BSEP-TIM2/dislinkt-backend/common v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.1.2
	github.com/joho/godotenv v1.4.0
	google.golang.org/grpc v1.45.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
