package utils

import "sync"

// StatusTracker 定义了一个状态跟踪器，使用泛型 T 来存储特定类型的状态。
// 泛型 T 必须是可比较的，以确保可以检查状态是否发生变化。
type StatusTracker[T comparable] struct {
	statusMap map[string]T // 存储状态的映射，键为 ApplicationID。
	mu        sync.Mutex   // 使用互斥锁保证并发安全
}

// NewStatusTracker 创建并返回一个新的 StatusTracker 实例，初始化状态映射。
// 返回值:
//
//	*StatusTracker[T] - 新创建的 StatusTracker 实例。
func NewStatusTracker[T comparable]() *StatusTracker[T] {
	return &StatusTracker[T]{
		statusMap: make(map[string]T),
	}
}

// UpdateStatus 试图更新给定 ApplicationID 的状态。
// 如果指定的 ApplicationID 的当前状态不存在或与新状态不同，则更新状态，并返回 true。
// 如果当前状态存在且与新状态相同，则不进行更新，返回 false。
//
// 参数:
//
//	key string - 状态跟踪的 ApplicationID。
//	newStatus T - 新的状态值。
//
// 返回值:
//
//	bool - 表示状态是否有变化。
//
// 示例:
//
//	tracker := NewStatusTracker[int]()
//	changed := tracker.UpdateStatus("app123", 1)
//	if changed {
//	    fmt.Println("状态更新成功")
//	} else {
//	    fmt.Println("状态未变化")
//	}
func (st *StatusTracker[T]) UpdateStatus(key string, newStatus T) bool {
	st.mu.Lock()         // 修改map前加
	defer st.mu.Unlock() // 在方法结束时解锁

	currentStatus, exists := st.statusMap[key]
	if !exists || currentStatus != newStatus {
		st.statusMap[key] = newStatus // 更新状态
		return true                   // 状态改变或者是新的状态
	}
	return false // 没有变化
}
