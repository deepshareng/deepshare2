package token

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"regexp"

	"github.com/MISingularity/deepshare2/api"
)

func TestGetNewToken(t *testing.T) {
	tokenServer := httptest.NewServer(NewTokenTestHandler(api.TokenPrefix))
	specificTokenURL := tokenServer.URL + api.TokenPrefix
	cli := http.DefaultClient
	token, err := GetNewToken(cli, specificTokenURL, "cookie")
	if err != nil {
		t.Fatal("Failed to GetNewToken, err:", err)
	}
	wPattern := `\w{10}`
	reg := regexp.MustCompile(wPattern)
	if !reg.Match([]byte(token)) {
		t.Fatal("Token is in the wrong format, token =", token)
	}
}
