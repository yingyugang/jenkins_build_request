package main

import (
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
						case "battle_server_release":
							url = "http://127.0.0.1:8080/job/climbers-battleserver-release/buildWithParameters?token=11db88c014135c00b7c5066c73c8ee9478"
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
						//data := `{}`
						fmt.Printf("url:%v\n", url)

						username := "yingyugang"
						password := "11db88c014135c00b7c5066c73c8ee9478"
						// 创建一个 HTTP 客户端
						client := &http.Client{}

						// 创建一个 HTTP 请求
						req, err3 := http.NewRequest("POST", url, nil)
						if err3 != nil {
							fmt.Println("Error creating HTTP request:", err3)
							return
						}

						// 设置用户名和密码
						req.SetBasicAuth(username, password)

						// 发送请求
						resp, err4 := client.Do(req)
						if err4 != nil {
							fmt.Println("Error sending HTTP request:", err4)
							return
						}
						defer resp.Body.Close()

						fmt.Println("Response status:", resp.Status)

					}
					//curl -v -X GET http://127.0.0.1:8080/crumbIssuer/api/json --user yingyugang:810412
					//-H 'Jenkins-Crumb: 0db38413bd7ec9e98974f5213f7ead8b'
					//-L --user your-user-name:apiToken
					//curl -X POST -L --user yingyugang:11db88c014135c00b7c5066c73c8ee9478 http://127.0.0.1:8080/job/unity-climber-client-ios/41/stop
					//curl -X POST http://127.0.0.1:8080/job/unity-climber-client-ios/40/stop --user yingyugang:810412 -H 'Jenkins-Crumb: 4cd88ceaa16bae965a1c73a709c525cb48c9f2a993495eb47979a7e6a7a24721'
				}
			}

		}
	}
}
