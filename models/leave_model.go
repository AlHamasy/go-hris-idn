package models

import (
	"database/sql"
	"hris-idn/config"
	"hris-idn/entities"
	"log"
	"strings"
	"time"
)

type LeaveModel struct {
	db *sql.DB
}

func NewLeaveModel() *LeaveModel {
	conn, err := config.DBConn()
	if err != nil {
		log.Println("Failed connect to database:", err)
	}
	return &LeaveModel{
		db: conn,
	}
}

func (model LeaveModel) FindAllLeaveType() ([]entities.LeaveType, error) {
	rows, err := model.db.Query("SELECT id, name, max_day FROM leave_type WHERE deleted_at IS NULL")
	if err != nil {
		return []entities.LeaveType{}, err
	}
	defer rows.Close()

	var leaveType []entities.LeaveType

	for rows.Next() {
		var leave entities.LeaveType
		err := rows.Scan(
			&leave.Id,
			&leave.Name,
			&leave.MaxDay,
		)
		if err != nil {
			return []entities.LeaveType{}, err
		}
		leaveType = append(leaveType, leave)
	}

	return leaveType, nil
}

func (model LeaveModel) InsertLeave(data entities.SubmitLeave) error {

	query := `
		INSERT INTO leave_employee 
		(nik, leave_type_id, leave_date, attachment, reason, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := model.db.Exec(
		query,
		data.NIK,
		data.LeaveTypeId,
		data.LeaveDateJoin,
		data.Attachment,
		data.Reason,
		data.Status,
	)

	return err
}

func (model LeaveModel) GetLeaveList(nik string, monthYear string) ([]entities.Leave, error) {
	var query string
	var args []interface{}

	baseQuery := `
		SELECT 
			le.id,
			le.nik,
			le.leave_type_id,
			lt.name AS leave_type_name,
			le.leave_date,
			le.attachment,
			le.reason,
			le.status,
			le.reason_status,
			le.created_at,
			le.updated_at,
			em.name
		FROM leave_employee le
		LEFT JOIN leave_type lt ON le.leave_type_id = lt.id
		LEFT JOIN employee em ON le.nik = em.nik
		WHERE le.deleted_at IS NULL
	`

	// Jika monthYear tidak kosong, tambahkan filter bulan dan tahun
	if strings.TrimSpace(monthYear) != "" {
		parsedDate, err := time.Parse("January 2006", monthYear)
		if err != nil {
			return nil, err
		}
		baseQuery += " AND MONTH(le.created_at) = ? AND YEAR(le.created_at) = ?"
		args = append(args, parsedDate.Month(), parsedDate.Year())
	}

	// Jika nik diberikan, tambahkan filter nik
	if nik != "" {
		baseQuery += " AND le.nik = ?"
		args = append(args, nik)
	}

	query = baseQuery + " ORDER BY le.created_at DESC"

	rows, err := model.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []entities.Leave
	for rows.Next() {
		var leave entities.Leave
		err := rows.Scan(
			&leave.Id,
			&leave.NIK,
			&leave.LeaveTypeId,
			&leave.LeaveTypeName,
			&leave.LeaveDateJoin,
			&leave.Attachment,
			&leave.Reason,
			&leave.Status,
			&leave.ReasonStatus,
			&leave.CreatedAt,
			&leave.UpdatedAt,
			&leave.Name,
		)
		if err != nil {
			return nil, err
		}

		// Parse LeaveDateJoin ke []time.Time
		if leave.LeaveDateJoin != "" {
			dateStrings := strings.Split(leave.LeaveDateJoin, ",")
			for _, ds := range dateStrings {
				parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(ds))
				if err == nil {
					leave.LeaveDate = append(leave.LeaveDate, parsedDate)
				}
			}
		}

		leaves = append(leaves, leave)
	}

	return leaves, nil
}

func (model LeaveModel) GetLeaveById(id int64) (*entities.Leave, error) {
	query := `
		SELECT 
			le.id,
			le.nik,
			le.leave_type_id,
			lt.name AS leave_type_name,
			le.leave_date,
			le.attachment,
			le.reason,
			le.status,
			le.reason_status,
			le.created_at,
			le.updated_at,
			em.name
		FROM leave_employee le
		LEFT JOIN leave_type lt ON le.leave_type_id = lt.id
		LEFT JOIN employee em ON le.nik = em.nik
		WHERE le.deleted_at IS NULL AND le.id = ?
		LIMIT 1
	`

	var leave entities.Leave

	err := model.db.QueryRow(query, id).Scan(
		&leave.Id,
		&leave.NIK,
		&leave.LeaveTypeId,
		&leave.LeaveTypeName,
		&leave.LeaveDateJoin,
		&leave.Attachment,
		&leave.Reason,
		&leave.Status,
		&leave.ReasonStatus,
		&leave.CreatedAt,
		&leave.UpdatedAt,
		&leave.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Tidak ditemukan
		}
		return nil, err
	}

	// Parse LeaveDateJoin ke []time.Time
	if leave.LeaveDateJoin != "" {
		dateStrings := strings.Split(leave.LeaveDateJoin, ",")
		for _, ds := range dateStrings {
			parsedDate, err := time.Parse("2006-01-02", strings.TrimSpace(ds))
			if err == nil {
				leave.LeaveDate = append(leave.LeaveDate, parsedDate)
			}
		}
	}

	return &leave, nil
}

func (model LeaveModel) UpdateLeaveStatus(approvalLeave entities.ApprovalLeave) error {
	query := `
		UPDATE leave_employee
		SET status = ?, reason_status = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := model.db.Exec(
		query,
		approvalLeave.Status,
		approvalLeave.ReasonStatus,
		time.Now(),
		approvalLeave.Id,
	)
	return err
}
