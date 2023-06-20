package web

import (
	"encoding/json"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
	"io"
)

type Headerer interface {
	Headers() http.Header
}

func EncodeJson(w http.ResponseWriter, v interface{}, code int) error {
	if headerer, ok := v.(Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}

	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return nil
	}

	var jsonData []byte
	var err error
	switch v := v.(type) {
	case []byte:
		jsonData = v
	case io.Reader:
		jsonData, err = io.ReadAll(v)
	default:
		jsonData, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(code)
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
