//go:build integration

package profilesrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	profilesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/profiles"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var (
	firstName = "Eric"
	lastName  = "Cartman"
)

type ProfileRepoAPISuite struct {
	testingh.DBSuite
	repo *profilesrepo.Repo
}

func TestProfileRepoAPISuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProfileRepoAPISuite{DBSuite: testingh.NewDBSuite("TestProfileRepoAPISuite")})
}

func (s *ProfileRepoAPISuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error
	s.repo, err = profilesrepo.New(profilesrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProfileRepoAPISuite) Test_CreateOrUpdate() {
	s.Run("profile not exists", func() {
		authorID := types.NewUserID()

		err := s.repo.CreateOrUpdate(s.Ctx, authorID, firstName, lastName)

		s.Require().NoError(err)

		dbProfile, err := s.Database.Profile(s.Ctx).Get(s.Ctx, authorID)
		s.Require().NoError(err)
		s.NotEmpty(dbProfile)
		s.Equal(firstName, dbProfile.FirstName)
		s.Equal(lastName, dbProfile.LastName)
		s.NotZero(dbProfile.UpdatedAt)
		s.NotZero(dbProfile.CreatedAt)
		s.NotZero(dbProfile.ID)

		count := s.Database.Profile(s.Ctx).Query().CountX(s.Ctx)
		s.Equal(1, count)
	})

	s.Run("update exists", func() {
		authorID := types.NewUserID()
		newFirstName, newLastName := "Stan", "Marsh"
		creatingTime := time.Now()

		err := s.Database.Profile(s.Ctx).Create().
			SetID(authorID).
			SetUpdatedAt(creatingTime).
			SetCreatedAt(creatingTime).
			Exec(s.Ctx)
		s.Require().NoError(err)

		err = s.repo.CreateOrUpdate(s.Ctx, authorID, newFirstName, newLastName)
		s.Require().NoError(err)

		dbProfile, err := s.Database.Profile(s.Ctx).Get(s.Ctx, authorID)

		s.Require().NoError(err)
		s.NotEmpty(dbProfile)
		s.Equal(newFirstName, dbProfile.FirstName)
		s.Equal(newLastName, dbProfile.LastName)
		s.True(dbProfile.UpdatedAt.After(creatingTime))
		s.True(dbProfile.CreatedAt.Equal(creatingTime))
	})
}

func (s *ProfileRepoAPISuite) Test_GetProfileByID() {
	s.Run("success", func() {
		authorID := types.NewUserID()

		err := s.Database.Profile(s.Ctx).Create().
			SetID(authorID).
			SetFirstName(firstName).
			SetLastName(lastName).
			SetUpdatedAt(time.Now()).
			Exec(s.Ctx)
		s.Require().NoError(err)

		profile, err := s.repo.GetProfileByID(s.Ctx, authorID)
		s.Require().NoError(err)
		s.NotEmpty(profile)
		s.IsType(profilesrepo.Profile{}, *profile)
		s.Equal(authorID, profile.ID)
		s.Equal(firstName, profile.FirstName)
		s.Equal(lastName, profile.LastName)
	})

	s.Run("not found", func() {
		profile, err := s.repo.GetProfileByID(s.Ctx, types.NewUserID())
		s.Require().Error(err)
		s.Empty(profile)
		s.ErrorIs(err, profilesrepo.ErrProfileNotFound)
	})
}
