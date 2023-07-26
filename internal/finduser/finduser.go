package finduser

import (
	"context"
	"errors"

	pagerduty "github.com/PagerDuty/go-pagerduty"
)

var UserNotFoundError = errors.New("Could not find specified user.")

type Client struct {
	pagerduty.Client
}

func (c *Client) FindAndValidate(ctx context.Context, in string) (*pagerduty.User, error) {
	u, e := c.GetUserWithContext(ctx, in, pagerduty.GetUserOptions{})
	if e == nil && u.ID == in {
		return u, e
	}
	return c.FindUser(ctx, in)
}
func (c *Client) FindUser(ctx context.Context, in string) (*pagerduty.User, error) {
	resp, e := c.ListUsersWithContext(ctx, pagerduty.ListUsersOptions{Query: in})
	if e != nil {
		return nil, e
	}

	if len(resp.Users) == 1 {
		return &resp.Users[0], nil
	}

	return nil, UserNotFoundError
}
