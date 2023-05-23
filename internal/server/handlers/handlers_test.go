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
		{
			name:    "missing login",
			prepare: func(f *fields) {},
			request: "/user/add",
			user: storage.User{
				Password: "testpassword",
				Email:    "testemail@test.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "missing password",
			prepare: func(f *fields) {},
			request: "/user/add",
			user: storage.User{
				Email: "testemail@test.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "missing email",
			prepare: func(f *fields) {},
			request: "/user/add",
			user: storage.User{
				Password: "testpassword",
			},
			expectedStatus: http.StatusBadRequest,
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

func TestHandler_Login(t *testing.T) {
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
				}

				gomock.InOrder(
					f.db.EXPECT().CheckUser(ctx, &user).Return(true, nil),
				)

			},
			request: "/user/login",
			user: storage.User{
				Login:    "testuser",
				Password: "testpassword",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "missing login",
			prepare: func(f *fields) {

				ctx := context.Background()
				user := storage.User{
					Password: "testpassword",
				}

				gomock.InOrder(
					f.db.EXPECT().CheckUser(ctx, &user).Return(false, nil),
				)

			},
			request: "/user/login",
			user: storage.User{
				Password: "testpassword",
			},

			expectedStatus: http.StatusForbidden,
		},
		{
			name: "missing password",
			prepare: func(f *fields) {
				ctx := context.Background()
				user := storage.User{
					Login: "testuser",
				}

				gomock.InOrder(
					f.db.EXPECT().CheckUser(ctx, &user).Return(false, nil),
				)
			},
			request: "/user/login",
			user: storage.User{
				Login: "testuser",
			},
			expectedStatus: http.StatusForbidden,
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

			handle := http.HandlerFunc(h.Login)

			handle(w, request)

			result := w.Result()

			assert.Equal(t, tt.expectedStatus, result.StatusCode)
		})
	}
}

func TestHandler_Add(t *testing.T) {
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
		pass           any
		dataType       string
	}{
		{
			name:    "empty header Data-Type",
			prepare: func(f *fields) {},

			request: "/user/add",
			pass: storage.Password{
				Service:    "yandex",
				LoginOwner: "testuser",
				Login:      "testuser",
				Password:   "testpassword",
			},
			expectedStatus: http.StatusBadRequest,
			dataType:       "",
		},
		{
			name: "success add password",
			prepare: func(f *fields) {

				ctx := context.Background()

				gomock.InOrder(
					f.db.EXPECT().Add(
						ctx,
						&storage.Password{
							Service:    "yandex",
							LoginOwner: "testuser",
							Login:      "testuser",
							Password:   "testpassword",
						},
						"testuser",
					).Return(nil),
				)
			},

			request: "/user/add",
			pass: storage.Password{
				Service:    "yandex",
				LoginOwner: "testuser",
				Login:      "testuser",
				Password:   "testpassword",
			},
			expectedStatus: http.StatusOK,
			dataType:       "password",
		},
		{
			name: "success add card",
			prepare: func(f *fields) {

				ctx := context.Background()

				gomock.InOrder(
					f.db.EXPECT().Add(
						ctx,
						&storage.Card{
							Bank:       "yandex",
							LoginOwner: "testuser",
							Number:     "123321123321",
							Owner:      "testpassword",
							SecretCode: "031",
							DataEnd:    "12/58",
						},
						"testuser",
					).Return(nil),
				)
			},

			request: "/user/add",
			pass: storage.Card{
				Bank:       "yandex",
				LoginOwner: "testuser",
				Number:     "123321123321",
				Owner:      "testpassword",
				SecretCode: "031",
				DataEnd:    "12/58",
			},
			expectedStatus: http.StatusOK,
			dataType:       "card",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			au := auth.NewAuth("some-secret")

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, err := json.Marshal(tt.pass)
			if err != nil {
				t.Errorf("err marshal: %s", err)
			}

			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(body))

			request.Header.Set("Data-Type", tt.dataType)

			f := &fields{
				db: mock_database.NewMockDatabase(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			w := httptest.NewRecorder()

			cook := &http.Cookie{
				Name:  "User",
				Value: "testuser",
			}

			request.AddCookie(cook)

			h := handlers.Handler{Db: f.db, Au: au}

			handle := http.HandlerFunc(h.Add)

			handle(w, request)

			result := w.Result()

			assert.Equal(t, tt.expectedStatus, result.StatusCode)
		})
	}
}

func TestHandler_Delete(t *testing.T) {
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
		pass           any
		dataType       string
	}{
		{
			name:    "empty header Data-Type",
			prepare: func(f *fields) {},

			request: "/user/delete",
			pass: storage.Password{
				Service:    "yandex",
				LoginOwner: "testuser",
			},
			expectedStatus: http.StatusBadRequest,
			dataType:       "",
		},
		{
			name: "success delete password",
			prepare: func(f *fields) {

				ctx := context.Background()

				gomock.InOrder(
					f.db.EXPECT().Delete(
						ctx,
						&storage.Password{
							Service:    "yandex",
							LoginOwner: "testuser",
						},
						"testuser",
					).Return(nil),
				)
			},

			request: "/user/delete",
			pass: storage.Password{
				Service:    "yandex",
				LoginOwner: "testuser",
			},
			expectedStatus: http.StatusOK,
			dataType:       "password",
		},
		{
			name: "success delete card",
			prepare: func(f *fields) {

				ctx := context.Background()

				gomock.InOrder(
					f.db.EXPECT().Delete(
						ctx,
						&storage.Card{
							Bank:       "yandex",
							LoginOwner: "testuser",
						},
						"testuser",
					).Return(nil),
				)
			},

			request: "/user/delete",
			pass: storage.Card{
				Bank:       "yandex",
				LoginOwner: "testuser",
			},
			expectedStatus: http.StatusOK,
			dataType:       "card",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			au := auth.NewAuth("some-secret")

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			body, err := json.Marshal(tt.pass)
			if err != nil {
				t.Errorf("err marshal: %s", err)
			}

			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(body))

			request.Header.Set("Data-Type", tt.dataType)

			f := &fields{
				db: mock_database.NewMockDatabase(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(f)
			}

			w := httptest.NewRecorder()

			cook := &http.Cookie{
				Name:  "User",
				Value: "testuser",
			}

			request.AddCookie(cook)

			h := handlers.Handler{Db: f.db, Au: au}

			handle := http.HandlerFunc(h.Delete)

			handle(w, request)

			result := w.Result()

			assert.Equal(t, tt.expectedStatus, result.StatusCode)
		})
	}
}
