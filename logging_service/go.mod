module github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service

go 1.18

replace github.com/XWS-BSEP-TIM2/dislinkt-backend/common => ../common

require (
	github.com/XWS-BSEP-TIM2/dislinkt-backend/common v0.0.0-20220617075150-48a21e06c92f
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.4.0
	google.golang.org/grpc v1.47.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
