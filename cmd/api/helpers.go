package main

import (
	"encoding/json"
	"net/http"
)

func jsonRes(w http.ResponseWriter, data interface{}, header http.Header, status ...int) error {
	sts := 200

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
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

func errRes(w http.ResponseWriter, errors map[string]string, header http.Header, status ...int) error {
	sts := 400

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
	if err := json.NewEncoder(w).Encode(errors); err != nil {
		return err
	}

	return nil
}

func jsonBody(r *http.Request, target_ptr interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(target_ptr); err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}
