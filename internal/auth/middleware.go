package auth

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/cache"
	"github.com/robitooS/api-service-go/internal/domain/user"
)

func AuthenticateHMAC(hmacKey []byte, repository user.UserRepository, cache cache.NonceStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// p fazer a nova verificação, precisa do método, path, ts, body e o nonce
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error":"body inválido"})
			return
		}
		// restaurar o body já q lemos ele antes 
		ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		tsStr, authHeader, nonce, userIDStr := ctx.GetHeader("X-Timestamp"), ctx.GetHeader("Authorization"), ctx.GetHeader("X-Nonce"), ctx.GetHeader("X-User-ID")

		if tsStr == "" || nonce == "" || authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "headers de autenticação ausentes"})
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

		if err := cache.CacheNonce(nonce); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"nonce já utilizado"})
			return
		}

		msg := BuildMessage(method, path, ts, string(bodyBytes), nonce)
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
