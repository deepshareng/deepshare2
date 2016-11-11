package devicecookier

import (
	"net/http/httptest"
	"testing"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"golang.org/x/net/context"
)

//end-to-end test against DeviceCookier
func TestDeviceCookieHelper(t *testing.T) {

	dc := NewDeviceCookier(storage.NewInMemSimpleKV())
	ts := httptest.NewServer(newDeviceCookieHandler(dc, nil, api.DeviceCookiePrefix))
	defer ts.Close()

	ctx := context.TODO()
	specificCookieUrl := ts.URL + "/v2/devicecookie/"
	cli := httputil.GetNewClient()

	PutCookieInfo(ctx, cli, specificCookieUrl, HardwareIDPrefix+"hh", "cc")
	PutCookieInfo(ctx, cli, specificCookieUrl, UniqueIDPrefix+"uu", "cc")
	PutCookieInfo(ctx, cli, specificCookieUrl, UniqueIDPrefixWCookie+"uu", "ww")
	PutDeviceInfo(ctx, cli, specificCookieUrl, CookieIDPrefix+"cc", "uu")
	PutDeviceInfo(ctx, cli, specificCookieUrl, WCookieIDPrefix+"ww", "uu")

	{
		c, err := GetCookieID(ctx, cli, specificCookieUrl, HardwareIDPrefix+"hh")
		if err != nil {
			t.Error(err)
		}
		if c != "cc" {
			t.Error("get cookieID by uniqueID uu; want = cc, got =", c)
		}
	}

	{
		c, err := GetCookieID(ctx, cli, specificCookieUrl, UniqueIDPrefix+"uu")
		if err != nil {
			t.Error(err)
		}
		if c != "cc" {
			t.Error("get cookieID by uniqueID uu; want = cc, got =", c)
		}
	}

	{
		w, err := GetCookieID(ctx, cli, specificCookieUrl, UniqueIDPrefixWCookie+"uu")
		if err != nil {
			t.Error(err)
		}
		if w != "ww" {
			t.Error("get wCookieID by uniqueID uu; want = ww, got =", w)
		}
	}

	{
		u, err := GetUniqueID(ctx, cli, specificCookieUrl, CookieIDPrefix+"cc")
		if err != nil {
			t.Error(err)
		}
		if u != "uu" {
			t.Error("get uniqueid by cookieID cc; want = uu, got =", u)
		}
	}

	{
		u, err := GetUniqueID(ctx, cli, specificCookieUrl, WCookieIDPrefix+"ww")
		if err != nil {
			t.Error(err)
		}
		if u != "uu" {
			t.Error("get uniqueid by wCookieID ww; want = uu, got =", u)
		}
	}
}
