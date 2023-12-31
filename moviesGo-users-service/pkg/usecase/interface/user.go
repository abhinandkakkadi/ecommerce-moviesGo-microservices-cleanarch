package interfaces

import (
	"github.com/abhinandkakkadi/moviesgo-user-service/pkg/utils/models"
)

type UserUseCase interface {
	UserSignUp(user models.UserDetails) (int, error)
	UserLogin(userSignUpDetails models.UserLogin) (string, error)
	ValidateUser(signedToken string) (int, error)
	AddAddress(addressInfo models.AddressInfo) (int, error)
	GetUserAddress(userID int) ([]models.AddressInfoResponse, error)
}
