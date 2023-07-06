package getmanagerprofile

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	profilesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/profiles"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=getmanagerprofilemocks

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrProfileNotFound = errors.New("profile not found")
)

type profilesRepository interface {
	GetProfileByID(ctx context.Context, id types.UserID) (*profilesrepo.Profile, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	profilesRepo profilesRepository `option:"mandatory" validate:"required"`
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

	profile, err := u.profilesRepo.GetProfileByID(ctx, req.ManagerID)
	if errors.Is(err, profilesrepo.ErrProfileNotFound) {
		return Response{}, fmt.Errorf("%w: %v", ErrProfileNotFound, err)
	}
	if err != nil {
		return Response{}, fmt.Errorf("get profile: %v", err)
	}

	return Response{
		ManagerID: profile.ID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}, nil
}
