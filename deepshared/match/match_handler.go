package match

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"bytes"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	in "github.com/MISingularity/deepshare2/pkg/instrumentation"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

type matchHandler struct {
	m        Matcher
	p        messaging.Producer
	endpoint string
}

// Used for unit testing handler core logic.
func AddHandler(mux *http.ServeMux, endpoint string, db storage.SimpleKV, p messaging.Producer, uaMatchValidSeconds int64) {
	mux.Handle(endpoint, newMatchHandler(NewSimpleMatcher(db, uaMatchValidSeconds), p, endpoint))
}

func newMatchHandler(m Matcher, mp messaging.Producer, endP string) *matchHandler {
	return &matchHandler{
		m:        m,
		p:        mp,
		endpoint: endP,
	}
}

func NewMatchTestHandler(endP string) *matchHandler {
	return &matchHandler{
		m:        NewSimpleMatcher(storage.NewInMemSimpleKV(), uaMatchExpireAfterSecDefault),
		p:        messaging.NewSimpleProducer(bytes.NewBuffer(nil)),
		endpoint: endP,
	}
}

func (mh *matchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !httputil.AllowMethod(w, r.Method, "GET", "PUT") {
		return
	}

	start := time.Now()

	switch r.Method {
	case "GET":
		defer in.PrometheusForMatch.HTTPGetDuration(start)

		fields := strings.Split(r.URL.Path[len(mh.endpoint):], "/")

		// TODO: we need to extract receiver info here so that we can allow
		// detailed attribution later. we also need to get cookieID from parameter
		if len(fields) < 1 || len(fields) > 2 {
			log.Error("Match; Invalid len of fields:", len(fields))
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}

		appID := fields[0]
		cookieID := ""
		// If client have forwarded cookieID, they access the /v1/matches/:appID/:cookieID
		// resource, this is when we need to exact match.
		if len(fields) > 1 {
			cookieID = fields[1]
		}

		//extract params from form
		ip := r.FormValue("client_ip")
		uaStr := r.FormValue("client_ua")

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

		//extract ua info out of uaStr
		ua := uainfo.ExtractUAInfoFromUAString(ip, uaStr)
		log.Debugf("ExtractUAInfoFromUAString uastr: %s uainfo: %#v\n", uaStr, ua)
		//for ios, uaStr is invalid, need to extract from receiverInfo param
		if ua.Os == "" || ua.Os == "ios" {
			ua = uainfo.ExtractUAInfoWithReceiverInfo(ip, uaStr, receiverInfo)
			log.Debugf("ExtractUAInfoFromReceiverInfo receiverInfo: %v uainfo: %#v\n", receiverInfo, ua)
		}

		if cookieID == "" && ua.Os == "" {
			log.Error(api.ErrMatchBadParameters.Message)
			httputil.WriteHTTPError(w, api.ErrMatchBadParameters)
			return
		}

		// This include both sender info and inapp_data
		mp, err := mh.m.Match(context.TODO(), appID, cookieID, uainfo.NewUAFingerPrinter(ua), receiverInfo.UniqueID)
		if err != nil {
			// TODO: define error message format in api.md
			// TODO: token error or appID error? or a server side error?
			log.Error(api.ErrMatchNotFound.Message, "err:", err)
			httputil.WriteHTTPError(w, api.ErrMatchNotFound)
			return
		}

		en := json.NewEncoder(w)
		resp := api.MatchResponse{
			InappData: mp.InappData,
			Channels:  mp.SenderInfo.Channels,
			SenderID:  mp.SenderInfo.SenderID,
		}
		if err := en.Encode(resp); err != nil {
			// TODO: use a logger pkg and change this to debug level
			log.Errorf("api: failed to encode match response to %s", r.RemoteAddr)
		}

		if err := produceMatchEvent(appID, mp, cookieID, tracking, receiverInfo, ua, mh.p); err != nil {
			//TODO should write an error log
			panic(err)
		}
	case "PUT":
		defer in.PrometheusForMatch.HTTPPostDuration(start)
		fields := strings.Split(r.URL.Path[len(mh.endpoint):], "/")
		if len(fields) != 2 {
			httputil.WriteHTTPError(w, api.ErrPathNotFound)
			return
		}

		appID := fields[0]
		cookieID := fields[1]

		decoder := json.NewDecoder(r.Body)
		matchReq := new(MatchPayload)
		if err := decoder.Decode(matchReq); err != nil {
			httputil.WriteHTTPError(w, api.ErrBadJSONBody)
			return
		}

		//TODO should we check if sender_id or channels is empty?

		if matchReq.ClientIP == "" {
			log.Error(api.ErrMatcPuthNeedIPAndUA.Message)
			httputil.WriteHTTPError(w, api.ErrMatcPuthNeedIPAndUA)
			return
		}
		ua := &uainfo.UAInfo{}
		if matchReq.ClientUA != "" {
			ua = uainfo.ExtractUAInfoFromUAString(matchReq.ClientIP, matchReq.ClientUA)
		}
		err := mh.m.Bind(context.TODO(), appID, cookieID, uainfo.NewUAFingerPrinter(ua), matchReq)
		if err != nil {
			log.Error(api.ErrInternalServer, "err:", err)
			httputil.WriteHTTPError(w, api.ErrInternalServer)
			return
		}

		if err := produceBindEvent(appID, matchReq, cookieID, ua, mh.p); err != nil {
			//TODO should write an error log for alerting
			panic(err)
		}
	}
}
