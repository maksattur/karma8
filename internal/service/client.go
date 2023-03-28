package service

import (
	"github.com/google/uuid"
	"github.com/maksattur/karma8/internal/client"
)

type AvailableServersGetter interface {
	GetAvailableServers(storageBalancerURL string) ([]client.StorageServer, error)
	UploadToServer(ip string, filePart []byte) (*client.FileMetaData, error)
	DownloadFromServer(ip string, filePartName uuid.UUID) ([]byte, error)
}
