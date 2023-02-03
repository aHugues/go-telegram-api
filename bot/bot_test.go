package bot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ahugues/go-telegram-api/structs"
)

func TestGetMeOK(t *testing.T) {
	t.Parallel()

	mockTGSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/bot110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw/getMe" {
			t.Errorf("Unexpected URL %v", r.URL)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true,"result":{"id":110201543,"is_bot":true,"first_name":"Test","username":"TestTheBot","can_join_groups":true,"can_read_all_group_messages":false,"supports_inline_queries":false}}`))
	}))
	defer mockTGSrv.Close()

	bot := ConcreteBot{
		httpClt: mockTGSrv.Client(),
		baseURL: mockTGSrv.URL,
		token:   "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw",
	}

	expectedUser := structs.User{
		ID:                      110201543,
		IsBot:                   true,
		FirstName:               "Test",
		LastName:                "",
		Username:                "TestTheBot",
		CanJoinGroups:           true,
		CanReadAllGroupMessages: false,
		SupportsInlineQueries:   false,
	}

	user, err := bot.GetMe(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}
	if !reflect.DeepEqual(user, expectedUser) {
		t.Fatalf("Unexpected user %+v", user)
	}
}

func TestGetMeInvalidToken(t *testing.T) {
	t.Parallel()

	mockTGSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"ok":false,"error_code":401,"description":"Unauthorized"}`))
	}))
	defer mockTGSrv.Close()

	bot := ConcreteBot{
		httpClt: mockTGSrv.Client(),
		baseURL: mockTGSrv.URL,
		token:   "110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw",
	}

	expectedErr := "Telegram API error (statuscode 401): Unauthorized"

	_, err := bot.GetMe(context.Background())
	if err == nil {
		t.Fatalf("Unexpected nil error")
	}
	if err.Error() != expectedErr {
		t.Fatalf("Unexpected error %s", err.Error())
	}
}
