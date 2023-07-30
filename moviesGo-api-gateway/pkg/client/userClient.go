package client

import (
	"context"
	"fmt"
	"net/http"

	services "github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/client/interface"
	config "github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/config"
	"github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/user/pb"
	"github.com/abhinandkakkadi/moviesgo-api-gateway/pkg/utils/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type userClient struct {
	client pb.AuthServiceClient
}

func NewUserClient(cfg config.Config) services.UserClient {

	grpcConnectoin, err := grpc.Dial(cfg.AuthSvcUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect", err)
	}

	grpcClient := pb.NewAuthServiceClient(grpcConnectoin)

	return &userClient{
		client: grpcClient,
	}

}

func (u *userClient) UserAuthRequired(c *gin.Context)  {

	token := c.GetHeader("authorization")

	userID, err := u.client.ValidateUser(context.Background(),&pb.ValidateRequest{
		Token: token,
	})

	fmt.Println("the error is: ",err)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("userID",userID.UserID)


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

func (u *userClient) LoginHandler(userAuthDetails models.UserLogin) (string,error) {

	res, err := u.client.UserLogin(context.Background(),&pb.LoginInRequest{
		Email: userAuthDetails.Email,
		Password: userAuthDetails.Password,
	})

	if err != nil {
		return "",err
	}

	return res.Token,nil
}