// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: report.proto

package report

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StatusReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TowerIP                       string                 `protobuf:"bytes,1,opt,name=towerIP,proto3" json:"towerIP,omitempty"`                                             // measure.Strings["host_name"]
	LastArduinoReachableTimestamp string                 `protobuf:"bytes,2,opt,name=lastArduinoReachableTimestamp,proto3" json:"lastArduinoReachableTimestamp,omitempty"` // time.Now().Add(measure.Strings["arduinoReached"]).String()
	LastTowerReachableTimestamp   string                 `protobuf:"bytes,3,opt,name=lastTowerReachableTimestamp,proto3" json:"lastTowerReachableTimestamp,omitempty"`     // testbed-monitor/report/receiver.go/receivedReports[ip]
	BootTimestamp                 string                 `protobuf:"bytes,4,opt,name=bootTimestamp,proto3" json:"bootTimestamp,omitempty"`                                 // measure.Strings["bootTime"]
	RebootsCurrentDay             int64                  `protobuf:"varint,5,opt,name=rebootsCurrentDay,proto3" json:"rebootsCurrentDay,omitempty"`                        // measure.Integers["Reboots_Today"]
	RAMUsed                       int64                  `protobuf:"varint,6,opt,name=RAMUsed,proto3" json:"RAMUsed,omitempty"`                                            // measure.Integers["vm_used"]
	RAMTotal                      int64                  `protobuf:"varint,7,opt,name=RAMTotal,proto3" json:"RAMTotal,omitempty"`                                          // measure.Integers["vm_total"]
	DiskUsed                      int64                  `protobuf:"varint,8,opt,name=diskUsed,proto3" json:"diskUsed,omitempty"`                                          // measure.Integers["DISK_USED"]
	DiskTotal                     int64                  `protobuf:"varint,9,opt,name=diskTotal,proto3" json:"diskTotal,omitempty"`                                        // measure.Integers["DISK_TOTAL"]
	CPUAvg                        int64                  `protobuf:"varint,10,opt,name=cpu,proto3" json:"cpu,omitempty"`                                             // measure.Integers["CPU_AVG"]
	Timestamp                     *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                                        // the time when this measure was collected
	Reachable                     bool                   `protobuf:"varint,12,opt,name=reachable,proto3" json:"reachable,omitempty"`                                       // has this tower been reached within the expected time?
	Temperature                   int64                  `protobuf:"varint,13,opt,name=temperature,proto3" json:"temperature,omitempty"`                                   // temperature from arduino controller mqtt report
}

func (x *StatusReport) Reset() {
	*x = StatusReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_report_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusReport) ProtoMessage() {}

func (x *StatusReport) ProtoReflect() protoreflect.Message {
	mi := &file_report_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusReport.ProtoReflect.Descriptor instead.
func (*StatusReport) Descriptor() ([]byte, []int) {
	return file_report_proto_rawDescGZIP(), []int{0}
}

func (x *StatusReport) GetTowerIP() string {
	if x != nil {
		return x.TowerIP
	}
	return ""
}

func (x *StatusReport) GetLastArduinoReachableTimestamp() string {
	if x != nil {
		return x.LastArduinoReachableTimestamp
	}
	return ""
}

func (x *StatusReport) GetLastTowerReachableTimestamp() string {
	if x != nil {
		return x.LastTowerReachableTimestamp
	}
	return ""
}

func (x *StatusReport) GetBootTimestamp() string {
	if x != nil {
		return x.BootTimestamp
	}
	return ""
}

func (x *StatusReport) GetRebootsCurrentDay() int64 {
	if x != nil {
		return x.RebootsCurrentDay
	}
	return 0
}

func (x *StatusReport) GetRAMUsed() int64 {
	if x != nil {
		return x.RAMUsed
	}
	return 0
}

func (x *StatusReport) GetRAMTotal() int64 {
	if x != nil {
		return x.RAMTotal
	}
	return 0
}

func (x *StatusReport) GetDiskUsed() int64 {
	if x != nil {
		return x.DiskUsed
	}
	return 0
}

func (x *StatusReport) GetDiskTotal() int64 {
	if x != nil {
		return x.DiskTotal
	}
	return 0
}

func (x *StatusReport) GetCPUAvg() int64 {
	if x != nil {
		return x.CPUAvg
	}
	return 0
}

func (x *StatusReport) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *StatusReport) GetReachable() bool {
	if x != nil {
		return x.Reachable
	}
	return false
}

