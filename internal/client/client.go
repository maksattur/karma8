package client

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetAvailableServers(storageBalancerURL string) ([]StorageServer, error) {
	resp, err := http.Get(storageBalancerURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var servers []StorageServer
	if err = json.Unmarshal(body, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func (c *Client) UploadToServer(ip string, filePart []byte) (*FileMetaData, error) {
	url := fmt.Sprintf("http://%s/upload_part", ip)
	resp, err := http.Post(url, "application/octet-stream", strings.NewReader(string(filePart)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmd := &FileMetaData{}
	if err = json.Unmarshal(body, fmd); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error uploading to server: %s", resp.Status)
	}

	return fmd, nil
}

func (c *Client) DownloadFromServer(ip string, filePartName uuid.UUID) ([]byte, error) {
	url := fmt.Sprintf("http://%s/%s", ip, filePartName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error downloading from server: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
