// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cyber/cron/v1beta1/types.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Params struct {
	MaxSlots uint32 `protobuf:"varint,1,opt,name=max_slots,json=maxSlots,proto3" json:"max_slots,omitempty" yaml:"max_slots"`
	MaxGas   uint32 `protobuf:"varint,2,opt,name=max_gas,json=maxGas,proto3" json:"max_gas,omitempty" yaml:"max_gas"`
	FeeTtl   uint32 `protobuf:"varint,3,opt,name=fee_ttl,json=feeTtl,proto3" json:"fee_ttl,omitempty" yaml:"fee_ttl"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_e79f11c3c745de7d, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetMaxSlots() uint32 {
	if m != nil {
		return m.MaxSlots
	}
	return 0
}

func (m *Params) GetMaxGas() uint32 {
	if m != nil {
		return m.MaxGas
	}
	return 0
}

func (m *Params) GetFeeTtl() uint32 {
	if m != nil {
		return m.FeeTtl
	}
	return 0
}

type Job struct {
	Creator  string  `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty" yaml:"creator"`
	Contract string  `protobuf:"bytes,2,opt,name=contract,proto3" json:"contract,omitempty" yaml:"contract"`
	Trigger  Trigger `protobuf:"bytes,3,opt,name=trigger,proto3" json:"trigger" yaml:"trigger"`
	Load     Load    `protobuf:"bytes,4,opt,name=load,proto3" json:"load" yaml:"load"`
	Label    string  `protobuf:"bytes,5,opt,name=label,proto3" json:"label,omitempty" yaml:"label"`
	Cid      string  `protobuf:"bytes,6,opt,name=cid,proto3" json:"cid,omitempty" yaml:"cid"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}
func (*Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_e79f11c3c745de7d, []int{1}
}
func (m *Job) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Job.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Job.Merge(m, src)
}
func (m *Job) XXX_Size() int {
	return m.Size()
}
func (m *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Job.DiscardUnknown(m)
}

var xxx_messageInfo_Job proto.InternalMessageInfo

type Trigger struct {
	Period uint64 `protobuf:"varint,1,opt,name=period,proto3" json:"period,omitempty" yaml:"period"`
	Block  uint64 `protobuf:"varint,2,opt,name=block,proto3" json:"block,omitempty" yaml:"block"`
}

func (m *Trigger) Reset()         { *m = Trigger{} }
func (m *Trigger) String() string { return proto.CompactTextString(m) }
func (*Trigger) ProtoMessage()    {}
func (*Trigger) Descriptor() ([]byte, []int) {
	return fileDescriptor_e79f11c3c745de7d, []int{2}
}
func (m *Trigger) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Trigger) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Trigger.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Trigger) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Trigger.Merge(m, src)
}
func (m *Trigger) XXX_Size() int {
	return m.Size()
}
func (m *Trigger) XXX_DiscardUnknown() {
	xxx_messageInfo_Trigger.DiscardUnknown(m)
}

var xxx_messageInfo_Trigger proto.InternalMessageInfo

type Load struct {
	CallData string                                  `protobuf:"bytes,1,opt,name=call_data,json=callData,proto3" json:"call_data,omitempty" yaml:"call_data"`
	GasPrice github_com_cosmos_cosmos_sdk_types.Coin `protobuf:"bytes,2,opt,name=gas_price,json=gasPrice,proto3,casttype=github.com/cosmos/cosmos-sdk/types.Coin" json:"gas_price" yaml:"gas_price"`
}

func (m *Load) Reset()         { *m = Load{} }
func (m *Load) String() string { return proto.CompactTextString(m) }
func (*Load) ProtoMessage()    {}
func (*Load) Descriptor() ([]byte, []int) {
	return fileDescriptor_e79f11c3c745de7d, []int{3}
}
func (m *Load) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Load) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Load.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Load) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Load.Merge(m, src)
}
func (m *Load) XXX_Size() int {
	return m.Size()
}
func (m *Load) XXX_DiscardUnknown() {
	xxx_messageInfo_Load.DiscardUnknown(m)
}

var xxx_messageInfo_Load proto.InternalMessageInfo

type JobStats struct {
	Creator   string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty" yaml:"creator"`
	Contract  string `protobuf:"bytes,2,opt,name=contract,proto3" json:"contract,omitempty" yaml:"contract"`
	Label     string `protobuf:"bytes,3,opt,name=label,proto3" json:"label,omitempty" yaml:"label"`
	Calls     uint64 `protobuf:"varint,4,opt,name=calls,proto3" json:"calls,omitempty" yaml:"calls"`
	Fees      uint64 `protobuf:"varint,5,opt,name=fees,proto3" json:"fees,omitempty" yaml:"fees"`
	Gas       uint64 `protobuf:"varint,6,opt,name=gas,proto3" json:"gas,omitempty" yaml:"gas"`
	LastBlock uint64 `protobuf:"varint,7,opt,name=last_block,json=lastBlock,proto3" json:"last_block,omitempty" yaml:"last_block"`
}

