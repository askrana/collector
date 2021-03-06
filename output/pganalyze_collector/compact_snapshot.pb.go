// Code generated by protoc-gen-go.
// source: compact_snapshot.proto
// DO NOT EDIT!

package pganalyze_collector

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type CompactSnapshot struct {
	// Basic information about this snapshot
	SnapshotVersionMajor int32                      `protobuf:"varint,1,opt,name=snapshot_version_major,json=snapshotVersionMajor" json:"snapshot_version_major,omitempty"`
	SnapshotVersionMinor int32                      `protobuf:"varint,2,opt,name=snapshot_version_minor,json=snapshotVersionMinor" json:"snapshot_version_minor,omitempty"`
	CollectorVersion     string                     `protobuf:"bytes,3,opt,name=collector_version,json=collectorVersion" json:"collector_version,omitempty"`
	SnapshotUuid         string                     `protobuf:"bytes,4,opt,name=snapshot_uuid,json=snapshotUuid" json:"snapshot_uuid,omitempty"`
	CollectedAt          *google_protobuf.Timestamp `protobuf:"bytes,5,opt,name=collected_at,json=collectedAt" json:"collected_at,omitempty"`
	BaseRefs             *CompactSnapshot_BaseRefs  `protobuf:"bytes,6,opt,name=base_refs,json=baseRefs" json:"base_refs,omitempty"`
	// Types that are valid to be assigned to Data:
	//	*CompactSnapshot_LogSnapshot
	//	*CompactSnapshot_SystemSnapshot
	Data isCompactSnapshot_Data `protobuf_oneof:"data"`
}

func (m *CompactSnapshot) Reset()                    { *m = CompactSnapshot{} }
func (m *CompactSnapshot) String() string            { return proto.CompactTextString(m) }
func (*CompactSnapshot) ProtoMessage()               {}
func (*CompactSnapshot) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

type isCompactSnapshot_Data interface {
	isCompactSnapshot_Data()
}

type CompactSnapshot_LogSnapshot struct {
	LogSnapshot *CompactLogSnapshot `protobuf:"bytes,10,opt,name=log_snapshot,json=logSnapshot,oneof"`
}
type CompactSnapshot_SystemSnapshot struct {
	SystemSnapshot *CompactSystemSnapshot `protobuf:"bytes,11,opt,name=system_snapshot,json=systemSnapshot,oneof"`
}

func (*CompactSnapshot_LogSnapshot) isCompactSnapshot_Data()    {}
func (*CompactSnapshot_SystemSnapshot) isCompactSnapshot_Data() {}

func (m *CompactSnapshot) GetData() isCompactSnapshot_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *CompactSnapshot) GetSnapshotVersionMajor() int32 {
	if m != nil {
		return m.SnapshotVersionMajor
	}
	return 0
}

func (m *CompactSnapshot) GetSnapshotVersionMinor() int32 {
	if m != nil {
		return m.SnapshotVersionMinor
	}
	return 0
}

func (m *CompactSnapshot) GetCollectorVersion() string {
	if m != nil {
		return m.CollectorVersion
	}
	return ""
}

func (m *CompactSnapshot) GetSnapshotUuid() string {
	if m != nil {
		return m.SnapshotUuid
	}
	return ""
}

func (m *CompactSnapshot) GetCollectedAt() *google_protobuf.Timestamp {
	if m != nil {
		return m.CollectedAt
	}
	return nil
}

func (m *CompactSnapshot) GetBaseRefs() *CompactSnapshot_BaseRefs {
	if m != nil {
		return m.BaseRefs
	}
	return nil
}

func (m *CompactSnapshot) GetLogSnapshot() *CompactLogSnapshot {
	if x, ok := m.GetData().(*CompactSnapshot_LogSnapshot); ok {
		return x.LogSnapshot
	}
	return nil
}

func (m *CompactSnapshot) GetSystemSnapshot() *CompactSystemSnapshot {
	if x, ok := m.GetData().(*CompactSnapshot_SystemSnapshot); ok {
		return x.SystemSnapshot
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*CompactSnapshot) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CompactSnapshot_OneofMarshaler, _CompactSnapshot_OneofUnmarshaler, _CompactSnapshot_OneofSizer, []interface{}{
		(*CompactSnapshot_LogSnapshot)(nil),
		(*CompactSnapshot_SystemSnapshot)(nil),
	}
}

func _CompactSnapshot_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*CompactSnapshot)
	// data
	switch x := m.Data.(type) {
	case *CompactSnapshot_LogSnapshot:
		b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.LogSnapshot); err != nil {
			return err
		}
	case *CompactSnapshot_SystemSnapshot:
		b.EncodeVarint(11<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SystemSnapshot); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("CompactSnapshot.Data has unexpected type %T", x)
	}
	return nil
}

