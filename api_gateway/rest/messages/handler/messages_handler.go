package handler

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbMessages "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageHandler struct {
	grpcClient *rest.ServiceClientGrpc
}

func InitMessageHandler() *MessageHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &MessageHandler{grpcClient: client}
}

func (handler *MessageHandler) GetMyContacts(ctx *gin.Context) {
	messageService := handler.grpcClient.MessageClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	myContacts, err := messageService.GetMyContacts(ctx, &pbMessages.GetMyContactsRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &myContacts)
}

func (handler *MessageHandler) GetChat(ctx *gin.Context) {
	messageService := handler.grpcClient.MessageClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	msgID := ctx.Param("chatId")
	getChatReq := pbMessages.GetChatRequest{UserID: dataFromToken.Id, MsgID: msgID}
	res, err := messageService.GetChat(ctx, &getChatReq)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *MessageHandler) SendMessage(ctx *gin.Context) {
	messageService := handler.grpcClient.MessageClient
	sendMessageDto := pbMessages.SendMessageRequest{}
	if err := ctx.BindJSON(&sendMessageDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := messageService.SendMessage(ctx, &sendMessageDto)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *MessageHandler) SetSeen(ctx *gin.Context) {
	messageService := handler.grpcClient.MessageClient
	setSeenDto := pbMessages.SetSeenRequest{}
	if err := ctx.BindJSON(&setSeenDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := messageService.SetSeen(ctx, &setSeenDto)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}
