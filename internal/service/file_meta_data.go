package service

import "github.com/google/uuid"

type FileMetaData struct {
	UserID         uuid.UUID `db:"user_id"`
	ServerIP       string    `db:"server_ip"`
	FileName       string    `db:"file_name"`
	FilePartName   uuid.UUID `db:"file_part_name"`
	FilePartNumber uint8     `db:"file_part_number"`
}
