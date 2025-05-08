package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Manage struct {
	client   *http.Client
	url      string
	user     string
	password string
}

func NewManage(url string, user string, password string) *Manage {
	m := &Manage{
		client:   &http.Client{},
		url:      "http://" + url,
		user:     user,
		password: password,
	}
	return m
}

func (m *Manage) ClearAllExchange(filterName string) error {
	size := 100
	resp, err := m.GetExchanges(filterName, 1, size)
	if err != nil {
		return err
	}
	countPage := resp.ItemCount / size
	if resp.ItemCount%size > 0 {
		countPage++
	}
	for page := 1; page <= countPage; page++ {
		if err != nil {
			return err
		}
		if err = m.DeleteExchanges(resp.Items); err != nil {
			return err
		}
		resp, err = m.GetExchanges(filterName, page, size)
	}

	return nil
}

func (m *Manage) ClearAllQueues(filterName string) error {
	size := 100
	resp, err := m.GetQueues(filterName, 1, size)
	if err != nil {
		return err
	}
	countPage := resp.ItemCount / size
	if resp.ItemCount%size > 0 {
		countPage++
	}
	for page := 1; page <= countPage; page++ {
		if err != nil {
			return err
		}
		if err = m.DeleteQueues(resp.Items); err != nil {
			return err
		}
		resp, err = m.GetQueues(filterName, page, size)
	}

	return nil
}

func (m *Manage) DeleteExchanges(items []*Exchange) error {
	for _, item := range items {
		if err := m.DeleteExchange(item.Vhost, item.Name); err != nil {
			fmt.Println("delete exchange:" + err.Error())
			if err.Error() != "404" {
				return err
			}
		}
	}
	return nil
}

func (m *Manage) GetExchanges(filterName string, page int, pageSize int) (*GetExchangeResponse, error) {
	api := fmt.Sprintf("/exchanges?page=%v&page_size=%v&name=%v&use_regex=false&pagination=true", page, pageSize, filterName)
	bytes, err := m.doGet(api, nil, func(r *http.Request) {

	})
	if err != nil {
		return nil, err
	}
	data := &GetExchangeResponse{}
	if err = json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (m *Manage) DeleteQueues(items []*Queue) error {
	for _, item := range items {
		if err := m.DeleteQueue(item.Vhost, item.Name); err != nil {
			fmt.Println("delete queue:" + err.Error())
			if err.Error() != "404" {
				return err
			}
		}
	}
	return nil
}

func (m *Manage) GetQueues(filterName string, page int, pageSize int) (*GetQueuesResponse, error) {
	api := fmt.Sprintf("/queues?page=%v&page_size=%v&name=%v&use_regex=false&pagination=true", page, pageSize, filterName)
	bytes, err := m.doGet(api, nil, func(r *http.Request) {})
	if err != nil {
		return nil, err
	}
	data := &GetQueuesResponse{}
	if err = json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (m *Manage) CreateExchange(req *CreateExchangeRequest) error {
	vhost := req.Vhost
	api := fmt.Sprintf("/exchanges/%s/%s", url.QueryEscape(vhost), url.QueryEscape(req.Name))
	_, err := m.doPut(api, req.GetData(), func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json")
	})
	return err
}

func (m *Manage) DeleteExchange(vhost, name string) error {
	if name == "" || strings.Contains(name, "amq") {
		return nil
	}
	api := fmt.Sprintf("/exchanges/%s/%s", url.QueryEscape(vhost), name)
	_, err := m.doDelete(api, nil, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json")
	})
	return err
}

func (m *Manage) DeleteQueue(vhost, name string) error {
	api := fmt.Sprintf("/queues/%s/%s", url.QueryEscape(vhost), name)
	_, err := m.doDelete(api, nil, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json")
	})
	return err
}

func (m *Manage) doPut(apiUrl string, data any, initReq func(req *http.Request)) ([]byte, error) {
	return m.do("PUT", apiUrl, data, initReq)
}

func (m *Manage) doGet(apiUrl string, data any, initReq func(req *http.Request)) ([]byte, error) {
	return m.do("GET", apiUrl, data, initReq)
}

func (m *Manage) doDelete(apiUrl string, data any, initReq func(req *http.Request)) ([]byte, error) {
	return m.do("DELETE", apiUrl, data, initReq)
}

func (m *Manage) do(methodType string, apiUrl string, data any, initReq func(req *http.Request)) ([]byte, error) {
	api := fmt.Sprintf("%v/api%v", m.url, apiUrl)

	var read io.Reader
	if data != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		str := string(bytes)
		read = strings.NewReader(str)
	}

	req, _ := http.NewRequest(methodType, api, read)
	req.SetBasicAuth(m.user, m.password)

	if initReq != nil {
		initReq(req)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 401 {
		fmt.Println(string(bytes))
	}
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 201 {
		return bytes, errors.New(fmt.Sprintf("%v", resp.Status))
	}

	return bytes, nil
}
