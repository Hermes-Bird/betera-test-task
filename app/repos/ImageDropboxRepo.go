package repos

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Hermes-Bird/betera-test-task/app/config"
	"github.com/Hermes-Bird/betera-test-task/app/helpers"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	UploadFileDropboxUrl = "https://content.dropboxapi.com/2/files/upload"
	CreateSharedLinkUrl  = "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings"
	GetAuthTokenUrl      = "https://api.dropboxapi.com/2/auth/token/from_oauth1"
)

var TokenExpiredErr = errors.New("token expired")

type ImageDropboxRepo struct {
	AuthToken string
}

func NewImageDropboxRepo(cfg *config.Config) ImageRepo {
	return &ImageDropboxRepo{
		AuthToken: cfg.DropboxAuthToken,
	}
}

func uploadDropboxFile(authToken, path string, r io.Reader) (string, error) {
	auth := fmt.Sprintf("Bearer %s", authToken)
	apiArgs := fmt.Sprintf("{\"autorename\": false, \"mode\": \"add\", \"mute\": false, \"path\": \"%s\",   \"strict_conflict\": false}", path)
	req, err := http.NewRequest(http.MethodPost, UploadFileDropboxUrl, r)
	if err != nil {
		log.Println(err)
		return "", errors.New("error while parsing upload url")
	}

	req.Header = map[string][]string{
		"Content-type":    {"application/octet-stream"},
		"Authorization":   {auth},
		"Dropbox-API-Arg": {apiArgs},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", TokenExpiredErr
		}
		data, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("error while uploading: %s %s", data)
	}

	m, err := helpers.GetJsonBodyMap(resp)
	if err != nil {
		return "", err
	}

	// TODO add check
	dropboxPath, _ := m["path_display"].(string)
	return dropboxPath, nil
}

func createDropboxSharedLink(authToken, path string) (string, error) {
	auth := fmt.Sprintf("Bearer %s", authToken)
	body := []byte(fmt.Sprintf(`{
		"path": "%s",
		"settings": {
			"access": "viewer",
			"allow_download": true,
			"audience": "public",
				 "requested_visibility": "public"
			}
	}`, path))

	req, err := http.NewRequest(http.MethodPost, CreateSharedLinkUrl, bytes.NewReader(body))
	if err != nil {
		log.Println(err)
		return "", errors.New("error while parsing create_shared_link url")
	}

	req.Header = map[string][]string{
		"Content-type":  {"application/json"},
		"Authorization": {auth},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	m, err := helpers.GetJsonBodyMap(resp)
	if err != nil {
		return "", err
	}

	// TODO add check
	url, _ := m["url"].(string)
	return url, nil
}

func (r *ImageDropboxRepo) UploadImage(name string, reader io.Reader) (string, error) {
	path, err := uploadDropboxFile(r.AuthToken, fmt.Sprintf("/apods/%s", name), reader)
	if err != nil {
		return "", err
	}

	return createDropboxSharedLink(r.AuthToken, path)
}
