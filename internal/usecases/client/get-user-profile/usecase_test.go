package getuserprofile_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	profilesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/profiles"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	getuserprofile "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-user-profile"
	getuserprofilemocks "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-user-profile/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl         *gomock.Controller
	profilesRepo *getuserprofilemocks.MockprofilesRepository
	uCase        getuserprofile.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.profilesRepo = getuserprofilemocks.NewMockprofilesRepository(s.ctrl)

	var err error
	s.uCase, err = getuserprofile.New(getuserprofile.NewOptions(s.profilesRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := getuserprofile.Request{}

	// Action
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getuserprofile.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestRequestProfileNotFoundError() {
	// Arrange.
	userID := types.NewUserID()

	req := getuserprofile.Request{
		ID:     types.NewRequestID(),
		UserID: userID,
	}

	s.profilesRepo.EXPECT().GetProfileByID(gomock.Any(), userID).Return(nil, profilesrepo.ErrProfileNotFound)

	// Action
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getuserprofile.ErrProfileNotFound)
}

func (s *UseCaseSuite) TestRequestProfileUnexpectedError() {
	// Arrange.
	userID := types.NewUserID()

	req := getuserprofile.Request{
		ID:     types.NewRequestID(),
		UserID: userID,
	}

	s.profilesRepo.EXPECT().GetProfileByID(gomock.Any(), userID).Return(nil, errors.New("unexpected"))

	// Action
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.NotErrorIs(err, getuserprofile.ErrProfileNotFound)
}

func (s *UseCaseSuite) TestRequestProfileSuccess() {
	// Arrange.
	userID := types.NewUserID()
	firstName, lastName := "Eric", "Cartman"

	req := getuserprofile.Request{
		ID:     types.NewRequestID(),
		UserID: userID,
	}

	s.profilesRepo.EXPECT().GetProfileByID(gomock.Any(), userID).Return(&profilesrepo.Profile{
		ID:        userID,
		FirstName: firstName,
		LastName:  lastName,
	}, nil)

	// Action
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.NotEmpty(resp)
	s.Equal(userID, resp.UserID)
	s.Equal(firstName, resp.FirstName)
	s.Equal(lastName, resp.LastName)
}
