package getchats_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
	getchats "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chats"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request getchats.Request
		wantErr bool
	}{
		// Positive
		{
			name: "all ID's is valid",
			request: getchats.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
			},
			wantErr: false,
		},

		// Negative
		{
			name: "req id is invalid",
			request: getchats.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.NewUserID(),
			},
			wantErr: true,
		},
		{
			name: "manager id is invalid",
			request: getchats.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.UserIDNil,
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
