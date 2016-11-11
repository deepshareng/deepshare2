package backup

import (
	"strconv"
	"time"

	pb "github.com/MISingularity/deepshare2/deepstats/backup/proto"
)

type BackupService interface {
	Insert(event pb.Event) error
	RetriveAllEvents() ([]pb.Event, error)
}

func convertTime(unixTime int64) string {
	generalTime := time.Unix(unixTime, 0)
	return strconv.Itoa(generalTime.Year()) + "-" + generalTime.Month().String()[0:3] + "-" + strconv.Itoa(generalTime.Day())
}
