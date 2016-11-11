package match

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"

	"fmt"
	"strconv"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

// UA match should only be available within 15min
const uaMatchExpireAfterSecDefault = int64(15 * 60)

const (
	keyData             = "data"
	keyCreateAt         = "create_at"
	keyAccessFlagPrefix = "r_"
	flagAccessed        = "1"
)

// simpleTokenMatcher matches user info to token based on exact same user info struct.
type simpleMatcher struct {
	skv                   storage.SimpleKV
	uaMatchExpireAfterSec int64
}

func NewSimpleMatcher(skv storage.SimpleKV, uaMatchValidSeconds int64) Matcher {
	return &simpleMatcher{
		skv: skv,
		uaMatchExpireAfterSec: uaMatchValidSeconds,
	}
}

func (tm *simpleMatcher) Bind(ctx context.Context, appID, cookieID string, u uainfo.UATransformer, mp *MatchPayload) error {
	log.Debugf("simpleMatcher.Bind, appid=%s, cookieID=%s, uakey=%v, payload=%v\n", appID, cookieID, u.Transform(), mp)

	if err := tm.bindByCookie(appID, cookieID, mp); err != nil {
		log.Error("bindByCookie failed, err:", err)
		return err
	}

	//if u.Transform() got empty, ua is not set
	//		for the situation for ios9 device with universal link
	//		UA is useless since in-app data is retrieved by cookieID for install and shortSeg for open
	if u.Transform() == "" {
		log.Debug("Bind, ua is empty, not to bind with UA")
		return nil
	}
	if err := tm.bindByUA(appID, u, mp); err != nil {
		log.Error("bindByUA failed, err:", err)
		return err
	}

	return nil
}

func (tm *simpleMatcher) Match(ctx context.Context, appID, cookieID string, u uainfo.UATransformer, receiverID string) (*MatchPayload, error) {
	log.Debugf("simpleMatcher.Match, appid=%s, cookieID=%s, u=%v, receiverID=%v\n", appID, cookieID, u.Transform(), receiverID)
	var result []byte
	var err error
	switch {
	//exact match with cookieID
	case len(cookieID) != 0:
		result, err = tm.matchWithCookie(ctx, appID, cookieID, u, receiverID)
	//match with UA info, should have 15min window and device repeatedly access limit
	default:
		result, err = tm.matchWithUA(ctx, appID, u, receiverID)
	}

	if err != nil {
		log.Error("Match Error:", err)
		return nil, err
	}

	mp := new(MatchPayload)
	if err := json.Unmarshal(result, mp); err != nil {
		return nil, err
	}
	log.Debugf("simpleMatcher.Match succeed, appid=%s, cookieID=%s, u=%v, receiverID=%v, got=%v\n",
		appID, cookieID, u.Transform(), receiverID, mp)
	return mp, nil
}

// get in-app data based on cookieID
// side effect: tag the corresponding UA binding as accessed
// here is why:
//	a click will cause two bindings, one with cookieID, the other with UA,
//	since when the binding happen, we can not predict which of the two will be used to match.
//	if a match is called based on cookieID, the two bindings should be both consumed
func (tm *simpleMatcher) matchWithCookie(ctx context.Context, appID string, cookieID string, u uainfo.UATransformer, receiverID string) ([]byte, error) {

	ck, err := formCookieKey(appID, cookieID)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	v, err := tm.skv.Get([]byte(ck))
	if err != nil {
		return nil, err
	}
	in.PrometheusForMatch.StorageGetDuration(start)

	if v == nil {
		return nil, NoMatchForCookieErr
	}

	//Don't delete in-app data binded under cookieID for the following two apps.
	//for TouTiao (fbbca54fd667867c):
	//	1. they use the new Android SDK which can avoid accessing with cookieID when app launched by tapping app icon.
	//	2. they use JS SDK, when the "open with app" button on TouTiao WAP page was clicked repeatedly,
	//	   only the first click has cookieID, but it's correct to return in-app data for every click.
	// for Linux Command (38CCA4C77072DDC9):
	//	We keep it the same as TouTiao, to reproduce any possible issues that maybe happen on TouTiao.
	if (appID != "fbbca54fd667867c" && appID != "38CCA4C77072DDC9") || u.Os() != "android" {
		startDelete := time.Now()
		err = tm.skv.Delete(ck)
		in.PrometheusForMatch.StorageDeleteDuration(startDelete)
	}

	if err := tm.markReceiverAccessed(appID, u, receiverID); err != nil {
		log.Error("simpleMatcher.matchWithCookie error when markReceiverAccessed:", err)
	}
	log.Debugf("simpleMatcher, mark ua key accessed. appID: %s, cookieID: %s, u: %s\n", appID, cookieID, u.Transform())

	log.Debugf("simpleMatcher.matchWithCookie: k: %s, v: %s, u: %s, receiverID: %s, err: %v\n",
		string(ck), string(v), u.Transform(), receiverID, err)
	return v, err
}

