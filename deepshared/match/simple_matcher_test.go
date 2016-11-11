package match

import (
	"reflect"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/storage"
	"golang.org/x/net/context"
)

type simpleUATransformer struct {
	u string
}

func (s *simpleUATransformer) Transform() string {
	return s.u
}
func (s *simpleUATransformer) Os() string {
	return ""
}

// TestSimpleTokenMatcher tests Bind() and Match() pair functions.
func TestSimpleMatcher(t *testing.T) {
	tests := []struct {
		appID       string
		cookieID    string
		u           *simpleUATransformer
		receiverID  string
		mp          *MatchPayload
		bindFirst   bool
		matchWithUA bool

		wmp  *MatchPayload
		werr error
	}{
		//good: match with cookie
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m1",
			u:          &simpleUATransformer{"u1"},
			receiverID: "r1",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id1", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk1"},
			bindFirst:   true,
			matchWithUA: false,

			wmp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id1", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk1"},
			werr: nil,
		},
		//bad: repeatedly match with UA after match with cookie
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m1",
			u:          &simpleUATransformer{"u1"},
			receiverID: "r1",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id1", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk1"},
			bindFirst:   false,
			matchWithUA: true,

			wmp:  nil,
			werr: MatchUABindAccessedErr,
		},
		//good: with another bind, match with UA should be OK
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m1",
			u:          &simpleUATransformer{"u1"},
			receiverID: "r1",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id1", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk1"},
			bindFirst:   true,
			matchWithUA: true,
			wmp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id1", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk1"},
			werr: nil,
		},
		//good: with another bind, match with cookie should be OK
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m1",
			u:          &simpleUATransformer{"u1"},
			receiverID: "r1",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id1", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk1"},
			bindFirst:   true,
			matchWithUA: false,
			wmp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id1", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk1"},
			werr: nil,
		},
		//good: match with UA
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m2",
			u:          &simpleUATransformer{"u2"},
			receiverID: "r2",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id2", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk2"},
			bindFirst:   true,
			matchWithUA: true,
			wmp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id2", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk2"},
			werr: nil,
		},
		//bad: repeatedly match with UA after match with UA
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m2",
			u:          &simpleUATransformer{"u2"},
			receiverID: "r2",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id2", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk2"},
			bindFirst:   false,
			matchWithUA: true,
			wmp:         nil,
			werr:        MatchUABindAccessedErr,
		},
		//bad: fails to get data with cookie without bind
		{
			appID:       "TestSimpleTokenMatcher-AppID",
			cookieID:    "m3",
			u:           &simpleUATransformer{"u3"},
			receiverID:  "r3",
			mp:          nil,
			bindFirst:   false,
			matchWithUA: false,
			wmp:         nil,
			werr:        NoMatchForCookieErr,
		},
		//bad: fails to get data with UA without bind
		{
			appID:       "TestSimpleTokenMatcher-AppID",
			cookieID:    "m4",
			u:           &simpleUATransformer{"u4"},
			receiverID:  "r4",
			mp:          nil,
			bindFirst:   false,
			matchWithUA: true,
			wmp:         nil,
			werr:        NoMatchForUAErr,
		},
		//good: match with cookie
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m5",
			u:          &simpleUATransformer{"u5"},
			receiverID: "r5",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id5", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk5"},
			bindFirst:   true,
			matchWithUA: false,

			wmp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id5", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk5"},
			werr: nil,
		},
		//bad: match with cookie repeatedly
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m5",
			u:          &simpleUATransformer{"u5"},
			receiverID: "r5",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id5", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk5"},
			bindFirst:   false,
			matchWithUA: false,

			wmp:  nil,
			werr: NoMatchForCookieErr,
		},
	}
	tm := NewSimpleMatcher(storage.NewInMemSimpleKV(), uaMatchExpireAfterSecDefault)
	//	tm := NewSimpleMatcher(storage.NewRedisSimpleKV("127.0.0.1:6379", ""))
	for i, tt := range tests {
		if tt.bindFirst && tt.mp != nil {
			err := tm.Bind(context.TODO(), tt.appID, tt.cookieID, tt.u, tt.mp)
			if err != nil {
				t.Fatalf("#%d: TokenMatcher.Bind failed: %v", i, err)
			}
		}

		switch tt.matchWithUA {
		case false:
			// test exact match here.
			payload, err := tm.Match(context.TODO(), tt.appID, tt.cookieID, tt.u, tt.receiverID)
			if err != tt.werr {
				t.Errorf("#%d: TokenMatcher.Match failed, err = %v, want = %v\n", i, err, tt.werr)
			}
			if !reflect.DeepEqual(payload, tt.wmp) {
				t.Errorf("#%d: payload = %v, want = %v", i, payload, tt.wmp)
			}
		case true:
			// test ua match here.
			payload, err := tm.Match(context.TODO(), tt.appID, "", tt.u, tt.receiverID)

			if err != tt.werr {
				t.Errorf("#%d: TokenMatcher.Match failed, err = %v, want = %v\n", i, err, tt.werr)
			}

			if !reflect.DeepEqual(payload, tt.wmp) {
				t.Errorf("#%d: payload = %v, want = %v", i, payload, tt.wmp)
			}
		}
	}
}

func TestUAMatchExpire(t *testing.T) {
	uaMatchExpireAfterSec := int64(1)
	tests := []struct {
		appID      string
		cookieID   string
		u          *simpleUATransformer
		receiverID string
		mp         *MatchPayload
	}{
		//good: match with cookie
		{
			appID:      "TestSimpleTokenMatcher-AppID",
			cookieID:   "m1",
			u:          &simpleUATransformer{"u1"},
			receiverID: "r1",
			mp: &MatchPayload{
				SenderInfo: api.SenderInfo{SenderID: "id1", Channels: []string{"ch_x", "ch_y"}},
				InappData:  "tk1"},
		},
	}
	tm := NewSimpleMatcher(storage.NewInMemSimpleKV(), uaMatchExpireAfterSec)
	for i, tt := range tests {
		err := tm.Bind(context.TODO(), tt.appID, tt.cookieID, tt.u, tt.mp)
		if err != nil {
			t.Fatalf("#%d: TokenMatcher.Bind failed: %v", i, err)
		}
		time.Sleep(time.Duration(uaMatchExpireAfterSec) * time.Second)
		// test ua match here.
		payload, err := tm.Match(context.TODO(), tt.appID, "", tt.u, tt.receiverID)
		if err != MatchUABindExpireErr {
			t.Errorf("#%d: TokenMatcher.Match failed, err = %v, want = %v\n", i, err, MatchUABindExpireErr)
		}

		if payload != nil {
			t.Errorf("#%d: payload = %v, want = %v", i, payload, nil)
		}
	}
}