func (m *JobStats) Reset()         { *m = JobStats{} }
func (m *JobStats) String() string { return proto.CompactTextString(m) }
func (*JobStats) ProtoMessage()    {}
func (*JobStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_e79f11c3c745de7d, []int{4}
}
func (m *JobStats) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *JobStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_JobStats.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *JobStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobStats.Merge(m, src)
}
func (m *JobStats) XXX_Size() int {
	return m.Size()
}
func (m *JobStats) XXX_DiscardUnknown() {
	xxx_messageInfo_JobStats.DiscardUnknown(m)
}

var xxx_messageInfo_JobStats proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "cyber.cron.v1beta1.Params")
	proto.RegisterType((*Job)(nil), "cyber.cron.v1beta1.Job")
	proto.RegisterType((*Trigger)(nil), "cyber.cron.v1beta1.Trigger")
	proto.RegisterType((*Load)(nil), "cyber.cron.v1beta1.Load")
	proto.RegisterType((*JobStats)(nil), "cyber.cron.v1beta1.JobStats")
}

func init() { proto.RegisterFile("cyber/cron/v1beta1/types.proto", fileDescriptor_e79f11c3c745de7d) }

var fileDescriptor_e79f11c3c745de7d = []byte{
	// 674 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x3b, 0x6f, 0xd4, 0x4c,
	0x14, 0x5d, 0x67, 0x9d, 0x7d, 0x4c, 0xbe, 0xbc, 0x26, 0xf9, 0xc0, 0x04, 0xc9, 0x8e, 0x06, 0x29,
	0x80, 0x20, 0xb6, 0x12, 0xa8, 0xd2, 0xe1, 0x20, 0x81, 0x22, 0x90, 0xa2, 0x49, 0x2a, 0x9a, 0xd5,
	0xd8, 0x9e, 0x35, 0x56, 0xec, 0x9d, 0xc8, 0x33, 0xa0, 0x8d, 0xf8, 0x03, 0x94, 0xf4, 0x34, 0xf9,
	0x29, 0x94, 0x11, 0x55, 0x4a, 0x2a, 0x0b, 0x25, 0x4d, 0x6a, 0x97, 0x54, 0x68, 0x1e, 0xbb, 0x71,
	0x00, 0xd1, 0x51, 0xad, 0xf7, 0x9e, 0x73, 0xee, 0xdc, 0xb9, 0xf7, 0xdc, 0x01, 0x6e, 0x7c, 0x12,
	0xd1, 0x32, 0x88, 0x4b, 0x36, 0x0a, 0xde, 0x6f, 0x45, 0x54, 0x90, 0xad, 0x40, 0x9c, 0x1c, 0x53,
	0xee, 0x1f, 0x97, 0x4c, 0x30, 0x08, 0x15, 0xee, 0x4b, 0xdc, 0x37, 0xf8, 0xda, 0x6a, 0xca, 0x52,
	0xa6, 0xe0, 0x40, 0x7e, 0x69, 0xe6, 0xda, 0xed, 0x98, 0xf1, 0x82, 0xf1, 0x81, 0x06, 0x62, 0x96,
	0x8d, 0x34, 0x80, 0x3e, 0x5b, 0xa0, 0xb3, 0x4f, 0x4a, 0x52, 0x70, 0xb8, 0x05, 0xfa, 0x05, 0x19,
	0x0f, 0x78, 0xce, 0x04, 0x77, 0xac, 0x75, 0xeb, 0xc1, 0x7c, 0xb8, 0x5a, 0x57, 0xde, 0xd2, 0x09,
	0x29, 0xf2, 0x1d, 0x34, 0x85, 0x10, 0xee, 0x15, 0x64, 0x7c, 0x20, 0x3f, 0xe1, 0x23, 0xd0, 0x95,
	0xf1, 0x94, 0x70, 0x67, 0x46, 0x09, 0x60, 0x5d, 0x79, 0x0b, 0xd7, 0x82, 0x94, 0x70, 0x84, 0x3b,
	0x05, 0x19, 0xbf, 0x20, 0x8a, 0x3c, 0xa4, 0x74, 0x20, 0x44, 0xee, 0xb4, 0x7f, 0x25, 0x1b, 0x00,
	0xe1, 0xce, 0x90, 0xd2, 0x43, 0x91, 0xef, 0xd8, 0x57, 0xa7, 0x9e, 0x85, 0xbe, 0xce, 0x80, 0xf6,
	0x1e, 0x8b, 0xe0, 0x63, 0xd0, 0x8d, 0x4b, 0x4a, 0x04, 0x2b, 0x55, 0x61, 0xfd, 0xa6, 0xd4, 0x00,
	0x08, 0x4f, 0x28, 0x30, 0x00, 0xbd, 0x98, 0x8d, 0x44, 0x49, 0x62, 0xa1, 0xca, 0xea, 0x87, 0x2b,
	0x75, 0xe5, 0x2d, 0x1a, 0xba, 0x41, 0x10, 0x9e, 0x92, 0xe0, 0x6b, 0xd0, 0x15, 0x65, 0x96, 0xa6,
	0xb4, 0x54, 0x95, 0xcd, 0x6d, 0xdf, 0xf5, 0x7f, 0xef, 0xac, 0x7f, 0xa8, 0x29, 0xe1, 0xad, 0xb3,
	0xca, 0x6b, 0x5d, 0x9f, 0x6f, 0x94, 0x08, 0x4f, 0x72, 0xc0, 0x67, 0xc0, 0xce, 0x19, 0x49, 0x1c,
	0x5b, 0xe5, 0x72, 0xfe, 0x94, 0xeb, 0x15, 0x23, 0x49, 0xb8, 0x62, 0x12, 0xcd, 0xe9, 0x44, 0x52,
	0x83, 0xb0, 0x92, 0xc2, 0x0d, 0x30, 0x9b, 0x93, 0x88, 0xe6, 0xce, 0xac, 0xaa, 0x7f, 0xa9, 0xae,
	0xbc, 0xff, 0x0c, 0x4b, 0x86, 0x11, 0xd6, 0x30, 0x5c, 0x07, 0xed, 0x38, 0x4b, 0x9c, 0x8e, 0x62,
	0x2d, 0xd4, 0x95, 0x07, 0xcc, 0x2d, 0xb3, 0x04, 0x61, 0x09, 0xed, 0xf4, 0x3e, 0x9e, 0x7a, 0xad,
	0xab, 0x53, 0xaf, 0x85, 0x86, 0xa0, 0x6b, 0xae, 0x00, 0x1f, 0x82, 0xce, 0x31, 0x2d, 0x33, 0x96,
	0xa8, 0x76, 0xda, 0xe1, 0x72, 0x5d, 0x79, 0xf3, 0x5a, 0xa9, 0xe3, 0x08, 0x1b, 0x82, 0xac, 0x24,
	0xca, 0x59, 0x7c, 0xa4, 0x3a, 0x69, 0x37, 0x2b, 0x51, 0x61, 0x84, 0x35, 0xdc, 0x38, 0xe7, 0x8b,
	0x05, 0x6c, 0x79, 0x3f, 0x69, 0xa8, 0x98, 0xe4, 0xf9, 0x20, 0x21, 0x82, 0x98, 0xb9, 0x35, 0x0c,
	0x35, 0x85, 0xe4, 0x24, 0x48, 0x9e, 0x3f, 0x27, 0x82, 0xc0, 0x0f, 0xa0, 0x9f, 0x12, 0x69, 0xd3,
	0x2c, 0xa6, 0xea, 0xc4, 0xb9, 0xed, 0x3b, 0xbe, 0xf6, 0xae, 0x1f, 0x11, 0x4e, 0xa7, 0x0d, 0xdc,
	0x65, 0xd9, 0x28, 0xdc, 0x35, 0x0d, 0x34, 0x19, 0xa7, 0x4a, 0xf4, 0xa3, 0xf2, 0xee, 0xa7, 0x99,
	0x78, 0xfb, 0x2e, 0xf2, 0x63, 0x56, 0x04, 0x3a, 0x81, 0xf9, 0xd9, 0xe4, 0xc9, 0x91, 0xd9, 0x22,
	0x99, 0x04, 0xf7, 0x52, 0xc2, 0xf7, 0xa5, 0xaa, 0x79, 0x85, 0x19, 0xd0, 0xdb, 0x63, 0xd1, 0x81,
	0x20, 0x82, 0xff, 0x6b, 0xf3, 0x4d, 0x47, 0xdd, 0xfe, 0xfb, 0xa8, 0x37, 0xc0, 0xac, 0x6c, 0x13,
	0x57, 0xb6, 0xba, 0x31, 0x08, 0x15, 0x46, 0x58, 0xc3, 0xf0, 0x1e, 0xb0, 0x87, 0x94, 0x72, 0xe5,
	0x1c, 0x3b, 0x5c, 0xbc, 0xf6, 0x97, 0x8c, 0x22, 0xac, 0x40, 0xe9, 0x1b, 0xb9, 0xb4, 0x1d, 0xc5,
	0x69, 0xf8, 0x46, 0x2d, 0xac, 0x84, 0xe0, 0x53, 0x00, 0x72, 0xc2, 0xc5, 0x40, 0x0f, 0xbf, 0xab,
	0x88, 0xff, 0xd7, 0x95, 0xb7, 0x3c, 0xa9, 0x6d, 0x82, 0x21, 0xdc, 0x97, 0x7f, 0xc2, 0x9b, 0x2e,
	0x08, 0x5f, 0x9e, 0x5d, 0xb8, 0xd6, 0xf9, 0x85, 0x6b, 0x7d, 0xbf, 0x70, 0xad, 0x4f, 0x97, 0x6e,
	0xeb, 0xfc, 0xd2, 0x6d, 0x7d, 0xbb, 0x74, 0x5b, 0x6f, 0xfc, 0xe6, 0x64, 0xe4, 0x6a, 0xc4, 0x6c,
	0x94, 0x96, 0x94, 0xf3, 0x20, 0x65, 0x9b, 0xfa, 0xc5, 0x1b, 0xeb, 0x37, 0x4f, 0x4d, 0x29, 0xea,
	0xa8, 0x97, 0xea, 0xc9, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x80, 0x5b, 0xcd, 0x0e, 0x05,
	0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.MaxSlots != that1.MaxSlots {
		return false
	}
	if this.MaxGas != that1.MaxGas {
		return false
	}
	if this.FeeTtl != that1.FeeTtl {
		return false
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.FeeTtl != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.FeeTtl))
		i--
		dAtA[i] = 0x18
	}
	if m.MaxGas != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.MaxGas))
		i--
		dAtA[i] = 0x10
	}
	if m.MaxSlots != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.MaxSlots))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Job) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Job) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Job) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Cid) > 0 {
		i -= len(m.Cid)
		copy(dAtA[i:], m.Cid)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Cid)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Label) > 0 {
		i -= len(m.Label)
		copy(dAtA[i:], m.Label)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Label)))
		i--
		dAtA[i] = 0x2a
	}
	{
		size, err := m.Load.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.Trigger.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Trigger) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Trigger) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Trigger) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Block != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Block))
		i--
		dAtA[i] = 0x10
	}
	if m.Period != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Period))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Load) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Load) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Load) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.GasPrice.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.CallData) > 0 {
		i -= len(m.CallData)
		copy(dAtA[i:], m.CallData)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.CallData)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *JobStats) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *JobStats) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *JobStats) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastBlock != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.LastBlock))
		i--
		dAtA[i] = 0x38
	}
	if m.Gas != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Gas))
		i--
		dAtA[i] = 0x30
	}
	if m.Fees != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Fees))
		i--
		dAtA[i] = 0x28
	}
	if m.Calls != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Calls))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Label) > 0 {
		i -= len(m.Label)
		copy(dAtA[i:], m.Label)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Label)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxSlots != 0 {
		n += 1 + sovTypes(uint64(m.MaxSlots))
	}
	if m.MaxGas != 0 {
		n += 1 + sovTypes(uint64(m.MaxGas))
	}
	if m.FeeTtl != 0 {
		n += 1 + sovTypes(uint64(m.FeeTtl))
	}
	return n
}

