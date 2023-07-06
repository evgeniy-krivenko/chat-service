package profilesrepo

import (
	"github.com/evgeniy-krivenko/chat-service/internal/store"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type Profile struct {
	ID        types.UserID
	FirstName string
	LastName  string
}

func adaptProfile(p *store.Profile) *Profile {
	return &Profile{
		ID:        p.ID,
		FirstName: p.FirstName,
		LastName:  p.LastName,
	}
}
