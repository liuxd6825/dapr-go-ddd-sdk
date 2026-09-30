package store_huge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
)

// AppendVertexLabelProperties 向已存在的 VertexLabel 追加 properties
// action=append 端点;若 label 不存在返回 error
func AppendVertexLabelProperties(ctx context.Context, client *hugegraph.CommonClient, name string, properties []string) error {
	if client == nil {
		return fmt.Errorf("client is nil")
	}
	if len(properties) == 0 {
		return nil
	}
	body, _ := json.Marshal(map[string]any{"name": name, "properties": properties})
	cfg := client.Transport.GetConfig()
	url := fmt.Sprintf("%s://%s/graphs/%s/schema/vertexlabels/%s?action=append",
		cfg.URL.Scheme, cfg.URL.Host, cfg.Graph, name)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.Username != "" {
		req.SetBasicAuth(cfg.Username, cfg.Password)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}()
	if resp.StatusCode == http.StatusAccepted ||
		resp.StatusCode == http.StatusCreated ||
		resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusConflict {
		return nil
	}
	rb, _ := io.ReadAll(resp.Body)
	bodyStr := string(rb)
	if containsCI(bodyStr, "already exists") || containsCI(bodyStr, "has existed") {
		return nil
	}
	return fmt.Errorf("append vertexlabel %s properties status=%d body=%s", name, resp.StatusCode, bodyStr)
}

// AppendEdgeLabelProperties 向已存在的 EdgeLabel 追加 properties
func AppendEdgeLabelProperties(ctx context.Context, client *hugegraph.CommonClient, name string, properties []string) error {
	if client == nil {
		return fmt.Errorf("client is nil")
	}
	if len(properties) == 0 {
		return nil
	}
	body, _ := json.Marshal(map[string]any{"name": name, "properties": properties})
	cfg := client.Transport.GetConfig()
	url := fmt.Sprintf("%s://%s/graphs/%s/schema/edgelabels/%s?action=append",
		cfg.URL.Scheme, cfg.URL.Host, cfg.Graph, name)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.Username != "" {
		req.SetBasicAuth(cfg.Username, cfg.Password)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}()
	if resp.StatusCode == http.StatusAccepted ||
		resp.StatusCode == http.StatusCreated ||
		resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusConflict {
		return nil
	}
	rb, _ := io.ReadAll(resp.Body)
	bodyStr := string(rb)
	if containsCI(bodyStr, "already exists") || containsCI(bodyStr, "has existed") {
		return nil
	}
	return fmt.Errorf("append edgelabel %s properties status=%d body=%s", name, resp.StatusCode, bodyStr)
}
