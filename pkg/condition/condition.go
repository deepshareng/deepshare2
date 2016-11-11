package condition

import (
	"reflect"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

const (
	DisplayPrefix = "match/"
)

func ConvertTimeToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func IsEmptyChannel(channels []string) bool {
	if len(channels) == 0 || len(channels) == 1 && channels[0] == "" {
		return true
	}
	return false
}

// only count the sharelink events with "shorturl_token"
func IsSharelink(dp *messaging.Event) bool {
	if !strings.Contains(dp.EventType, "sharelink:/d/") {
		return false
	}
	v, ok := dp.KVs["shorturl_token"]
	if !ok || v == "" {
		return false
	}

	b, ok := dp.KVs["shorturl_token_valid"]
	if !ok || b != true {
		return false
	}
	return true
}

func HasDstag(dp *messaging.Event) bool {
	v, ok := dp.KVs["ds_tag"]
	if !ok || v == "" {
		return false
	}
	return true
}

func ByDeepshare(dp *messaging.Event) bool {
	if dp.SenderID != "" || !IsEmptyChannel(dp.Channels) || HasInappdata(dp) {
		return true
	}
	return false
}

func SubstitutePrefix(event string) string {
	event = strings.Replace(event, api.GetInAppDataPrefix, DisplayPrefix, 1)
	return event
}

func HasInappdata(dp *messaging.Event) bool {
	s, ok := dp.KVs["inapp_data"]
	if ok {
		m := reflect.ValueOf(s).Interface().(string)
		if s != nil && m != "{}" && m != "" {
			return true
		}
	}
	return false
}
func IsOpenEvent(dp *messaging.Event) bool {
	return strings.HasPrefix(dp.EventType, api.GetInAppDataPrefix) && strings.HasSuffix(dp.EventType, "open")
}

func IsInstallEvent(dp *messaging.Event) bool {
	return strings.HasPrefix(dp.EventType, api.GetInAppDataPrefix) && strings.HasSuffix(dp.EventType, "install")
}
