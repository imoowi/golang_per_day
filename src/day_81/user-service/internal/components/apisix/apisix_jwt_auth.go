package apisix

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"user-service/internal/global"
)

type ApisixJwt struct {
	Addr           string
	XApiKey        string
	ConsumerKey    string
	ConsumerSecret string
	Timeout        time.Duration
}

func NewImoowiApisix() *ApisixJwt {
	return &ApisixJwt{
		Addr:           global.Config.GetString("apisix.addr"),
		XApiKey:        global.Config.GetString("apisix.xapikey"),
		ConsumerKey:    global.Config.GetString("apisix.consumer_key"),
		ConsumerSecret: global.Config.GetString("apisix.consumer_secret"),
		Timeout:        global.Config.GetDuration("apisix.timeout"),
	}
}

func (a *ApisixJwt) CreateConsumer(name string) error {
	// 调用apisix创建consumer
	url := `http://` + a.Addr + `/apisix/admin/consumers`
	method := "PUT"
	payload := strings.NewReader(`{
		"username": "` + name + `"
	}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {

		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X_API_KEY", a.XApiKey)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
func (a *ApisixJwt) CreateConsumerCredentials(name string, jwtAuthKey string, jwtAuthSecret string) error {
	// 调用apisix创建consumer credentials
	url := `http://` + a.Addr + `/apisix/admin/consumers/` + name + `/credentials`
	method := "PUT"
	payload := strings.NewReader(`{
		"id": "cred-` + name + `-jwt-auth",
		"plugins": {
			"jwt-auth": {
				"key": "` + a.ConsumerKey + `",
				"secret": "` + a.ConsumerSecret + `"
			}
		}
	}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {

		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X_API_KEY", a.XApiKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}

func (a *ApisixJwt) CreateRoute(name string, uri string, jwtAuthKey string) error {
	// 调用apisix创建route
	url := `http://` + a.Addr + `/apisix/admin/routes/` + name
	method := "PUT"
	payload := strings.NewReader(`{
		"uri": "` + uri + `",
		"plugins": {
			"key-auth": {
				"key": "` + a.ConsumerKey + `"
			}
		}
	}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {

		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X_API_KEY", a.XApiKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
