package middlewares

import (
	"github.com/ak-karimzai/bank-api/internel/token"
	"github.com/gin-gonic/gin"
)

func GetPayload(ctx *gin.Context) *token.Payload {
	authPayload := ctx.MustGet(
		AuthorizationPayloadKey).(*token.Payload)
	return authPayload
}
