package service

import (
	"context"
	"github.com/google/uuid"
)

type PostgresService interface {
	InsertFileMetaData(ctx context.Context, data FileMetaData) error
	UpdateFileMetaData(ctx context.Context, data FileMetaData) error
	CheckFileIsExists(ctx context.Context, userID uuid.UUID, fileName string) (bool, error)
	SelectFileMetaDataList(ctx context.Context, userID uuid.UUID, fileName string) ([]FileMetaData, error)
}
