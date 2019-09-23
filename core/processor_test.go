package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestHandler(t *testing.T) {
	type TestCase struct {
		Response, Request string
		StatusCode        int
	}
	cases := []TestCase{
		TestCase{
			Request:    `{"payload_data": "v=1&tid=UA-XXXXX-Y&cid=555&t=pageview&dp=%2Fhome"}`,
			StatusCode: http.StatusOK,
		},
		TestCase{
			Request:    ``,
			Response:   "EOF",
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Request:    `{"payload_data": "v=1&tid=UA-XXXXX-Y"}`,
			Response:   "missing or empty parameter: cid",
			StatusCode: http.StatusBadRequest,
		},
		TestCase{
			Request:    `{"payload_data": ""}`,
			Response:   "payload_data is empty",
			StatusCode: http.StatusBadRequest,
		},
	}
	for caseNum, item := range cases {
		url := "http://localhost:8080/collect"
		req := httptest.NewRequest("POST", url, bytes.NewReader([]byte(item.Request)))
		w := httptest.NewRecorder()
		RequestHandler(w, req)
		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] wrong Response: got %s, expected %s",
				caseNum, bodyStr, item.Response)
		}
	}
}

func TestWriteInFile(t *testing.T) {
	type TestCase struct {
		FileName, ErrorText string
	}
	cases := []TestCase{
		TestCase{
			FileName:  "../test.txt",
			ErrorText: "",
		},
		TestCase{
			FileName:  "../test.txt",
			ErrorText: "",
		},
		TestCase{
			FileName:  "",
			ErrorText: "filename is not assigned",
		},
	}
	for caseNum, item := range cases {
		err := writeInFile(item.FileName)
		errText := ""
		if err != nil {
			errText = err.Error()
		}
		if errText != item.ErrorText {
			t.Errorf("[%d] wrong StatusCode: got %s, expected %s",
				caseNum, errText, item.ErrorText)
		}
	}
}

func TestIsValid(t *testing.T) {
	type TestCase struct {
		Payload, ErrorText string
	}
	cases := []TestCase{
		TestCase{
			Payload:   "v=1&tid=UA-XXXXX-Y&cid=555&t=pageview&dp=%2Fhome",
			ErrorText: "",
		},
		TestCase{
			Payload:   `{"json": "test}`,
			ErrorText: "missing or empty parameter: v",
		},
		TestCase{
			Payload:   "2",
			ErrorText: "missing or empty parameter: v",
		},
		TestCase{
			Payload:   "v=1",
			ErrorText: "missing or empty parameter: tid",
		},
	}
	for caseNum, item := range cases {
		err := isValid(item.Payload)
		errText := ""
		if err != nil {
			errText = err.Error()
		}
		if errText != item.ErrorText {
			t.Errorf("[%d] wrong StatusCode: got %s, expected %s",
				caseNum, errText, item.ErrorText)
		}
	}
}
