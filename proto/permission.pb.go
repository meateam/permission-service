// Code generated by protoc-gen-go. DO NOT EDIT.
// source: permission.proto

package permission

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Role int32

const (
	Role_NONE  Role = 0
	Role_OWNER Role = 1
	Role_WRITE Role = 2
	Role_READ  Role = 3
)

var Role_name = map[int32]string{
	0: "NONE",
	1: "OWNER",
	2: "WRITE",
	3: "READ",
}

var Role_value = map[string]int32{
	"NONE":  0,
	"OWNER": 1,
	"WRITE": 2,
	"READ":  3,
}

func (x Role) String() string {
	return proto.EnumName(Role_name, int32(x))
}

func (Role) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{0}
}

type CreatePermissionRequest struct {
	// The ID of the file which is being permitted.
	FileID string `protobuf:"bytes,1,opt,name=fileID,proto3" json:"fileID,omitempty"`
	// The ID of the user that's given the permission.
	UserID string `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	// The role of the permission.
	Role                 Role     `protobuf:"varint,3,opt,name=role,proto3,enum=permission.Role" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreatePermissionRequest) Reset()         { *m = CreatePermissionRequest{} }
func (m *CreatePermissionRequest) String() string { return proto.CompactTextString(m) }
func (*CreatePermissionRequest) ProtoMessage()    {}
func (*CreatePermissionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{0}
}

func (m *CreatePermissionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreatePermissionRequest.Unmarshal(m, b)
}
func (m *CreatePermissionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreatePermissionRequest.Marshal(b, m, deterministic)
}
func (m *CreatePermissionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreatePermissionRequest.Merge(m, src)
}
func (m *CreatePermissionRequest) XXX_Size() int {
	return xxx_messageInfo_CreatePermissionRequest.Size(m)
}
func (m *CreatePermissionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreatePermissionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreatePermissionRequest proto.InternalMessageInfo

func (m *CreatePermissionRequest) GetFileID() string {
	if m != nil {
		return m.FileID
	}
	return ""
}

func (m *CreatePermissionRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *CreatePermissionRequest) GetRole() Role {
	if m != nil {
		return m.Role
	}
	return Role_NONE
}

type DeletePermissionRequest struct {
	// The ID of the file which is being permitted.
	FileID string `protobuf:"bytes,1,opt,name=fileID,proto3" json:"fileID,omitempty"`
	// The ID of the user that's given the permission.
	UserID               string   `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeletePermissionRequest) Reset()         { *m = DeletePermissionRequest{} }
func (m *DeletePermissionRequest) String() string { return proto.CompactTextString(m) }
func (*DeletePermissionRequest) ProtoMessage()    {}
func (*DeletePermissionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{1}
}

func (m *DeletePermissionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeletePermissionRequest.Unmarshal(m, b)
}
func (m *DeletePermissionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeletePermissionRequest.Marshal(b, m, deterministic)
}
func (m *DeletePermissionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeletePermissionRequest.Merge(m, src)
}
func (m *DeletePermissionRequest) XXX_Size() int {
	return xxx_messageInfo_DeletePermissionRequest.Size(m)
}
func (m *DeletePermissionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeletePermissionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeletePermissionRequest proto.InternalMessageInfo

func (m *DeletePermissionRequest) GetFileID() string {
	if m != nil {
		return m.FileID
	}
	return ""
}

func (m *DeletePermissionRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type PermissionObject struct {
	// The ID of the permission.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The ID of the file which is being pemitted.
	FileID string `protobuf:"bytes,2,opt,name=fileID,proto3" json:"fileID,omitempty"`
	// The ID of the user that's given the permission.
	UserID string `protobuf:"bytes,3,opt,name=userID,proto3" json:"userID,omitempty"`
	// The role of the permission.
	Role                 Role     `protobuf:"varint,4,opt,name=role,proto3,enum=permission.Role" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PermissionObject) Reset()         { *m = PermissionObject{} }
func (m *PermissionObject) String() string { return proto.CompactTextString(m) }
func (*PermissionObject) ProtoMessage()    {}
func (*PermissionObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{2}
}

func (m *PermissionObject) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PermissionObject.Unmarshal(m, b)
}
func (m *PermissionObject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PermissionObject.Marshal(b, m, deterministic)
}
func (m *PermissionObject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PermissionObject.Merge(m, src)
}
func (m *PermissionObject) XXX_Size() int {
	return xxx_messageInfo_PermissionObject.Size(m)
}
func (m *PermissionObject) XXX_DiscardUnknown() {
	xxx_messageInfo_PermissionObject.DiscardUnknown(m)
}

