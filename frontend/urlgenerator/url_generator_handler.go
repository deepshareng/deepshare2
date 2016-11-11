package urlgenerator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/shorturl"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

type generateUrlHandler struct {
	gu               *UrlGenerator
	endPoint         string
	mProducer        messaging.Producer
	client           *http.Client
	specificTokenUrl string
}

func newGenerateUrlHandler(cli *http.Client, s shorturl.UrlShortener, urlBase *url.URL, mp messaging.Producer, ep string, tokenUrl string) http.Handler {
	genUrlHandler := &generateUrlHandler{
		gu:               NewUrlGenerator(s, urlBase),
		endPoint:         ep,
		mProducer:        mp,
		client:           cli,
		specificTokenUrl: tokenUrl,
	}
	return genUrlHandler
}

// NewGenerateUrlTestHandler sets up simple urlGenerator backend (storage, etc.)
func NewGenerateUrlTestHandler(tokenUrl string) http.Handler {
	s := shorturl.NewUrlShortener(http.DefaultClient, storage.NewInMemSimpleKV(), tokenUrl)
	urlBase, _ := url.Parse("http://example.com")
	p := messaging.NewSimpleProducer(bytes.NewBuffer(nil))
	return newGenerateUrlHandler(http.DefaultClient, s, urlBase, p, api.GenerateUrlPrefix, tokenUrl)
}

func CreateHandler(endpoint, urlGeneratorBase string, db storage.SimpleKV, tokenUrl string, mp messaging.Producer) http.Handler {
	client := httputil.GetNewClient()
	s := shorturl.NewUrlShortener(client, db, tokenUrl)
	urlBase, err := url.Parse(urlGeneratorBase)
	if err != nil {
		log.Fatalf("[Error]. URL generator; url Base %s is passed by docker, should not in wrong format: %v", urlGeneratorBase, err)
	}
	return newGenerateUrlHandler(client, s, urlBase, mp, endpoint, tokenUrl)
}

func (guh *generateUrlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "POST", "GET") {
		return
	}

	start := time.Now()
	switch r.Method {
	case "POST":
		defer in.PrometheusForShorturl.HTTPPostDuration(start)
		decoder := json.NewDecoder(r.Body)
		//set the default value of isshort is true
		genBody := &GenURLPostBody{IsShort: true}
		if err := decoder.Decode(genBody); err != nil {
			log.Errorf("Url generator post body decode error: %v, r.ContentLength: %v, r.Body: %#v, r", err, r.ContentLength, r.Body)
			httputil.WriteHTTPError(w, api.ErrBadJSONBody)
			return
		}
		if v, err := StringfyInappData(genBody); err != nil {
			log.Error("Url generator; got invalid inappdata, err:", err)
			httputil.WriteHTTPError(w, api.ErrBadJSONBody)
			return
		} else {
			genBody = v
		}

		appID := path.Base(r.URL.Path)
		ua := uainfo.ExtractUAInfoFromUAString(httputil.ParseClientIP(r), r.UserAgent())
		u, err := guh.gu.GenerateUrl(appID, genBody, ua, guh.endPoint, guh.mProducer)

		//for ios9, wechat browser, append wcookie in the path (will be sent to server with POST InappData request when app opened via universal link)
		// http://fds.so/appid/shortseg?wcookie=www
		//TODO check universal link? need request appinfo which will cause higher latency
		if ua.IsWechat && ua.IosMajorVersion() >= 9 {
			u = guh.gu.AppendWCookie(w, r, guh.client, guh.specificTokenUrl, u)
		}

		urlResp := GenURLResponseBody{Url: u, Path: strings.TrimPrefix(u, guh.gu.urlGeneratorBase.String())}
		if err != nil {
			//The error means system clock is moving backwards, which means the system works abnormal
			//So we should kill this image
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			log.Error("Url generator", err)
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(urlResp); err != nil {
			// TODO: use a logger pkg and change this to debug level
			log.Errorf("Url generator: failed to encode url response to %s", r.RemoteAddr)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
		}
	case "GET":
		defer in.PrometheusForShorturl.HTTPGetDuration(start)
		fields := strings.Split(r.URL.Path[len(guh.endPoint):], "/")
		if len(fields) != 2 {
			log.Error("urlgenerator; Invalid len of fields:", len(fields))
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}

		appID := fields[0]
		shortSeg := fields[1]
		tracking := strings.ToLower(r.FormValue("tracking"))
		ri := r.FormValue("receiver_info")
		if tracking != "install" && tracking != "open" {
			//TODO return an error code?
			log.Error("invalid tracking:", tracking)
		}
		receiverInfo := api.MatchReceiverInfo{}
		if err := json.Unmarshal([]byte(ri), &receiverInfo); err != nil {
			//TODO return an error code?
			log.Error("invalid receiverInfo:", ri)
		}
		u := guh.gu.urlPrefix(appID)
		u.Path = path.Join(u.Path, shortSeg)
		rawUrl, err := guh.gu.urlShortener.ToRawURL(u, appID)
		if err != nil || rawUrl == nil {
			//TODO return an error code
			log.Error(err, u.String())
			httputil.WriteHTTPError(w, api.ErrMatchBadParameters)
			return
		}
		resp := api.MatchResponse{}
		values := rawUrl.Query()
		if contexts, ok := values["inapp_data"]; ok {
			resp.InappData = contexts[0]
		} else {
			log.Debug("failed to extract inappData from shortseg")
		}
		if contexts, ok := values["channels"]; ok {
			resp.Channels = pkg.DecodeStringSlice(contexts[0])
		} else {
			log.Debug("failed to extract channels from shortseg")
		}
		if contexts, ok := values["sender_id"]; ok {
			resp.SenderID = contexts[0]
		} else {
			log.Debug("failed to extract sender_id from shortseg")
		}
		en := json.NewEncoder(w)
		if err := en.Encode(resp); err != nil {
			log.Error("api: failed to encode GetUrlDataResponse response, shortseg:", shortSeg)
		}
		log.Debugf("[GET] urlgenerator, response: %#v\n", resp)

		ip := r.FormValue("client_ip")
		uaStr := r.FormValue("client_ua")
		ua := uainfo.ExtractUAInfoFromUAString(ip, uaStr)
		log.Debugf("ExtractUAInfoFromUAString uastr: %s uainfo: %#v\n", uaStr, ua)
		//for ios, uaStr is invalid, need to extract from receiverInfo param
		if ua.Os == "" {
			ua = uainfo.ExtractUAInfoWithReceiverInfo(ip, uaStr, receiverInfo)
			log.Debugf("ExtractUAInfoFromReceiverInfo receiverInfo: %v uainfo: %#v\n", receiverInfo, ua)
		}
		if err := produceGetInappDataEvent(appID, guh.endPoint, shortSeg, tracking, rawUrl, receiverInfo, ua, guh.mProducer); err != nil {
			panic(err)
		}
	}
}