func _CompactSnapshot_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*CompactSnapshot)
	switch tag {
	case 10: // data.log_snapshot
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CompactLogSnapshot)
		err := b.DecodeMessage(msg)
		m.Data = &CompactSnapshot_LogSnapshot{msg}
		return true, err
	case 11: // data.system_snapshot
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(CompactSystemSnapshot)
		err := b.DecodeMessage(msg)
		m.Data = &CompactSnapshot_SystemSnapshot{msg}
		return true, err
	default:
		return false, nil
	}
}

func _CompactSnapshot_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*CompactSnapshot)
	// data
	switch x := m.Data.(type) {
	case *CompactSnapshot_LogSnapshot:
		s := proto.Size(x.LogSnapshot)
		n += proto.SizeVarint(10<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CompactSnapshot_SystemSnapshot:
		s := proto.Size(x.SystemSnapshot)
		n += proto.SizeVarint(11<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type CompactSnapshot_BaseRefs struct {
	RoleReferences     []*RoleReference     `protobuf:"bytes,1,rep,name=role_references,json=roleReferences" json:"role_references,omitempty"`
	DatabaseReferences []*DatabaseReference `protobuf:"bytes,2,rep,name=database_references,json=databaseReferences" json:"database_references,omitempty"`
	QueryReferences    []*QueryReference    `protobuf:"bytes,3,rep,name=query_references,json=queryReferences" json:"query_references,omitempty"`
	QueryInformations  []*QueryInformation  `protobuf:"bytes,4,rep,name=query_informations,json=queryInformations" json:"query_informations,omitempty"`
	RelationReferences []*RelationReference `protobuf:"bytes,5,rep,name=relation_references,json=relationReferences" json:"relation_references,omitempty"`
}

func (m *CompactSnapshot_BaseRefs) Reset()                    { *m = CompactSnapshot_BaseRefs{} }
func (m *CompactSnapshot_BaseRefs) String() string            { return proto.CompactTextString(m) }
func (*CompactSnapshot_BaseRefs) ProtoMessage()               {}
func (*CompactSnapshot_BaseRefs) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0, 0} }

func (m *CompactSnapshot_BaseRefs) GetRoleReferences() []*RoleReference {
	if m != nil {
		return m.RoleReferences
	}
	return nil
}

func (m *CompactSnapshot_BaseRefs) GetDatabaseReferences() []*DatabaseReference {
	if m != nil {
		return m.DatabaseReferences
	}
	return nil
}

func (m *CompactSnapshot_BaseRefs) GetQueryReferences() []*QueryReference {
	if m != nil {
		return m.QueryReferences
	}
	return nil
}

func (m *CompactSnapshot_BaseRefs) GetQueryInformations() []*QueryInformation {
	if m != nil {
		return m.QueryInformations
	}
	return nil
}

func (m *CompactSnapshot_BaseRefs) GetRelationReferences() []*RelationReference {
	if m != nil {
		return m.RelationReferences
	}
	return nil
}

func init() {
	proto.RegisterType((*CompactSnapshot)(nil), "pganalyze.collector.CompactSnapshot")
	proto.RegisterType((*CompactSnapshot_BaseRefs)(nil), "pganalyze.collector.CompactSnapshot.BaseRefs")
}

func init() { proto.RegisterFile("compact_snapshot.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 480 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x92, 0x5f, 0x6f, 0xd3, 0x30,
	0x14, 0xc5, 0xc9, 0xfa, 0x87, 0xcd, 0x29, 0xeb, 0xe6, 0xa1, 0xc9, 0x8a, 0x84, 0xa8, 0x36, 0x01,
	0x15, 0x88, 0x4c, 0x1a, 0xbc, 0xf2, 0xc0, 0xe0, 0x81, 0x3f, 0x03, 0x09, 0xb3, 0xc1, 0x63, 0xe4,
	0x26, 0x6e, 0x16, 0xe4, 0xc4, 0xa9, 0xaf, 0x83, 0x54, 0xbe, 0x14, 0xcf, 0x7c, 0x3b, 0x54, 0x27,
	0x4e, 0xb2, 0x28, 0xeb, 0x5b, 0xee, 0xbd, 0xe7, 0xfc, 0xec, 0x7b, 0x1c, 0x74, 0x1c, 0xca, 0x34,
	0x67, 0xa1, 0x0e, 0x20, 0x63, 0x39, 0xdc, 0x48, 0xed, 0xe7, 0x4a, 0x6a, 0x89, 0x8f, 0xf2, 0x98,
	0x65, 0x4c, 0xac, 0xff, 0x70, 0x3f, 0x94, 0x42, 0xf0, 0x50, 0x4b, 0xe5, 0x3d, 0x8e, 0xa5, 0x8c,
	0x05, 0x3f, 0x33, 0x92, 0x45, 0xb1, 0x3c, 0xd3, 0x49, 0xca, 0x41, 0xb3, 0x34, 0x2f, 0x5d, 0x9e,
	0x67, 0x69, 0x42, 0xc6, 0x1d, 0xa2, 0xf7, 0xa8, 0x3e, 0x69, 0x0d, 0x9a, 0xa7, 0xdd, 0xf1, 0x04,
	0x6e, 0x98, 0xe2, 0x51, 0x59, 0x9d, 0xfc, 0xbb, 0x8f, 0xa6, 0xef, 0x4a, 0xfd, 0xf7, 0x4a, 0x87,
	0x5f, 0xa3, 0x63, 0xeb, 0x09, 0x7e, 0x73, 0x05, 0x89, 0xcc, 0x82, 0x94, 0xfd, 0x92, 0x8a, 0x38,
	0x33, 0x67, 0x3e, 0xa2, 0x0f, 0xed, 0xf4, 0x47, 0x39, 0xfc, 0xb2, 0x99, 0xf5, 0xbb, 0x92, 0x4c,
	0x2a, 0xb2, 0xd3, 0xef, 0xda, 0xcc, 0xf0, 0x0b, 0x74, 0x58, 0xaf, 0x6d, 0x6d, 0x64, 0x30, 0x73,
	0xe6, 0x7b, 0xf4, 0xa0, 0x1e, 0x54, 0x0e, 0x7c, 0x8a, 0x1e, 0xd4, 0x47, 0x14, 0x45, 0x12, 0x91,
	0xa1, 0x11, 0x4e, 0x6c, 0xf3, 0xba, 0x48, 0x22, 0xfc, 0x06, 0x4d, 0x2a, 0x23, 0x8f, 0x02, 0xa6,
	0xc9, 0x68, 0xe6, 0xcc, 0xdd, 0x73, 0xcf, 0x2f, 0x23, 0xf5, 0x6d, 0xa4, 0xfe, 0x95, 0x8d, 0x94,
	0xba, 0xb5, 0xfe, 0xad, 0xc6, 0x9f, 0xd0, 0xde, 0x82, 0x01, 0x0f, 0x14, 0x5f, 0x02, 0x19, 0x1b,
	0xef, 0x4b, 0xbf, 0xe7, 0x8d, 0xfc, 0x4e, 0x6a, 0xfe, 0x05, 0x03, 0x4e, 0xf9, 0x12, 0xe8, 0xee,
	0xa2, 0xfa, 0xc2, 0x97, 0x68, 0xd2, 0x7e, 0x1f, 0x82, 0x0c, 0xee, 0xd9, 0x36, 0xdc, 0xa5, 0x8c,
	0x2d, 0xf1, 0xc3, 0x3d, 0xea, 0x8a, 0xa6, 0xc4, 0xd7, 0x68, 0xda, 0x79, 0x51, 0xe2, 0x1a, 0xe0,
	0xf3, 0xad, 0xf7, 0x33, 0x96, 0x16, 0x73, 0x1f, 0x6e, 0x75, 0xbc, 0xbf, 0x03, 0xb4, 0x6b, 0xef,
	0x8e, 0x3f, 0xa3, 0xa9, 0x92, 0xc2, 0x6c, 0xcf, 0x15, 0xcf, 0x42, 0x0e, 0xc4, 0x99, 0x0d, 0xe6,
	0xee, 0xf9, 0x49, 0xef, 0x19, 0x54, 0x8a, 0x8d, 0xaf, 0x94, 0xd2, 0x7d, 0xd5, 0x2e, 0x01, 0xff,
	0x44, 0x47, 0x11, 0xd3, 0xcc, 0xc6, 0x69, 0x81, 0x3b, 0x06, 0xf8, 0xb4, 0x17, 0xf8, 0xbe, 0xd2,
	0x37, 0x50, 0x1c, 0x75, 0x5b, 0x80, 0xbf, 0xa2, 0x83, 0x55, 0xc1, 0xd5, 0xba, 0x4d, 0x1d, 0x18,
	0xea, 0x69, 0x2f, 0xf5, 0xdb, 0x46, 0xdc, 0x20, 0xa7, 0xab, 0x5b, 0x35, 0xe0, 0x2b, 0x84, 0x4b,
	0x5e, 0x92, 0x2d, 0xa5, 0x4a, 0x99, 0x4e, 0x64, 0x06, 0x64, 0x68, 0x88, 0x4f, 0xee, 0x26, 0x7e,
	0x6c, 0xd4, 0xf4, 0x70, 0xd5, 0xe9, 0x98, 0xf5, 0x15, 0x17, 0xa6, 0x68, 0x5f, 0x74, 0xb4, 0x65,
	0x7d, 0x5a, 0xe9, 0x5b, 0xeb, 0xab, 0x6e, 0x0b, 0x2e, 0xc6, 0x68, 0x68, 0x42, 0x19, 0x9b, 0x7f,
	0xf9, 0xd5, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0x67, 0xda, 0x4c, 0x88, 0x5b, 0x04, 0x00, 0x00,
}
