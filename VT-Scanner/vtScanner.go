package main

import (
	"log"

	cmn "modules/scanmod/common"
	cfg "modules/scanmod/config"
	scn "modules/scanmod/scanner"
)

func ScanFile(op *cfg.UserOptions) {
	data, _ := scn.FS(op)
	log.Println(data)
	return
}

func ScanUrl(op *cfg.UserOptions) {
	log.Println(op.File)
}

func main() {
	log.Println("VT Scanner is Running.....")
	op := &cfg.UserOptions{}
	op.GetInputs()
	if op.Apikey == "" {
		cmn.Usage()
	}

	switch {
	case op.File != "":
		ScanFile(op)

	case op.Url != "":
		ScanUrl(op)

	default:
		cmn.Usage()
	}
	log.Println("====================== SCAN COMPLETED ===================")
}
