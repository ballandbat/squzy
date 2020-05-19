package database

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"time"
)

type postgres struct {
	db *gorm.DB
}

type Snapshot struct {
	gorm.Model
	SchedulerID string    `gorm:"schedulerId"`
	Code        string    `gorm:"code"`
	Type        string    `gorm:"column:type"`
	Error       string    `gorm:"error"`
	Meta        *MetaData `gorm:"meta"`
}

type MetaData struct {
	gorm.Model
	SnapshotID uint           `gorm:"snapshotId"`
	StartTime  time.Time      `gorm:"startTime"`
	EndTime    time.Time      `gorm:"endTime"`
	Value      *_struct.Value `gorm:"value"` //TODO: google
}

//Agent gorm description
type StatRequest struct {
	gorm.Model
	AgentID    string      `gorm:"agentID"`
	AgentName  string      `gorm:"agentName"`
	CPUInfo    []*CPUInfo  `gorm:"cpuInfo"`
	MemoryInfo *MemoryInfo `gorm:"memoryInfo"`
	DiskInfo   []*DiskInfo `gorm:"diskInfo"`
	NetInfo    []*NetInfo  `gorm:"netInfo"`
	Time       time.Time   `gorm:"time"`
}

const (
	cpuInfoKey    = "cpuInfo"
	memoryInfoKey = "memoryInfo"
	diskInfoKey   = "diskInfo"
	netInfoKey    = "netInfo"
)

type CPUInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"statRequestId"`
	Load          float64 `gorm:"load"`
}

type MemoryInfo struct {
	gorm.Model
	StatRequestID uint    `gorm:"statRequestId"`
	Mem           *Memory `gorm:"mem"`
	Swap          *Memory `gorm:"swap"`
}

type Memory struct { ///Will we need check for Total = used + free + shared?
	gorm.Model
	MemoryInfoID uint    `gorm:"memoryInfoId"`
	Total        uint64  `gorm:"total"`
	Used         uint64  `gorm:"used"`
	Free         uint64  `gorm:"free"`
	Shared       uint64  `gorm:"shared"`
	UsedPercent  float64 `gorm:"usedPercent"`
}

type DiskInfo struct { ///Will we need check for Total = used + free?
	gorm.Model
	StatRequestID uint    `gorm:"statRequestId"`
	Name          string  `gorm:"name"`
	Total         uint64  `gorm:"total"`
	Free          uint64  `gorm:"free"`
	Used          uint64  `gorm:"used"`
	UsedPercent   float64 `gorm:"usedPercent"`
}

type NetInfo struct {
	gorm.Model
	StatRequestID uint   `gorm:"statRequestId"`
	Name          string `gorm:"name"`
	BytesSent     uint64 `gorm:"bytesSent"`
	BytesRecv     uint64 `gorm:"bytesRecv"`
	PacketsSent   uint64 `gorm:"packetsSent"`
	PacketsRecv   uint64 `gorm:"packetsRecv"`
	ErrIn         uint64 `gorm:"errIn"`
	ErrOut        uint64 `gorm:"errOut"`
	DropIn        uint64 `gorm:"dropIn"`
	DropOut       uint64 `gorm:"dropOut"`
}

const (
	dbSnapshotCollection    = "snapshots"     //TODO: check
	dbStatRequestCollection = "stat_requests" //TODO: check
)

var (
	errorDataBase = errors.New("ERROR_DATABASE_OPERATION")
)

func (p *postgres) Migrate() error {
	models := []interface{}{
		&Snapshot{},
		&MetaData{},
		&StatRequest{},
		&CPUInfo{},
		&MemoryInfo{},
		&Memory{},
		&DiskInfo{},
		&NetInfo{},
	}

	var err error
	for _, model := range models {
		err = p.db.AutoMigrate(model).Error // migrate models one-by-one
	}

	return err
}

func (p *postgres) InsertSnapshot(data *apiPb.SchedulerResponse) error {
	snapshot, err := ConvertToPostgresSnapshot(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbSnapshotCollection).Create(snapshot).Error; err != nil {
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetSnapshots(schedulerId string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.SchedulerSnapshot, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int
	err = p.db.Table(dbSnapshotCollection).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), schedulerId).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var dbSnapshots []*Snapshot
	err = p.db.Table(dbSnapshotCollection).
		Where(fmt.Sprintf(`"%s"."schedulerId" = ?`, dbSnapshotCollection), schedulerId).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbSnapshotCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&dbSnapshots).Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgresSnapshots(dbSnapshots), int32(count), nil
}

func (p *postgres) InsertStatRequest(data *apiPb.Metric) error {
	pgData, err := ConvertToPostgressStatRequest(data)
	if err != nil {
		return err
	}
	if err := p.db.Table(dbStatRequestCollection).Create(pgData).Error; err != nil {
		//TODO: log?
		return errorDataBase
	}
	return nil
}

func (p *postgres) GetStatRequest(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func (p *postgres) GetCPUInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, cpuInfoKey)
}

func (p *postgres) GetMemoryInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, memoryInfoKey)
}

func (p *postgres) GetDiskInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, diskInfoKey)
}

func (p *postgres) GetNetInfo(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	return p.getSpecialRecords(agentID, pagination, filter, netInfoKey)
}

func (p *postgres) getSpecialRecords(agentID string, pagination *apiPb.Pagination, filter *apiPb.TimeFilter, key string) ([]*apiPb.GetAgentInformationResponse_Statistic, int32, error) {
	timeFrom, timeTo, err := getTime(filter)
	if err != nil {
		return nil, -1, err
	}

	var count int
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Count(&count).Error
	if err != nil {
		return nil, -1, err
	}

	offset, limit := getOffsetAndLimit(count, pagination)

	//TODO: test if it works
	var statRequests []*StatRequest
	err = p.db.Table(dbStatRequestCollection).
		Where(fmt.Sprintf(`"%s"."agentID" = ?`, dbStatRequestCollection), agentID).
		Where(fmt.Sprintf(`"%s"."time" BETWEEN ? and ?`, dbStatRequestCollection), timeFrom, timeTo).
		Select(fmt.Sprintf("%s, time", key)).
		Order("time").
		Offset(offset).
		Limit(limit).
		Find(&statRequests).Error

	if err != nil {
		return nil, -1, errorDataBase
	}

	return ConvertFromPostgressStatRequests(statRequests), int32(count), nil
}

func getTime(filter *apiPb.TimeFilter) (time.Time, time.Time, error) {
	timeFrom := time.Unix(0, 0)
	timeTo := time.Now()
	var err error
	if filter != nil {
		if filter.GetFrom() != nil {
			timeFrom, err = ptypes.Timestamp(filter.From)
		}
		if filter.GetTo() != nil {
			timeFrom, err = ptypes.Timestamp(filter.From)
		}
	}
	return timeFrom, timeTo, err
}

func getOffsetAndLimit(count int, pagination *apiPb.Pagination) (int32, int32) {
	offset := int32(0)
	limit := int32(count)
	if pagination != nil {
		offset = pagination.GetLimit() * pagination.GetPage()
		limit = pagination.GetLimit()
	}
	return offset, limit
}
