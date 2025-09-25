package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}

	return "NoAck"
}

type Server struct {
	sync.RWMutex
	*threading.TaskRunner
	authentication Authentication
	routes         map[string]HandlerFunc
	addr           string
	connMeta       map[*Conn]connIdentity
	userConns      map[string]map[string]*Conn
	upgrader       websocket.Upgrader
	logx.Logger
	opt *serverOption
}

type connIdentity struct {
	uid      string
	deviceId string
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routes:         make(map[string]HandlerFunc),
		authentication: opt.Authentication,
		addr:           addr,
		upgrader:       websocket.Upgrader{},
		Logger:         logx.WithContext(context.Background()),
		opt:            &opt,
		connMeta:       make(map[*Conn]connIdentity),
		userConns:      make(map[string]map[string]*Conn),
		TaskRunner:     threading.NewTaskRunner(opt.concurrency),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("ServerWs")
	defer func() {
		if err := recover(); err != nil {
			s.Logger.Errorf("Error: %v", err)
		}
	}()
	if !s.authentication.Auth(w, r) {
		s.Info("Authentication failed")
		return
	}
	// conn, err := s.upgrader.Upgrade(w, r, nil)
	conn := NewConn(s, w, r)
	if conn == nil {
		s.Logger.Errorf("Failed to upgrade connection")
		return
	}
	s.addConn(conn, r)
	go s.handlerConn(conn)

}
func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)
	deviceId := s.authentication.DeviceId(req)
	if deviceId == "" {
		deviceId = fmt.Sprintf("%s:%d", uid, time.Now().UnixNano())
	}
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	if s.userConns[uid] == nil {
		s.userConns[uid] = make(map[string]*Conn)
	}
	if c := s.userConns[uid][deviceId]; c != nil {
		s.Logger.Infof("User %s device %s already connected, closing old connection", uid, deviceId)
		c.Close()
	}
	s.Logger.Infof("User %s device %s connected", uid, deviceId)
	conn.Uid = uid
	conn.DeviceId = deviceId
	s.connMeta[conn] = connIdentity{uid: uid, deviceId: deviceId}
	s.userConns[uid][deviceId] = conn

}

// 任务的处理
func (s *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done:
			// 连接关闭
			return
		case message := <-conn.message:
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{FrameType: FramePing}, conn)
			case FrameData:
				// 根据请求的method分发路由并执行
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)}, conn)
					//conn.WriteMessage(&Message{}, []byte(fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method)))
				}
			}

			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}
func (s *Server) readAck(conn *Conn) {
	for {
		select {
		case <-conn.done:
			s.Infof("close message ack uid %v ", conn.Uid)
			return
		default:
		}
		fmt.Println("read ack message...")
		// 从队列中读取新的消息
		conn.messageMu.Lock()
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			fmt.Println("")
			// 增加睡眠
			time.Sleep(1 * time.Second)
			continue
		}

		// 读取第一条
		message := conn.readMessage[0]
		fmt.Println("ack handler...")
		// 判断ack的方式
		switch s.opt.ack {
		case OnlyAck:
			// 直接给客户端回复
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			// 进行业务处理
			// 把消息从队列中移除
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()

			conn.message <- message
		case RigorAck:
			// 先回
			if message.AckSeq == 0 {
				// 还未确认
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].ackTime = time.Now()
				s.Send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)
				s.Infof("message ack RigorAck send mid %v, seq %v , time%v", message.Id, message.AckSeq,
					message.ackTime)
				conn.messageMu.Unlock()
				continue
			}

			// 再验证

			// 1. 客户端返回结果，再一次确认
			// 得到客户端的序号
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 确认
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				s.Infof("message ack RigorAck success mid %v", message.Id)
				continue
			}

			// 2. 客户端没有确认，考虑是否超过了ack的确认时间
			val := s.opt.ackTimeout - time.Since(message.ackTime)
			if !message.ackTime.IsZero() && val <= 0 {
				//		2.2 超过结束确认
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}
			//		2.1 未超过，重新发送
			conn.messageMu.Unlock()
			s.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)
			// 睡眠一定的时间
			time.Sleep(3 * time.Second)
		}
	}
}
func (s *Server) handlerConn(conn *Conn) {
	// 处理任务
	go s.handlerWrite(conn)

	if s.isAck(nil) {
		fmt.Println("s.isAck", s.isAck(nil))
		go s.readAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			s.Close(conn)
			return
		}
		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(msg))
			s.Close(conn)
			return
		}

		// 依据消息进行处理
		if s.isAck(&message) {
			s.Infof("conn message read ack msg %v", message)
			conn.appendMsgMq(&message)
		} else {
			conn.message <- &message
		}
	}
}

func (s *Server) isAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck
}
func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}
func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	if conns := s.userConns[uid]; conns != nil {
		for _, conn := range conns {
			return conn
		}
	}
	return nil
}

func (s *Server) GetUserConns(uid string) []*Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	if conns := s.userConns[uid]; conns != nil {
		res := make([]*Conn, 0, len(conns))
		for _, conn := range conns {
			res = append(res, conn)
		}
		return res
	}
	return nil
}
func (s *Server) GetConns(ids ...string) []*Conn {
	if len(ids) == 0 {
		return nil
	}
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	res := make([]*Conn, 0)
	for _, id := range ids {
		if conns, ok := s.userConns[id]; ok {
			for _, conn := range conns {
				res = append(res, conn)
			}
		}
	}
	return res

}
func (s *Server) GetUsers(conns ...*Conn) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var res []string
	if len(conns) == 0 {
		res = make([]string, 0, len(s.userConns))
		for uid := range s.userConns {
			res = append(res, uid)
		}
	} else {
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			if meta, ok := s.connMeta[conn]; ok {
				res = append(res, meta.uid)
			}
		}
	}
	return res
}
func (s *Server) Close(conn *Conn) {

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	meta, ok := s.connMeta[conn]
	if !ok {
		return
	}

	delete(s.connMeta, conn)
	if conns := s.userConns[meta.uid]; conns != nil {
		delete(conns, meta.deviceId)
		if len(conns) == 0 {
			delete(s.userConns, meta.uid)
		}
	}

}
func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds...)...)
}
func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for _, conn := range conns {
		conn.WriteMessage(websocket.TextMessage, data)
	}
	return nil
}
func (s *Server) Start() {
	http.HandleFunc("/ws", s.ServerWs)
	http.ListenAndServe(s.addr, nil)
}

func (s *Server) Stop() {
	// Implement graceful shutdown logic if needed
	fmt.Println("Stopping server...")
}
