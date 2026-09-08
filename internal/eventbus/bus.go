package eventbus

import (
	"sync"
)

// EventBus 事件总线
type EventBus struct {
	handlers map[string][]func(interface{})
	mu       sync.RWMutex
}

var (
	bus  *EventBus
	once sync.Once
)

// New 创建新的事件总线
func New() *EventBus {
	return &EventBus{
		handlers: make(map[string][]func(interface{})),
	}
}

// GetBus 获取全局单例事件总线
func GetBus() *EventBus {
	once.Do(func() {
		bus = New()
	})
	return bus
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(event string, handler func(interface{})) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[event] = append(eb.handlers[event], handler)
}

// Publish 发布事件
func (eb *EventBus) Publish(event string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	handlers, ok := eb.handlers[event]
	if !ok {
		return
	}

	// 异步执行所有处理器
	for _, handler := range handlers {
		go handler(data)
	}
}
