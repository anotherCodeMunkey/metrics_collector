package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var payloadBatch batch

type batch struct {
	sync.Mutex
	data string
}

func (b *batch) set(data string) {
	b.Lock()
	defer b.Unlock()
	b.data += "\n" + data
}

// RequestHandler handler from metrics
func RequestHandler(resp http.ResponseWriter, req *http.Request) {
	var reqData = struct {
		PayloadData string `json:"payload_data"`
	}{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqData)

	if err != nil {
		response(http.StatusBadRequest, []byte(err.Error()), resp)
		return
	}

	if err := isValid(reqData.PayloadData); err != nil {
		response(http.StatusBadRequest, []byte(err.Error()), resp)
		return
	}

	payloadBatch.set(reqData.PayloadData)

	response(http.StatusOK, nil, resp)
	return
}

// WriteManager ticker for data writing
func WriteManager() {
	c := time.Tick(time.Second * time.Duration(viper.GetInt("SaveTimeout")))
	for range c {
		if payloadBatch.data != "" {
			if err := writeInFile(viper.GetString("Filename")); err != nil {
				log.Println("write in file failed: ", err)
			}
		}
	}
}

func writeInFile(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename is not assigned")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return payloadBatch.createAndWrite(filename)
	}
	return payloadBatch.openAndWrite(filename)
}

func (b *batch) createAndWrite(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	b.Lock()
	defer func() {
		file.Close()
		b.Unlock()
	}()
	_, err = file.WriteString(b.data)
	b.data = ""
	return err
}

func (b *batch) openAndWrite(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	b.Lock()
	defer func() {
		file.Close()
		b.Unlock()
	}()
	_, err = fmt.Fprintln(file, b.data)
	b.data = ""
	return err
}

func isValid(payloadData string) error {
	if payloadData == "" {
		return fmt.Errorf("payload_data is empty")
	}
	keys := []string{"v", "tid", "cid", "t"}
	values, err := url.ParseQuery(payloadData)
	if err != nil {
		return err
	}
	if values == nil {
		return fmt.Errorf("unable to parse payload_data")
	}
	for _, k := range keys {
		if values.Get(k) == "" {
			return fmt.Errorf("missing or empty parameter: %s", k)
		}
	}
	return nil
}

func response(errCode int, payload []byte, resp http.ResponseWriter) {
	resp.WriteHeader(errCode)
	resp.Write(payload)
}
