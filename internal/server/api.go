package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/maksattur/karma8/internal/service"
	"io"
)

const (
	storageBalancerURL = "http://localhost:8081/available_servers"
	countServers       = 5
	userID             = "9c4a5bc1-e314-4969-80ec-c3997305e10c"
)

type API struct {
	storage service.PostgresService
	client  service.AvailableServersGetter
}

func NewAPI(storage service.PostgresService, client service.AvailableServersGetter) *API {
	return &API{
		storage: storage,
		client:  client,
	}
}

func (a *API) UploadFile(ctx *fiber.Ctx) error {
	// uuid нашего дефолтного и единственного пользователя
	// аутентификация и авторизация в реализации отсутствует
	// uuid пользователя исключительно нужен для БД
	uid, err := uuid.Parse(userID)
	if err != nil {
		return NewError(fmt.Errorf("form file: %v", err), fiber.NewError(fiber.StatusUnauthorized, "unauthorized"))
	}

	// получаем загружаемый файл
	inFile, err := ctx.FormFile("file")
	if err != nil {
		return NewError(fmt.Errorf("form file: %v", err), fiber.NewError(fiber.StatusBadRequest, "invalid argument"))
	}

	// запрашиваем доступные сервера
	servers, err := a.client.GetAvailableServers(storageBalancerURL)
	if err != nil || len(servers) < countServers {
		return NewError(fmt.Errorf("get available servers: %v", err), fiber.NewError(fiber.StatusInternalServerError, "no available servers to upload"))
	}

	// проверяем есть ли метаданные файла юзера у нас в БД
	isExists, err := a.storage.CheckFileIsExists(ctx.UserContext(), uid, inFile.Filename)
	if err != nil {
		return NewError(fmt.Errorf("check file is exists: %v", err), fiber.ErrInternalServerError)
	}

	// открываем файл для чтения
	src, err := inFile.Open()
	if err != nil {
		return NewError(fmt.Errorf("open file: %v", err), fiber.ErrInternalServerError)
	}
	defer src.Close()

	// читаем файл
	file, err := io.ReadAll(src)
	if err != nil {
		return NewError(fmt.Errorf("read file: %v", err), fiber.ErrInternalServerError)
	}

	parts := splitFile(file)

	// загружаем части файла и записываем в БД метаданные файла
	for i, part := range parts {
		ip := servers[i%len(servers)].IP
		fileMetaData, err := a.client.UploadToServer(ip, part)
		if err != nil {
			return NewError(fmt.Errorf("uploading file part: %v", err), fiber.ErrInternalServerError)
		}

		meta := service.FileMetaData{
			UserID:         uid,
			ServerIP:       ip,
			FileName:       inFile.Filename,
			FilePartName:   fileMetaData.Name,
			FilePartNumber: uint8(i),
		}

		if !isExists {
			if err := a.storage.InsertFileMetaData(ctx.UserContext(), meta); err != nil {
				return NewError(fmt.Errorf("write file meta data: %v", err), fiber.ErrInternalServerError)
			}
		} else {
			if err := a.storage.UpdateFileMetaData(ctx.UserContext(), meta); err != nil {
				return NewError(fmt.Errorf("update file meta data: %v", err), fiber.ErrInternalServerError)
			}
		}
	}

	return nil
}

func (a *API) DownloadFile(ctx *fiber.Ctx) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return NewError(fmt.Errorf("form file: %v", err), fiber.NewError(fiber.StatusUnauthorized, "unauthorized"))
	}

	// получаем метаданные файла с сортировкой по порядку загруженности по серверам
	metaList, err := a.storage.SelectFileMetaDataList(ctx.UserContext(), uid, ctx.Params("file_name"))
	if err != nil {
		return NewError(fmt.Errorf("read metadata: %v", err), fiber.ErrInternalServerError)
	}

	// нет метаданных для запрашиваемого файла
	if len(metaList) == 0 {
		return NewError(fmt.Errorf("file metadata not found: %v", err), fiber.ErrNotFound)
	}

	// количество метаданных не верное, что то не так при записи
	if len(metaList) != countServers {
		return NewError(fmt.Errorf("wrong metadata: %v", err), fiber.ErrInternalServerError)
	}

	var parts [][]byte

	for _, m := range metaList {
		part, err := a.client.DownloadFromServer(m.ServerIP, m.FilePartName)
		if err != nil {
			return NewError(fmt.Errorf("downloading file part: %v", err), fiber.ErrInternalServerError)
		}

		parts = append(parts, part)
	}

	return ctx.Status(fiber.StatusOK).Send(joinFile(parts))
}

func splitFile(file []byte) [][]byte {
	partSize := len(file) / countServers
	var parts [][]byte
	for i := 0; i < countServers; i++ {
		start := i * partSize
		end := start + partSize
		if i == countServers-1 {
			end = len(file)
		}
		parts = append(parts, file[start:end])
	}
	return parts
}

func joinFile(parts [][]byte) []byte {
	var result []byte
	for i, v := range parts {
		result = append(result, v[i])
	}
	return result
}
