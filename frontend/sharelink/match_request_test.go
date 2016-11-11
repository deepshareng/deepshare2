package sharelink

import "testing"

func TestSetupMatchBody(t *testing.T) {
	tests := []struct {
		inAppData     string
		senderID      string
		channels      string
		sdkInfo       string
		ip            string
		ua            string
		wantMatchBody string
	}{
		{
			`{"k1": "v1","k2": "v2"}`,
			`senderID1`,
			"channel1|chanel2",
			"android1.2",
			"ip1",
			"ua1",
			`{"inapp_data":"{\"k1\": \"v1\",\"k2\": \"v2\"}","sender_info":{"sender_id":"senderID1","channels":["channel1","chanel2"],"sdk_info":"android1.2"},"client_ip":"ip1","client_ua":"ua1"}`,
		},
	}
	sl := &Sharelink{}
	for i, tt := range tests {
		matchBody := sl.setupMatchBody(tt.inAppData, tt.senderID, tt.channels, tt.sdkInfo, tt.ip, tt.ua)
		if tt.wantMatchBody != matchBody {
			t.Errorf("#%d: Match Body = %s, want = %s", i, matchBody, tt.wantMatchBody)
		}
	}
}
