package tokenutil

import (
	"errors"
	"sync"
	"time"

	"github.com/MISingularity/deepshare2/pkg/log"
)

const (
	workerIdBits = 5
	maxWorkerId  = (1 << workerIdBits) - 1

	dataCenterIdBits = 5
	maxDataCenterId  = (1 << dataCenterIdBits) - 1

	sequenceBits = 12
	sequenceMask = int64(-1 ^ (-1 << sequenceBits))

	workerIdShift     = sequenceBits
	dataCenterIdShift = sequenceBits + workerIdBits
	timestampShift    = sequenceBits + workerIdBits + dataCenterIdBits

	//the epoch is a number of milliseconds since the UNIX UTC epoch.
	//the value is a time when I log timeStamp()
	epoch = int64(1445500382851)
)

type snowflakeTokenGenerator struct {
	workerID      int64
	dataCenterID  int64
	lastTimestamp int64
	sequence      int64
	mutex         sync.Mutex
}

//Twitter Snowflake arithmetic to generate Unique ID
//Token would not repeat given different Worker ID and dataCenter ID
//Worker ID and data center ID must be in the range [0 : 31].
//Namespace is not used in Snowflake implementation
func NewSnowflakeTokenGenerator(workerID, dataCenterID int64) (TokenGenerator, error) {
	if workerID < 0 || workerID > maxWorkerId {
		log.Errorf("snowflakeTokenGenerator dataCenterID can't be greater than %d or less than 0", maxWorkerId)
		err := errors.New("workerID is illegal")
		return nil, err
	}

	if dataCenterID < 0 || dataCenterID > maxDataCenterId {
		log.Errorf("snowflakeTokenGenerator dataCenterID can't be greater than %d or less than 0", maxDataCenterId)
		err := errors.New("dataCenterID is illegal")
		return nil, err
	}
	return &snowflakeTokenGenerator{
		workerID:      workerID,
		dataCenterID:  dataCenterID,
		lastTimestamp: 0,
		sequence:      0,
	}, nil
}

//Generate next token
func (tg *snowflakeTokenGenerator) Generate(namespace string) (string, error) {
	tg.mutex.Lock()
	defer tg.mutex.Unlock()
	// Get the current timestamp.
	timestamp := timeStamp()
	if tg.lastTimestamp > timestamp {
		err := errors.New("system clock is moving backwards")
		return "", err
	}
	if timestamp == tg.lastTimestamp {
		tg.sequence = (tg.sequence + 1) & sequenceMask
		if tg.sequence == 0 {
			// do get timestamp until we get the next millisecond.
			for timestamp == tg.lastTimestamp {
				timestamp = timeStamp()
			}
			tg.lastTimestamp = timestamp
		}
	} else {
		tg.sequence = 0
		tg.lastTimestamp = timestamp
	}
	tokenInt := ((timestamp - epoch) << timestampShift) | (tg.dataCenterID << dataCenterIdShift) | (tg.workerID << workerIdShift) | tg.sequence
	tokenStr := Encode(tokenInt)
	log.Debugf("TokenGenerator, SnowflakeTokenGenerator generate toke = %s", tokenStr)
	return tokenStr, nil
}

func timeStamp() int64 {
	return int64(time.Now().UTC().UnixNano() / 1000000)
}
