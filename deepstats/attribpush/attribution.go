package attribpush

import (
	"strings"

	"github.com/MISingularity/deepshare2/deepstats/attribution/sender"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"gopkg.in/mgo.v2"
)

type AttributionRetriver interface {
	//AttributionRetriver holds the core logic for attribution retrieving
	// accepts a event and retrive senderID and tag
	Retrive(e messaging.Event) (senderID string, tag string)
}

//simpleAttributionRetriver is a simple attribution retriver
// we maintain <receiver_id, last sender_id> pairs
// for a specific receiver R, when an open or install is introduced by a sender S1, we set <R, S1>
// when the app on R's device is closed, delete <R, S1> from the pairs
// when a new open or install is introduced by sender S2. we update the pair to <R, S2>
// when receiver R pushed a event with tag and value, simply add the (tag,value) to the last sender_id of R
type simpleAttributionRetriver struct {
	rsm *sender.ReceiverSenderMapping
}

func NewAttributionParser(mgoColl *mgo.Collection) AttributionRetriver {
	return &simpleAttributionRetriver{rsm: sender.NewReceiverSenderMapping(mgoColl)}
}

func (sa *simpleAttributionRetriver) Retrive(e messaging.Event) (senderID string, tag string) {
	switch parseEventType(e.EventType) {
	case EventInstall:
		tag = "ds/install"
		sa.rsm.OnInstall(e.UniqueID, e.SenderID)
	case EventOpen:
		tag = "ds/open"
		sa.rsm.OnOpen(e.UniqueID, e.SenderID)
	case EventClose:
		sa.rsm.OnClose(e.UniqueID)
		return
	case EventCounter:
		tag = e.EventType[strings.Index(e.EventType, "/counters/")+len("/counters/"):]
	default:
		log.Debug("irrelevant event: ", e.EventType)
		return
	}

	senderID = e.SenderID
	if senderID == "" {
		senderID = sa.rsm.GetSenderID(e.UniqueID)
	}

	return senderID, tag
}

const (
	EventOther = iota
	EventInstall
	EventOpen
	EventClose
	EventCounter
)

func parseEventType(s string) int {
	if strings.HasSuffix(s, "/install") {
		return EventInstall
	}
	if strings.HasSuffix(s, "/open") {
		return EventOpen
	}
	if strings.HasSuffix(s, "app/close") {
		return EventClose
	}
	if strings.Contains(s, "counters/") {
		return EventCounter
	}
	return EventOther
}
