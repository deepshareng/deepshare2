package appcookiedevice

import (
	"net/http"

	"golang.org/x/net/context"

	"bytes"
	"encoding/json"

	"errors"
	"fmt"

	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
)

func RefreshCookie(ctx context.Context, cli *http.Client, specificAppCookieURL, cookieId, newCookieID string) error {
	if specificAppCookieURL == "" || cookieId == "" || newCookieID == "" {
		return nil
	}
	appCookieURL, err := httputil.AppendPath(specificAppCookieURL, RefreshCookiePath)
	if err != nil {
		log.Error("Request AppCookieDevice; RefreshCookie AppendPath failed, err:", err)
		return err
	}
	prcb := PostRefreshCookieBody{Cookie: cookieId, NewCookie: newCookieID}
	b, err := json.Marshal(prcb)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", appCookieURL, bytes.NewReader(b))
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("Request AppCookieDevice; failed, code = %d", resp.StatusCode)
		log.Error(errStr)
		return errors.New(errStr)
	}

	return nil
}

// TODO put all requests to appcookiedevice in request_helper
