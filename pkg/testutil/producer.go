package testutil

import (
	"encoding/json"

	"github.com/MISingularity/deepshare2/pkg/messaging"
)

type TestUtilProducer struct {
	C map[string]int
}

func NewTestUtilProducer() *TestUtilProducer {
	return &TestUtilProducer{make(map[string]int)}
}

func (tp *TestUtilProducer) Produce(topic []byte, e *messaging.Event) {
	s, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	tp.C[string(topic)+"#"+string(s)]++
}

func (tp *TestUtilProducer) Clear() {
	tp.C = make(map[string]int)
}
