package finduser

import (
	"errors"

	"github.com/quiffman/go-pagerduty/pagerduty"
)

var UserNotFoundError = errors.New("Could not find specified user.")

type Service struct {
	pagerduty.UsersService
}

func (s *Service) FindAndValidate(in string) (*pagerduty.User, error) {
	u, _, e := s.Get(in)
	if e == nil && u.ID == in {
		return u, e
	}
	return s.FindUser(in)
}
func (s *Service) FindUser(in string) (*pagerduty.User, error) {
	users, _, e := s.List(&pagerduty.UsersOptions{Query: in})
	if e != nil {
		return nil, e
	}

	if len(users) == 1 {
		return &users[0], nil
	}

	return nil, UserNotFoundError
}