var xxx_messageInfo_PermissionObject proto.InternalMessageInfo

func (m *PermissionObject) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PermissionObject) GetFileID() string {
	if m != nil {
		return m.FileID
	}
	return ""
}

func (m *PermissionObject) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *PermissionObject) GetRole() Role {
	if m != nil {
		return m.Role
	}
	return Role_NONE
}

type GetFilePermissionsRequest struct {
	// The ID of the file which is being permitted.
	FileID               string   `protobuf:"bytes,1,opt,name=fileID,proto3" json:"fileID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFilePermissionsRequest) Reset()         { *m = GetFilePermissionsRequest{} }
func (m *GetFilePermissionsRequest) String() string { return proto.CompactTextString(m) }
func (*GetFilePermissionsRequest) ProtoMessage()    {}
func (*GetFilePermissionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{3}
}

func (m *GetFilePermissionsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFilePermissionsRequest.Unmarshal(m, b)
}
func (m *GetFilePermissionsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFilePermissionsRequest.Marshal(b, m, deterministic)
}
func (m *GetFilePermissionsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFilePermissionsRequest.Merge(m, src)
}
func (m *GetFilePermissionsRequest) XXX_Size() int {
	return xxx_messageInfo_GetFilePermissionsRequest.Size(m)
}
func (m *GetFilePermissionsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFilePermissionsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetFilePermissionsRequest proto.InternalMessageInfo

func (m *GetFilePermissionsRequest) GetFileID() string {
	if m != nil {
		return m.FileID
	}
	return ""
}

type GetFilePermissionsResponse struct {
	// Array of user roles.
	Permissions          []*GetFilePermissionsResponse_UserRole `protobuf:"bytes,1,rep,name=permissions,proto3" json:"permissions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                               `json:"-"`
	XXX_unrecognized     []byte                                 `json:"-"`
	XXX_sizecache        int32                                  `json:"-"`
}

func (m *GetFilePermissionsResponse) Reset()         { *m = GetFilePermissionsResponse{} }
func (m *GetFilePermissionsResponse) String() string { return proto.CompactTextString(m) }
func (*GetFilePermissionsResponse) ProtoMessage()    {}
func (*GetFilePermissionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{4}
}

func (m *GetFilePermissionsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFilePermissionsResponse.Unmarshal(m, b)
}
func (m *GetFilePermissionsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFilePermissionsResponse.Marshal(b, m, deterministic)
}
func (m *GetFilePermissionsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFilePermissionsResponse.Merge(m, src)
}
func (m *GetFilePermissionsResponse) XXX_Size() int {
	return xxx_messageInfo_GetFilePermissionsResponse.Size(m)
}
func (m *GetFilePermissionsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFilePermissionsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetFilePermissionsResponse proto.InternalMessageInfo

func (m *GetFilePermissionsResponse) GetPermissions() []*GetFilePermissionsResponse_UserRole {
	if m != nil {
		return m.Permissions
	}
	return nil
}

// The role of a user.
type GetFilePermissionsResponse_UserRole struct {
	// The user ID.
	UserID string `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	// The role of the user.
	Role                 Role     `protobuf:"varint,2,opt,name=role,proto3,enum=permission.Role" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFilePermissionsResponse_UserRole) Reset()         { *m = GetFilePermissionsResponse_UserRole{} }
func (m *GetFilePermissionsResponse_UserRole) String() string { return proto.CompactTextString(m) }
func (*GetFilePermissionsResponse_UserRole) ProtoMessage()    {}
func (*GetFilePermissionsResponse_UserRole) Descriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{4, 0}
}

func (m *GetFilePermissionsResponse_UserRole) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFilePermissionsResponse_UserRole.Unmarshal(m, b)
}
func (m *GetFilePermissionsResponse_UserRole) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFilePermissionsResponse_UserRole.Marshal(b, m, deterministic)
}
func (m *GetFilePermissionsResponse_UserRole) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFilePermissionsResponse_UserRole.Merge(m, src)
}
func (m *GetFilePermissionsResponse_UserRole) XXX_Size() int {
	return xxx_messageInfo_GetFilePermissionsResponse_UserRole.Size(m)
}
func (m *GetFilePermissionsResponse_UserRole) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFilePermissionsResponse_UserRole.DiscardUnknown(m)
}

var xxx_messageInfo_GetFilePermissionsResponse_UserRole proto.InternalMessageInfo

