package scanner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	cmn "modules/scanmod/common"
	cfg "modules/scanmod/config"
)

func FS(op *cfg.UserOptions) (string, bool) {

	if cfg.DEBUG {
		log.Println("Checking Existence of File.....")
	}
	if err, ok := cmn.IsFileExists(op.File); !ok {
		return err, ok
	}

	if cfg.DEBUG {
		log.Println("File exists. Submitting to scanner....")
	}
	return SubmitFile(op)
}

func SubmitFile(op *cfg.UserOptions) (string, bool) {

	var RespData map[string]interface{}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", op.File)
	if err != nil {
		return err.Error(), false
	}
	file, err := os.Open(op.File)
	if err != nil {
		return err.Error(), false
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return err.Error(), false
	}

	writer.Close()

	headers := map[string]string{"Content-Type": writer.FormDataContentType(), "X-Apikey": op.Apikey}

	res, ok := cmn.SendRequest("POST", cfg.VTAPIS["FILE_SUBMITION_URL"], body, headers)
	if !ok {
		return res, ok
	}

	json.Unmarshal([]byte(res), &RespData)

	rdata := RespData["data"].(map[string]interface{})

	tmpdata := make(map[string]string)

	for key, value := range rdata {
		// Each value is an interface{} type, that is type asserted as a string
		tmpdata[key] = value.(string)
	}

	if tmpdata["type"] != "analysis" {
		return "", false
	}

	if cfg.DEBUG {
		log.Println(tmpdata)
		log.Println("File Submitted successfully. Analyzing Scan Report....")
	}

	return RetriveReport(tmpdata["id"], op)
}

func RetriveReport(analysis_id string, op *cfg.UserOptions) (errstr string, rflag bool) {

	var RespData map[string]map[string]map[string]interface{}

	headers := map[string]string{"X-Apikey": op.Apikey}
	res, ok := cmn.SendRequest("GET", fmt.Sprintf(cfg.VTAPIS["FILE_ANALYSIS_URL"], analysis_id), nil, headers)
	if !ok {
		return res, ok
	}

	json.Unmarshal([]byte(res), &RespData)

	rdata := RespData["data"]["attributes"]["stats"].(map[string]interface{})
	tmpdata := make(map[string]int)

	for key, value := range rdata {
		// Each value is an interface{} type, that is type asserted as a string
		tmpdata[key] = int(value.(float64))
	}

	if tmpdata["malicious"] > 0 {
		res = op.File + ": Malicious file!"
	} else if tmpdata["malicious"] == 0 && tmpdata["suspicious"] > 0 {
		res = op.File + ": Suspicious file. May contain a virus!"
	} else {
		res = op.File + ": Clean File!"
	}

	if cfg.DEBUG {
		log.Println(tmpdata)
		log.Println("Scan Report Analysis Completed.....")
	}
	return res, true
}
