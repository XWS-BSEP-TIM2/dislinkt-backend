package security

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/services"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"github.com/casbin/casbin"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type AccessDetails struct {
	TokenUuid string
	UserId    string
	UserName  string
}

type AuthMiddleware struct {
	authClient authService.AuthServiceClient
}

func NewAuthMiddleware(authClientAddress string) *AuthMiddleware {

	return &AuthMiddleware{
		authClient: services.NewAuthClient(authClientAddress),
	}
}

func (auth *AuthMiddleware) AuthRequired(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("authorization")

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	res, err := auth.authClient.Validate(context.Background(), &authService.ValidateRequest{
		Token: token[1],
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("userId", res.UserId)

	ctx.Next()
}

// Authorize determines if current subject has been authorized to take an action on an object.
func (auth *AuthMiddleware) Authorize(obj string, act string, isApiMethod bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		metadata, err := ExtractTokenDataFromContext(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(401, "error while trying to extract token data")
			return
		}
		if isApiMethod {
			response, _ := auth.authClient.ValidateApiToken(ctx, &authService.ValidateApiTokenRequest{TokenCode: metadata.ApiCode})
			if response.Error != nil {
				ctx.AbortWithStatusJSON(401, "Invalid Api Token")
			}

		}
		// casbin enforces policy
		var ok = false
		if isApiMethod {
			ok, err = enforce(metadata.TokenType, obj, act)
		} else {
			ok, err = enforce(metadata.Role, obj, act)
		}

		if err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(500, "error occurred when authorizing user")
			return
		}
		if !ok {
			ctx.AbortWithStatusJSON(403, "forbidden")
			return
		}
		ctx.Next()
	}
}

func ExtractTokenDataFromContext(ctx *gin.Context) (*authService.ExtractDataFromTokenResponse, error) {
	authorization := ctx.Request.Header.Get("authorization")
	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return nil, nil
	}
	token := strings.Split(authorization, "Bearer ")
	authS := services.NewAuthClient(fmt.Sprintf("%s:%s", config.NewConfig().AuthHost, config.NewConfig().AuthPort))
	metadata, err := authS.ExtractDataFromToken(ctx, &authService.ExtractDataFromTokenRequest{Token: token[1]})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return nil, err
	}
	return metadata, nil
}

func enforce(sub string, obj string, act string) (bool, error) {
	fileAdapter := fileadapter.NewAdapter("./security/rbac_policy.csv")
	enforcer := casbin.NewEnforcer("./security/rbac_model.conf", fileAdapter)
	err := enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	ok := enforcer.Enforce(sub, obj, act)
	return ok, nil
}
