package aggregate

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type CounterEvent struct {
	AppID     string    `bson:"appid,omitempty"`
	Channel   string    `bson:"channel,omitempty"`
	Timestamp time.Time `bson:"timestamp,omitempty"`
	Event     string    `bson:"event,omitempty"`
	Count     int       `bson:"count,omitempty"`
	Os        string    `bson:"os,omitempty"`
}

type AggregateResult struct {
	Event  string
	Counts []*AggregateCount
}

type AggregateCount struct {
	Count int
}

type AggregateResultWithTimestamp struct {
	Event  string
	Counts []*AggregateCountWithTimestamp
}

type AggregateCountWithTimestamp struct {
	Timestamp time.Time
	Count     int
}

func AggregateResultsToString(aggrs []*AggregateResult) string {
	buf := new(bytes.Buffer)
	buf.WriteString("\n")
	for _, aggr := range aggrs {
		buf.WriteString(fmt.Sprintf("=> Event: %s\n", aggr.Event))
		for _, count := range aggr.Counts {
			buf.WriteString(fmt.Sprintf("=> => Count: %d\n", count.Count))
		}
	}
	return string(buf.Bytes())
}

func AggregateResultsWithTImeStampToString(aggrs []*AggregateResultWithTimestamp) string {
	buf := new(bytes.Buffer)
	buf.WriteString("\n")
	for _, aggr := range aggrs {
		buf.WriteString(fmt.Sprintf("=> Event: %s\n", aggr.Event))
		for _, count := range aggr.Counts {
			buf.WriteString(fmt.Sprintf("=> => Timestamp: %s\n", count.Timestamp))
			buf.WriteString(fmt.Sprintf("=> => Count: %d\n", count.Count))
		}
	}
	return string(buf.Bytes())
}

type aggrResults []*AggregateResult

func (a aggrResults) Len() int {
	return len(a)
}

func (a aggrResults) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a aggrResults) Less(i, j int) bool { return strings.Compare(a[i].Event, a[j].Event) == 1 }
