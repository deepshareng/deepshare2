package match

import (
	"errors"

	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"golang.org/x/net/context"
)

var (
	// TODO: Define standard errors for no match, timeout, etc.
	// 			We should have standard errors for package purpose.
	NoMatchForUAErr        = errors.New("match: No such match to given user information")
	NoMatchForCookieErr    = errors.New("match: No such match to given cookieID")
	MatchUABindExpireErr   = errors.New("match: The data binded by UA has expired")
	MatchUABindAccessedErr = errors.New("match: The data binded by UA has been accessed by the device")
)

// Matcher helps bind in-app data to given payload and
// hides matching logic in its API.
type Matcher interface {
	// Bind binds in-app data to userinfo and cookieID under appID.
	Bind(ctx context.Context, appID string, cookieID string, u uainfo.UATransformer, mp *MatchPayload) error

	// Match tries to find a payload that matches given AppID and cookieID or user info.
	// If the cookieID is presented, we have an exact match. Otherwise, we do a probalistic
	// match based on user info.
	// If no such match could be found, it returns nil and NoMatchErr.
	Match(ctx context.Context, appID string, cookieID string, u uainfo.UATransformer, receiverID string) (mp *MatchPayload, err error)
}
