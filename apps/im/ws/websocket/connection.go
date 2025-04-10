package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	*websocket.Conn
	s                 *Server
	idleMu                sync.Mutex
	idle              time.Time
	maxConnectionIdle time.Duration
	done              chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Errorf("Failed to upgrade connection: %v", err)
		return nil
	}
	conn :=  &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		done:              make(chan struct{}),
	}
	go conn.keepalive()
	return conn
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTicker(c.maxConnectionIdle)
	defer idleTimer.Stop()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			fmt.Printf("idle %v, maxIdle %c \n", idle,c.maxConnectionIdle)
			// if idle.IsZero() {
			// 	c.idleMu.Unlock()
			// 	idleTimer.Reset(c.maxConnectionIdle)
			// 	continue
			// }
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			fmt.Printf("val %v \n", val)
			if val <= 0 {
				c.s.Logger.Infof("Connection idle timeout, closing connection")
				c.Close()
				return
			}
			idleTimer.Reset(val);
			case <- c.done:
				fmt.Println("client connection finished")
				return
		}

	}
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Now()
	fmt.Printf("idle %v, maxIdle %c \n", c.idle,c.maxConnectionIdle)
	return
}
func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	// 方法是并不安全
	err := c.Conn.WriteMessage(messageType, data)
	c.idle = time.Now()
	return err
}
func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	return c.Conn.Close()

}