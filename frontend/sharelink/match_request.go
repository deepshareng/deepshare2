package sharelink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/MISingularity/deepshare2/deepshared/match"
	"github.com/MISingularity/deepshare2/pkg"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
)

func (sl *Sharelink) RequestMatch(client *http.Client, appID, inAppData, senderID, channels, sdkInfo, cookieID, deeplinkID, ip, ua string) {
	matchUrlStr, err := httputil.AppendPath(sl.specificMatchUrl, appID)
	if err != nil {
		log.Errorf("jsApiHandler; matchUrl %s is constructed by us, should not in wrong format: %v", sl.specificMatchUrl, err)
		panic(err)
	}
	matchurl, err := url.Parse(matchUrlStr)
	if err != nil {
		log.Errorf("[error]; requestMatch match url is illegal")
		return
	}
	pathMatch := matchurl.Path
	newPathMatch := path.Join(pathMatch, cookieID)
	if deeplinkID != "" {
		newPathMatch = path.Join(pathMatch, cookieID+"_"+deeplinkID)
	}
	matchurl.Path = newPathMatch
	requestMatchUrlStr := matchurl.String()
	log.Debugf("Share link; Request Match URL string is %s", requestMatchUrlStr)
	jsonBody := sl.setupMatchBody(inAppData, senderID, channels, sdkInfo, ip, ua)
	log.Debugf("Share link; MatchBody is %s", jsonBody)
	req, err := http.NewRequest("PUT", requestMatchUrlStr, strings.NewReader(jsonBody))
	if err != nil {
		log.Panic(fmt.Sprintln("[Error],Sharelink requestMatch setup new request failed:", err))
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(fmt.Sprintln("[Error],Sharelink requestMatch do request failed:", err))
	}
	defer resp.Body.Close()
}

func (sl *Sharelink) setupMatchBody(inAppData, senderID, channels, sdkInfo, ip, ua string) string {
	channelArray := pkg.DecodeStringSlice(channels)
	senderInfo := match.SenderInfoObj{
		SenderID: senderID,
		Channels: channelArray,
		SDKInfo:  sdkInfo,
	}
	matchRequestBody := match.MatchRequestBody{
		InAppData:  inAppData,
		SenderInfo: senderInfo,
		ClientIP:   ip,
		ClientUA:   ua,
	}
	if b, err := json.Marshal(matchRequestBody); err != nil {
		log.Panic(fmt.Sprintln("[Error],Sharelink setupMatchBody failed:", err))
		return ""
	} else {
		return string(b)
	}
}
