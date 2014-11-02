package titanium

import (
	"errors"
	"fmt"

	"github.com/atomosio/common"
)

var _ = fmt.Printf

type CreateTokenRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type CreateTokenResponse struct {
	Response
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
}

func (client *HttpClient) Login(user, password string) error {
	request := CreateTokenRequest{
		User:     user,
		Password: password,
	}

	//send request
	response := CreateTokenResponse{}
	err := client.postAndUnmarshal(TokensEndpoint, &request, &response)
	if err != nil {
		return err
	}

	if response.Code != common.Success {
		return errors.New(response.Description)
	}
	client.Token = response.Token

	return nil
}
