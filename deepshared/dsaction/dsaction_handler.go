package dsaction

import (
	"encoding/json"
	"net/http"

	"io/ioutil"
	"path"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/httputil"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

type dsactionHandler struct {
	p        messaging.Producer
	endpoint string
}

// Used for unit testing handler core logic.
func CreateHandler(endpoint string, mp messaging.Producer) http.Handler {
	return NewDSActionHandler(mp, endpoint)
}

func NewDSActionHandler(p messaging.Producer, endP string) http.Handler {
	return &dsactionHandler{
		p,
		endP,
	}
}

func (ch *dsactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("dsactionHandler is called, request:", r.Method, r.URL.Path)
	if !httputil.AllowMethod(w, r.Method, "POST") {
		return
	}
	var body []byte
	var err error
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		log.Errorf("DSActionHandler; Read Request body error: %s", err)
		httputil.WriteHTTPError(w, api.ErrBadRequestBody)
	}
	log.Debugf("DSActionHandler; Post Body : %s", string(body))
	//Extract app ID from: /dsactions/:appID
	appID := path.Base(r.URL.Path)
	if appID == "/" {
		log.Errorf("DSActionHandler; Illegal Request URL: %s", r.URL.String())
		httputil.WriteHTTPError(w, api.ErrPathNotFound)
	}

	dsAction := DSActionRequest{}
	if err := json.Unmarshal(body, &dsAction); err != nil {
		log.Errorf("DSActionHandler; Get DSAction post body decode error: body = %s; err = %v", body, err)
		httputil.WriteHTTPError(w, api.ErrBadJSONBody)
		return
	}
	ua := uainfo.ExtractUAInfoFromUAString(httputil.ParseClientIP(r), r.UserAgent())
	if err := produceDSActionEvent(appID, &dsAction, ua, ch.p, ch.endpoint); err != nil {
		//TODO should write an error log for alerting
		panic(err)
	}

}
