package storage

import "time"

type mockSimpleKV struct {
}

var mockData = map[string]string{
	"app:7713337217A6E150": `{"AppID":"7713337217A6E150","UserConf":{"BgWeChatAndroidTipUrl":"www.androidwechattip","BgWeChatIosTipUrl":"www.ioswechattip"},"Android":{"Scheme":"deepshare","Host":"com.singulariti.deepsharedemo","Pkg":"com.singulariti.deepsharedemo","DownloadUrl":"","IsDownloadDirectly":true,"YYBEnable":false},"Ios":{"Scheme":"deepsharedemo","DownloadUrl":"", "BundleID":"bundleID1", "TeamID":"teamID1", "YYBEnableBelow9":true, "YYBEnableAbove9":true},"YYBUrl":"","YYBEnable":false,"Theme":"0"}`,
}

func NewMockSimpleKV() SimpleKV {
	return &mockSimpleKV{}
}
func (mockKV *mockSimpleKV) Get(k []byte) ([]byte, error) {
	return []byte(mockData[string(k)]), nil
}

func (mockKV *mockSimpleKV) Set(k []byte, v []byte) error {
	return nil
}
func (mockKV *mockSimpleKV) SetEx(k []byte, v []byte, expiration time.Duration) error {
	return nil
}
func (db *mockSimpleKV) Delete(k []byte) error {
	return nil
}
func (db *mockSimpleKV) HSet(k []byte, hk string, v []byte) error {
	return nil
}
func (db *mockSimpleKV) HGet(k []byte, hk string) ([]byte, error) {
	return nil, nil
}
func (db *mockSimpleKV) HGetAll(k []byte) (map[string]string, error) {
	return nil, nil
}
func (db *mockSimpleKV) HDel(k []byte, hk string) error {
	return nil
}
func (db *mockSimpleKV) HIncrBy(k []byte, hk string, n int) error {
	return nil
}
func (db *mockSimpleKV) Exists(k []byte) bool {
	return false
}

func (db *mockSimpleKV) SAdd(k []byte, v string) error {
	return nil
}

func (db *mockSimpleKV) SRem(k []byte, v string) error {
	return nil
}

func (db *mockSimpleKV) SCard(k []byte) (int64, error) {
	return 0, nil
}

func (db *mockSimpleKV) SMembers(k []byte) ([]string, error) {
	return nil, nil
}
