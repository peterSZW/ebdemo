package ebgame

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xiaomi-tc/log15"
)

type GSRsp struct {
	Code  int    `json:"code,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Msgex string `json:"msgex,omitempty"`
	GSRsp string `json:"status,omitempty"`
}

type GSRspSignup struct {
	GSRsp
	Token string `json:"token,omitempty"`
	Uuid  string `json:"uuid,omitempty"`
}

type GSRspSignin struct {
	GSRsp
	Token string `json:"token,omitempty"`
	Uuid  string `json:"uuid,omitempty"`
}

type GSRspSignout struct {
	GSRsp
}

type GSRspGetrooms struct {
	GSRsp
	Ids []string `json:"ids,omitempty"`
}

type GSRspJoinroom struct {
	GSRsp
	Roomid string `json:"roomid,omitempty"`
}

type GSRspJoinnewroom struct {
	GSRsp
	Roomid string `json:"roomid,omitempty"`
}

type GSRspLeaveroom struct {
	GSRsp
}

// type Rooms struct {
// 	Listeners        []string `json:"listeners"`
// 	ListenersCount   int      `json:"listeners_count"`
// 	Name             string   `json:"name"`
// 	Subscribers      []string `json:"subscribers"`
// 	SubscribersCount int      `json:"subscribers_count"`
// 	Type             string   `json:"type"`
// 	CreatedAt        int      `json:"created_at"`
// 	UpdatedAt        int      `json:"updated_at"`
// }

// type ClientResp struct {
// 	Roomss    []string `json:"channels"`
// 	ID        string   `json:"id"`
// 	Token     string   `json:"token"`
// 	CreatedAt int      `json:"created_at"`
// 	UpdatedAt int      `json:"updated_at"`
// }

var gs *GsClient

type GsClient struct {
	Token string
	URL   string
}

func init() {
	gs = NewGSConnect("", "http://127.0.0.1:7403")
}
func NewGSConnect(token string, url string) *GsClient {

	client := GsClient{
		Token: token,
		URL:   url,
	}

	return &client
}

func (c *GsClient) command(method string, url string, payload string) ([]byte, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var req *http.Request
	var err error
	fullurl := c.URL + url

	if payload == "" {
		req, err = http.NewRequest(method, fullurl, nil)
	} else {
		payday := strings.NewReader(payload)
		req, err = http.NewRequest(method, fullurl, payday)
	}

	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", c.Token)

	res, err := client.Do(req)
	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}
func (c *GsClient) GetRooms() (*GSRspGetrooms, error) {
	method := "POST"
	url := "/getrooms"

	body, err := c.command(method, url, "")

	if err != nil {
		return nil, err
	}
	var rsp *GSRspGetrooms
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	return rsp, nil
}

func (c *GsClient) Signup(acc, pass string) (*GSRspSignup, error) {
	method := "POST"
	url := "/signup"
	var rsp *GSRspSignup
	type GSRepSignup struct {
		Account string `json:"account,omitempty"`
		Pass    string `json:"pass,omitempty"`
		Debug   int    `json:"debug,omitempty"`
	}
	var req GSRepSignup
	req.Account = acc
	req.Debug = 1
	req.Pass = pass

	jsonreq, _ := json.Marshal(req)

	body, err := c.command(method, url, string(jsonreq))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	return rsp, nil
}
func (c *GsClient) Signin(acc, pass string) (*GSRspSignin, error) {
	method := "POST"
	url := "/signin"
	var rsp *GSRspSignin
	type GSRepSignin struct {
		Account string `json:"account,omitempty"`
		Pass    string `json:"pass,omitempty"`
		Debug   int    `json:"debug,omitempty"`
	}
	var req GSRepSignin
	req.Account = acc
	req.Debug = 1
	req.Pass = pass

	jsonreq, _ := json.Marshal(req)

	body, err := c.command(method, url, string(jsonreq))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	return rsp, nil
}

func (c *GsClient) Joinroom(token, roomid string) (*GSRspJoinroom, error) {
	method := "POST"
	url := "/joinroom"

	type GSReqJoinroom struct {
		Token  string `json:"token,omitempty"`
		Roomid string `json:"roomid,omitempty"`
		Debug  int    `json:"debug,omitempty"`
	}
	var req GSReqJoinroom
	req.Token = token
	req.Debug = 1
	req.Roomid = roomid

	jsonreq, _ := json.Marshal(req)

	body, err := c.command(method, url, string(jsonreq))

	if err != nil {
		return nil, err
	}
	var rsp *GSRspJoinroom
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	return rsp, nil
}

func (c *GsClient) Joinnewroom(token string) (*GSRspJoinnewroom, error) {
	method := "POST"
	url := "/joinnewroom"

	type GSReqJoinnewroom struct {
		Token string `json:"token,omitempty"`

		Debug int `json:"debug,omitempty"`
	}
	var req GSReqJoinnewroom
	req.Token = token
	req.Debug = 1

	jsonreq, _ := json.Marshal(req)

	body, err := c.command(method, url, string(jsonreq))

	if err != nil {
		return nil, err
	}
	var rsp *GSRspJoinnewroom
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	return rsp, nil
}

func (c *GsClient) Leaveroom(token, roomid string) (*GSRspLeaveroom, error) {
	method := "POST"
	url := "/leaveroom"
	var rsp *GSRspLeaveroom
	type GSReqLeaveroom struct {
		Token  string `json:"token,omitempty"`
		Roomid string `json:"roomid,omitempty"`
	}
	var req GSReqLeaveroom
	req.Token = token
	req.Roomid = roomid
	jsonreq, _ := json.Marshal(req)
	body, err := c.command(method, url, string(jsonreq))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	return rsp, nil
}

func (c *GsClient) HealthCheck() (*GSRsp, error) {
	method := "GET"
	url := "/_health"
	body, err := c.command(method, url, "")

	if err != nil {
		return nil, err
	}
	var status *GSRsp
	err = json.Unmarshal(body, &status)
	if err != nil {
		log15.Error("", "err", err, string(body), url)
		return nil, err
	}
	return status, nil
}

/*
func (c *GsClient) CreateConfig(key string, value string) (bool, error) {
	method := "POST"
	url := "/api/config"
	payload := "{\n	\"key\" : \"" + key + "\",\n	\"value\" : \"" + value + "\"\n}"

	_, err := c.command(method, url, payload)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *GsClient) GetConfig(key string) (*Config, error) {
	method := "GET"
	url := "/api/config/"

	body, err := c.command(method, url+key, "")
	if err != nil {
		return nil, err
	}
	log15.Debug("[", string(body), "]")
	var config *Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		log15.Error("", "err", err)
		return nil, err
	}
	return config, nil
}

func (c *GsClient) UpdateConfig(value string) (bool, error) {
	method := "PUT"
	url := "/api/config/"
	payload := "{\n	\"value\" : \"" + value + "\"\n}"
	_, err := c.command(method, url+value, payload)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *GsClient) DeleteConfig(key string) (bool, error) {
	method := "DELETE"
	url := "/api/config/"

	_, err := c.command(method, url+key, "")

	if err != nil {
		return false, err
	}
	return true, nil
}
*/

// func (c *GsClient) CreateRooms(channel string, ctype string) (bool, error) {
// 	method := "POST"
// 	url := "/api/channel"
// 	payload := "{\n	\"name\" : \"" + channel + "\",\n	\"type\" : \"" + ctype + "\"\n}"
// 	_, err := c.command(method, url, payload)

// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func (c *GsClient) UpdateRooms(channel string, ctype string) (bool, error) {
// 	method := "POST"
// 	url := "/api/channel/"
// 	payload := "{\n	\"type\" : \"" + ctype + "\"}"
// 	_, err := c.command(method, url+channel, payload)

// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// // has to be json string
// func (c *GsClient) PublishRooms(channel string, data string) (bool, error) {
// 	method := "POST"
// 	url := "/api/publish"
// 	quoteddata := strconv.Quote(data)
// 	payload := "{\n	\"channel\" : \"" + channel + "\",\n	\"data\" : " + quoteddata + "\n}"
// 	_, err := c.command(method, url, payload)

// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// // has to be json string
// func (c *GsClient) BroadcastRooms(channels []string, data string) (bool, error) {
// 	method := "POST"
// 	url := "/api/broadcast"
// 	urlsJSON, _ := json.Marshal(channels)
// 	quoteddata := strconv.Quote(data)
// 	payload := "{\n	\"channels\" : " + string(urlsJSON) + ",\n	\"data\" : " + quoteddata + "\n}"
// 	_, err := c.command(method, url, payload)

// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// func (c *GsClient) DeleteRooms(channel string) (bool, error) {
// 	method := "DELETE"
// 	url := "/api/channel/"

// 	body, err := c.command(method, url+channel, "")

// 	if err != nil {

// 		return false, err

// 	}
// 	log15.Error("", "err", string(body))
// 	return true, nil
// }

// func (c *GsClient) CreateClient(channel []string) (*ClientResp, error) {
// 	method := "POST"
// 	url := "/api/client"
// 	urlsJSON, _ := json.Marshal(channel)
// 	payload := "{\n	\"channels\" : " + string(urlsJSON) + "\n}"
// 	body, err := c.command(method, url, payload)

// 	if err != nil {
// 		return nil, err
// 	}
// 	log15.Error("", "err", string(body))
// 	var clientresp *ClientResp
// 	err = json.Unmarshal(body, &clientresp)
// 	if err != nil {
// 		log15.Error("", "err", err)
// 		return nil, err
// 	}
// 	return clientresp, nil
// }

// func (c *GsClient) GetClient(id string) (*ClientResp, error) {
// 	method := "GET"
// 	url := "/api/client/"

// 	body, err := c.command(method, url+id, "")

// 	if err != nil {
// 		return nil, err
// 	}
// 	var clientresp *ClientResp
// 	err = json.Unmarshal(body, &clientresp)
// 	if err != nil {
// 		log15.Error("", "err", err)
// 		return nil, err
// 	}
// 	return clientresp, nil
// }

// func (c *GsClient) SubscribeClient(channel []string, id string) (bool, error) {
// 	method := "PUT"
// 	url := "/api/client/"
// 	urlsJSON, _ := json.Marshal(channel)

// 	payload := "{\n	\"channels\" : " + string(urlsJSON) + "\n}"
// 	_, err := c.command(method, url+id+"/subscribe", payload)

// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
// func (c *GsClient) UnsubscribeClient(channel []string, id string) (bool, error) {
// 	method := "PUT"
// 	url := "/api/client/"
// 	urlsJSON, _ := json.Marshal(channel)
// 	payload := "{\n	\"channels\" : " + string(urlsJSON) + "\n}"
// 	_, err := c.command(method, url+id+"/unsubscribe", payload)

// 	if err != nil {
// 		return false, err
// 	}
// 	return true, err
// }

// func (c *GsClient) DeleteClient(id string) (bool, error) {
// 	method := "DELETE"
// 	url := "/api/client/"

// 	_, err := c.command(method, url+id, "")

// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// /* for future use */
// func (c *GsClient) Metrics() (bool, error) {
// 	method := "GET"
// 	url := "/api/metrics"
// 	_, err := c.command(method, url, "")
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// /* for future use */
// func (c *GsClient) Node() (bool, error) {
// 	method := "GET"
// 	url := "/api/node"
// 	_, err := c.command(method, url, "")
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
