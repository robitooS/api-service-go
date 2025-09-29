package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/robitooS/api-service-go/internal/domain/user"
	"github.com/gin-gonic/gin"
)

func AuthenticateHMAC(hmacKey []byte, repository user.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDStr, tsStr, authHeader := ctx.GetHeader("X-User-ID"), ctx.GetHeader("X-Timestamp"), ctx.GetHeader("Authorization")
		
		if userIDStr == "" || tsStr == "" || authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error":"headers ausentes"})
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error":"X-User-ID inválido"})
			return
		}
		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error":"X-Timestamp inválido"})
			return
		}

		if err := ValidateTimeStamp(ts); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"timestamp inválido ou expirado"})
			return
		}

		msg := BuildMessage(userID, ts)
		ok, err := ValidateSignature(msg, authHeader, hmacKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":"erro ao validar assinatura"})
			return
		}
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"assinatura inválida"})
			return
		}

		if !verifyUserExists(repository, ctx.Request.Context(), userID) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"usuário não encontrado"})
			return

		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func verifyUserExists(r user.UserRepository, ctx context.Context, userID int64) bool {
	// Futuramente é legal adicionar mais clareza para o erro
	// caso der erro de conexão, banco caiu, etc
	// vai continuar lançando que o user não foi encontrado
	_, err := r.FindByID(ctx, userID)
	return err == nil
}


