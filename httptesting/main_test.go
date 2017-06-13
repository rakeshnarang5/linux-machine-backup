package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"string"
	"testing"
)

const (
	multiServerConnect = "multi-server-resiproc/"
	multiServer        = "multi-server/"
	helperScriptPath   = "helperScripts/build"
)

var (
	urls                   = make(map[string]string)
	connectBuildNames      = []string{"multi-server", "caseta"}
	casetaBuilsNames       = []string{"multi-server-connect", "connect"}
	localBuildNames        = []string{"system", "computer", "local"}
	crossCompileBuildNames = []string{"machine", "cross platform", "cross compile"}
	goSrcPath              string
)

func init() {
	urls["local"] = "http://localhost:5000/system"
	urls["remote"] = "https://shrouded-beyond-17924.herokuapp.com/system"
	goSrcPath = os.Getenv("GOPATH") + "/src/"
}

type Parameters struct {
	Status    int
	Code      string `json:"build-source-type"`
	BuildType string `json:"build-type"`
}

func TestHandleRequest(t *testing.T) {
	params := &Parameters{
		Code:      "multi-server-connect",
		BuildType: "machine",
	}
	handleRequest(params)
}

func handleRequest(p *Parameters) {
	path := scriptPath(p.Code)
	//buildFlag := buildFlag(p.BuildType)
	fmt.Printf("path: %v\n", path)

	cmd := exec.Command(path, "-a")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.Path = path
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to read serial number. Run error %v, StdErr %v, StdOut %v", err, stderr.String(), stdout.String())
	}

}

func scriptPath(code string) string {
	for _, names := range connectBuildNames {
		if str {
			return goSrcPath + multiServerConnect + helperScriptPath
		}
	}
	return goSrcPath + multiServerConnect + helperScriptPath
}

func buildFlag(buildType string) []string {
	for _, name := range localBuildNames {
		if name == buildType {
			return nil
		}
	}
	return []string{"-a"}
}
