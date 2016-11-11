package integration

import (
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/MISingularity/deepshare2/api"
	"github.com/MISingularity/deepshare2/deepshared/counter"
	"github.com/MISingularity/deepshare2/deepshared/uainfo"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/testutil"
	"github.com/nsqio/go-nsq"
)

var nullLogger = log.New(ioutil.Discard, "", log.LstdFlags)

const (
	nsqdAddr = "127.0.0.1:4150"
)

func TestNSQSetup(t *testing.T) {
	p := mustNewNSQProducer()

	err := p.Publish("write_test", []byte("test"))
	if err != nil {
		t.Fatalf("should lazily connect - %s", err)
	}

	p.Stop()

	err = p.Publish("write_test", []byte("fail test"))
	if err != nsq.ErrStopped {
		t.Fatalf("should not be able to write after Stop()")
	}
}

func TestCounterPostWithNSQ(t *testing.T) {
	// Counter handler will call counter service and write counter events to
	// producer. Consumer decodes messages and returns it via eventChan.
	p := mustNewNSQProducer()
	eventChan := make(chan *messaging.Event, 10)
	hf := func(message *nsq.Message) error {
		event := new(messaging.Event)
		event.Unmarshal(message.Body)
		select {
		case eventChan <- event:
		case <-time.After(10 * time.Second):
			panic("")
		}
		return nil
	}
	c := mustNewNSQConsumer(string(messaging.CounterTopic), hf)
	defer p.Stop()
	defer c.Stop()

	handler := counter.NewCounterHandler(messaging.NewTestingNSQProducer(p), api.CounterPrefix)
	url := "http://" + path.Join("example.com", api.CounterPrefix, "AppID")

	tests := []struct {
		code   int
		header map[string]string
		body   string
		events []*messaging.Event
	}{
		{
			http.StatusOK,
			map[string]string{},
			`{"receiver_info":{"unique_id":"rec_dev"},` +
				`"counters":[{"event":"e1","count":1},{"event":"e2","count":10}]}` + "\n",
			[]*messaging.Event{
				{AppID: "AppID", UniqueID: "rec_dev", EventType: api.CounterPrefix + "e1", Count: 1, UAInfo: uainfo.UAInfo{}},
				{AppID: "AppID", UniqueID: "rec_dev", EventType: api.CounterPrefix + "e2", Count: 10, UAInfo: uainfo.UAInfo{}},
			},
		},
	}

	for i, tt := range tests {
		w := testutil.HandleWithRequestInfo(handler, "POST", url, tt.body, tt.header, "ip1:port1")
		if w.Code != tt.code {
			t.Errorf("#%d: HTTP status code = %d, want = %d", i, w.Code, tt.code)
		}
		for _, event := range tt.events {
			ce := <-eventChan
			ce.UAInfo = uainfo.UAInfo{}
			ce.TimeStamp = 0
			if !reflect.DeepEqual(ce, event) {
				t.Errorf("#%d, event=%v, want=%v", i, ce, event)
			}
		}
	}
}

func mustNewNSQProducer() *nsq.Producer {
	config := nsq.NewConfig()
	w, _ := nsq.NewProducer(nsqdAddr, config)
	w.SetLogger(nullLogger, nsq.LogLevelInfo)
	if err := w.Ping(); err != nil {
		panic(err)
	}
	return w
}

func mustNewNSQConsumer(topic string, hf nsq.HandlerFunc) *nsq.Consumer {
	config := nsq.NewConfig()
	config.DefaultRequeueDelay = 0
	config.MaxBackoffDuration = 50 * time.Millisecond
	c, _ := nsq.NewConsumer(topic, "ch", config)
	c.SetLogger(nullLogger, nsq.LogLevelInfo)
	c.AddHandler(hf)
	c.ConnectToNSQD(nsqdAddr)
	return c
}
