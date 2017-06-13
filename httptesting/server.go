package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"httptesting/telnet"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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
	telnetConn             *telnet.Conn
)

func init() {
	urls["local"] = "http://localhost:5000/system"
	urls["remote"] = "https://shrouded-beyond-17924.herokuapp.com/system"
	goSrcPath = os.Getenv("GOPATH") + "/src/"
}

type request struct {
	Result *Result `json:"result"`
}

type Result struct {
	ServiceType   string
	ResolvedQuery string      `json:"resolvedQuery"`
	MetaData      *MetaData   `json:"metadata"`
	Parameters    *Parameters `json:"parameters"`
}

type MetaData struct {
	IntentName string `json:"intentName"`
}

type Parameters struct {
	Status       int
	Code         string `json:"build-source-type"`
	BuildType    string `json:"build-type"`
	Action       string
	Location     string
	AlexaMessage string
	SceneName    string
}

type Fulfillment struct {
	Speech      string `json:"speech"`
	Data        string `json:"data"`
	DisplayText string `json:"displayText"`
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
	conn, err := telnet.NewConnection()
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
	}
	telnetConn = conn

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
		req := &request{}
		err = json.Unmarshal(body, req)
		if err != nil {
			fmt.Printf("ERROR:%v\n", err)
		}
		// if params.Status == 0 {
		// 	go handleRequest(params)
		// }
		switch req.Result.ServiceType {
		case "Google Home":
			ffm := handleRequestGoogleHome(req.Result)
			ffmbytes, err := json.Marshal(ffm)
			if err != nil {
				fmt.Printf("got error while marhalling the json %v\n", err)
				ffmbytes = []byte("sorry not able to fulfill this request")
			}
			_, err = http.Post(serverurl, "application/json", bytes.NewBuffer(ffmbytes))
			if err != nil {
				fmt.Printf("ERROR:%v\n", err)
			}
		case "Alexa":
			go handleRequestAlexa(req.Result.Parameters)
		default:
			fmt.Println("nothing to process")
		}
		fmt.Printf("Req:%+v\n", req)
		time.Sleep(1 * time.Second)
	}
}

func handleRequestGoogleHome(res *Result) *Fulfillment {
	switch res.MetaData.IntentName {
	case "LightAssistant":
		return handleLightAssistantIntent(res.Parameters)
	case "build-helper":
		return handleBuildIntent(res.Parameters)
	default:
		return &Fulfillment{
			Speech: "This intent in not supported yet",
		}
	}
}

func handleLightAssistantIntent(params *Parameters) *Fulfillment {
	if telnetConn == nil || !telnetConn.Connected() {
		conn, err := telnet.NewConnection()
		if err != nil {
			fmt.Printf("ERROR:%v\n", err)
			return &Fulfillment{
				Speech: "Sorry not able to fulfill request at this time",
			}
		}
		telnetConn = conn
	}
	err := handleLightAssistant(params.Location, params.Action)
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
		return &Fulfillment{
			Speech: "Sorry not able to fulfill request at this time",
		}
	}
	return &Fulfillment{
		Speech: "Your " + params.Location + " lights are " + params.Action,
	}
}

func handleLightAssistant(location, action string) error {
	fmt.Printf("Handling location :%v , action: %v\n", location, action)
	cmd, err := formIntegrationCommand(location, action)
	if err != nil {
		return err
	}
	if telnetConn == nil || !telnetConn.Connected() {
		conn, err := telnet.NewConnection()
		if err != nil {
			return err
		}
		telnetConn = conn
	}
	fmt.Printf("Writing: %s\v", cmd)
	_, err = telnetConn.Write([]byte(cmd))
	if err != nil {
		return err
	}
	return nil
}

func formIntegrationCommand(location, action string) (string, error) {
	var level string
	switch action {
	case "On":
		level = "100\r\n"
	case "Off":
		level = "0\r\n"
	}
	switch location {
	case "Kitchen":
		return "#OUTPUT,107,1," + level, nil
	case "Bedroom":
		return "#OUTPUT,95,1," + level, nil
	case "All":
		return "#OUTPUT,107,1," + level + "#OUTPUT,95,1," + level, nil
	}

	return "", fmt.Errorf("format Sorry this is not supported")
}

func handleBuildIntent(params *Parameters) *Fulfillment {
	go func() {
		path := scriptPath(params.Code)
		buildFlag := buildFlag(params.BuildType)
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
	}()
	return &Fulfillment{
		Speech: "your " + params.Code + " is compiled for " + params.BuildType,
	}
}

func handleRequestAlexa(p *Parameters) {
	fmt.Println(p.AlexaMessage)
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
