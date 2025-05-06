package msgTransfer

import (
	"sync"
	"time"

	"github.com/wujunhui99/easy-chat/apps/im/ws/ws"
	"github.com/wujunhui99/easy-chat/pkg/constants"
	"github.com/zeromicro/go-zero/core/logx"
)

// 群聊消息已读的推送合并

type groupMsgRead struct {
	mu             sync.Mutex
	conversationId string
	push           *ws.Push      // 记录推送的消息
	pushCh         chan *ws.Push // 最终要推送的
	count          int
	// 上次推送时间
	pushTime time.Time
	done     chan struct{}
}

func newGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushCh:         pushCh,
		count:          1,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}

	// 这个位置挂一个 for 循环处理的，不会有问题吗
	go m.transfer()
	return m
}

func (m *groupMsgRead) transfer() {
	// 1，超时发送
	// 2. 超量发送

	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop()

	for {
		select {
		case <-m.done:
			return
		case <-timer.C:
			m.mu.Lock()

			pushTime := m.pushTime
			val := GroupMsgReadRecordDelayTime - time.Since(pushTime)
			push := m.push
			logx.Infof("timer.C %v val %v", time.Now(), val)
			if val > 0 && m.count < GroupMsgReadRecordDelayCount || push == nil {
				if val > 0 {
					timer.Reset(val)
				}

				// 未达标
				m.mu.Unlock()
				continue
			}

			m.pushTime = time.Now()
			m.push = nil
			m.count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			m.mu.Unlock()

			// 推送
			logx.Infof("超过 合并的条件推送 %v ", push)
			m.pushCh <- push
		default:
			m.mu.Lock()

			if m.count >= GroupMsgReadRecordDelayCount {
				push := m.push
				m.push = nil
				m.count = 0
				m.mu.Unlock()

				// 推送
				logx.Infof("default 超过 合并的条件推送 %v ", push)
				m.pushCh <- push
				continue
			}

			if m.isIdle() {
				m.mu.Unlock()
				// 使得msgReadTransfer释放
				m.pushCh <- &ws.Push{
					ChatType:       constants.GroupChatType,
					ConversationId: m.conversationId,
				}
				continue
			}
			m.mu.Unlock()

			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}
			time.Sleep(tempDelay)
		}
	}
}

// 合并消息
func (m *groupMsgRead) mergePush(push *ws.Push) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.count++
	// 消息id 和 消息读取记录
	for msgId, read := range push.ReadRecords {
		m.push.ReadRecords[msgId] = read
	}
}

// 判断是否为空消息
func (m *groupMsgRead) IsIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isIdle()
}

func (m *groupMsgRead) isIdle() bool {
	pushTime := m.pushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)

	if val <= 0 && m.push == nil && m.count == 0 {
		return true
	}

	return false
}

func (m *groupMsgRead) Clear() {
	select {
	case <-m.done:
	default:
		close(m.done)
	}

	m.push = nil
}
