package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/robitooS/api-service-go/internal/domain/address"
	"modernc.org/sqlite"
)

type SQLiteAddressRepository struct {
	DB *sql.DB
}

func NewAddressRepository(db *sql.DB) address.AddressRepository {
	return &SQLiteAddressRepository{DB: db}
}

// Create implements address.AddressRepository.
func (r *SQLiteAddressRepository) Create(ctx context.Context, address *address.Address) (*address.Address, error) {
	query := `
        INSERT INTO address (
            address_address, 
            address_number, 
            address_neighborhood, 
            address_city, 
            address_state, 
            address_cep, 
            user_id
        ) VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := r.DB.ExecContext(ctx, query,
		address.Address,
		address.Number,
		address.Neighborhood,
		address.City,
		address.State,
		address.CEP,
		address.UserID,
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == 2067 { // SQLITE_CONSTRAINT_UNIQUE
			return nil, fmt.Errorf("o usuário já possui um endereço cadastrado em seu nome")
		}

		return nil, fmt.Errorf("não foi possível inserir o endereço no banco - %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("não foi possivel obter o ID do último endereço inserido - %w", err)
	}

	return r.FindByID(ctx, id)
}

// FindByID implements address.AddressRepository.
func (r *SQLiteAddressRepository) FindByID(ctx context.Context, id int64) (*address.Address, error) {
	query := `
		SELECT 
			address_id,
			address_address,
			address_number,
			address_neighborhood,
			address_city,
			address_state,
			address_cep,
			user_id,
			address_created_at
		FROM address WHERE address_id = ?`

	var (
		addressID        int64
		addressAddress   string
		addressNumber    string
		addressNeighbor  string
		addressCity      string
		addressState     string
		addressCEP       string
		userID           int64
		addressCreatedAt time.Time
	)

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&addressID,
		&addressAddress,
		&addressNumber,
		&addressNeighbor,
		&addressCity,
		&addressState,
		&addressCEP,
		&userID,
		&addressCreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("endereço não encontrado")
		}
		return nil, fmt.Errorf("não foi possível buscar o endereço: %w", err)
	}

	addr := address.Address{
		ID:           addressID,
		Address:      addressAddress,
		Number:       addressNumber,
		Neighborhood: addressNeighbor,
		City:         addressCity,
		State:        addressState,
		CEP:          addressCEP,
		UserID:       userID,
		CreatedAt:    addressCreatedAt,
	}

	return &addr, nil
}

// Update implements address.AddressRepository.
func (r *SQLiteAddressRepository) Update(ctx context.Context, userID int64, newAddress *address.Address) (*address.Address, error) {
	query := `
		UPDATE address
		SET address_address = ?,
			address_number = ?,
			address_neighborhood = ?,
			address_city = ?,
			address_state = ?,
			address_cep = ?
		WHERE user_id = ?`

	res, err := r.DB.ExecContext(ctx, query, newAddress.Address, newAddress.Number, newAddress.Neighborhood, newAddress.City, newAddress.State, newAddress.CEP, userID)
	if err != nil {
		return nil, fmt.Errorf("não foi possível atualizar o endereço do usuário: %w", err)
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("nenhum endereço foi alterado, verificar se o usuário possui um endereço cadastrado")
	}
	if err != nil {
		return nil, err
	}

	add, err := r.FindByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return add, nil
}

// FindByUserId implements address.AddressRepository.
func (r *SQLiteAddressRepository) FindByUserId(ctx context.Context, userId int64) (*address.Address, error) {
	query := `
		SELECT 
			address_id,
			address_address,
			address_number,
			address_neighborhood,
			address_city,
			address_state,
			address_cep,
			user_id,
			address_created_at
		FROM address WHERE user_id = ?`

	var (
		addressID        int64
		addressAddress   string
		addressNumber    string
		addressNeighbor  string
		addressCity      string
		addressState     string
		addressCEP       string
		userID           int64
		addressCreatedAt time.Time
	)

	err := r.DB.QueryRowContext(ctx, query, userId).Scan(
		&addressID,
		&addressAddress,
		&addressNumber,
		&addressNeighbor,
		&addressCity,
		&addressState,
		&addressCEP,
		&userID,
		&addressCreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("endereço não encontrado")
		}
		return nil, fmt.Errorf("não foi possível buscar o endereço: %w", err)
	}

	addr := address.Address{
		ID:           addressID,
		Address:      addressAddress,
		Number:       addressNumber,
		Neighborhood: addressNeighbor,
		City:         addressCity,
		State:        addressState,
		CEP:          addressCEP,
		UserID:       userID,
		CreatedAt:    addressCreatedAt,
	}

	return &addr, nil
}