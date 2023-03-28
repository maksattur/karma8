package client

import "github.com/google/uuid"

type FileMetaData struct {
	Name uuid.UUID `json:"name"`
}
