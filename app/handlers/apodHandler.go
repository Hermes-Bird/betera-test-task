package handlers

import (
	"bytes"
	"fmt"
	"github.com/Hermes-Bird/betera-test-task/app/domain"
	"github.com/Hermes-Bird/betera-test-task/app/helpers"
	"github.com/Hermes-Bird/betera-test-task/app/repos"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ApodsHandler struct {
	ApodRepo      repos.ApodRepo
	ImageApodRepo repos.ImageRepo
}

func NewApodHandler(apodRepo repos.ApodRepo, imageRepo repos.ImageRepo) *ApodsHandler {
	return &ApodsHandler{
		ApodRepo:      apodRepo,
		ImageApodRepo: imageRepo,
	}
}

func (h *ApodsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.HandleGet(w, r)
	case "POST":
		h.HandlePost(w, r)
	}
}

func (h *ApodsHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	queryDate := r.URL.Query().Get("date")

	var err error
	var res []domain.Apod

	if queryDate != "" {
		parsedTime, e := time.Parse("2006-01-02", queryDate)
		if e != nil {
			helpers.WriteHttpStringError(http.StatusBadRequest, "query data format is invalid", w)
			return
		}
		res, err = h.ApodRepo.GetApodsForDate(parsedTime)
	} else {
		res, err = h.ApodRepo.GetApodsList()
	}

	if err != nil {
		helpers.WriteHttpError(http.StatusBadRequest, err, w)
		return
	}

	helpers.WriteHttpJsonResponse(http.StatusOK, res, w)

}

func (h *ApodsHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	fileName := r.Header.Get("File-Name")
	err := helpers.ValidateImgFileNameHeader(fileName)
	if err != nil {
		helpers.WriteHttpError(http.StatusBadRequest, err, w)
		return
	}

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil || (err == nil && len(bd) == 0) {
		helpers.WriteHttpStringError(http.StatusBadRequest, "no body data", w)
		return
	}

	_, _, err = image.Decode(bytes.NewReader(bd))
	if err != nil {
		helpers.WriteHttpStringError(http.StatusBadRequest, "specified body is not an image of supported format (png, jpeg gif)", w)
		return
	}

	url, err := h.ImageApodRepo.UploadImage(fileName, bytes.NewReader(bd))
	if err != nil {
		helpers.WriteHttpStringError(http.StatusInternalServerError, "error while uploading image", w)
		return
	}

	apod := domain.NewApod(url)
	err = h.ApodRepo.SaveApod(apod)
	if err != nil {
		log.Println(err)
		helpers.WriteHttpStringError(http.StatusInternalServerError, "error while saving to db", w)
		return
	}

	helpers.WriteHttpStringJsonResponse(http.StatusCreated, fmt.Sprintf(`{ "previewUrl": "%s" }`, url), w)
}
