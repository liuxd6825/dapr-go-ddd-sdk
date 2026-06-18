package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

const companyIndex = "company"

type ElasticClient struct {
	client *elasticsearch.Client
}

func NewElasticClient() (*ElasticClient, error) {
	c := &ElasticClient{}
	c.init()
	return c, nil
}

func (c *ElasticClient) init() {
	_env := env.GetEnv()
	c.client = _env.GetDB("elastic").GetElastic()
}

func (c *ElasticClient) Ping(ctx context.Context) error {
	res, err := c.client.Ping(c.client.Ping.WithContext(ctx))
	if err != nil {
		return errors.NewErr(err, "Ping failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("Ping error: status=%d, body=%s", res.StatusCode, string(body)))
	}
	return nil
}

func (c *ElasticClient) Info(ctx context.Context) (map[string]interface{}, error) {
	res, err := c.client.Info(c.client.Info.WithContext(ctx))
	if err != nil {
		return nil, errors.NewErr(err, "Info failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, errors.New(fmt.Sprintf("Info error: status=%d, body=%s", res.StatusCode, string(body)))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, errors.NewErr(err, "Decode Info response failed")
	}
	return result, nil
}

func (c *ElasticClient) CreateIndex(ctx context.Context, index string, mappings map[string]interface{}) error {
	body, err := json.Marshal(map[string]interface{}{
		"mappings": mappings,
	})
	if err != nil {
		return errors.NewErr(err, "Marshal mappings failed")
	}

	res, err := c.client.Indices.Create(
		index,
		c.client.Indices.Create.WithContext(ctx),
		c.client.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return errors.NewErr(err, "CreateIndex failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("CreateIndex error: index=%s, status=%d, body=%s", index, res.StatusCode, string(body)))
	}
	return nil
}

func (c *ElasticClient) DeleteIndex(ctx context.Context, index string) error {
	res, err := c.client.Indices.Delete(
		[]string{index},
		c.client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return errors.NewErr(err, "DeleteIndex failed")
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("DeleteIndex error: index=%s, status=%d, body=%s", index, res.StatusCode, string(body)))
	}
	return nil
}

func (c *ElasticClient) IndexExists(ctx context.Context, index string) (bool, error) {
	res, err := c.client.Indices.Exists(
		[]string{index},
		c.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, errors.NewErr(err, "IndexExists failed")
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}

func (c *ElasticClient) Index(ctx context.Context, index, id string, document interface{}) error {
	body, err := json.Marshal(document)
	if err != nil {
		return errors.NewErr(err, "Marshal document failed")
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "false",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return errors.NewErr(err, "Index document failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("Index error: index=%s, id=%s, status=%d, body=%s", index, id, res.StatusCode, string(body)))
	}
	return nil
}

func (c *ElasticClient) Get(ctx context.Context, index, id string, result interface{}) error {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return errors.NewErr(err, "Get document failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == http.StatusNotFound {
			return errors.New(fmt.Sprintf("document not found: index=%s, id=%s", index, id))
		}
		body, _ := io.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("Get error: index=%s, id=%s, status=%d, body=%s", index, id, res.StatusCode, string(body)))
	}

	var response struct {
		Source_ json.RawMessage `json:"_source"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return errors.NewErr(err, "Decode Get response failed")
	}

	if result != nil && response.Source_ != nil {
		if err := json.Unmarshal(response.Source_, result); err != nil {
			return errors.NewErr(err, "Unmarshal document failed")
		}
	}
	return nil
}

func (c *ElasticClient) Delete(ctx context.Context, index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return errors.NewErr(err, "Delete document failed")
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("Delete error: index=%s, id=%s, status=%d, body=%s", index, id, res.StatusCode, string(body)))
	}
	return nil
}

func (c *ElasticClient) Search(ctx context.Context, index string, query string) (*SearchResult, error) {
	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(index),
		c.client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, errors.NewErr(err, "Search failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, errors.New(fmt.Sprintf("Search error: index=%s, status=%d, body=%s", index, res.StatusCode, string(body)))
	}

	var response struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source_ json.RawMessage `json:"_source"`
				Id_     string          `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.NewErr(err, "Decode Search response failed")
	}

	result := &SearchResult{
		Total: response.Hits.Total.Value,
		Hits:  make([]map[string]interface{}, 0, len(response.Hits.Hits)),
	}

	for _, hit := range response.Hits.Hits {
		var doc map[string]interface{}
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			continue
		}
		doc["_id"] = hit.Id_
		result.Hits = append(result.Hits, doc)
	}

	return result, nil
}

func (c *ElasticClient) SearchWithDSL(ctx context.Context, index string, query io.Reader) (*SearchResult, error) {
	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(index),
		c.client.Search.WithBody(query),
	)
	if err != nil {
		return nil, errors.NewErr(err, "SearchWithDSL failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, errors.New(fmt.Sprintf("SearchWithDSL error: index=%s, status=%d, body=%s", index, res.StatusCode, string(body)))
	}

	var response struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source_ json.RawMessage `json:"_source"`
				Id_     string          `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.NewErr(err, "Decode Search response failed")
	}

	result := &SearchResult{
		Total: response.Hits.Total.Value,
		Hits:  make([]map[string]interface{}, 0, len(response.Hits.Hits)),
	}

	for _, hit := range response.Hits.Hits {
		var doc map[string]interface{}
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			continue
		}
		doc["_id"] = hit.Id_
		result.Hits = append(result.Hits, doc)
	}

	return result, nil
}

func (c *ElasticClient) Bulk(ctx context.Context, index string, documents []map[string]interface{}) (int, error) {
	if len(documents) == 0 {
		return 0, nil
	}

	var buf bytes.Buffer
	for _, doc := range documents {
		id, _ := doc["_id"].(string)
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": index,
			},
		}
		if id != "" {
			meta["index"].(map[string]interface{})["_id"] = id
		}

		metaBytes, _ := json.Marshal(meta)
		buf.Write(metaBytes)
		buf.WriteByte('\n')

		docBytes, _ := json.Marshal(doc)
		buf.Write(docBytes)
		buf.WriteByte('\n')
	}

	res, err := c.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		c.client.Bulk.WithContext(ctx),
		c.client.Bulk.WithIndex(index),
	)
	if err != nil {
		return 0, errors.NewErr(err, "Bulk failed")
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return 0, errors.New(fmt.Sprintf("Bulk error: status=%d, body=%s", res.StatusCode, string(body)))
	}

	var bulkResponse struct {
		Errors bool `json:"errors"`
		Items  []struct {
			Index struct {
				Error *struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"error,omitempty"`
			} `json:"index"`
		} `json:"items"`
	}

	if err := json.NewDecoder(res.Body).Decode(&bulkResponse); err != nil {
		return 0, errors.NewErr(err, "Decode Bulk response failed")
	}

	if bulkResponse.Errors {
		var errCount int
		for _, item := range bulkResponse.Items {
			if item.Index.Error != nil {
				errCount++
			}
		}
		return errCount, errors.New(fmt.Sprintf("Bulk partially failed: %d errors", errCount))
	}

	return 0, nil
}

func (c *ElasticClient) Close() error {
	if c.client != nil {
		ctx := context.Background()
		return c.client.Close(ctx)
	}
	return nil
}

func (c *ElasticClient) GetClient() *elasticsearch.Client {
	return c.client
}

type SearchResult struct {
	Total int64
	Hits  []map[string]interface{}
}
