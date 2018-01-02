package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Result struct {
	Code   int
	Result interface{}
}

func Failed(errCode int) *Result {
	return &Result{
		Code:   http.StatusInternalServerError,
		Result: errCode,
	}
}

func Succeed() *Result {
	return &Result{
		Code:   http.StatusOK,
		Result: nil,
	}
}

func SucceedWithResult(result interface{}) *Result {
	return &Result{
		Code:   http.StatusOK,
		Result: result,
	}
}

func EncodeResult(w http.ResponseWriter, result *Result) {
	w.WriteHeader(int(result.Code))
	if result.Result == nil {
		return
	}

	var body []byte
	if result.Code == http.StatusInternalServerError {
		body, _ = json.Marshal(struct {
			ErrCode int `json:"err_code"`
		}{result.Result.(int)})
	} else {
		body, _ = json.Marshal(result.Result)
	}
	w.Write(body)
}

func DecodeRequestBody(req *http.Request, params interface{}) error {
	reqBody, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(reqBody, params)
}
