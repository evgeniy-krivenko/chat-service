package profilesrepo

import (
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
)

//go:generate options-gen -out-filename=repo_options.gen.go -from-struct=Options
type Options struct {
	db *store.Database `option:"mandatory" validate:"required"`
}

type Repo struct {
	Options
}

func New(opts Options) (*Repo, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate problem repo options: %v", err)
	}
	return &Repo{Options: opts}, nil
}
