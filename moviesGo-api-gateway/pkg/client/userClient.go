package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/client/interceptors"
	interfaces "github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/client/interface"
	services "github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/client/interface"
	config "github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/config"
	"github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/user/pb"
	"github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/utils/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type userClient struct {
	client     pb.AuthServiceClient
	cartClient interfaces.CartClient
}

func NewUserClient(cfg config.Config, cartClient interfaces.CartClient) services.UserClient {

	grpcConnectoin, err := grpc.Dial(cfg.AuthSvcUrl, grpc.WithInsecure(), interceptors.UnaryClientInterceptor())
	if err != nil {
		fmt.Println("Could not connect", err)
	}

	grpcClient := pb.NewAuthServiceClient(grpcConnectoin)

	return &userClient{
		client:     grpcClient,
		cartClient: cartClient,
	}

}

func (u *userClient) UserAuthRequired(c *gin.Context) {

	token := c.GetHeader("authorization")

	userID, err := u.client.ValidateUser(context.Background(), &pb.ValidateRequest{
		Token: token,
	})

	fmt.Println("the error is: ", err)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("userID", userID.UserID)

}

func (u *userClient) SampleRequest(request string) (string, error) {

	res, err := u.client.SampleRequest(context.Background(), &pb.RegisterRequest{Request: request})
	if err != nil {
		return "", err
	}

	return res.Response, nil

}

func (u *userClient) SignUpRequest(userDetails models.UserDetails) (int, error) {

	res, err := u.client.UserSignUp(context.Background(), &pb.SingUpRequest{
		Name:            userDetails.Name,
		Email:           userDetails.Email,
		Phone:           userDetails.Phone,
		Password:        userDetails.Password,
		Confirmpassword: userDetails.ConfirmPassword,
	})

	if err != nil {
		return 0, err
	}

	return int(res.Status), nil

}

func (u *userClient) LoginHandler(userAuthDetails models.UserLogin) (string, error) {

	res, err := u.client.UserLogin(context.Background(), &pb.LoginInRequest{
		Email:    userAuthDetails.Email,
		Password: userAuthDetails.Password,
	})

	if err != nil {
		return "", err
	}

	return res.Token, nil
}

func (u *userClient) AddAddress(address models.AddressInfo, userID int) (int, error) {

	res, err := u.client.AddAddress(context.Background(),
		&pb.AddAddressRequest{
			Userid:    int64(userID),
			Name:      address.Name,
			Housename: address.HouseName,
			City:      address.City,
		})

	if err != nil {
		return 400, err
	}

	return int(res.Status), nil
}

func (u *userClient) GetAllAddress(userID int) ([]models.AddressInfoResponse, error) {

	res, err := u.client.GetAddress(context.Background(),
		&pb.GetAddressRequest{Userid: int64(userID)})

	if err != nil {
		return []models.AddressInfoResponse{}, err
	}

	var addressResponse []models.AddressInfoResponse
	for _, address := range res.Addresses {
		addressResponse = append(addressResponse, models.AddressInfoResponse{
			ID:        uint(address.Id),
			UserID:    uint(address.Userid),
			Name:      address.Name,
			HouseName: address.Housename,
			City:      address.City,
		})
	}

	return addressResponse, nil
}

func (u *userClient) Checkout(userID int) (models.CheckoutDetails, error) {

	address, err := u.GetAllAddress(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	cartsDetails, err := u.cartClient.DisplayCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	return models.CheckoutDetails{
		AddressInfoResponse: address,
		Cart:                cartsDetails.Cart,
		Grand_Total:         cartsDetails.TotalPrice,
	}, nil

}
