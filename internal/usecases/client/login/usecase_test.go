package login_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	keycloakclient "github.com/evgeniy-krivenko/chat-service/internal/clients/keycloak"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/usecases/client/login"
	loginmocks "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/login/mocks"
)

const (
	token               = "token"
	firstName, lastName = "Eric", "Cartman"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl         *gomock.Controller
	authClient   *loginmocks.MockauthClient
	usrGetter    *loginmocks.MockuserGetter
	problemsRepo *loginmocks.MockprofilesRepository
	uCase        login.UseCase

	login    string
	password string
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.authClient = loginmocks.NewMockauthClient(s.ctrl)
	s.usrGetter = loginmocks.NewMockuserGetter(s.ctrl)
	s.problemsRepo = loginmocks.NewMockprofilesRepository(s.ctrl)

	s.login, s.password = "client", "password"

	var err error
	s.uCase, err = login.New(login.NewOptions(s.authClient, s.usrGetter, s.problemsRepo, "test", "test"))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := login.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	s.Require().Error(err)
	s.ErrorIs(err, login.ErrInvalidRequest)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestAuthError() {
	// Arrange.
	req := login.Request{
		Login:    s.login,
		Password: s.password,
	}

	s.authClient.EXPECT().Auth(gomock.Any(), req.Login, req.Password).
		Return(nil, errors.New("unexpected"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	s.Require().Error(err)
	s.ErrorIs(err, login.ErrAuthClient)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestGetUserInfoError() {
	// Arrange.
	req := login.Request{
		Login:    s.login,
		Password: s.password,
	}

	s.authClient.EXPECT().Auth(gomock.Any(), req.Login, req.Password).
		Return(&keycloakclient.RPT{AccessToken: token}, nil)
	s.usrGetter.EXPECT().GetUserInfoFromToken(token).
		Return(nil, errors.New("unexpected"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	s.Require().Error(err)
	s.ErrorIs(err, login.ErrAuthClient)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestClientNoResourceAccessError() {
	// Arrange.
	req := login.Request{
		Login:    s.login,
		Password: s.password,
	}
	clientID := types.NewUserID()

	s.authClient.EXPECT().Auth(gomock.Any(), req.Login, req.Password).
		Return(&keycloakclient.RPT{AccessToken: token}, nil)
	s.usrGetter.EXPECT().GetUserInfoFromToken(token).
		Return(
			&keycloakclient.User{
				ID:        clientID,
				FirstName: firstName,
				LastName:  lastName,
				ResourcesAccess: map[string]struct {
					Roles []string `json:"roles"`
				}{},
			},
			nil,
		)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, login.ErrNoResourceAccess)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestClientCreateOrUpdateError() {
	// Arrange.
	req := login.Request{
		Login:    s.login,
		Password: s.password,
	}
	clientID := types.NewUserID()

	s.authClient.EXPECT().Auth(gomock.Any(), req.Login, req.Password).
		Return(&keycloakclient.RPT{AccessToken: token}, nil)
	s.usrGetter.EXPECT().GetUserInfoFromToken(token).
		Return(
			&keycloakclient.User{
				ID:        clientID,
				FirstName: firstName,
				LastName:  lastName,
				ResourcesAccess: map[string]struct {
					Roles []string `json:"roles"`
				}{"test": {
					[]string{"test"},
				}},
			},
			nil,
		)
	s.problemsRepo.EXPECT().CreateOrUpdate(gomock.Any(), clientID, firstName, lastName).
		Return(errors.New("unexpected"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestSuccess() {
	// Arrange.
	req := login.Request{
		Login:    s.login,
		Password: s.password,
	}
	clientID := types.NewUserID()

	s.authClient.EXPECT().Auth(gomock.Any(), req.Login, req.Password).
		Return(&keycloakclient.RPT{AccessToken: token}, nil)
	s.usrGetter.EXPECT().GetUserInfoFromToken(token).
		Return(
			&keycloakclient.User{
				ID:        clientID,
				FirstName: firstName,
				LastName:  lastName,
				ResourcesAccess: map[string]struct {
					Roles []string `json:"roles"`
				}{"test": {
					[]string{"test"},
				}},
			},
			nil,
		)
	s.problemsRepo.EXPECT().CreateOrUpdate(gomock.Any(), clientID, firstName, lastName).
		Return(nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Equal(token, resp.Token)
	s.Equal(clientID, resp.ClientID)
	s.Equal(firstName, resp.FirstName)
	s.Equal(lastName, resp.LastName)
}