func (m *Job) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.Trigger.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Load.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = len(m.Label)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Cid)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *Trigger) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Period != 0 {
		n += 1 + sovTypes(uint64(m.Period))
	}
	if m.Block != 0 {
		n += 1 + sovTypes(uint64(m.Block))
	}
	return n
}

func (m *Load) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CallData)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.GasPrice.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *JobStats) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Label)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Calls != 0 {
		n += 1 + sovTypes(uint64(m.Calls))
	}
	if m.Fees != 0 {
		n += 1 + sovTypes(uint64(m.Fees))
	}
	if m.Gas != 0 {
		n += 1 + sovTypes(uint64(m.Gas))
	}
	if m.LastBlock != 0 {
		n += 1 + sovTypes(uint64(m.LastBlock))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxSlots", wireType)
			}
			m.MaxSlots = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxSlots |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxGas", wireType)
			}
			m.MaxGas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxGas |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeTtl", wireType)
			}
			m.FeeTtl = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FeeTtl |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Job) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Job: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Job: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Trigger", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Trigger.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Load", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Load.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Label", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Label = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Cid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Trigger) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Trigger: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Trigger: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Period", wireType)
			}
			m.Period = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Period |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Block", wireType)
			}
			m.Block = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Block |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Load) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Load: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Load: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CallData", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CallData = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasPrice", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GasPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *JobStats) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: JobStats: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: JobStats: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Label", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Label = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Calls", wireType)
			}
			m.Calls = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Calls |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Fees", wireType)
			}
			m.Fees = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Fees |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gas", wireType)
			}
			m.Gas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Gas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastBlock", wireType)
			}
			m.LastBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
