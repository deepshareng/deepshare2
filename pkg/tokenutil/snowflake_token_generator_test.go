package tokenutil

import (
	"sync"
	"testing"
)

func TestNewSnowflakeTokenGeneratorParams(t *testing.T) {
	testData := []struct {
		workerID     int64
		dataCenterID int64
		error        bool
	}{
		{
			31,
			31,
			false,
		},
		{
			32,
			31,
			true,
		},
		{
			31,
			32,
			true,
		},
	}

	for i, testToken := range testData {
		_, err := NewSnowflakeTokenGenerator(testToken.workerID, testToken.dataCenterID)
		if err != nil && !testToken.error {
			t.Errorf("#%d Test New SnowflakeTokenGenerator failed, should got error", i)
		} else if err == nil && testToken.error {
			t.Errorf("#%d Test New SnowflakeTokenGenerator failed, should not got error", i)
		}
	}

}

func TestSnowflakeTokenGeneratorAync(t *testing.T) {
	tg, _ := NewSnowflakeTokenGenerator(1, 1)
	l := make(map[string]struct{})
	var mutex sync.Mutex
	for i := 0; i < 10; i++ {
		go func() {
			token, err := tg.Generate("1")
			go func(string, error) {
				if err != nil {
					t.Fatalf("Generate failed: %v", err)
				}
				mutex.Lock()
				_, ok := l[token]
				if ok {
					t.Errorf("SnowflakeTokenGenerator generates idential tokens")
				} else {
					l[token] = struct{}{}
				}
				mutex.Unlock()
			}(token, err)
		}()
	}

}
