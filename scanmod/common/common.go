package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func Usage() {
	_, file, _, _ := runtime.Caller(0)

	log.Println("Improper parameter used!")
	log.Printf("USAGE:\ngo run " + file + " -f <file_to_scan> -a <apikey> \nOR\ngo run " + file + " -u <url_to_scan> -a <apikey>\n")
	log.Printf("Example:\ngo run " + file + " -f /home/test.txt -a kjsdjknusckjsdcec\nOR\ngo run " + file + " -u www.test.com -a kjsdjknusckjsdcec\n")
	os.Exit(1)
}

func IsFileExists(filename string) (errstr string, rflag bool) {
	errstr = ""
	rflag = true
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		errstr = filename + ": Does not exists!"
		rflag = false
	} else if err != nil {
		errstr = filename + ": Something went wrong!"
		rflag = false
	}

	return
}

func SendRequest(rtype string, url string, data io.Reader, headers map[string]string) (errstr string, rflag bool) {
	errstr = ""
	rflag = false

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(rtype, url, data)
	if err != nil {
		errstr = fmt.Sprintf("Got error %s", err.Error())
		return
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		errstr = fmt.Sprintf("Got error %s", err.Error())
		return
	}

	defer func() {
		e := response.Body.Close()
		if e != nil {
			log.Fatal(e)
		}
	}()

	bytes, errRead := ioutil.ReadAll(response.Body)

	if errRead != nil {
		errstr = fmt.Sprintf("Failed to read response %s", err.Error())
		return
	}

	return string(bytes), true
}
