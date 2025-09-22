package grafcookier

import (
	"log/slog"
	"os"
	"testing"
)

func TestGetCookie(t *testing.T) {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	testCookier, err := New("https://172.17.134.91:3000", "/login/generic_oauth", *l)
	if err != nil {
		t.Fatalf("Error create client: %v", err)
	}
	// Пример использования
	username := "admin"
	password := "admin"

	cookies, err := testCookier.GetCookies(username, password)
	if err != nil {
		t.Fatalf("Error getting cookies: %v", err)
	}
	t.Logf("\ngrafana_session: %s \ngrafana_session_expiry: %s\n",
		cookies.GrafanaSession,
		cookies.GrafanaSessionExpiry)
}
