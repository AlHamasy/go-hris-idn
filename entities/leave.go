package entities

import (
	"database/sql"
	"time"
)

type LeaveType struct {
	Id     int64
	Name   string
	MaxDay int64
}

type SubmitLeave struct {
	Id            int64
	NIK           string
	LeaveDate     []string `validate:"required,gte=1" label:"Tanggal Cuti"`
	LeaveDateJoin string
	Attachment    string
	Reason        string `validate:"required" label:"Alasan"`
	LeaveTypeId   string `validate:"required" label:"Tipe Cuti"`
	Status        int64
	ReasonStatus  sql.NullString
	CreatedAt     time.Time
	UpdatedAt     sql.NullTime
}

type Leave struct {
	Id            int64
	LeaveTypeId   int64
	LeaveTypeName string
	NIK           string
	LeaveDateJoin string
	Attachment    string
	Reason        string
	Status        int64
	ReasonStatus  sql.NullString
	CreatedAt     time.Time
	UpdatedAt     sql.NullTime
}
