package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robitooS/api-service-go/internal/domain/address"
	"github.com/robitooS/api-service-go/internal/service"
)

type AddressRequest struct {
	Street       string `json:"address_street"`
	Number       string `json:"address_number"`
	Neighborhood string `json:"address_neighborhood"`
	City         string `json:"address_city"`
	State        string `json:"address_state"`
	CEP          string `json:"address_cep"`
}

type AddressHandler struct {
	AddressService service.AddressService
}

func NewAddressHandler(addressService *service.AddressService) *AddressHandler {
	return &AddressHandler{AddressService: *addressService}
}

func (h *AddressHandler) CreateAddress(ctx *gin.Context)  {
	var request AddressRequest // Usar var para declarar
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"body da requisição inválido"})
		return
	}
	
	userIDFromCtx, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado no contexto"})
		return
	}

	userID, ok := userIDFromCtx.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "formato de ID de usuário inválido no contexto"})
		return
	}

	address, err := h.AddressService.Create(
		ctx, 
		request.Street, 
		request.Number, 
		request.Neighborhood, 
		request.City, 
		request.State, 
		request.CEP, 
		userID, 
	)

	if err != nil {
		log.Printf("Erro ao criar endereço no banco: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, address)
}

type UpdateAddressRequest struct {
	Street       string `json:"address_street"`
	Number       string `json:"address_number"`
	Neighborhood string `json:"address_neighborhood"`
	City         string `json:"address_city"`
	State        string `json:"address_state"`
	CEP          string `json:"address_cep"`
	UserID       int64  `json:"user_id"`
}


func (h *AddressHandler) UpdateAddress(ctx *gin.Context)  {
	var request UpdateAddressRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"body da requisição inválido"})
		return
	}

	addressToUpdate, err := address.NewAddress(request.Street, request.Number, request.Neighborhood, request.City, request.State, request.CEP, request.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	addressUpdated, err := h.AddressService.Update(ctx, request.UserID, addressToUpdate)
	if err != nil {
		log.Printf("Erro ao atualizar endereço no banco: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, addressUpdated)
}