func (m *GetFilePermissionsResponse_UserRole) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *GetFilePermissionsResponse_UserRole) GetRole() Role {
	if m != nil {
		return m.Role
	}
	return Role_NONE
}

type IsPermittedRequest struct {
	// The ID of the file which is being permitted.
	FileID string `protobuf:"bytes,1,opt,name=fileID,proto3" json:"fileID,omitempty"`
	// The ID of the user that's given the permission.
	UserID string `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	// The role of the permission.
	Role                 Role     `protobuf:"varint,3,opt,name=role,proto3,enum=permission.Role" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IsPermittedRequest) Reset()         { *m = IsPermittedRequest{} }
func (m *IsPermittedRequest) String() string { return proto.CompactTextString(m) }
func (*IsPermittedRequest) ProtoMessage()    {}
func (*IsPermittedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{5}
}

func (m *IsPermittedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IsPermittedRequest.Unmarshal(m, b)
}
func (m *IsPermittedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IsPermittedRequest.Marshal(b, m, deterministic)
}
func (m *IsPermittedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IsPermittedRequest.Merge(m, src)
}
func (m *IsPermittedRequest) XXX_Size() int {
	return xxx_messageInfo_IsPermittedRequest.Size(m)
}
func (m *IsPermittedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IsPermittedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IsPermittedRequest proto.InternalMessageInfo

func (m *IsPermittedRequest) GetFileID() string {
	if m != nil {
		return m.FileID
	}
	return ""
}

func (m *IsPermittedRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *IsPermittedRequest) GetRole() Role {
	if m != nil {
		return m.Role
	}
	return Role_NONE
}

type IsPermittedResponse struct {
	Permitted            bool     `protobuf:"varint,1,opt,name=permitted,proto3" json:"permitted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IsPermittedResponse) Reset()         { *m = IsPermittedResponse{} }
func (m *IsPermittedResponse) String() string { return proto.CompactTextString(m) }
func (*IsPermittedResponse) ProtoMessage()    {}
func (*IsPermittedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c837ef01cbda0ad8, []int{6}
}

func (m *IsPermittedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IsPermittedResponse.Unmarshal(m, b)
}
func (m *IsPermittedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IsPermittedResponse.Marshal(b, m, deterministic)
}
func (m *IsPermittedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IsPermittedResponse.Merge(m, src)
}
func (m *IsPermittedResponse) XXX_Size() int {
	return xxx_messageInfo_IsPermittedResponse.Size(m)
}
func (m *IsPermittedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IsPermittedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IsPermittedResponse proto.InternalMessageInfo

func (m *IsPermittedResponse) GetPermitted() bool {
	if m != nil {
		return m.Permitted
	}
	return false
}

func init() {
	proto.RegisterEnum("permission.Role", Role_name, Role_value)
	proto.RegisterType((*CreatePermissionRequest)(nil), "permission.CreatePermissionRequest")
	proto.RegisterType((*DeletePermissionRequest)(nil), "permission.DeletePermissionRequest")
	proto.RegisterType((*PermissionObject)(nil), "permission.PermissionObject")
	proto.RegisterType((*GetFilePermissionsRequest)(nil), "permission.GetFilePermissionsRequest")
	proto.RegisterType((*GetFilePermissionsResponse)(nil), "permission.GetFilePermissionsResponse")
	proto.RegisterType((*GetFilePermissionsResponse_UserRole)(nil), "permission.GetFilePermissionsResponse.UserRole")
	proto.RegisterType((*IsPermittedRequest)(nil), "permission.IsPermittedRequest")
	proto.RegisterType((*IsPermittedResponse)(nil), "permission.IsPermittedResponse")
}

func init() { proto.RegisterFile("permission.proto", fileDescriptor_c837ef01cbda0ad8) }

var fileDescriptor_c837ef01cbda0ad8 = []byte{
	// 387 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x4d, 0x4f, 0xc2, 0x40,
	0x10, 0xa5, 0x1f, 0x12, 0x18, 0x12, 0xb2, 0x59, 0x13, 0xa9, 0x0d, 0x51, 0x52, 0x3f, 0x42, 0x3c,
	0xa0, 0x81, 0x5f, 0x60, 0x6c, 0xd5, 0x5e, 0x00, 0x37, 0x1a, 0xce, 0x22, 0x63, 0x52, 0x52, 0xd9,
	0xda, 0x5d, 0x12, 0xff, 0x9a, 0xbf, 0xcc, 0xab, 0x69, 0x01, 0xbb, 0x7c, 0x14, 0x3c, 0xe8, 0xad,
	0xfb, 0xa6, 0xf3, 0xde, 0xbc, 0xed, 0x9b, 0x02, 0x89, 0x30, 0x7e, 0x0b, 0x84, 0x08, 0xf8, 0xa4,
	0x15, 0xc5, 0x5c, 0x72, 0x0a, 0x19, 0xe2, 0x70, 0xa8, 0xdd, 0xc4, 0xf8, 0x2c, 0xb1, 0xff, 0x83,
	0x31, 0x7c, 0x9f, 0xa2, 0x90, 0xf4, 0x00, 0x8a, 0xaf, 0x41, 0x88, 0xbe, 0x6b, 0x69, 0x0d, 0xad,
	0x59, 0x66, 0xf3, 0x53, 0x82, 0x4f, 0x05, 0xc6, 0xbe, 0x6b, 0xe9, 0x33, 0x7c, 0x76, 0xa2, 0xa7,
	0x60, 0xc6, 0x3c, 0x44, 0xcb, 0x68, 0x68, 0xcd, 0x6a, 0x9b, 0xb4, 0x14, 0x5d, 0xc6, 0x43, 0x64,
	0x69, 0xd5, 0xf1, 0xa1, 0xe6, 0x62, 0x88, 0x7f, 0x20, 0xe8, 0x7c, 0x00, 0xc9, 0x48, 0x7a, 0xc3,
	0x31, 0xbe, 0x48, 0x5a, 0x05, 0x3d, 0x18, 0xcd, 0xfb, 0xf5, 0x60, 0xa4, 0x70, 0xea, 0x39, 0x9c,
	0xc6, 0x46, 0x13, 0xe6, 0x56, 0x13, 0x1d, 0x38, 0xbc, 0x43, 0x79, 0x1b, 0x84, 0x8a, 0x0b, 0xb1,
	0xc3, 0x86, 0xf3, 0xa9, 0x81, 0xbd, 0xa9, 0x4b, 0x44, 0x7c, 0x22, 0x90, 0x3e, 0x40, 0x25, 0x13,
	0x13, 0x96, 0xd6, 0x30, 0x9a, 0x95, 0xf6, 0xa5, 0x3a, 0x40, 0x7e, 0x73, 0xeb, 0x49, 0x60, 0x9c,
	0xce, 0xa7, 0x72, 0xd8, 0xf7, 0x50, 0x5a, 0x14, 0x14, 0xc3, 0xda, 0x46, 0xc3, 0xfa, 0x56, 0xc3,
	0x63, 0xa0, 0xbe, 0x48, 0x85, 0xa5, 0xc4, 0xd1, 0xff, 0x26, 0xa4, 0x03, 0xfb, 0x4b, 0x5a, 0xf3,
	0xfb, 0xa9, 0x43, 0x39, 0x5a, 0x80, 0xa9, 0x5e, 0x89, 0x65, 0xc0, 0xc5, 0x15, 0x98, 0xa9, 0xcd,
	0x12, 0x98, 0xdd, 0x5e, 0xd7, 0x23, 0x05, 0x5a, 0x86, 0xbd, 0xde, 0xa0, 0xeb, 0x31, 0xa2, 0x25,
	0x8f, 0x03, 0xe6, 0x3f, 0x7a, 0x44, 0x4f, 0xea, 0xcc, 0xbb, 0x76, 0x89, 0xd1, 0xfe, 0xd2, 0x01,
	0xb2, 0xab, 0xa4, 0x03, 0x20, 0xab, 0x8b, 0x40, 0x4f, 0xd4, 0x09, 0x73, 0xd6, 0xc4, 0xae, 0xab,
	0x2f, 0xad, 0xe6, 0xd1, 0x29, 0x24, 0xc4, 0xab, 0x81, 0x5f, 0x26, 0xce, 0x59, 0x87, 0x9d, 0xc4,
	0x08, 0x74, 0x3d, 0x11, 0xf4, 0x6c, 0x57, 0x62, 0x66, 0xe4, 0xe7, 0xbf, 0x0b, 0x96, 0x53, 0xa0,
	0x7d, 0xa8, 0x28, 0x9f, 0x83, 0x1e, 0xa9, 0x8d, 0xeb, 0x99, 0xb0, 0x8f, 0x73, 0xeb, 0x0b, 0xc6,
	0x61, 0x31, 0xfd, 0x0d, 0x75, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0xd6, 0xf4, 0x6e, 0x23, 0x9a,
	0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PermissionClient is the client API for Permission service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PermissionClient interface {
	CreatePermission(ctx context.Context, in *CreatePermissionRequest, opts ...grpc.CallOption) (*PermissionObject, error)
	DeletePermission(ctx context.Context, in *DeletePermissionRequest, opts ...grpc.CallOption) (*PermissionObject, error)
	GetFilePermissions(ctx context.Context, in *GetFilePermissionsRequest, opts ...grpc.CallOption) (*GetFilePermissionsResponse, error)
	IsPermitted(ctx context.Context, in *IsPermittedRequest, opts ...grpc.CallOption) (*IsPermittedResponse, error)
}

type permissionClient struct {
	cc *grpc.ClientConn
}

func NewPermissionClient(cc *grpc.ClientConn) PermissionClient {
	return &permissionClient{cc}
}

func (c *permissionClient) CreatePermission(ctx context.Context, in *CreatePermissionRequest, opts ...grpc.CallOption) (*PermissionObject, error) {
	out := new(PermissionObject)
	err := c.cc.Invoke(ctx, "/permission.Permission/CreatePermission", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) DeletePermission(ctx context.Context, in *DeletePermissionRequest, opts ...grpc.CallOption) (*PermissionObject, error) {
	out := new(PermissionObject)
	err := c.cc.Invoke(ctx, "/permission.Permission/DeletePermission", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) GetFilePermissions(ctx context.Context, in *GetFilePermissionsRequest, opts ...grpc.CallOption) (*GetFilePermissionsResponse, error) {
	out := new(GetFilePermissionsResponse)
	err := c.cc.Invoke(ctx, "/permission.Permission/GetFilePermissions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) IsPermitted(ctx context.Context, in *IsPermittedRequest, opts ...grpc.CallOption) (*IsPermittedResponse, error) {
	out := new(IsPermittedResponse)
	err := c.cc.Invoke(ctx, "/permission.Permission/IsPermitted", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PermissionServer is the server API for Permission service.
type PermissionServer interface {
	CreatePermission(context.Context, *CreatePermissionRequest) (*PermissionObject, error)
	DeletePermission(context.Context, *DeletePermissionRequest) (*PermissionObject, error)
	GetFilePermissions(context.Context, *GetFilePermissionsRequest) (*GetFilePermissionsResponse, error)
	IsPermitted(context.Context, *IsPermittedRequest) (*IsPermittedResponse, error)
}

// UnimplementedPermissionServer can be embedded to have forward compatible implementations.
type UnimplementedPermissionServer struct {
}

func (*UnimplementedPermissionServer) CreatePermission(ctx context.Context, req *CreatePermissionRequest) (*PermissionObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePermission not implemented")
}
func (*UnimplementedPermissionServer) DeletePermission(ctx context.Context, req *DeletePermissionRequest) (*PermissionObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePermission not implemented")
}
func (*UnimplementedPermissionServer) GetFilePermissions(ctx context.Context, req *GetFilePermissionsRequest) (*GetFilePermissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilePermissions not implemented")
}
func (*UnimplementedPermissionServer) IsPermitted(ctx context.Context, req *IsPermittedRequest) (*IsPermittedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsPermitted not implemented")
}

func RegisterPermissionServer(s *grpc.Server, srv PermissionServer) {
	s.RegisterService(&_Permission_serviceDesc, srv)
}

func _Permission_CreatePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).CreatePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/permission.Permission/CreatePermission",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).CreatePermission(ctx, req.(*CreatePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_DeletePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).DeletePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/permission.Permission/DeletePermission",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).DeletePermission(ctx, req.(*DeletePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_GetFilePermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFilePermissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).GetFilePermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/permission.Permission/GetFilePermissions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).GetFilePermissions(ctx, req.(*GetFilePermissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_IsPermitted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsPermittedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).IsPermitted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/permission.Permission/IsPermitted",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).IsPermitted(ctx, req.(*IsPermittedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Permission_serviceDesc = grpc.ServiceDesc{
	ServiceName: "permission.Permission",
	HandlerType: (*PermissionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePermission",
			Handler:    _Permission_CreatePermission_Handler,
		},
		{
			MethodName: "DeletePermission",
			Handler:    _Permission_DeletePermission_Handler,
		},
		{
			MethodName: "GetFilePermissions",
			Handler:    _Permission_GetFilePermissions_Handler,
		},
		{
			MethodName: "IsPermitted",
			Handler:    _Permission_IsPermitted_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "permission.proto",
}
