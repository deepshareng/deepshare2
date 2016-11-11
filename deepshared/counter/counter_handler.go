package counter

import (
	"encoding/json"
	"net/http"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

type counterHandler struct {
	p        messaging.Producer
	endpoint string
}

func AddHandler(mux *http.ServeMux, endpoint string, mp messaging.Producer) {
	mux.Handle(endpoint, NewCounterHandler(mp, endpoint))
}

func NewCounterHandler(p messaging.Producer, endP string) http.Handler {
	return &counterHandler{
		p,
		endP,
	}
}

func (ch *counterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("counterHandler is called, request:", r.Method, r.URL.Path)
	if !httputil.AllowMethod(w, r.Method, "POST") {
		return
	}
	// Endpoint:
	//  /counters/:appID
	appID := r.URL.Path[len(ch.endpoint):]
	cp := new(CounterRequest)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(cp); err != nil {
		log.Error(api.ErrBadJSONBody.Message, "err:", err)
		httputil.WriteHTTPError(w, api.ErrBadJSONBody)
		return
	}

	ua := uainfo.ExtractUAInfoFromUAString(httputil.ParseClientIP(r), r.UserAgent())
	if err := produceCounterEvents(appID, cp, ua, ch.p, ch.endpoint); err != nil {
		panic(err)
	}
}
