package login

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrAuthClient       = errors.New("auth client")
	ErrNoResourceAccess = errors.New("no access to resource")
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=loginmocks

type authClient interface {
	Auth(ctx context.Context, username, password string) (*keycloakclient.RPT, error)
}

type userGetter interface {
	GetUserInfoFromToken(token string) (*keycloakclient.User, error)
}

type profilesRepository interface {
	CreateOrUpdate(ctx context.Context, id types.UserID, firstName, lastName string) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	authClient   authClient         `option:"mandatory" validate:"required"`
	usrGetter    userGetter         `option:"mandatory" validate:"required"`
	profilesRepo profilesRepository `option:"mandatory" validate:"required"`

	resource string `option:"mandatory" validate:"required"`
	role     string `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	return UseCase{Options: opts}, opts.Validate()
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrInvalidRequest, err)
	}

	token, err := u.authClient.Auth(ctx, req.Login, req.Password)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrAuthClient, err)
	}

	user, err := u.usrGetter.GetUserInfoFromToken(token.AccessToken)
	if err != nil {
		return Response{}, fmt.Errorf("%w: %v", ErrAuthClient, err)
	}

	if !user.ResourcesAccess.HasResourceRole(u.resource, u.role) {
		return Response{}, ErrNoResourceAccess
	}

	if err := u.profilesRepo.CreateOrUpdate(ctx, user.ID, user.FirstName, user.LastName); err != nil {
		return Response{}, fmt.Errorf("create or update profile: %v", err)
	}

	return Response{
		Token:     token.AccessToken,
		ClientID:  user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
