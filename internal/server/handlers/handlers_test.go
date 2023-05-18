package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EgorKo25/GophKeeper/internal/server/handlers"

	"github.com/stretchr/testify/assert"

	"github.com/EgorKo25/GophKeeper/pkg/auth"

	mock_database "github.com/EgorKo25/GophKeeper/internal/database/mocks"
	"github.com/EgorKo25/GophKeeper/internal/storage"

	"github.com/golang/mock/gomock"
)

func TestHandler_Register(t *testing.T) {
	type fields struct {
		db *mock_database.MockDatabase
	}

	tests := []struct {
		name           string
		fields         fields
		prepare        func(f *fields)
		expectedStatus int
		au             *auth.Auth
		request        string
		user           storage.User
	}{
		{
			name: "success",
			prepare: func(f *fields) {

				ctx := context.Background()
				user := storage.User{
					Login:    "testuser",
					Password: "testpassword",
					Email:    "testemail@test.com",
				}

				gomock.InOrder(
					f.db.EXPECT().Read(ctx, &user, "testuser").Return(nil, nil),
					f.db.EXPECT().Add(ctx, &user, "testuser").Return(nil),
				)

			},
			request: "/user/add",
			user: storage.User{
				Login:    "testuser",
				Password: "testpassword",
				Email:    "testemail@test.com",
			},
			expectedStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			au := auth.NewAuth("some-secret")

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, err := json.Marshal(tt.user)
			if err != nil {
				t.Errorf("err unmarshal: %s", err)
			}

			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(body))

			f := &fields{
				db: mock_database.NewMockDatabase(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			w := httptest.NewRecorder()

			h := handlers.Handler{Db: f.db, Au: au}

			handle := http.HandlerFunc(h.Register)

			handle(w, request)

			result := w.Result()

			assert.Equal(t, tt.expectedStatus, result.StatusCode)
		})
	}
}
