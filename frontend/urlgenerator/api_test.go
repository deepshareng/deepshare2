package urlgenerator

import (
	"errors"
	"reflect"
	"testing"
)

func TestGenURLInAppDataStringfy(t *testing.T) {
	tests := []struct {
		body  *GenURLPostBody
		wbody *GenURLPostBody
		werr  error
	}{
		{ // ok, InAppDataReq is string
			body:  &GenURLPostBody{InAppDataReq: "sss"},
			wbody: &GenURLPostBody{InAppDataReq: "sss", InAppData: "sss"},
			werr:  nil,
		},
		{ // ok, InAppDataReq is nil
			body:  &GenURLPostBody{InAppDataReq: nil},
			wbody: &GenURLPostBody{InAppDataReq: nil, InAppData: ""},
			werr:  nil,
		},
		{ // ok, InAppDataReq is a map (in http request, a json object)
			body:  &GenURLPostBody{InAppDataReq: map[string]interface{}{"a": "aaa", "b": "bbb"}},
			wbody: &GenURLPostBody{InAppDataReq: map[string]interface{}{"a": "aaa", "b": "bbb"}, InAppData: `{"a":"aaa","b":"bbb"}`},
			werr:  nil,
		},
		{ // bad, InAppDataReq can not be types other than string and map[string]interface{}
			body:  &GenURLPostBody{InAppDataReq: 1},
			wbody: &GenURLPostBody{InAppDataReq: 1, InAppData: ""},
			werr:  errors.New("InAppData type error"),
		},
	}

	for i, tt := range tests {
		gupb, err := StringfyInappData(tt.body)
		if !reflect.DeepEqual(*gupb, *tt.wbody) {
			t.Errorf("#%d: Inappdata Stringfied genBody = %v, want = %v\n", i, gupb, tt.wbody)
		}
		if !reflect.DeepEqual(err, tt.werr) {
			t.Errorf("#%d: StringfyInappData() returned error = %v, want = %v\n", i, err, tt.werr)
		}
	}
}
