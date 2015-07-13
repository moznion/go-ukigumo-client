package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// TODO followings should be kicked out to common library
const (
	StatusSuccess = "1"
	StatusFail    = "2"
	StatusNA      = "3"
	StatusSkip    = "4"
	StatusPending = "5"
	StatusTimeout = "6"
)

type Client struct {
	VC             VC
	ServerURL      string
	WorkDir        string
	Project        string
	CompareURL     string
	ElapsedTimeSec int64

	revFrom string
}

func NewClient(serverURL string, vc VC) *Client {
	c := new(Client)

	c.ServerURL = serverURL
	c.VC = vc

	c.WorkDir = ""
	c.Project = ""
	c.CompareURL = ""

	return c
}

func (c *Client) Run() {
	origDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Chdir(origDir)

	vc := c.VC
	vcInfo := vc.GetVCInfo()

	if c.WorkDir == "" {
		var err error
		c.WorkDir, err = ioutil.TempDir(os.TempDir(), "")
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	workDir := path.Join(c.WorkDir, c.Project, vcInfo.Branch)

	log.Printf("start testing : %s", vcInfo.Description)
	log.Printf("working directory : %s", workDir)

	os.MkdirAll(workDir, 0777)
	os.Chdir(workDir)

	c.revFrom = vc.GetRevision()
	c.VC.Update()
	currentRev := vcInfo.revision

	if vcInfo.SkipIfUnmodified && c.revFrom == currentRev {
		log.Printf("Skip testing")
		return
	}

	conf := NewYAMLConfig()

	conf.applyEnvironmentVariables()

	if conf.ProjectName != "" {
		c.Project = conf.ProjectName
	}

	// TODO notify here

	c.ElapsedTimeSec = int64(0)
	r := NewRunner(conf, &c.ElapsedTimeSec)

	r.Run("BeforeInstall")
	r.Run("Install")
	r.Run("BeforeScript")

	// TODO support to switch executor
	var executor Executor = NewCommandExecutor(conf.Script)
	status := executor.Exec()
	log.Printf("Finished testing: %s", status)

	r.Run("AfterScript")

	log.Printf("End testing")
}

func (c *Client) submitResult(status string, logFilename string) {
	c.sendResultToServer(status, logFilename)
	// TODO notify here
}

type responseJSON struct {
	report report `json:"report"`
}

type report struct {
	url        string `json:"url"`
	lastStatus string `json:"last_status"`
}

func (c *Client) sendResultToServer(status string, logFilename string) (string, string) {

	vc := c.VC
	vcInfo := vc.GetVCInfo()

	revFrom := c.revFrom
	currentRev := vcInfo.revision

	vcLog, _ := vc.GetLog(revFrom, currentRev)

	serverURL := strings.TrimRight(c.ServerURL, "/")
	log.Printf("Sending result to server at %s (status: %s)", serverURL, status)
	res, err := http.PostForm(serverURL+"/api/v1/report/add", url.Values{
		"project":          {c.Project},
		"branch":           {vcInfo.Branch},
		"repo":             {vcInfo.Repository},
		"revision":         {currentRev[:10]},
		"status":           {status},
		"vc_log":           {vcLog},
		"body":             {logFilename},
		"compare_url":      {c.CompareURL},
		"elapsed_time_sec": {strconv.FormatInt(c.ElapsedTimeSec, 10)},
	})
	if err != nil {
		log.Fatal(err)
		return "", ""
	}

	resStatusCode := res.StatusCode
	if !(200 <= resStatusCode && resStatusCode <= 299) {
		log.Fatalf("Failed to send result to server (status: %d)", resStatusCode)
	}

	var dat responseJSON
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	json.Unmarshal(buf.Bytes(), &dat)

	reportURL := dat.report.url
	log.Printf("report url: %s", reportURL)
	return reportURL, dat.report.lastStatus
}
