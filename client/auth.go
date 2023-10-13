package sonic_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// SignUp Sign up - Create new user, return user token upon successful creation
func (c *Client) SignUp(auth AuthStruct) (*AuthResponse, error) {
	if auth.Username == "" || auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}
	rb, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signup", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

// SignIn - Get a new token for user
func (c *Client) SignIn() (*AuthResponse, error) {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}
	rb, err := json.Marshal(c.Auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/authenticate", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

// GetUserTokenSignIn SignIn - Get a new token for user
func (c *Client) GetUserTokenSignIn(auth AuthStruct) (*AuthResponse, error) {
	if auth.Username == "" || auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}
	rb, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signin", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, errors.New("Unable to login")
	}

	ar := AuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

// SignOut - Revoke the token for a user
func (c *Client) SignOut(authToken *string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/signout", c.HostURL), strings.NewReader(string("")))
	if err != nil {
		return err
	}

	body, err := c.doRequest(req, authToken)
	if err != nil {
		return err
	}

	if string(body) != "Signed out user" {
		return errors.New(string(body))
	}

	return nil
}
