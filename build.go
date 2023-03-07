package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Empty struct {
}

type BuildPhases struct {
	Builds []*BuildPhase
}

type BuildPhase struct {
	Build  string
	Target string
	Brunch string
}

func main() {
	ticker := time.Tick(1 * time.Second)
	go RequestBuild(ticker)
	select {}
}

func RequestBuild(ch <-chan time.Time) {
	for {
		select {
		case <-ch:
			resp, err := http.Get("http://13.231.44.197:8080/get_build")
			if err != nil {
				fmt.Printf("RequestBuild err:%v\n", err)
			}
			content, err1 := ioutil.ReadAll(resp.Body)
			if err1 != nil {
				fmt.Printf("RequestBuild err:%v\n", err1)
			}
			fmt.Printf("RequestBuild:%v\n", string(content))
			var baseMessage BuildPhases
			err2 := json.Unmarshal(content, &baseMessage)
			if err2 != nil {
				fmt.Printf("RequestBuild err:%v\n", err2)
			}
			if baseMessage.Builds != nil {
				for i := 0; i < len(baseMessage.Builds); i++ {
					var build = baseMessage.Builds[i]
					var url = ""
					if build.Build == "compile" {
						switch build.Target {
						case "ios":
							url = "http://127.0.0.1:8080/job/unity-climber-client-ios/buildWithParameters?token=11db88c014135c00b7c5066c73c8ee9478&brunch=" + build.Brunch
						case "battle_server":
							url = "http://127.0.0.1:8080/job/climbers-battleserver/buildWithParameters?token=11db88c014135c00b7c5066c73c8ee9478"
						}
						if build.Target == "ios" {
							url = "http://127.0.0.1:8080/job/unity-climber-client-ios/buildWithParameters?token=11db88c014135c00b7c5066c73c8ee9478&brunch=" + build.Brunch
						}
						resp1, err3 := http.Get(url)
						if err3 != nil {
							fmt.Printf("RequestBuild err:%v\n", err3)
						}
						content1, err4 := ioutil.ReadAll(resp1.Body)
						if err4 != nil {
							fmt.Printf("RequestBuild err:%v\n", err4)
						}
						fmt.Printf("RequestBuild:%v\n", string(content1))
					} else if build.Build == "stop" {
						switch build.Target {
						case "ios":
							url = "http://127.0.0.1:8080/job/unity-climber-client-ios/" + build.Brunch + "/stop"
						}
						data := `{}`
						fmt.Printf("url:%v\n", url)
						resp1, err3 := http.Post(url, "application/json", bytes.NewBufferString(data))
						if err3 != nil {
							fmt.Printf("RequestBuild err:%v\n", err3)
						}
						content1, err4 := ioutil.ReadAll(resp1.Body)
						if err4 != nil {
							fmt.Printf("RequestBuild err:%v\n", err4)
						}
						fmt.Printf("RequestBuild:%v\n", string(content1))
					}

				}
			}

		}
	}
}
