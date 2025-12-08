package event

// TaskEvent 定义了在 Dapr 中传输的消息结构
type TaskEvent struct {
	UserID     string `json:"user_id"`
	WorkflowID string `json:"workflow_id"`
	Status     string `json:"status"` // e.g., "COMPLETED", "FAILED"
	Message    string `json:"message"`
	ResultData string `json:"result_data,omitempty"`
}

const (
	PubSubName = "notify-pubsub"
	TopicName  = "task.status"
)
