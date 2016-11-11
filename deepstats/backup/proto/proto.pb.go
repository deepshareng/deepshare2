// Code generated by protoc-gen-go.
// source: proto.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	proto.proto

It has these top-level messages:
	Event
	UAInfo
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Event struct {
	AppID     string   `protobuf:"bytes,1,opt,name=AppID" json:"AppID,omitempty"`
	EventType string   `protobuf:"bytes,2,opt,name=EventType" json:"EventType,omitempty"`
	Channels  []string `protobuf:"bytes,3,rep,name=Channels" json:"Channels,omitempty"`
	SenderID  string   `protobuf:"bytes,4,opt,name=SenderID" json:"SenderID,omitempty"`
	CookieID  string   `protobuf:"bytes,5,opt,name=CookieID" json:"CookieID,omitempty"`
	UniqueID  string   `protobuf:"bytes,6,opt,name=UniqueID" json:"UniqueID,omitempty"`
	Count     int64    `protobuf:"varint,7,opt,name=Count" json:"Count,omitempty"`
	UAInfo    *UAInfo  `protobuf:"bytes,8,opt,name=UAInfo" json:"UAInfo,omitempty"`
	KVs       []byte   `protobuf:"bytes,9,opt,name=KVs,proto3" json:"KVs,omitempty"`
	Timestamp int64    `protobuf:"varint,10,opt,name=Timestamp" json:"Timestamp,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto1.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Event) GetUAInfo() *UAInfo {
	if m != nil {
		return m.UAInfo
	}
	return nil
}

type UAInfo struct {
	Ua          string `protobuf:"bytes,1,opt,name=Ua" json:"Ua,omitempty"`
	Ip          string `protobuf:"bytes,2,opt,name=Ip" json:"Ip,omitempty"`
	Os          string `protobuf:"bytes,3,opt,name=Os" json:"Os,omitempty"`
	OsVersion   string `protobuf:"bytes,4,opt,name=OsVersion" json:"OsVersion,omitempty"`
	Brand       string `protobuf:"bytes,5,opt,name=Brand" json:"Brand,omitempty"`
	Browser     string `protobuf:"bytes,6,opt,name=Browser" json:"Browser,omitempty"`
	IsWechat    bool   `protobuf:"varint,7,opt,name=IsWechat" json:"IsWechat,omitempty"`
	IsWeibo     bool   `protobuf:"varint,8,opt,name=IsWeibo" json:"IsWeibo,omitempty"`
	IsQQ        bool   `protobuf:"varint,9,opt,name=IsQQ" json:"IsQQ,omitempty"`
	ChromeMajor int64  `protobuf:"varint,10,opt,name=ChromeMajor" json:"ChromeMajor,omitempty"`
}

func (m *UAInfo) Reset()                    { *m = UAInfo{} }
func (m *UAInfo) String() string            { return proto1.CompactTextString(m) }
func (*UAInfo) ProtoMessage()               {}
func (*UAInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto1.RegisterType((*Event)(nil), "proto.Event")
	proto1.RegisterType((*UAInfo)(nil), "proto.UAInfo")
}

var fileDescriptor0 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x3c, 0x8f, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0x87, 0x95, 0xa6, 0x49, 0x93, 0x4b, 0x2b, 0xc0, 0x2c, 0x5e, 0x90, 0xaa, 0x4e, 0x4c, 0x1d,
	0xe0, 0x09, 0xda, 0xc0, 0x10, 0x21, 0x54, 0x55, 0x34, 0x65, 0x76, 0xe9, 0xa1, 0x04, 0x88, 0x6d,
	0xec, 0x14, 0xc4, 0x03, 0xf1, 0x1c, 0xbc, 0x1a, 0xfe, 0xd7, 0x2e, 0xf6, 0xdd, 0xe7, 0x3b, 0xf9,
	0xfb, 0x41, 0x21, 0x95, 0xe8, 0xc5, 0xdc, 0x9d, 0x24, 0x71, 0xd7, 0xec, 0x2f, 0x82, 0xe4, 0xfe,
	0x0b, 0x79, 0x4f, 0x26, 0x90, 0x2c, 0xa4, 0xac, 0xee, 0x68, 0x34, 0x8d, 0xae, 0x73, 0x72, 0x01,
	0xb9, 0xe3, 0x9b, 0x1f, 0x89, 0x74, 0xe0, 0xd0, 0x39, 0x64, 0x65, 0xc3, 0x38, 0xc7, 0x0f, 0x4d,
	0xe3, 0x69, 0xec, 0xc9, 0x13, 0xf2, 0x3d, 0x2a, 0xb3, 0x36, 0x3c, 0xcd, 0x08, 0xf1, 0xde, 0xa2,
	0x21, 0xc9, 0x91, 0xd4, 0xbc, 0xfd, 0x3c, 0x58, 0x92, 0x3a, 0x62, 0x7e, 0x2a, 0xc5, 0x81, 0xf7,
	0x74, 0x64, 0xda, 0x98, 0x5c, 0x41, 0x5a, 0x2f, 0x2a, 0xfe, 0x2a, 0x68, 0x66, 0xfa, 0xe2, 0x66,
	0xe2, 0x0d, 0xe7, 0x1e, 0x92, 0x02, 0xe2, 0x87, 0xad, 0xa6, 0xb9, 0x79, 0x1b, 0x5b, 0xab, 0x4d,
	0xdb, 0xa1, 0xee, 0x59, 0x27, 0x29, 0xd8, 0xf5, 0xd9, 0x6f, 0x74, 0xdc, 0x27, 0x00, 0x83, 0x9a,
	0x05, 0x7f, 0x53, 0x57, 0x32, 0x88, 0x9b, 0x7a, 0x65, 0x95, 0x43, 0xae, 0x95, 0xde, 0xa2, 0xd2,
	0xad, 0xe0, 0xc1, 0xd9, 0xf8, 0x2c, 0x15, 0xe3, 0xfb, 0x20, 0x7c, 0x06, 0xa3, 0xa5, 0x12, 0xdf,
	0x1a, 0x55, 0xf0, 0x35, 0x09, 0x2a, 0xfd, 0x8c, 0x2f, 0x0d, 0xf3, 0xca, 0x99, 0x1d, 0xb1, 0xa4,
	0xdd, 0x79, 0xe7, 0x8c, 0x8c, 0x61, 0x58, 0xe9, 0xf5, 0xda, 0x59, 0x66, 0xe4, 0x12, 0x8a, 0xb2,
	0x51, 0xa2, 0xc3, 0x47, 0xf6, 0x26, 0x94, 0xf7, 0xdc, 0xa5, 0x2e, 0xd5, 0xed, 0x7f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x2b, 0x10, 0xa3, 0xce, 0x86, 0x01, 0x00, 0x00,
}
