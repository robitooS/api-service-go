package address

import "context"

type AddressRepository interface {
	Create(ctx context.Context, address *Address) (*Address, error)
	Update(ctx context.Context, userID int64, newAddress *Address) (*Address, error)
	FindByID(ctx context.Context, id int64) (*Address, error)
	FindByUserId(ctx context.Context, userID int64) (*Address, error)
}