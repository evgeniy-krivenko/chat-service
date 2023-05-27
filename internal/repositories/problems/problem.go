package problemsrepo

import (
	"github.com/evgeniy-krivenko/chat-service/internal/store"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type Problem struct {
	ID     types.ProblemID
	ChatID types.ChatID
}

func adaptProblem(p *store.Problem) Problem {
	return Problem{
		ID:     p.ID,
		ChatID: p.ChatID,
	}
}
