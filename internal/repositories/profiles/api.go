package profilesrepo

import (
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"golang.org/x/net/context"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	storeprofile "github.com/evgeniy-krivenko/chat-service/internal/store/profile"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var ErrProfileNotFound = errors.New("profile not found")

func (r *Repo) CreateOrUpdate(ctx context.Context, id types.UserID, firstName, lastName string) error {
	err := r.db.Profile(ctx).Create().
		SetID(id).
		SetFirstName(firstName).
		SetLastName(lastName).
		SetUpdatedAt(time.Now()).
		OnConflict(
			sql.ConflictColumns(storeprofile.FieldID),
			sql.ResolveWith(func(set *sql.UpdateSet) {
				set.SetIgnore(storeprofile.FieldCreatedAt)
			}),
		).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("upsert profile: %v", err)
	}

	return nil
}

func (r *Repo) GetProfileByID(ctx context.Context, id types.UserID) (profile *Profile, err error) {
	p, err := r.db.Profile(ctx).Get(ctx, id)
	if nil == err {
		profile = adaptProfile(p)
		return
	}
	if store.IsNotFound(err) {
		return nil, fmt.Errorf("%w: %v", ErrProfileNotFound, err)
	}
	if err != nil {
		return nil, fmt.Errorf("get profile: %v", err)
	}
	return
}
