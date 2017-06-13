package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"ss.go"
	"strings"
	"time"
)

const (
	multiServerConnect = "multi-server-resiproc/"
	multiServer        = "multi-server-casseta/"
	helperScriptPath   = "helperScripts/build"
)

var (
	urls                   = make(map[string]string)
	casetaBuildNames       = []string{"multi-server", "caseta"}
	connectBuildNames      = []string{"multi-server-connect", "connect"}
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
	ServiceType  string
	Status       int
	Code         string `json:"build-source-type"`
	BuildType    string `json:"build-type"`
	AlexaMessage string
}

//http://localhost:5000
func main() {
	reader := bufio.NewReader(os.Stdin)
	var serverurl string
	fmt.Print("Enter the  address of client: remote or local\n")
	for {
		add, _ := reader.ReadString('\n')
		add = strings.TrimSpace(add)
		url, ok := urls[add]
		if ok {
			serverurl = url
			break
		}
	}

	for {
		resp, err := http.Get(serverurl)
		if err != nil {
			fmt.Printf("ERROR:%v\n", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ERROR:%v\n", err)
		}
		fmt.Printf("Body:%s\n", body)
		resp.Body.Close()
		params := &Parameters{}
		err = json.Unmarshal(body, params)
		if err != nil {
			fmt.Printf("ERROR:%v\n", err)
		}
		// if params.Status == 0 {
		// 	go handleRequest(params)
		// }
		switch params.ServiceType {
		case "Google Home":
			go handleRequestGoogleHome(params)
		case "Alexa":
			go handleRequestAlexa(params)
		default:
			fmt.Println("nothing to process")
		}
		fmt.Printf("Parameters:%+v\n", params)
		time.Sleep(1 * time.Second)
	}
}

func handleRequestGoogleHome(p *Parameters) {
	path := scriptPath(p.Code)
	buildFlag := buildFlag(p.BuildType)
	fmt.Printf("path: %v\n", path)

	cmd := exec.Command(path, buildFlag...)
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

func handleRequestAlexa(p *Parameters) {
	fmt.Println(p.AlexaMessage)
	conn, err := ssh.Dial("tcp", server, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	defer conn.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := conn.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var w io.Writer
	session.Stdout = w
	w.Write(p.AlexaMessage)

}

func scriptPath(code string) string {
	for _, names := range connectBuildNames {
		fmt.Printf("origional: %v, Current: %v\n", code, names)
		if code == names {
			return goSrcPath + multiServerConnect + helperScriptPath
		}
	}
	return goSrcPath + multiServer + helperScriptPath
}

func buildFlag(buildType string) []string {
	for _, name := range localBuildNames {
		if name == buildType {
			return nil
		}
	}
	return []string{"-a"}
}
