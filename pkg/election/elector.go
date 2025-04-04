package election

import (
	"context"
	daprsdk "github.com/dapr/go-sdk/client"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type LeaderElector struct {
	client      daprsdk.Client
	storeName   string
	instanceID  string
	isLeader    bool
	stopChan    chan struct{}
	lockTTL     time.Duration // 锁持有时间
	renewPeriod time.Duration // 续期间隔
	resourceID  string        // 锁资源标识

	healthCheckURL string      // 健康检查端点（可选）
	lastHealthOK   atomic.Bool // 原子操作的健康状态
}

func NewLeaderElector(instanceID string, lockTTL time.Duration, storeName string, client daprsdk.Client) (*LeaderElector, error) {
	return &LeaderElector{
		client:      client,
		storeName:   storeName,
		instanceID:  instanceID,
		lockTTL:     lockTTL,
		renewPeriod: lockTTL / 2, // 续期间隔为TTL的一半
		resourceID:  "leader-election-lock",
		stopChan:    make(chan struct{}),
	}, nil
}

func (l *LeaderElector) Run(ctx context.Context) {
	ticker := time.NewTicker(l.renewPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			l.releaseLockIfLeader()
			return

		case <-ticker.C:
			if !l.isLeader {
				l.tryAcquireLock(ctx)
			} else {
				// 通过重新获取锁实现续期
				resp, err := l.client.TryLockAlpha1(ctx, l.storeName, &daprsdk.LockRequest{
					ResourceID:      l.resourceID,
					LockOwner:       l.instanceID,
					ExpiryInSeconds: int32(l.lockTTL.Seconds()),
				})

				if err != nil || resp == nil || !resp.Success {
					l.loseLeadership()
				}
			}
		}
	}
}

// 尝试获取锁
func (l *LeaderElector) tryAcquireLock(ctx context.Context) {
	resp, err := l.client.TryLockAlpha1(ctx, l.storeName, &daprsdk.LockRequest{
		ResourceID:      l.resourceID,
		LockOwner:       l.instanceID,
		ExpiryInSeconds: int32(l.lockTTL.Seconds()),
	})

	if err != nil {
		log.Printf("锁获取失败: %v", err)
		return
	}

	if resp != nil && resp.Success {
		log.Println("当选为Leader")
		l.isLeader = true
		go l.startLeaderTasks()
	}
}

// 失去领导权处理
func (l *LeaderElector) loseLeadership() {
	log.Println("失去Leader身份")
	l.isLeader = false
	close(l.stopChan)
	l.stopChan = make(chan struct{})
}

// 释放锁（优雅退出时）
func (l *LeaderElector) releaseLockIfLeader() {
	if l.isLeader {
		_, err := l.client.UnlockAlpha1(context.Background(), l.storeName, &daprsdk.UnlockRequest{
			ResourceID: l.resourceID,
			LockOwner:  l.instanceID,
		})
		if err != nil {
			log.Printf("锁释放失败: %v", err)
		}
		l.loseLeadership()
	}
}

func (l *LeaderElector) startLeaderTasks() {
	// 示例任务：处理发件箱消息
	go l.processOutboxMessages()

	// 示例任务：定期健康检查
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.doHealthCheck()
		case <-l.stopChan:
			return
		}
	}
}

func (l *LeaderElector) processOutboxMessages() {
	// 实现发件箱处理逻辑
	for l.isLeader {
		select {
		case <-l.stopChan:
			return
		default:
			// 处理消息...
			time.Sleep(1 * time.Second)
		}
	}
}

// 新增健康检查方法
func (l *LeaderElector) doHealthCheck() {
	// 示例1：简单自检
	l.lastHealthOK.Store(true) // 假设默认健康

	// 示例2：实际HTTP检查（如检查数据库连接）
	if l.healthCheckURL != "" {
		_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := http.Get(l.healthCheckURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			l.lastHealthOK.Store(false)
			log.Printf("健康检查失败: %v", err)
			return
		}
		l.lastHealthOK.Store(true)
	}

	// 示例3：检查关键资源（如发件箱表）
	if err := l.checkDatabase(); err != nil {
		l.lastHealthOK.Store(false)
		log.Printf("数据库检查失败: %v", err)
	}
}

func (l *LeaderElector) checkDatabase() error {
	// 实现实际数据库检查逻辑
	_, err := l.client.GetState(context.Background(), "statestore", "health-check", nil)
	return err
}