func (x *StatusReport) GetTemperature() int64 {
	if x != nil {
		return x.Temperature
	}
	return 0
}

var File_report_proto protoreflect.FileDescriptor

var file_report_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x86, 0x04, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x6f, 0x77, 0x65,
	0x72, 0x49, 0x50, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x6f, 0x77, 0x65, 0x72,
	0x49, 0x50, 0x12, 0x44, 0x0a, 0x1d, 0x6c, 0x61, 0x73, 0x74, 0x41, 0x72, 0x64, 0x75, 0x69, 0x6e,
	0x6f, 0x52, 0x65, 0x61, 0x63, 0x68, 0x61, 0x62, 0x6c, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1d, 0x6c, 0x61, 0x73, 0x74, 0x41,
	0x72, 0x64, 0x75, 0x69, 0x6e, 0x6f, 0x52, 0x65, 0x61, 0x63, 0x68, 0x61, 0x62, 0x6c, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x40, 0x0a, 0x1b, 0x6c, 0x61, 0x73, 0x74,
	0x54, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x61, 0x63, 0x68, 0x61, 0x62, 0x6c, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1b, 0x6c,
	0x61, 0x73, 0x74, 0x54, 0x6f, 0x77, 0x65, 0x72, 0x52, 0x65, 0x61, 0x63, 0x68, 0x61, 0x62, 0x6c,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x24, 0x0a, 0x0d, 0x62, 0x6f,
	0x6f, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x62, 0x6f, 0x6f, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x2c, 0x0a, 0x11, 0x72, 0x65, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x43, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x44, 0x61, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x72, 0x65, 0x62,
	0x6f, 0x6f, 0x74, 0x73, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x79, 0x12, 0x18,
	0x0a, 0x07, 0x52, 0x41, 0x4d, 0x55, 0x73, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x52, 0x41, 0x4d, 0x55, 0x73, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x52, 0x41, 0x4d, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x52, 0x41, 0x4d, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x6b, 0x55, 0x73, 0x65, 0x64,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x69, 0x73, 0x6b, 0x55, 0x73, 0x65, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x73, 0x6b, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x64, 0x69, 0x73, 0x6b, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x16,
	0x0a, 0x06, 0x43, 0x50, 0x55, 0x41, 0x76, 0x67, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x43, 0x50, 0x55, 0x41, 0x76, 0x67, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x63, 0x68, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x09, 0x72, 0x65, 0x61, 0x63, 0x68, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0b, 0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x42, 0x09, 0x5a, 0x07, 0x2f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_report_proto_rawDescOnce sync.Once
	file_report_proto_rawDescData = file_report_proto_rawDesc
)

func file_report_proto_rawDescGZIP() []byte {
	file_report_proto_rawDescOnce.Do(func() {
		file_report_proto_rawDescData = protoimpl.X.CompressGZIP(file_report_proto_rawDescData)
	})
	return file_report_proto_rawDescData
}

var file_report_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_report_proto_goTypes = []interface{}{
	(*StatusReport)(nil),          // 0: report.StatusReport
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.timestamp
}
var file_report_proto_depIdxs = []int32{
	1, // 0: report.StatusReport.timestamp:type_name -> google.protobuf.timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_report_proto_init() }
func file_report_proto_init() {
	if File_report_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_report_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusReport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_report_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_report_proto_goTypes,
		DependencyIndexes: file_report_proto_depIdxs,
		MessageInfos:      file_report_proto_msgTypes,
	}.Build()
	File_report_proto = out.File
	file_report_proto_rawDesc = nil
	file_report_proto_goTypes = nil
	file_report_proto_depIdxs = nil
}
