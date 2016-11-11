package device_stat

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"github.com/MISingularity/deepshare2/pkg/messaging"
	"github.com/MISingularity/deepshare2/pkg/storage"
)

const RedisPrefixDevice = "device_stat_"

type DeviceStatService interface {
	Insert(device *messaging.Event) error
	Count(os string) (int64, error)
}

type GeneralDeviceStat struct {
	skv                                               storage.SimpleKV
	androidDeviceColl, iosDeviceColl, otherDeviceColl string
}

func NewGeneralDeviceStat(skv storage.SimpleKV, devicePrefix string) DeviceStatService {
	return &GeneralDeviceStat{
		skv:               skv,
		androidDeviceColl: devicePrefix + "android",
		iosDeviceColl:     devicePrefix + "ios",
		otherDeviceColl:   devicePrefix + "other",
	}
}

func (gs *GeneralDeviceStat) Insert(device *messaging.Event) error {
	key := gs.otherDeviceColl
	if device.UAInfo.IsAndroid() {
		key = gs.androidDeviceColl
	} else if device.UAInfo.IsIos() {
		key = gs.iosDeviceColl
	}
	if device.UniqueID != "" {
		log.Infof("SAdd %s into Set", device.UniqueID)
		err := gs.skv.SAdd([]byte(key), device.UniqueID)
		if err != nil {
			return err
		}
	}
	if device.SenderID != "" {
		log.Infof("SAdd %s into Set", device.SenderID)
		err := gs.skv.SAdd([]byte(key), device.SenderID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *GeneralDeviceStat) Count(os string) (int64, error) {
	iosDeviceNum, err := gs.skv.SCard([]byte(gs.iosDeviceColl))
	if err != nil {
		return 0, err
	}
	androidDeviceNum, err := gs.skv.SCard([]byte(gs.androidDeviceColl))
	if err != nil {
		return 0, err
	}
	otherDeviceNum, err := gs.skv.SCard([]byte(gs.otherDeviceColl))
	if err != nil {
		return 0, err
	}
	switch os {
	case "all":
		return iosDeviceNum + androidDeviceNum + otherDeviceNum, nil
	case "android":
		return androidDeviceNum, nil
	case "ios":
		return iosDeviceNum, nil
	case "other":
		return otherDeviceNum, nil
	}
	return 0, nil
}
