package attribpush

import (
	"github.com/MISingularity/deepshare2/pkg/log"
	"sync"
)

type AttributionBuffer interface {
	PopAttributions(appID string) []*AttributionPushInfo
	PutAttribution(appID string, info *AttributionPushInfo)
	ListAppIDs() []string
}

type AttrInfos struct {
	mutex     sync.Mutex
	attrInfos []*AttributionPushInfo
}

func putAttrInfo(infos *AttrInfos, info *AttributionPushInfo) {
	infos.mutex.Lock()
	if infos.attrInfos == nil {
		infos.attrInfos = []*AttributionPushInfo{info}
	} else {
		infos.attrInfos = append(infos.attrInfos, info)
	}
	infos.mutex.Unlock()
	log.Debugf("putAttrInfo: %v %d\n", info, len(infos.attrInfos))
}

type InMemBuffer struct {
	m      sync.Mutex
	buffer map[string]*AttrInfos
}

func NewInMemBuffer() *InMemBuffer {
	return &InMemBuffer{
		buffer: make(map[string]*AttrInfos),
	}
}

func (b *InMemBuffer) PopAttributions(appID string) []*AttributionPushInfo {
	b.m.Lock()
	defer b.m.Unlock()
	if data, ok := b.buffer[appID]; ok {
		delete(b.buffer, appID)
		return data.attrInfos
	}
	return nil
}

func (b *InMemBuffer) PutAttribution(appID string, info *AttributionPushInfo) {
	b.m.Lock()
	defer b.m.Unlock()
	if _, ok := b.buffer[appID]; !ok {
		b.buffer[appID] = new(AttrInfos)
	}
	putAttrInfo(b.buffer[appID], info)
	log.Debugf("PutAttribution, cur buffer: %#v\n", b.buffer[appID].attrInfos)
}

func (b *InMemBuffer) ListAppIDs() []string {
	var ids []string
	for id, _ := range b.buffer {
		ids = append(ids, id)
	}
	return ids
}
