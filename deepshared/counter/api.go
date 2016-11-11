package counter

import "github.com/MISingularity/deepshare2/api"

type Counter struct {
	Event string `json:"event"`
	Count int    `json:"count"`
}

type CounterRequest struct {
	ReceiverInfo api.ReceiverInfo `json:"receiver_info"`
	Counters     []*Counter       `json:"counters"`
}
