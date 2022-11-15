package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func WriteHttpStringError(code int, err string, rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err)))
}

func WriteHttpError(code int, err error, rw http.ResponseWriter) {
	WriteHttpStringError(code, err.Error(), rw)
}

func WriteHttpJsonResponse(code int, data any, rw http.ResponseWriter) {
	bs, err := json.Marshal(data)
	if err != nil {
		WriteHttpError(http.StatusInternalServerError, errors.New("error while converting to json"), rw)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(code)
	rw.Write(bs)
}

func WriteHttpStringJsonResponse(code int, data string, rw http.ResponseWriter) {
	rw.Header().Set("Content-type", "application/json")
	rw.Write([]byte(data))
}

func GetJsonBodyMap(r *http.Response) (map[string]any, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	m := map[string]any{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.New("error while unmarshal dropbox data to map")
	}

	return m, nil
}
