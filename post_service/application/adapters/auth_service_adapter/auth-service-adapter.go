package auth_service_adapter

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/helper"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters"
	lsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/logging_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthServiceAdapter struct {
	address               string
	loggingServiceAdapter lsa.ILoggingServiceAdapter
}

func NewAuthServiceAdapter(address string, loggingServiceAdapter lsa.ILoggingServiceAdapter) *AuthServiceAdapter {
	return &AuthServiceAdapter{address: address, loggingServiceAdapter: loggingServiceAdapter}
}

func (auth *AuthServiceAdapter) GetRequesterId(ctx context.Context) primitive.ObjectID {
	span := tracer.StartSpanFromContext(ctx, "GetRequesterId")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	authClient := adapters.NewAuthClient(auth.address)
	token := getRequesterToken(ctx2)
	if token == "" { // user is not logged in
		return primitive.NilObjectID
	}
	currentUserData, dataExtractionErr := authClient.ExtractDataFromToken(ctx2, &authService.ExtractDataFromTokenRequest{Token: token})
	if dataExtractionErr != nil {
		panic(fmt.Errorf("Error during data extraction from token : Auth Service"))
	}

	requesterId, IdFromHexErr := primitive.ObjectIDFromHex(currentUserData.Id)
	if IdFromHexErr != nil {
		panic(fmt.Errorf("Error during userId extraction from token : Auth Service"))
	}
	return requesterId
}

func (auth *AuthServiceAdapter) ValidateCurrentUser(ctx context.Context, userId primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "ValidateCurrentUser")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	requesterId := auth.GetRequesterId(ctx2)
	if requesterId != userId {
		message := fmt.Sprintf("Current user (id: %s) is trying to take action on behalf user with id %s", requesterId.Hex(), userId.Hex())
		auth.loggingServiceAdapter.Log(ctx2, "WARNING", "ValidateCurrentUser", requesterId.Hex(), message)
		panic(errors.NewEntityForbiddenError(message))
	}
}

func getRequesterToken(ctx context.Context) string {
	span := tracer.StartSpanFromContext(ctx, "getRequesterToken")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	token, tokenExtractionErr := helper.ExtractTokenFromContext(ctx2)
	if tokenExtractionErr != nil {
		if token == "invalid_auth_header" {
			panic(fmt.Errorf("Error during token extraction"))
		} else {
			return ""
		}
	}
	return token
}
