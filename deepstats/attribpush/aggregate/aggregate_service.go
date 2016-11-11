package aggregate

import (
	"github.com/nsqio/go-nsq"
	"gopkg.in/mgo.v2"
)

type AggregateService struct {
	aggregator  AttributionAggregator
	mgoSession  *mgo.Session
	nsqConsumer *nsq.Consumer
}

func newAggregateService(mmongoAddr string, nsqlookupAddr string) *AggregateService {
	return nil
}

//Start the aggregation service loop
func (as *AggregateService) Start() {

}
