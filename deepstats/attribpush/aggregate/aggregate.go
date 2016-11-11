package aggregate

import (
	"github.com/MISingularity/deepshare2/deepstats/attribpush"
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AttributionAggregator interface {
	//Aggregate holds the core logic for attribution aggregation
	// accepts a event and aggregates the event to some place that holds the aggregated data
	Aggregate(e messaging.Event) error
}

//simpleAggregator is a simple attribution aggregator
// we maintain <receiver_id, last sender_id> pairs
// for a specific receiver R, when an open or install is introduced by a sender S1, we set <R, S1>
// when the app on R's device is closed, delete <R, S1> from the pairs
// when a new open or install is introduced by sender S2. we update the pair to <R, S2>
// when receiver R pushed a event with tag and value, simply add the (tag,value) to the last sender_id of R
type simpleAggregator struct {
	ar      attribpush.AttributionRetriver
	mgocoll *mgo.Collection
}

func (sa *simpleAggregator) Aggregate(e messaging.Event) error {
	senderID, tag := sa.ar.Retrive(e)
	if tag == "" {
		log.Debug("irrelevant event, dropped. event:", e)
		return nil
	}
	if senderID == "" {
		log.Debug("event is not introduced by any sender, dropped. event:", e)
		return nil
	}
	return sa.aggregate(e.AppID, senderID, tag, e.Count, e.TimeStamp)
}

/*
{
  "sender_id": "S",
  "tagged_values": {
    "tag1": [
      {
        "t": "time1",
        "v": 100
      }
    ],
    "tag2": [
      {
        "v": 1,
        "t": "time2"
      }
    ]
  }
}
aggregate (S, tag2, 50, time3)
->
{
  "sender_id": "S",
  "tagged_values": {
    "tag1": [
      {
        "t": "time1",
        "v": 100
      }
    ],
    "tag2": [
      {
        "v": 1,
        "t": "time2"
      },
      {
        "v": 50,
        "t": "time3"
      }
    ]
  }
}
*/
func (sa *simpleAggregator) aggregate(appID, senderID, tag string, value int, timestamp int64) error {
	_, err := sa.mgocoll.Upsert(bson.M{"app_id": appID, "sender_id": senderID},
		bson.M{"$addToSet": bson.M{"tagged_values." + tag: bson.M{"v": value, "t": timestamp}}})
	if err != nil {
		log.Errorf("upsert to mongo failed, err:%v, appID:%s, senderID:%s, tag:%s, value:%d, timestamp:%v\n",
			err, appID, senderID, tag, value, timestamp)
	}
	return err
}
