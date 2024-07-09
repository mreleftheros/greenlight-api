package main

import (
	"encoding/json"
	"net/http"
)

type Err map[string]string

func jsonRes(w http.ResponseWriter, data interface{}, header http.Header, status ...int) error {
	sts := 200
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if header != nil {
		for k, v := range header {
			w.Header()[k] = v
		}
	}

	if len(status) > 0 {
		sts = status[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sts)
	if _, err = w.Write(json); err != nil {
		return err
	}

	return nil
}

func errRes(w http.ResponseWriter, error interface{}, header http.Header, status ...int) error {
	sts := 400
	var errStruct Err = nil

	if error != nil {
		errStruct = Err{"error": error.(string)}
	}

	json, err := json.Marshal(errStruct)
	if err != nil {
		return err
	}

	if header != nil {
		for k, v := range header {
			w.Header()[k] = v
		}
	}

	if len(status) > 0 {
		sts = status[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sts)
	if _, err = w.Write(json); err != nil {
		return err
	}

	return nil
}
