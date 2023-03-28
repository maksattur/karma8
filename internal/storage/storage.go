package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/maksattur/karma8/internal/service"
	"time"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) InsertFileMetaData(ctx context.Context, data service.FileMetaData) error {
	query := `insert into files_meta_data(user_id, server_ip, file_name, file_part_name, file_part_number) values ($1, $2, $3, $4, $5)`
	if _, err := s.db.ExecContext(ctx, query, data.UserID, data.ServerIP, data.FileName, data.FilePartName, data.FilePartNumber); err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				return service.ErrForeignKeyOrUniqueViolation
			}
		}
		return fmt.Errorf("insert file meta data: %w", err)
	}
	return nil
}

func (s *Storage) UpdateFileMetaData(ctx context.Context, data service.FileMetaData) error {
	query := `update files_meta_data set 
                user_id = $1,
                server_ip = $2,
                file_name = $3,
                file_part_name = $4,
                file_part_number = $5,
                updated_at = $6
            where user_id = $1 and file_name = $3`
	if _, err := s.db.ExecContext(ctx, query, data.UserID, data.ServerIP, data.FileName, data.FilePartName, data.FilePartName, time.Now()); err != nil {
		return fmt.Errorf("update file meta data: %w", err)
	}

	return nil
}

func (s *Storage) CheckFileIsExists(ctx context.Context, userID uuid.UUID, fileName string) (bool, error) {
	query := `select fmd.file_name from files_meta_data as fmd where fmd.user_id = $1 and fmd.file_name = $2`
	var fm string
	if err := s.db.GetContext(ctx, &fm, query, userID, fileName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("get file meta data: %w", err)
	}
	return true, nil
}

func (s *Storage) SelectFileMetaDataList(ctx context.Context, userID uuid.UUID, fileName string) ([]service.FileMetaData, error) {
	query := `select * from files_meta_data as fmd where fmd.user_id = $1 and fmd.file_name = $2 order by fmd.file_part_number desc`
	rows, err := s.db.QueryContext(ctx, query, userID, fileName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metaList []service.FileMetaData

	for rows.Next() {
		var meta service.FileMetaData
		if err := rows.Scan(&meta); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		metaList = append(metaList, meta)
	}
	return metaList, nil
}
