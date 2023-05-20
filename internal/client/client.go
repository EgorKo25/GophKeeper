package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/EgorKo25/GophKeeper/internal/storage"
)

type Client struct {
	urlServer string
}

func NewClient(urlServer string) *Client {

	return &Client{
		urlServer: urlServer,
	}
}

func (c *Client) Send(src any, dataType string, cookie []*http.Cookie, path string) (int, any, []*http.Cookie, error) {

	var data []byte
	var err error

	client := &http.Client{}

	data, err = c.myMarshal(src)
	if err != nil {
		return 0, nil, nil, err
	}

	req, err := http.NewRequest("POST", c.urlServer+path, bytes.NewBuffer(data))
	if err != nil {
		return 0, nil, nil, err
	}

	for _, cook := range cookie {
		req.AddCookie(cook)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Data-Type", dataType)

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}

	body, _ := io.ReadAll(resp.Body)

	res, err := c.myUnmarshal(resp.Header.Get("Data-Type"), body)
	if err != nil {
		return 0, nil, nil, err
	}

	return resp.StatusCode, res, resp.Cookies(), nil
}

func (c *Client) myMarshal(body any) ([]byte, error) {
	switch t := body.(type) {
	case *storage.Card:
		res, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		return res, nil
	case *storage.Password:
		res, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		return res, nil
	case *storage.BinaryData:
		res, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		return res, nil
	case *storage.User:
		res, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, errors.New("unknown type")
}

func (c *Client) myUnmarshal(t string, body []byte) (any, error) {

	if len(body) == 0 {
		return nil, nil
	}

	switch t {
	case "card":
		res := storage.Card{}
		err := json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "password":
		res := storage.Password{}
		err := json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "bin":
		res := storage.BinaryData{}
		err := json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, errors.New("unknown type")
}
