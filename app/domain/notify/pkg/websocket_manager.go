package pkg

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 使用 map[string]bool 模拟 Set，方便删除特定连接
type ConnectionSet map[*websocket.Conn]bool

type WebsocketManager struct {
	// 核心结构：UserID -> 多个连接
	clients map[string]ConnectionSet
	// 读写锁，保证并发安全
	lock sync.RWMutex
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		clients: make(map[string]ConnectionSet),
	}
}

// 定义 Upgrader
// 建议作为包级变量或者放在结构体中
var upgrader = websocket.Upgrader{
	// 1. 设置读取缓冲区大小（字节）
	ReadBufferSize: 1024,

	// 2. 设置写入缓冲区大小（字节）
	WriteBufferSize: 1024,

	// 3. 跨域检查（非常重要！）
	// 在开发环境下（例如前端在 localhost:3000，后端在 localhost:8080），
	// 浏览器会发起跨域连接。默认情况下 Upgrader 会拒绝。
	// 这里返回 true 表示允许所有来源连接。
	// 生产环境建议根据 r.Header.Get("Origin") 进行白名单校验。
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// RegisterClient 注册连接
func (m *WebsocketManager) RegisterClient(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	m.lock.Lock()
	if _, ok := m.clients[userID]; !ok {
		m.clients[userID] = make(ConnectionSet)
	}
	m.clients[userID][conn] = true
	m.lock.Unlock()

	// 监听连接断开，进行清理
	go m.listenForClose(userID, conn)
}

// listenForClose 监听连接关闭信号并清理内存
func (m *WebsocketManager) listenForClose(userID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		m.lock.Lock()
		if connections, ok := m.clients[userID]; ok {
			delete(connections, conn) // 从 Set 中删除该连接
			if len(connections) == 0 {
				delete(m.clients, userID) // 如果该用户无连接，删除 UserID Key
			}
		}
		m.lock.Unlock()
	}()

	// 阻塞读取，直到出错或断开
	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
}

// LocalPush 仅向当前服务实例中已连接的用户推送
func (m *WebsocketManager) LocalPush(userID string, msg interface{}) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if connections, ok := m.clients[userID]; ok {
		for conn := range connections {
			// 对该用户的每个连接（每个 Tab/设备）都发送消息
			go conn.WriteJSON(msg)
		}
	}
}
