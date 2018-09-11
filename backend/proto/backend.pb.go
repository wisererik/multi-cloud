// Code generated by protoc-gen-go. DO NOT EDIT.
// source: backend.proto

package backend

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CreateBackendRequest struct {
	Backend              *BackendDetail `protobuf:"bytes,1,opt,name=backend,proto3" json:"backend,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CreateBackendRequest) Reset()         { *m = CreateBackendRequest{} }
func (m *CreateBackendRequest) String() string { return proto.CompactTextString(m) }
func (*CreateBackendRequest) ProtoMessage()    {}
func (*CreateBackendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{0}
}
func (m *CreateBackendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateBackendRequest.Unmarshal(m, b)
}
func (m *CreateBackendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateBackendRequest.Marshal(b, m, deterministic)
}
func (dst *CreateBackendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateBackendRequest.Merge(dst, src)
}
func (m *CreateBackendRequest) XXX_Size() int {
	return xxx_messageInfo_CreateBackendRequest.Size(m)
}
func (m *CreateBackendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateBackendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateBackendRequest proto.InternalMessageInfo

func (m *CreateBackendRequest) GetBackend() *BackendDetail {
	if m != nil {
		return m.Backend
	}
	return nil
}

type CreateBackendResponse struct {
	Backend              *BackendDetail `protobuf:"bytes,1,opt,name=backend,proto3" json:"backend,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CreateBackendResponse) Reset()         { *m = CreateBackendResponse{} }
func (m *CreateBackendResponse) String() string { return proto.CompactTextString(m) }
func (*CreateBackendResponse) ProtoMessage()    {}
func (*CreateBackendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{1}
}
func (m *CreateBackendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateBackendResponse.Unmarshal(m, b)
}
func (m *CreateBackendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateBackendResponse.Marshal(b, m, deterministic)
}
func (dst *CreateBackendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateBackendResponse.Merge(dst, src)
}
func (m *CreateBackendResponse) XXX_Size() int {
	return xxx_messageInfo_CreateBackendResponse.Size(m)
}
func (m *CreateBackendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateBackendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateBackendResponse proto.InternalMessageInfo

func (m *CreateBackendResponse) GetBackend() *BackendDetail {
	if m != nil {
		return m.Backend
	}
	return nil
}

type GetBackendRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBackendRequest) Reset()         { *m = GetBackendRequest{} }
func (m *GetBackendRequest) String() string { return proto.CompactTextString(m) }
func (*GetBackendRequest) ProtoMessage()    {}
func (*GetBackendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{2}
}
func (m *GetBackendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetBackendRequest.Unmarshal(m, b)
}
func (m *GetBackendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetBackendRequest.Marshal(b, m, deterministic)
}
func (dst *GetBackendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetBackendRequest.Merge(dst, src)
}
func (m *GetBackendRequest) XXX_Size() int {
	return xxx_messageInfo_GetBackendRequest.Size(m)
}
func (m *GetBackendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetBackendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetBackendRequest proto.InternalMessageInfo

func (m *GetBackendRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetBackendResponse struct {
	Backend              *BackendDetail `protobuf:"bytes,1,opt,name=backend,proto3" json:"backend,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetBackendResponse) Reset()         { *m = GetBackendResponse{} }
func (m *GetBackendResponse) String() string { return proto.CompactTextString(m) }
func (*GetBackendResponse) ProtoMessage()    {}
func (*GetBackendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{3}
}
func (m *GetBackendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetBackendResponse.Unmarshal(m, b)
}
func (m *GetBackendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetBackendResponse.Marshal(b, m, deterministic)
}
func (dst *GetBackendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetBackendResponse.Merge(dst, src)
}
func (m *GetBackendResponse) XXX_Size() int {
	return xxx_messageInfo_GetBackendResponse.Size(m)
}
func (m *GetBackendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetBackendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetBackendResponse proto.InternalMessageInfo

func (m *GetBackendResponse) GetBackend() *BackendDetail {
	if m != nil {
		return m.Backend
	}
	return nil
}

type ListBackendRequest struct {
	Limit                int32    `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset               string   `protobuf:"bytes,2,opt,name=offset,proto3" json:"offset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListBackendRequest) Reset()         { *m = ListBackendRequest{} }
func (m *ListBackendRequest) String() string { return proto.CompactTextString(m) }
func (*ListBackendRequest) ProtoMessage()    {}
func (*ListBackendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{4}
}
func (m *ListBackendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListBackendRequest.Unmarshal(m, b)
}
func (m *ListBackendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListBackendRequest.Marshal(b, m, deterministic)
}
func (dst *ListBackendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListBackendRequest.Merge(dst, src)
}
func (m *ListBackendRequest) XXX_Size() int {
	return xxx_messageInfo_ListBackendRequest.Size(m)
}
func (m *ListBackendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListBackendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListBackendRequest proto.InternalMessageInfo

func (m *ListBackendRequest) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *ListBackendRequest) GetOffset() string {
	if m != nil {
		return m.Offset
	}
	return ""
}

type ListBackendResponse struct {
	Entries              []*ListBackendResponse_Entry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
	Next                 string                       `protobuf:"bytes,2,opt,name=next,proto3" json:"next,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *ListBackendResponse) Reset()         { *m = ListBackendResponse{} }
func (m *ListBackendResponse) String() string { return proto.CompactTextString(m) }
func (*ListBackendResponse) ProtoMessage()    {}
func (*ListBackendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{5}
}
func (m *ListBackendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListBackendResponse.Unmarshal(m, b)
}
func (m *ListBackendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListBackendResponse.Marshal(b, m, deterministic)
}
func (dst *ListBackendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListBackendResponse.Merge(dst, src)
}
func (m *ListBackendResponse) XXX_Size() int {
	return xxx_messageInfo_ListBackendResponse.Size(m)
}
func (m *ListBackendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListBackendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListBackendResponse proto.InternalMessageInfo

func (m *ListBackendResponse) GetEntries() []*ListBackendResponse_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *ListBackendResponse) GetNext() string {
	if m != nil {
		return m.Next
	}
	return ""
}

type ListBackendResponse_Entry struct {
	Backend              *BackendDetail `protobuf:"bytes,1,opt,name=backend,proto3" json:"backend,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ListBackendResponse_Entry) Reset()         { *m = ListBackendResponse_Entry{} }
func (m *ListBackendResponse_Entry) String() string { return proto.CompactTextString(m) }
func (*ListBackendResponse_Entry) ProtoMessage()    {}
func (*ListBackendResponse_Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{5, 0}
}
func (m *ListBackendResponse_Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListBackendResponse_Entry.Unmarshal(m, b)
}
func (m *ListBackendResponse_Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListBackendResponse_Entry.Marshal(b, m, deterministic)
}
func (dst *ListBackendResponse_Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListBackendResponse_Entry.Merge(dst, src)
}
func (m *ListBackendResponse_Entry) XXX_Size() int {
	return xxx_messageInfo_ListBackendResponse_Entry.Size(m)
}
func (m *ListBackendResponse_Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_ListBackendResponse_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_ListBackendResponse_Entry proto.InternalMessageInfo

func (m *ListBackendResponse_Entry) GetBackend() *BackendDetail {
	if m != nil {
		return m.Backend
	}
	return nil
}

type UpdateBackendRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Access               string   `protobuf:"bytes,2,opt,name=access,proto3" json:"access,omitempty"`
	Security             string   `protobuf:"bytes,3,opt,name=security,proto3" json:"security,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateBackendRequest) Reset()         { *m = UpdateBackendRequest{} }
func (m *UpdateBackendRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateBackendRequest) ProtoMessage()    {}
func (*UpdateBackendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{6}
}
func (m *UpdateBackendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateBackendRequest.Unmarshal(m, b)
}
func (m *UpdateBackendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateBackendRequest.Marshal(b, m, deterministic)
}
func (dst *UpdateBackendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateBackendRequest.Merge(dst, src)
}
func (m *UpdateBackendRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateBackendRequest.Size(m)
}
func (m *UpdateBackendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateBackendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateBackendRequest proto.InternalMessageInfo

func (m *UpdateBackendRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UpdateBackendRequest) GetAccess() string {
	if m != nil {
		return m.Access
	}
	return ""
}

func (m *UpdateBackendRequest) GetSecurity() string {
	if m != nil {
		return m.Security
	}
	return ""
}

type UpdateBackendResponse struct {
	Backend              *BackendDetail `protobuf:"bytes,1,opt,name=backend,proto3" json:"backend,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *UpdateBackendResponse) Reset()         { *m = UpdateBackendResponse{} }
func (m *UpdateBackendResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateBackendResponse) ProtoMessage()    {}
func (*UpdateBackendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{7}
}
func (m *UpdateBackendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateBackendResponse.Unmarshal(m, b)
}
func (m *UpdateBackendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateBackendResponse.Marshal(b, m, deterministic)
}
func (dst *UpdateBackendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateBackendResponse.Merge(dst, src)
}
func (m *UpdateBackendResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateBackendResponse.Size(m)
}
func (m *UpdateBackendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateBackendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateBackendResponse proto.InternalMessageInfo

func (m *UpdateBackendResponse) GetBackend() *BackendDetail {
	if m != nil {
		return m.Backend
	}
	return nil
}

type DeleteBackendRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteBackendRequest) Reset()         { *m = DeleteBackendRequest{} }
func (m *DeleteBackendRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteBackendRequest) ProtoMessage()    {}
func (*DeleteBackendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{8}
}
func (m *DeleteBackendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteBackendRequest.Unmarshal(m, b)
}
func (m *DeleteBackendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteBackendRequest.Marshal(b, m, deterministic)
}
func (dst *DeleteBackendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteBackendRequest.Merge(dst, src)
}
func (m *DeleteBackendRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteBackendRequest.Size(m)
}
func (m *DeleteBackendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteBackendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteBackendRequest proto.InternalMessageInfo

func (m *DeleteBackendRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DeleteBackendResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteBackendResponse) Reset()         { *m = DeleteBackendResponse{} }
func (m *DeleteBackendResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteBackendResponse) ProtoMessage()    {}
func (*DeleteBackendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{9}
}
func (m *DeleteBackendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteBackendResponse.Unmarshal(m, b)
}
func (m *DeleteBackendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteBackendResponse.Marshal(b, m, deterministic)
}
func (dst *DeleteBackendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteBackendResponse.Merge(dst, src)
}
func (m *DeleteBackendResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteBackendResponse.Size(m)
}
func (m *DeleteBackendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteBackendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteBackendResponse proto.InternalMessageInfo

type BackendDetail struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TenantId             string   `protobuf:"bytes,2,opt,name=tenantId,proto3" json:"tenantId,omitempty"`
	UserId               string   `protobuf:"bytes,3,opt,name=userId,proto3" json:"userId,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Type                 string   `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	Endpoint             string   `protobuf:"bytes,6,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	BucketName           string   `protobuf:"bytes,7,opt,name=bucketName,proto3" json:"bucketName,omitempty"`
	Access               string   `protobuf:"bytes,8,opt,name=access,proto3" json:"access,omitempty"`
	Security             string   `protobuf:"bytes,9,opt,name=security,proto3" json:"security,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BackendDetail) Reset()         { *m = BackendDetail{} }
func (m *BackendDetail) String() string { return proto.CompactTextString(m) }
func (*BackendDetail) ProtoMessage()    {}
func (*BackendDetail) Descriptor() ([]byte, []int) {
	return fileDescriptor_backend_ac18f7211c31e45d, []int{10}
}
func (m *BackendDetail) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackendDetail.Unmarshal(m, b)
}
func (m *BackendDetail) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackendDetail.Marshal(b, m, deterministic)
}
func (dst *BackendDetail) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackendDetail.Merge(dst, src)
}
func (m *BackendDetail) XXX_Size() int {
	return xxx_messageInfo_BackendDetail.Size(m)
}
func (m *BackendDetail) XXX_DiscardUnknown() {
	xxx_messageInfo_BackendDetail.DiscardUnknown(m)
}

var xxx_messageInfo_BackendDetail proto.InternalMessageInfo

func (m *BackendDetail) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *BackendDetail) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *BackendDetail) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *BackendDetail) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *BackendDetail) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *BackendDetail) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *BackendDetail) GetBucketName() string {
	if m != nil {
		return m.BucketName
	}
	return ""
}

func (m *BackendDetail) GetAccess() string {
	if m != nil {
		return m.Access
	}
	return ""
}

func (m *BackendDetail) GetSecurity() string {
	if m != nil {
		return m.Security
	}
	return ""
}

func init() {
	proto.RegisterType((*CreateBackendRequest)(nil), "CreateBackendRequest")
	proto.RegisterType((*CreateBackendResponse)(nil), "CreateBackendResponse")
	proto.RegisterType((*GetBackendRequest)(nil), "GetBackendRequest")
	proto.RegisterType((*GetBackendResponse)(nil), "GetBackendResponse")
	proto.RegisterType((*ListBackendRequest)(nil), "ListBackendRequest")
	proto.RegisterType((*ListBackendResponse)(nil), "ListBackendResponse")
	proto.RegisterType((*ListBackendResponse_Entry)(nil), "ListBackendResponse.Entry")
	proto.RegisterType((*UpdateBackendRequest)(nil), "UpdateBackendRequest")
	proto.RegisterType((*UpdateBackendResponse)(nil), "UpdateBackendResponse")
	proto.RegisterType((*DeleteBackendRequest)(nil), "DeleteBackendRequest")
	proto.RegisterType((*DeleteBackendResponse)(nil), "DeleteBackendResponse")
	proto.RegisterType((*BackendDetail)(nil), "BackendDetail")
}

func init() { proto.RegisterFile("backend.proto", fileDescriptor_backend_ac18f7211c31e45d) }

var fileDescriptor_backend_ac18f7211c31e45d = []byte{
	// 469 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x5f, 0x8b, 0xd3, 0x40,
	0x10, 0x6f, 0x7a, 0xd7, 0x3f, 0x37, 0xa5, 0x07, 0x4e, 0x93, 0x1a, 0xf6, 0x41, 0x8e, 0x15, 0xa4,
	0x4f, 0x0b, 0x56, 0x41, 0xf0, 0x41, 0xce, 0xf3, 0x44, 0x0e, 0xc4, 0x87, 0x80, 0x2f, 0xbe, 0xa5,
	0xc9, 0x1c, 0x2c, 0xd7, 0x6e, 0x62, 0x76, 0x0b, 0xf6, 0x6b, 0xf8, 0x99, 0xfc, 0x36, 0x7e, 0x09,
	0x49, 0xb2, 0x39, 0x9b, 0x64, 0x45, 0xfa, 0x96, 0x99, 0x9d, 0xfd, 0xcd, 0xfc, 0x66, 0x7f, 0xbf,
	0xc0, 0x7c, 0x13, 0x27, 0x0f, 0xa4, 0x52, 0x91, 0x17, 0x99, 0xc9, 0xf8, 0x35, 0xf8, 0x1f, 0x0a,
	0x8a, 0x0d, 0xdd, 0xd4, 0xe9, 0x88, 0xbe, 0xef, 0x49, 0x1b, 0x5c, 0xc1, 0xc4, 0x16, 0x86, 0xde,
	0x95, 0xb7, 0x9a, 0xad, 0x2f, 0x85, 0xad, 0xb8, 0x25, 0x13, 0xcb, 0x6d, 0xd4, 0x1c, 0xf3, 0xf7,
	0x10, 0x74, 0x10, 0x74, 0x9e, 0x29, 0x4d, 0x27, 0x40, 0x3c, 0x87, 0x27, 0x9f, 0xc8, 0x74, 0x26,
	0xb8, 0x84, 0xa1, 0xac, 0x6f, 0x5e, 0x44, 0x43, 0x99, 0xf2, 0x77, 0x80, 0xc7, 0x45, 0x27, 0x37,
	0xb9, 0x01, 0xfc, 0x2c, 0x75, 0xb7, 0x8b, 0x0f, 0xa3, 0xad, 0xdc, 0x49, 0x53, 0xdd, 0x1e, 0x45,
	0x75, 0x80, 0x4b, 0x18, 0x67, 0xf7, 0xf7, 0x9a, 0x4c, 0x38, 0xac, 0xfa, 0xdb, 0x88, 0xff, 0xf4,
	0x60, 0xd1, 0x02, 0xb1, 0x53, 0xbc, 0x86, 0x09, 0x29, 0x53, 0x48, 0xd2, 0xa1, 0x77, 0x75, 0xb6,
	0x9a, 0xad, 0x99, 0x70, 0x94, 0x89, 0x8f, 0xca, 0x14, 0x87, 0xa8, 0x29, 0x45, 0x84, 0x73, 0x45,
	0x3f, 0x9a, 0x1e, 0xd5, 0x37, 0x7b, 0x09, 0xa3, 0xaa, 0xea, 0x04, 0x62, 0xdf, 0xc0, 0xff, 0x9a,
	0xa7, 0xfd, 0x27, 0xec, 0x2c, 0xb0, 0x24, 0x15, 0x27, 0x09, 0x69, 0xdd, 0x90, 0xaa, 0x23, 0x64,
	0x30, 0xd5, 0x94, 0xec, 0x0b, 0x69, 0x0e, 0xe1, 0x59, 0x75, 0xf2, 0x18, 0x97, 0x8f, 0xdb, 0xc1,
	0x3e, 0x79, 0xef, 0x2f, 0xc0, 0xbf, 0xa5, 0x2d, 0xfd, 0x6f, 0x3c, 0xfe, 0x14, 0x82, 0x4e, 0x5d,
	0xdd, 0x8a, 0xff, 0xf6, 0x60, 0xde, 0xc2, 0xee, 0x31, 0x63, 0x30, 0x35, 0xa4, 0x62, 0x65, 0xee,
	0x52, 0xcb, 0xed, 0x31, 0x2e, 0x59, 0xef, 0x35, 0x15, 0x77, 0xa9, 0xe5, 0x66, 0xa3, 0x6a, 0xf9,
	0xf1, 0x8e, 0xc2, 0x73, 0xbb, 0xfc, 0x78, 0x47, 0x65, 0xce, 0x1c, 0x72, 0x0a, 0x47, 0x75, 0xae,
	0xfc, 0x2e, 0xb1, 0x49, 0xa5, 0x79, 0x26, 0x95, 0x09, 0xc7, 0x35, 0x76, 0x13, 0xe3, 0x33, 0x80,
	0xcd, 0x3e, 0x79, 0x20, 0xf3, 0xa5, 0x44, 0x9a, 0x54, 0xa7, 0x47, 0x99, 0xa3, 0x8d, 0x4f, 0xff,
	0xb9, 0xf1, 0x8b, 0xf6, 0xc6, 0xd7, 0xbf, 0x86, 0x30, 0xb1, 0x6c, 0xf1, 0x1a, 0xe6, 0x2d, 0x6b,
	0x61, 0x20, 0x5c, 0x66, 0x65, 0x4b, 0xe1, 0x74, 0x20, 0x1f, 0xe0, 0x1b, 0x80, 0xbf, 0xa6, 0x41,
	0x14, 0x3d, 0x9b, 0xb1, 0x85, 0xe8, 0xbb, 0x8a, 0x0f, 0xf0, 0x2d, 0xcc, 0x8e, 0x14, 0x8c, 0x0b,
	0xd1, 0xf7, 0x0e, 0xf3, 0x5d, 0x22, 0xe7, 0x83, 0x72, 0xec, 0x96, 0x68, 0x30, 0x10, 0x2e, 0x81,
	0xb2, 0xa5, 0x70, 0x6a, 0xab, 0x46, 0x68, 0x69, 0x01, 0x03, 0xe1, 0xd2, 0x10, 0x5b, 0x0a, 0xb7,
	0x64, 0x06, 0x9b, 0x71, 0xf5, 0x7b, 0x7b, 0xf5, 0x27, 0x00, 0x00, 0xff, 0xff, 0xa9, 0xfb, 0xe6,
	0x69, 0xef, 0x04, 0x00, 0x00,
}