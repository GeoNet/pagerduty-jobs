package finduser

import (
	"errors"

	pagerduty "github.com/PagerDuty/go-pagerduty"
)

var UserNotFoundError = errors.New("Could not find specified user.")

type Client struct {
	pagerduty.Client
}

func (c *Client) FindAndValidate(in string) (*pagerduty.User, error) {
	u, e := c.GetUser(in, pagerduty.GetUserOptions{})
	if e == nil && u.ID == in {
		return u, e
	}
	return c.FindUser(in)
}
func (c *Client) FindUser(in string) (*pagerduty.User, error) {
	resp, e := c.ListUsers(pagerduty.ListUsersOptions{Query: in})
	if e != nil {
		return nil, e
	}

	if len(resp.Users) == 1 {
		return &resp.Users[0], nil
	}

	return nil, UserNotFoundError
}
