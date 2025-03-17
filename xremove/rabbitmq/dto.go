package rabbitmq

type GetExchangeResponse struct {
	FilteredCount int         `json:"filtered_count"`
	ItemCount     int         `json:"item_count"`
	Items         []*Exchange `json:"items"`
}

type Exchange struct {
	Arguments struct {
	} `json:"arguments"`
	AutoDelete             bool   `json:"auto_delete"`
	Durable                bool   `json:"durable"`
	Internal               bool   `json:"internal"`
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	UserWhoPerformedAction string `json:"user_who_performed_action"`
	Vhost                  string `json:"vhost"`
}

type GetQueuesResponse struct {
	FilteredCount int      `json:"filtered_count"`
	ItemCount     int      `json:"item_count"`
	Page          int      `json:"page"`
	PageCount     int      `json:"page_count"`
	TotalCount    int      `json:"total_count"`
	Items         []*Queue `json:"items"`
}

type Queue struct {
	Arguments                     map[string]string `json:"arguments"`
	AutoDelete                    bool              `json:"auto_delete"`
	BackingQueueStatus            map[string]any    `json:"backing_queue_status"`
	ConsumerCapacity              any               `json:"consumer_capacity"`
	ConsumerUtilisation           any               `json:"consumer_utilisation"`
	Consumers                     any               `json:"consumers"`
	Durable                       bool              `json:"durable"`
	EffectivePolicyDefinition     map[string]any    `json:"effective_policy_definition"`
	Exclusive                     bool              `json:"exclusive"`
	ExclusiveConsumerTag          interface{}       `json:"exclusive_consumer_tag"`
	GarbageCollection             map[string]any    `json:"garbage_collection"`
	HeadMessageTimestamp          interface{}       `json:"head_message_timestamp"`
	IdleSince                     string            `json:"idle_since"`
	Memory                        int               `json:"memory"`
	MessageBytes                  int               `json:"message_bytes"`
	MessageBytesPagedOut          int               `json:"message_bytes_paged_out"`
	MessageBytesPersistent        int               `json:"message_bytes_persistent"`
	MessageBytesRam               int               `json:"message_bytes_ram"`
	MessageBytesReady             int               `json:"message_bytes_ready"`
	MessageBytesUnacknowledged    int               `json:"message_bytes_unacknowledged"`
	Messages                      int               `json:"messages"`
	MessagesDetails               map[string]any    `json:"messages_details"`
	MessagesPagedOut              int               `json:"messages_paged_out"`
	MessagesPersistent            int               `json:"messages_persistent"`
	MessagesRam                   int               `json:"messages_ram"`
	MessagesReady                 int               `json:"messages_ready"`
	MessagesReadyDetails          map[string]any    `json:"messages_ready_details"`
	MessagesReadyRam              int               `json:"messages_ready_ram"`
	MessagesUnacknowledged        int               `json:"messages_unacknowledged"`
	MessagesUnacknowledgedDetails map[string]any    `json:"messages_unacknowledged_details"`
	MessagesUnacknowledgedRam     int               `json:"messages_unacknowledged_ram"`
	Name                          string            `json:"name"`
	Node                          string            `json:"node"`
	OperatorPolicy                interface{}       `json:"operator_policy"`
	Policy                        interface{}       `json:"policy"`
	RecoverableSlaves             interface{}       `json:"recoverable_slaves"`
	Reductions                    int               `json:"reductions"`
	ReductionsDetails             map[string]any    `json:"reductions_details"`
	SingleActiveConsumerTag       any               `json:"single_active_consumer_tag"`
	State                         string            `json:"state"`
	Type                          string            `json:"type"`
	Vhost                         string            `json:"vhost"`
}

type CreateExchangeRequest struct {
	Vhost      string            `json:"vhost"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Durable    bool              `json:"durable"`
	AutoDelete bool              `json:"auto_delete"`
	Internal   bool              `json:"internal"`
	Arguments  map[string]string `json:"arguments"`
}

type createExchangeRequestData struct {
	Vhost      string            `json:"vhost"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Durable    string            `json:"durable"`
	AutoDelete string            `json:"auto_delete"`
	Internal   string            `json:"internal"`
	Arguments  map[string]string `json:"arguments"`
}

func (r *CreateExchangeRequest) GetData() *createExchangeRequestData {
	arg := r.Arguments
	if arg == nil {
		arg = make(map[string]string)
	}
	return &createExchangeRequestData{
		Vhost:      r.Vhost,
		Name:       r.Name,
		Type:       r.Type,
		Durable:    boolToStr(r.Durable),
		Internal:   boolToStr(r.Internal),
		AutoDelete: boolToStr(r.AutoDelete),
		Arguments:  arg,
	}
}

func boolToStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
