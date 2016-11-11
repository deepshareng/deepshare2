package attribution

import (
	"strings"

	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
)

const (
	EventOther = iota
	EventInstall
	EventOpen
	EventClose
	EventCounter
)

type AttributionService interface {
	OnEvent(e *messaging.Event) error
}

type attributionService struct {
	ar       AttributionRetriever
	producer messaging.Producer
}

func NewAttributionService(ar AttributionRetriever, producer messaging.Producer) AttributionService {
	return &attributionService{
		ar:       ar,
		producer: producer,
	}
}

func (as *attributionService) OnEvent(e *messaging.Event) error {
	if e == nil {
		return nil
	}
	switch parseEventType(e.EventType) {
	// install/open/close events are used to retrieve attribution given appID and uniqueID
	case EventInstall:
		as.ar.OnInstall(e.AppID, e.UniqueID, e.SenderID, e.Channels, e.TimeStamp)
	case EventOpen:
		as.ar.OnOpen(e.AppID, e.UniqueID, e.SenderID, e.Channels, e.TimeStamp)
	case EventClose:
		as.ar.OnClose(e.AppID, e.UniqueID, e.TimeStamp)

	// counter event: fill sender_id and channels then push back to NSQ
	// so that aggregate service will calculate stats data by channels
	case EventCounter:
		attr := as.ar.GetAttribution(e.AppID, e.UniqueID)
		log.Debug("AttributionService; counter event. parsed attribution info:", attr)
		if attr != nil && attr.SenderID != "" {
			e.SenderID = attr.SenderID
		}
		if attr != nil && attr.Channels != nil {
			e.Channels = attr.Channels
		}
		as.producer.Produce(nsqProduceTopic, e)
	default:
		log.Debug("AttributionService; irrelevant event: ", e.EventType)
		return nil
	}

	return nil
}

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
