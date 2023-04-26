package clientv1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

var stub = GetHistoryResponse{
	Data: MessagesPage{Messages: []Message{
		{
			AuthorId:  types.NewUserID(),
			Body:      "Здравствуйте! Разберёмся.",
			CreatedAt: time.Now(),
			Id:        types.NewMessageID(),
		},
		{
			AuthorId:  types.MustParse[types.UserID]("bc7e300b-29e4-47d5-bc90-90ca0046f9f7"),
			Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
			CreatedAt: time.Now().Add(-time.Minute),
			Id:        types.NewMessageID(),
		},
	}},
	Next: "",
}

func (h Handlers) PostGetHistory(
	eCtx echo.Context,
	_ PostGetHistoryParams,
) error {
	return eCtx.JSON(http.StatusOK, &stub)
}
