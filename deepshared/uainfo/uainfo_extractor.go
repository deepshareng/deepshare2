package uainfo

import (
	"strconv"
	"strings"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/path"
	"github.com/MISingularity/uap-go/uaparser"
)

var regexParser *uaparser.Parser
var specialCases map[string]bool

func Init() error {
	var err error
	curdir, err := path.Getcurdir()
	if err != nil {
		return err
	}
	regexParser, err = uaparser.New(curdir + "/regexes.yaml")
	if err != nil {
		return err
	}
	return nil
}

// Get the client ip of the request, parse the user-agent to form UAInfo
//TODO add metrics
func ExtractUAInfoFromUAString(ip, uaStr string) *UAInfo {
	if regexParser == nil {
		err := Init()
		if err != nil {
			log.Fatalf("UA Parser load regex failed! Err Msg=%v", err)
			return nil
		}
	}
	client := regexParser.Parse(uaStr)
	family := strings.ToLower(client.Os.Family)
	brand := "-"
	if family == "other" {
		family = ""
	}
	if family != "ios" {
		brand = strings.ToLower(client.Device.ToString())
		brand = strings.Replace(brand, " ", "#", -1)
		brand = strings.Replace(brand, "-", "#", -1)
		brand = strings.Replace(brand, "_", "#", -1)
	}
	if brand == "other" {
		brand = "-"
	}

	chromeMajor := 0
	if strings.Contains(client.UserAgent.Family, "Chrome") {
		chromeMajor, _ = strconv.Atoi(client.UserAgent.Major)

	}
	isWechat := client.Channel.Params["wechat"]
	isWeibo := client.Channel.Params["weibo"]
	isQQ := client.Channel.Params["qq"]
	isQQBrowser := client.Channel.Params["qq browser"]
	isTwitter := client.Channel.Params["twitter"]
	isFaceBook := strings.Contains(client.UserAgent.ToString(), "Facebook")
	isFirefox := client.Channel.Params["firefox"]
	cannotDeeplink := client.Channel.Params["cannotdeeplink"]
	cannotGoMarket := client.Channel.Params["cannot go market"]
	forceUserScheme := client.Channel.Params["sogou search"]
	cannotGetWindowsEvent := client.Channel.Params["windows cannot catch"]
	return &UAInfo{
		Ua:                    uaStr,
		Ip:                    ip,
		Os:                    family,
		OsVersion:             strings.ToLower(client.Os.ToVersionString()),
		Brand:                 brand,
		Browser:               client.UserAgent.ToString(),
		IsWechat:              isWechat,
		IsWeibo:               isWeibo,
		IsQQ:                  isQQ,
		IsTwitter:             isTwitter,
		IsFacebook:            isFaceBook,
		IsQQBrowser:           isQQBrowser,
		IsFirefox:             isFirefox,
		ChromeMajor:           chromeMajor,
		CannotDeeplink:        cannotDeeplink,
		CannotGoMarket:        cannotGoMarket,
		CannotGetWindowsEvent: cannotGetWindowsEvent,
		ForceUseScheme:        forceUserScheme,
	}
}

func ExtractUAInfoWithReceiverInfo(ip, uaStr string, receiverInfo api.MatchReceiverInfo) *UAInfo {
	u := ExtractUAInfoFromUAString(ip, uaStr)
	//use os and os_version in receiverInfo (parsed from SDK, more reliable than UAString)
	u.Os = strings.ToLower(receiverInfo.OS)
	u.OsVersion = strings.Replace(receiverInfo.OSVersion, "_", ".", -1)
	return u
}
