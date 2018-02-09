package proto

import (
	"log"
	"os"
	"testing"

	"encoding/json"

	"reflect"

	"github.com/emicklei/proto"
)

//func TestCamelCase(t *testing.T) {
//	str := `syntax="proto3";
//package ctob;
//
//import "proto/service/ctob/common.proto";
//
//service CarSource{
//    /// GetCarSourceList 车源列表
//    rpc GetCarSourceList(CarSourceListRequest) returns(CarSourceListResponse);
//    /// GetCarSourceDetail 车源详情
//    rpc GetCarSourceDetail(CarSourceDetailRequest) returns(CarSourceDetailResponse);
//    /// GetAppointFailReason 获取失败原因
//    rpc GetAppointFailReason(FailReasonRequest) returns(FailReasonResponse);
//    /// GetStopSaleReason 获取停售原因
//    rpc GetStopSaleReason(EmptyRequest) returns(StopSaleReasonResponse);
//    /// GetAppointSchedule 获取排班信息
//    rpc GetAppointSchedule(AppointScheduleRequest) returns(AppointScheduleResponse);
//    /// GetAppointStoreSchedule 获取门店排班信息
//    rpc GetAppointStoreSchedule(AppointStoreScheduleRequest) returns(AppointScheduleResponse);
//    /// GetCarSourceListPermission 获取列表权限
//    rpc GetCarSourceListPermission(EmptyRequest) returns(CarSourceListPermissionResponse);
//    ///GetCarSourceInfoPermission 获取详情权限
//    rpc GetCarSourceInfoPermission(CarSourceInfoPermissionRequest) returns(CarSourceInfoPermissionResponse);
//}
//
/////FundSource 资方接口
//service FundSource {
//    ///获取贷款合同列表
//    rpc GetFundSourceMap (common.DefaultEmpty) returns (FundSourceResponse) {
//    }
//}
//
//service Preaudit {
//    ///获取预审结果
//    rpc Get(GetRequest) returns (GetResponse) {
//    }
//}
//service Order {
//    /// 订单列表
//    rpc List (OrderListRequest) returns (OrderListResponse) {}
//    /// 预审页面
//    rpc PreAuditDetail (OrderDetailRequest) returns (OrderPreAuditDetailResponse) {}
//    /// 材料审核页面
//    rpc AuditDetail (OrderDetailRequest) returns (OrderAuditDetailResponse) {}
//    /// 提交审核结果
//    rpc AuditSubmit (AuditSubmitRequest) returns (OrderAuditDetailResponse) {}
//    /// 初审提交
//    //rpc FirstCheckSubmit (AuditSubmitRequest) returns (OrderAuditDetailResponse) {}
//    /// 复审提交
//    //rpc RecheckSubmit (AuditSubmitRequest) returns (OrderAuditDetailResponse) {}
//    /// 订单处理记录
//    rpc AuditRecord (AuditRecordRequest) returns (AuditRecordResponse) {}
//    /// 订单状态数量统计
//    rpc StatusCount (OrderStatusCountRequest) returns (OrderStatusCountResponse) {}
//    /// 重置回预审核
//    rpc ResetPreCheck(ResetPreCheckRequest) returns (common.DefaultEmpty) {}
//}
//
/////GPS信息服务
//service Gps {
//    ///获取工单信息
//    rpc GetBySheetId(GetGpsRequest) returns (WorksheetInfo) {
//    }
//    ///获取GPS 信息，包含工单信息，GPS信息，安装进度，图片信息
//    rpc GetGpsInfoBySheeId (GetGpsRequest)returns(GpsInfoResponse) {
//
//    }
//}
//
//service QueryData {
//    /// 获取通话录音数据
//    rpc GetCallRecord(CallRecordRequest) returns (CallRecordResponse) {}
//    /// 批量获取通话记录数据
//    rpc GetBatchCallRecord(BatchCallRecordRequest) returns (BatchCallRecordResponse) {}
//    /// 根据commit_id获取通话记录数据
//    rpc GetCallRecordsByCommitID(RecordsByCommitIDRequest) returns (BatchCallRecordResponse) {}
//    /// 根据source_id获取通话记录数据
//    rpc GetCallRecordsBySourceID(RecordsBySourceIDRequest) returns (BatchCallRecordResponse) {}
//    /// 根据手机号和时间戳批量获取通话记录数据
//    rpc GetBatchRecordsByPhone(BatchRecordsByPhoneRequest) returns (BatchCallRecordResponse) {}
//}
//
//service CallSaleProcessor {
//    // 处理提交业务的接口: 新建带看工单
//    rpc SubmitAppointTask(SubmitAppointTaskRequest) returns (google.protobuf.Empty) {}
//    // 核验保卖可下单
//    rpc CheckConsignTaskCanSubmit(ConsignTaskCanSubmitRequest) returns (ConsignTaskCanSubmitResponse) {}
//    // 处理提交业务的接口: 保卖工单
//    rpc SubmitConsignTask(SubmitConsignTaskRequest) returns (google.protobuf.Empty) {}
//    // 处理提交业务的接口: 回访
//    rpc StashAssignment(StashAssignmentRequest) returns (google.protobuf.Empty) {}
//    ///结束任务接口
//    rpc FinishTask(FinishTaskRequest) returns(google.protobuf.Empty){}
//}
//
//// 目标管理服务
//service StaffPlansService {
//    /// GetPlanList 获取目标列表数据.
//    rpc GetPlanList (StaffPlansPlanListRequest) returns (StaffPlansPlanListResponse);
//    /// GetPlanStat 获取人员目标数据.
//    rpc GetPlanStat (StaffPlansPlanStatRequest) returns (StaffPlansPlanStatResponse);
//}
//
//// 目标管理服务
//service StaffPlansService {
//    /// GetPlanList 获取目标列表数据.
//    rpc GetPlanList (StaffPlansPlanListRequest) returns (StaffPlansPlanListResponse);
//    /// GetPlanStat 获取人员目标数据.
//    rpc GetPlanStat (StaffPlansPlanStatRequest) returns (StaffPlansPlanStatResponse);
//}
///// MyCustomerService 我的客服仪表盘服务
//service MyCustomerService {
//    /// 获取我的客服仪表盘按天薪酬
//    rpc GetDailySalary        (Request) returns (DailySalaryResponse);
//    /// 获取我的客服仪表盘按月薪酬
//    rpc GetMonthlySalary      (Request) returns (MonthlySalaryResponse);
//    /// 获取我的客服仪表盘按天业绩
//    rpc GetDailyPerformance   (Request) returns (DailyPerformanceResponse);
//    /// 获取我的客服仪表盘按月业绩
//    rpc GetMonthlyPerformance (Request) returns (MonthlyPerformanceResponse);
//
//    /// 获取我的客服仪表盘额外信息(通话时长, 等级等）
//    rpc GetAdditionalItems(AdditionalItemsRequest) returns (AdditionalItemsResponse);
//}
//
//service AssignmentList {
//    /// 获取客服逾期任务列表接口
//    rpc ListDelayAssignments(ListRequest) returns (DelayAssignmentListResponse) {}
//    /// 获取客服未处理线索类任务列表接口
//    rpc ListClueAssignments(ListRequest) returns (ClueAssignmentListResponse) {}
//    /// 获取客服未处理回访类任务列表接口
//    rpc ListNextCallAssignments(ListRequest) returns (NextCallAssignmentListResponse) {}
//    /// 获取客服未处理失败工单任务列表接口
//    rpc ListFailedTaskAssignments(ListRequest) returns (FailedTaskAssignmentListResponse) {}
//    /// 获取列表页统计信息数据
//    rpc GetStatistics(google.protobuf.Empty) returns (StatisticsResponse) {}
//    /// 获取客户的相关任务列表接口
//    rpc ListCustomerRelatedAssignments(ListCustomerRelatedAssignmentRequest) returns (CustomerRelatedAssignmentListResponse) {}
//}
///**
// * Assignment list related messages.
// *
// * 任务列表相关定义
// *
// * Author: Stella
// */
//syntax = "proto3";
//package assignment;
//
//import "google/protobuf/empty.proto";
//import "proto/common/pagination.proto";
//import "proto/common/filter.proto";
//import "proto/service/assignment/base.proto";
//
///**
// * 任务列表服务
// */
//service AssignmentList {
//    /// 获取客服逾期任务列表接口
//    rpc ListDelayAssignments(ListRequest) returns (DelayAssignmentListResponse) {}
//    /// 获取客服未处理线索类任务列表接口
//    rpc ListClueAssignments(ListRequest) returns (ClueAssignmentListResponse) {}
//    /// 获取客服未处理回访类任务列表接口
//    rpc ListNextCallAssignments(ListRequest) returns (NextCallAssignmentListResponse) {}
//    /// 获取客服未处理失败工单任务列表接口
//    rpc ListFailedTaskAssignments(ListRequest) returns (FailedTaskAssignmentListResponse) {}
//    /// 获取列表页统计信息数据
//    rpc GetStatistics(google.protobuf.Empty) returns (StatisticsResponse) {}
//    /// 获取客户的相关任务列表接口
//    rpc ListCustomerRelatedAssignments(ListCustomerRelatedAssignmentRequest) returns (CustomerRelatedAssignmentListResponse) {}
//}
//
///**
// * 任务统计信息服务
// */
//service AssignmentStatistics {
//    /// 获取任务列表中各项数量
//    rpc GetAssignmentNum(google.protobuf.Empty) returns (AssignmentNumResponse){}
//}
//
///**
// * 表格排序
// */
//message TableSorter {
//    /**
//     * 排序的字段
//     */
//    enum Prop {
//        Unknown = 0;
//        CreatedAt = 1;          ///  创建时间
//        DueDate = 2;            ///  到期时间，即期望完成时间
//        ScBusinessTypeId = 3;   ///  线索类型
//        CallBackType = 4;       ///  回访类型
//        AssignmentCategory = 5; ///  任务类别
//        Weight = 6;             ///  权重
//        NextCallTime = 7;       ///  回访时间
//    }
//
//    Prop prop = 1;              ///  按表格中哪个字段来排序，范围参考枚举值，并根据具体列表中需要排序项做设置
//    common.Order order = 2;     ///  升序or降序
//}
//
///**
// * 请求：列表
// */
//message ListRequest {
//    common.TimePeriod time_period = 1;  /// 时间段
//    common.Pagination pagination = 2;   /// 分页
//    TableSorter table_sorter = 3;       /// 表格排序
//}
//
///**
// * 逾期任务
// */
//message DelayAssignment {
//    int32 id = 1;                               /// 任务号
//    int32 customer_id = 2;                      /// 客户号
//    string customer_name = 3;                   /// 客户姓名
//    Category category = 4;                      /// 任务类别
//    int64 created_at = 5;                       /// 创建时间
//    int64 due_date = 6;                         /// 到期时间，即期望完成时间
//    int32 car_id = 7;                           /// 车源号
//    int64 delay_duration = 8;                   /// 逾期时长
//    string remark = 9;                          /// 备注
//    bool starred = 12;                          /// 是否加星标
//    bool is_pre_audit_passed = 13;              /// 是否金融预审通过
//}
//
///**
// * 线索类任务
// */
//message ClueAssignment {
//    int32 id = 1;                       /// 任务号
//    int32 customer_id = 2;              /// 客户号
//    string customer_name = 3;           /// 客户姓名
//    bool is_regular_customer = 4;       /// 是否为老客户，true:老客户，false:新客户
//    int32 sc_business_type_id = 5;      /// 线索类型ID
//    string sc_business_type_name = 6;   /// 线索类型名称
//    int64 created_at = 7;               /// 创建时间
//    int64 due_date = 8;                 /// 到期时间，即期望完成时间
//    int64 relaimed_at = 9;              /// 回收时间
//    int32 car_id = 10;                  /// 车源号
//    string remark = 11;                 /// 备注
//    bool starred = 12;                  /// 是否加星标
//    bool is_pre_audit_passed = 13;              /// 是否金融预审通过
//}
//
///**
// * 回访类任务
// */
//message NextCallAssignment {
//    int32 id = 1;                               /// 任务号
//    int32 customer_id = 2;                      /// 客户号
//    string customer_name = 3;                   /// 客户姓名
//    CallBackType call_back_type = 4;            /// 回访类别
//    repeated CallBackSubReason call_back_sub_reason = 6; /// 回访子原因
//    int64 created_at = 7;                       /// 创建时间
//    int64 due_date = 8;                         /// 到期时间，即期望完成时间
//    string remark = 9;                          /// 备注
//    bool starred = 10;                          /// 是否加星标
//    int64 next_call_time = 11;                  /// 下次回访时间
//    bool is_pre_audit_passed = 12;              /// 是否金融预审通过
//}
//
///**
// * 失败工单类任务
// */
//message FailedTaskAssignment {
//    int32 id = 1;                                   /// 任务号
//    int32 customer_id = 2;                          /// 客户号
//    string customer_name = 3;                       /// 客户姓名
//    int32 task_id = 4;                              /// 工单号
//    int64 task_created_at = 5;                      /// 工单创建时间
//    int64 task_failed_at = 6;                       /// 工单失败时间
//    string task_status = 7;                         /// 工单状态
//    string task_result = 8;                         /// 工单结果，失败原因
//    string sales_remark = 9;                        /// 销售备注
//    string telesales_remark = 10;                   /// 电销备注
//    int64 due_date = 11;                            /// 到期时间，即期望完成时间
//    bool starred = 12;                              /// 是否加星标
//    bool is_pre_audit_passed = 13;              /// 是否金融预审通过
//}
//
///**
// * 响应：逾期任务列表
// */
//message DelayAssignmentListResponse {
//    repeated DelayAssignment assignments = 1;
//    common.Pagination pagination = 2;
//}
//
///**
// * 响应：线索类任务列表
// */
//message ClueAssignmentListResponse {
//    repeated ClueAssignment assignments = 1;
//    common.Pagination pagination = 2;
//}
//
///**
// * 响应：回访类任务列表
// */
//message NextCallAssignmentListResponse {
//    repeated NextCallAssignment assignments = 1;
//    common.Pagination pagination = 2;
//}
//
///**
// * 响应：失败工单类任务列表
// */
//message FailedTaskAssignmentListResponse {
//    repeated FailedTaskAssignment assignments = 1;
//    common.Pagination pagination = 2;
//}
//
///**
// * 响应：客服任务列表数据统计
// */
//message StatisticsResponse {
//    int32 delay_assignment_num = 1;                ///逾期任务数量
//    int32 clue_assignment_num = 2;                 ///线索类任务数量
//    int32 today_next_call_assignment_num = 3;      ///今日回访任务数量
//    int32 tomorrow_next_call_assignment_num = 4;   ///明日回访任务数量
//    int32 others_next_call_assignment_num = 5;     ///其他回访任务数量
//    int32 failed_task_assignment_num = 6;          ///失败工单任务数量
//
//    int32 today_failed_task_assignment_num = 7;    ///今日失败工单任务数量
//    int32 today_customer_num = 8;                  ///今日处理客户量,已处理任务量按照手机号码去重
//    int32 today_call_out_times = 9;                ///今日外呼接听次数
//    int32 today_call_in_times = 10;                ///今日接听次数
//    int32 today_create_task_num = 11;              ///建单量,当日当前员工工单创建量,包含保卖工单
//}
//
///**
// * 响应：客服任务列表数据统计, 只统计任务数量
// */
//message AssignmentNumResponse {
//    int32 delay_assignment_num = 1;                ///逾期任务数量
//    int32 clue_assignment_num = 2;                 ///线索类任务数量
//    int32 today_next_call_assignment_num = 3;      ///今日回访任务数量
//    int32 tomorrow_next_call_assignment_num = 4;   ///明日回访任务数量
//    int32 others_next_call_assignment_num = 5;     ///其他回访任务数量
//    int32 failed_task_assignment_num = 6;          ///失败工单任务数量
//}
//
///**
// * 请求：客户的相关任务列表接口需要的入参
// */
//message ListCustomerRelatedAssignmentRequest {
//    string customer_phone_encrypt = 1;  /// 客户手机号密文
//    common.Pagination pagination = 2;   /// 分页
//}
//
///**
// * 客户的相关任务
// */
//message CustomerRelatedAssignment {
//    int32 id = 1;                               /// 任务ID
//    Category category = 2;                      /// 任务类型
//    int64 created_at = 3;                       /// 任务创建时间
//    Status status = 4;                          /// 任务状态
//    string sc_business_type_name = 5;           /// 线索类型名称
//    repeated CallBackSubReason call_back_sub_reason = 7; /// 回访子原因
//    int64 updated_at = 8;                       /// 处理时间
//    string result = 9;                          /// 提交结果
//    int32 operator_id = 10;                     /// 处理客服的坐席号
//    string operator_name = 11;                  /// 处理客服的姓名
//}
//
///**
// * 响应：客户的相关任务列表接口的返回数据
// */
//message CustomerRelatedAssignmentListResponse {
//    repeated CustomerRelatedAssignment assignments = 1; /// 客户相关任务
//    common.Pagination pagination = 2;   /// 分页
//}
//
//`
//	//rpc [A-Za-z]+ \([A-Za-z.]+\) returns \([A-Za-z]+\)
//	//service [A-Za-z]+ {[A-Za-z\p{Han}/\\\n{} ().]*}
//	//service [A-Za-z]+ {(rpc [A-Za-z]+ \([A-Za-z.]+\) returns \([A-Za-z]+\)({(\n)*}))+|\n|(/{2,3}(+))}
//r, err := regexp.Compile(`service[A-Za-z ]+{([A-Za-z_:.）（(), \\\n\p{Han}/，]*(rpc[A-Za-z ]+\([A-Za-z.]+\)[ ]*returns[ ]*\([A-Za-z.]+\)([; {}\\\n]*)))+}`)
//	if err != nil {
//		log.Println(err)
//	}
//	arr := r.FindAllString(str, -1)
//	if len(arr) > 0 {
//		fmt.Println(len(arr))
//		for _, v := range arr {
//			fmt.Printf("%v", v)
//			fmt.Println()
//		}
//	}
//}

func TestCamelCasea(t *testing.T) {
	reader, _ := os.Open("/Users/luan/go/src/protobuf-schema/proto/finance/service/schedule/schedule_record.proto")
	defer reader.Close()
	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()
	for _, v := range definition.Elements {
		ty := reflect.TypeOf(v)
		log.Println(ty)
		//	if a, ok := v.(*proto.Service); ok {
		//		for _, b := range a.Elements {
		//			c := reflect.TypeOf(b)
		//			log.Println(c)
		//		}
		//	}
	}
	d, _ := json.Marshal(definition)
	log.Println(string(d))
}

func TestGeneratorProto(t *testing.T) {
	g := GeneratorProto{}
	err := g.Generator("/Users/luan/go/src/protobuf-schema/proto/finance/service/borrow/borrow.proto")
	if err != nil {
		log.Println(err)
		return
	}
	data, _ := json.Marshal(g)
	log.Println(string(data))
}