func (tm *simpleMatcher) matchWithUA(ctx context.Context, appID string, u uainfo.UATransformer, receiverID string) ([]byte, error) {

	if !tm.existsBinding(appID, u) {
		return nil, NoMatchForUAErr
	}
	if tm.isExpired(appID, u) {
		key, _ := formUserInfoKey(appID, u)
		//Delete expired binding
		start := time.Now()
		tm.skv.Delete(key)
		in.PrometheusForMatch.StorageDeleteDuration(start)
		return nil, MatchUABindExpireErr
	}
	if tm.hasAccessed(appID, u, receiverID) {
		return nil, MatchUABindAccessedErr
	}

	log.Debugf("simpleMatcher.matchWithUA, first time Access. appid: %s u: %s receiverID: %s\n", appID, u.Transform(), receiverID)
	if err := tm.markReceiverAccessed(appID, u, receiverID); err != nil {
		log.Error("simpleMatcher.matchWithUA error when markReceiverAccessed:", err)
		return nil, err
	}

	key, _ := formUserInfoKey(appID, u)
	start := time.Now()
	b, err := tm.skv.HGet(key, keyData)
	if err != nil {
		return nil, err
	}
	in.PrometheusForMatch.StorageHGetDuration(start)
	if b == nil {
		return nil, NoMatchForUAErr
	}

	log.Debugf("simpleMatcher.matchWithUA: k: %s, v: %s, \n", string(key), string(b))
	return b, nil
}

func (tm *simpleMatcher) bindByCookie(appID string, cookieID string, mp *MatchPayload) error {
	if cookieID == api.CookiePlaceHolder {
		log.Debug("bindByCookie skipped, cookieID =", cookieID)
		return nil
	}
	payload, err := json.Marshal(mp)
	if err != nil {
		return err
	}
	ck, err := formCookieKey(appID, cookieID)
	if err != nil {
		return err
	}

	start := time.Now()
	if err := tm.skv.Set(ck, payload); err != nil {
		return err
	}
	in.PrometheusForMatch.StorageSaveDuration(start)
	log.Debugf("Bind with cookieID succeed. k = %s, len(v) = %v\n", string(ck), len(payload))

	return nil
}

func (tm *simpleMatcher) bindByUA(appID string, u uainfo.UATransformer, mp *MatchPayload) error {
	payload, err := json.Marshal(mp)
	if err != nil {
		return err
	}
	uik, err := formUserInfoKey(appID, u)
	if err != nil {
		return err
	}
	startDelete := time.Now()
	err = tm.skv.Delete(uik)
	in.PrometheusForMatch.StorageDeleteDuration(startDelete)
	start := time.Now()
	if err := tm.skv.HSet(uik, keyData, payload); err != nil {
		return err
	}
	in.PrometheusForMatch.StorageHSetDuration(start)
	start = time.Now()
	if err := tm.skv.HSet(uik, keyCreateAt, []byte(fmt.Sprintf("%d", time.Now().Unix()))); err != nil {
		return err
	}
	in.PrometheusForMatch.StorageHSetDuration(start)
	log.Debugf("Bind with uainfo succeed. k = %v, len(v) = %v\n", string(uik), len(payload))

	return nil
}

func (tm *simpleMatcher) existsBinding(appID string, u uainfo.UATransformer) bool {
	uik, _ := formUserInfoKey(appID, u)
	return tm.skv.Exists(uik)
}

func (tm *simpleMatcher) isExpired(appID string, u uainfo.UATransformer) bool {
	uik, _ := formUserInfoKey(appID, u)
	b, err := tm.skv.HGet(uik, keyCreateAt)
	if err != nil {
		panic(err)
	}
	n, err := strconv.Atoi(string(b))
	if err != nil {
		log.Error("isExpired Atoi err:", err, n)
	}
	if time.Now().Unix()-int64(n) >= tm.uaMatchExpireAfterSec {
		return true
	}
	return false
}

func (tm *simpleMatcher) hasAccessed(appID string, u uainfo.UATransformer, receiverID string) bool {
	uik, _ := formUserInfoKey(appID, u)
	start := time.Now()
	b, err := tm.skv.HGet(uik, keyAccessFlagPrefix+receiverID)
	if err != nil {
		panic(err)
	}
	in.PrometheusForMatch.StorageHGetDuration(start)
	if string(b) == flagAccessed {
		return true
	}
	return false
}

func (tm *simpleMatcher) markReceiverAccessed(appID string, u uainfo.UATransformer, receiverID string) error {
	if receiverID == "" {
		log.Error("markReceiverAccessed failed, receiverID should not be emoty")
		return nil
	}
	key, _ := formUserInfoKey(appID, u)
	start := time.Now()
	if err := tm.skv.HSet(key, keyAccessFlagPrefix+receiverID, []byte(flagAccessed)); err != nil {
		return err
	}
	in.PrometheusForMatch.StorageHSetDuration(start)
	return nil
}

// Assuming appID is url encoded, this implementaiton make sure
// combined appID and UserInfo (unique per app) still unique.
func formUserInfoKey(appID string, u uainfo.UATransformer) ([]byte, error) {
	return []byte(appID + ":" + u.Transform()), nil
}

func formCookieKey(appID, cookieID string) ([]byte, error) {
	return []byte(appID + ":" + cookieID), nil
}
