package livyserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	http2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/http"
)

type LivyServer struct {
	client *http2.Client
	ip     string
	port   int
	url    string
}

// NewLivyServer
//
//	@Description: Livy
//	@param url "http://192.168.120.200:8998/batches"
func NewLivyServer(ip string, port int) *LivyServer {
	return &LivyServer{
		client: http2.NewClient(),
		ip:     ip,
		port:   port,
		url:    fmt.Sprintf("http://%s:%d", ip, port),
	}
}

func (s *LivyServer) GetURL() string {
	return s.url
}

/*
func (s *LivyServer) httpDelete(url string) int {
	// 2. 创建自定义的 http.Client（建议设置超时时间）
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 3. 构建 DELETE 请求对象
	// 如果需要发送请求体，可以将最后一个参数替换为 bytes.NewBuffer(jsonData)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
	}

	// 4. 设置请求头（根据 API 需求可选）
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer YOUR_TOKEN_HERE")

	// 5. 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return
	}
	// 务必记得关闭 Body
	defer resp.Body.Close()

	// 6. 读取并打印响应状态码与内容
	fmt.Printf("响应状态码: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体失败: %v\n", err)
		return
	}

	fmt.Printf("响应内容: %s\n", string(body))
}

func (s *LivyServer) httpGet(url string, response any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, response); err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}

func (s *LivyServer) httpPost(url string, request any, response any) (status *http.Result, err error) {
	// 2. 将结构体序列化为 JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("JSON 序列化失败: %v\n", err)
		return nil, err
	}

	// 3. 发送 POST 请求到 Livy
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("请求 Livy 失败: %v\n", err)
		return nil, err
	}
	// 4. 读取并解析返回结果
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Livy 报错，状态码: %d, 原因: %s\n", resp.StatusCode, string(body))
		return nil, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("解析响应 JSON 失败: %v\n", err)
		return nil, err
	}
	return resp, nil
}
*/

// CreateBatch
//
//	@Description: 批量导入
//	@receiver s
//	@param ctx
//	@param request CreateTaskResult
//	@return error
func (s *LivyServer) CreateBatch(request CreateBatchRequest) (*CreateTaskResult, error) {
	livyURL := s.url + "/batches"

	// 2. 将结构体序列化为 JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("JSON 序列化失败: %v\n", err)
		return nil, err
	}

	// 3. 发送 POST 请求到 Livy
	resp, err := http.Post(livyURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("请求 Livy 失败: %v\n", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 4. 读取并解析返回结果
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Livy 报错，状态码: %d, 原因: %s\n", resp.StatusCode, string(body))
		return nil, err
	}

	var livyResp CreateTaskResult
	if err = json.Unmarshal(body, &livyResp); err != nil {
		fmt.Printf("解析响应 JSON 失败: %v\n", err)
		return nil, err
	}

	fmt.Printf("成功提交 Spark 任务！分配的 Batch ID: %d, 当前状态: %s\n", livyResp.ID, livyResp.State)

	// 5. 异步轮询任务状态（可选流程）
	go s.trackStatus(livyResp.ID)

	// 阻塞主线程方便查看演示结果
	time.Sleep(60 * time.Second)
	return &livyResp, err
}

func (s *LivyServer) GetSessions(from, size int) (*GetSessionResult, error) {
	url := fmt.Sprintf("%s/sessions", s.GetURL())
	var resp GetSessionResult
	_, err := s.client.Get(url, nil, &resp)
	return &resp, err
}

func (s *LivyServer) GetBatchState(batchId int) (*GetBatchStateResult, error) {
	url := fmt.Sprintf("%s/batches/%d/state", s.GetURL(), batchId)
	var resp GetBatchStateResult
	_, err := s.client.Get(url, nil, &resp)
	return &resp, err
}

func (s *LivyServer) DeleteBatch(batchId int) error {
	url := fmt.Sprintf("%s/batches/%d", s.GetURL(), batchId)
	_, err := s.client.Get(url, nil, nil)
	return err
}

func (s *LivyServer) GetBatchId(batchId int) (*GetBatchResult, error) {
	statusURL := fmt.Sprintf("%s/batches/%d", s.GetURL(), batchId)
	var resp GetBatchResult
	_, err := s.client.Get(statusURL, nil, &resp)
	if err != nil {
		return nil, errors.New("获取状态失败: %v\n", err)
	}

	return &resp, err
}

func (s *LivyServer) GetBatchLogs(batchId int, size int) (*GetBatchLogsResult, error) {
	statusURL := fmt.Sprintf("%s/batches/%d/log?size=%d", s.GetURL(), batchId, size)
	var resp GetBatchLogsResult
	_, err := s.client.Get(statusURL, nil, &resp)
	if err != nil {
		return nil, errors.New("获取状态失败: %v\n", err)
	}
	return &resp, nil
}

// trackStatus 根据批处理 ID 轮询任务状态
func (s *LivyServer) trackStatus(batchId int) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		res, err := s.GetBatchId(batchId)
		if err != nil {
			fmt.Printf("获取状态失败: %v\n", err)
			continue
		}

		fmt.Printf("BatchID [%d] AppID [%s] 当前状态: %s\n", batchId, res.AppId, res.State)
		for _, logInfo := range res.Log {
			fmt.Println(logInfo)
		}
		// 如果状态达到终态则退出轮询
		if res.State == "success" || res.State == "dead" || res.State == "killed" {
			fmt.Printf("任务结束，最终状态: %s\n", res.State)
			break
		}
	}
}
