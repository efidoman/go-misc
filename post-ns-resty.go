package main

import (
	"flag"
	"fmt"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"log"
	"os"
	"time"
)

//
// exit codes
// -1 ==> INPUT file doesn't exist
//  0 ==> success
// -2 ==> err on http request

var curlStatus int = -1

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s inputjsonfile type\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "          inputjsonfile = Nightscout json record file\n")
	fmt.Fprintf(os.Stderr, "          type = Nightscout record type, default is \"entries\"\n")
	flag.PrintDefaults()
	os.Exit(curlStatus)
}

func main() {

	flag.Parse()
	flag.Usage = usage

	var nsUrl string = os.Getenv("NIGHTSCOUT_HOST")
	var nsSecret string = os.Getenv("API_SECRET")

	//fmt.Fprintf(os.Stderr, "nsUrl=%s, nsSecret=%s\n", nsUrl, nsSecret)

	if flag.NArg() < 2 {
		usage()
	}
	//fmt.Fprintf(os.Stderr, "arg0=%s\n", flag.Arg(0))
	//fmt.Fprintf(os.Stderr, "arg1=%s\n", flag.Arg(1))

	jsonFile := flag.Arg(0)
	nsType := flag.Arg(1)

	b, err := ioutil.ReadFile(jsonFile) // just pass the file name
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/api/v1/%s/.json", nsUrl, nsType)

	resty.SetHTTPMode()
	resty.SetTimeout(time.Duration(10 * time.Second))

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("API-SECRET", nsSecret).
		SetBody(b).
		Post(url)

	if err != nil {
		curlStatus = -2
	}

	fmt.Printf("\n%v\n", resp)

	os.Exit(curlStatus)
}
