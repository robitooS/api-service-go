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
	UserID       int64  `json:"user_id"`
}

type AddressHandler struct {
	AddressService service.AddressService
}

func NewAddressHandler(addressService *service.AddressService) *AddressHandler {
	return &AddressHandler{AddressService: *addressService}
}

func (h *AddressHandler) CreateAddress(ctx *gin.Context)  {
	request := AddressRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"body da requisição inválido"})
		return
	}

	address, err := h.AddressService.Create(ctx, request.Street, request.Number, request.Neighborhood, request.City, request.State, request.CEP, request.UserID)
	if err != nil {
		log.Printf("Erro ao criar endereço no banco: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"não foi possível criar o endereço"})
		return
	}

	ctx.JSON(http.StatusCreated, address)
}

func (h *AddressHandler) UpdateAddress(ctx *gin.Context)  {
	request := AddressRequest{} // struct que pega da requisição os novos dados do endereço (reutilizou o de criar por ser os mesmos campos)
	if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error":"body da requisição inválido"})
			return
		}

	addressToUpdate, err := address.NewAddress(request.Street, request.Number, request.Neighborhood, request.City, request.State, request.CEP, request.UserID)
	if err != nil {
		log.Printf("Erro ao atualizar endereço no banco: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"não foi possível atualizar o endereço"})
		return
	}
	
	addressUpdated, err := h.AddressService.Update(ctx, request.UserID, addressToUpdate)
	if err != nil {
		log.Printf("Erro ao atualizar endereço no banco: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"não foi possível atualizar o endereço"})
		return
	}

	ctx.JSON(http.StatusOK, addressUpdated)
}