package api

import (
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/kavanahuang/config"
	"github.com/kavanahuang/log"
	"github.com/mozillazg/request"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	contentType = "application/json; charset=utf-8"
	Timeout     = 30 * time.Second
)

type client struct {
	conn   *http.Client
	buffer *bytes.Buffer
	Uri    string
}

var Http = new(client)

func (c *client) Post(url string, data any) *simplejson.Json {
	conn := new(http.Client)
	req := request.NewRequest(conn)

	req.Json = data
	req.Client.Timeout = Timeout
	response, err := req.Post(url)
	if err != nil {
		log.Logs.Error("Post error: ", err)
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Logs.Error("Response body close error: ", err)
		}
	}()

	var result *simplejson.Json
	result, err = response.Json()
	if err != nil {
		log.Logs.Error("Response context error: ", err)
	}

	statusCode, err := result.Get("Code").Int()
	if err != nil {
		log.Logs.Error("Get code error: ", err)
	}

	if statusCode == 200 {
		return result.Get("Data")
	}

	msg, err := result.Get("Msg").String()
	if err != nil {
		log.Logs.Error("Get msg error: ", err)
	}

	log.Logs.Warning("Response code: ", strconv.Itoa(statusCode)+", msg: "+msg)
	return nil
}

func (c *client) NewJsonClient(url string) *client {
	c.Uri = url
	return c
}

func (c *client) NewFormDataClient(url string) *client {
	contentType = "form-data, charset=utf-8"
	c.Uri = url
	return c
}

func (c *client) NewFormEncodeClient(url string) *client {
	contentType = "application/x-www-form-urlencoded, charset=utf-8"
	c.Uri = url
	return c
}

func (c *client) PostString(data string, stu any) (a any) {
	req := bytes.NewBufferString(data)
	res, err := http.Post(c.Uri, contentType, req)
	if err != nil {
		log.Logs.Error("Post error: ", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Logs.Error("Read body error: ", err)
	}

	defer func() { _ = res.Body.Close() }()

	err = json.Unmarshal(body, &stu)
	if err != nil {
		log.Logs.Error("Unmarshal response body error: ", err)
	}

	return
}

func (c *client) PostByte(data []byte, stu any) any {
	req := bytes.NewBuffer(data)
	res, err := http.Post(c.Uri, contentType, req)
	if err != nil {
		log.Logs.Error("Post error: ", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Logs.Error("Read body error: ", err)
	}

	defer func() { _ = res.Body.Close() }()

	err = json.Unmarshal(body, &stu)
	if err != nil {
		log.Logs.Error("Unmarshal response body error: ", err)
	}

	return stu
}
