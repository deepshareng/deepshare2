package cookieutil

import (
	"net/http"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/token"
)

const (
	expireAfterYears = 2
)

func GetCookie(r *http.Request, cli *http.Client, specificTokenUrl string) (cookie *http.Cookie, isNew bool, err error) {
	cookie, err = r.Cookie(api.CookieName)
	if err == http.ErrNoCookie {
		cookieNew, err := newCookie(cli, specificTokenUrl)
		return cookieNew, true, err
	}
	return cookie, false, nil
}

func newCookie(cli *http.Client, specificTokenUrl string) (*http.Cookie, error) {
	cookieValue, err := token.GetNewToken(cli, specificTokenUrl, "cookie")
	expire := time.Now().AddDate(expireAfterYears, 0, 0)
	if err != nil {
		return nil, err
	}
	cookieNew := &http.Cookie{
		Name:    api.CookieName,
		Value:   cookieValue,
		Path:    "/",
		Expires: expire,
	}
	return cookieNew, nil
}
