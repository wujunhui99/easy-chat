package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	sync.RWMutex
	authentication Authentication
	routers  map[string]HandlerFunc
	addr     string
	connToUser map[*Conn]string
	userToConn map[string]*Conn
	upgrader websocket.Upgrader
	logx.Logger
	opt *serverOption
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routers:  make(map[string]HandlerFunc),
		authentication: opt.Authentication,
		addr:     addr,
		upgrader: websocket.Upgrader{},
		Logger:   logx.WithContext(context.Background()),
		opt:      &opt,
		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServerWs")
	defer func() {
		if err := recover(); err != nil {
			s.Logger.Errorf("Error: %v", err)
		}
	}()
	fmt.Println("before auth")
	if !s.authentication.Auth(w,r) {
		s.Info("Authentication failed")
		return 
	}
	fmt.Println("Authentication passed")
	// conn, err := s.upgrader.Upgrade(w, r, nil)
	conn := NewConn(s, w, r)
	if conn == nil {
		s.Logger.Errorf("Failed to upgrade connection")
		return
	}
	s.addConn(conn, r)
	go s.handlerConn(conn)

}
func (s *Server) addConn( conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	fmt.Println("uid is ", uid)
	if c := s.userToConn[uid]; c != nil {
		s.Logger.Infof("User %s already connected, closing old connection", uid)
		c.Close()
	}
	s.Logger.Infof("User %s connected", uid)
	s.connToUser[conn] = uid
	s.userToConn[uid] = conn

}
func (s *Server) handlerConn(conn *Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Logger.Errorf("Failed to read message: %v", err)
			s.Close(conn)
			return
		}
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			s.Send(NewErrMessage(err),conn);
		}

		switch message.FrameType {
		case FramePing:
			s.Send(&Message{FrameType: FramePing},conn)
		case FrameData:
			if handler, ok := s.routers[message.Method]; ok {
				handler(s, conn, &message)
			} else {
				s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在方法 %v \n",message.Method)}, conn)
			}
	}
	}
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routers[r.Method] = r.Handler
	}
}
func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.userToConn[uid]
}
func (s *Server) GetConns(ids... string) []*Conn {
	if(len(ids) == 0) {
		return nil
	}
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	res := make([]*Conn, 0, len(ids))
	for _, id := range ids {
		if conn, ok := s.userToConn[id]; ok {
			res = append(res, conn)
		}
	}
	return res

}
func (s *Server) GetUsers(conns ...*Conn) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var res []string
	if len(conns) == 0 {
		res = make([]string, 0, len(s.userToConn))
		for uid := range s.userToConn {
			res = append(res, uid)
		}
	}else{
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}
	return res
}
func (s *Server) Close(conn *Conn) {
	
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	uid := s.connToUser[conn]
	if uid == ""{
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

}
func (s *Server) SendByUserId(msg interface{},sendIds ... string) error {
	if( len(sendIds) == 0) {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds ...)...)
}
func (s *Server) Send(msg interface{}, conns...*Conn) error {
	if(len(conns) == 0) {
		return nil
	}
	data, err := json.Marshal(msg)
	if err!= nil {
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
