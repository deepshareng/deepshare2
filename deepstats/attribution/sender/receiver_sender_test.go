package sender

import (
	"testing"

	"github.com/MISingularity/deepshare2/pkg/testutil"
)

func TestReceiverSender(t *testing.T) {
	tests := []struct {
		eventType string
		senderID  string

		receiverID string
		wSenderID  string
	}{
		{"install", "s1", "r", "s1"},
		{"close", "", "r", ""},
		{"open", "", "r", ""},
		{"open", "s2", "r", "s2"},
		{"close", "", "r", ""},
	}

	c := testutil.MustNewMongoColl("testRS", "rs")
	c.DropCollection()
	rsm := &ReceiverSenderMapping{c}

	for i, tt := range tests {
		switch tt.eventType {
		case "install":
			rsm.OnInstall(tt.receiverID, tt.senderID)
		case "open":
			rsm.OnOpen(tt.receiverID, tt.senderID)
		case "close":
			rsm.OnClose(tt.receiverID)
		}
		if s := rsm.GetSenderID(tt.receiverID); s != tt.wSenderID {
			t.Errorf("#%d Expected senderID: %s, got: %s", i, tt.wSenderID, s)
		}
	}
}
