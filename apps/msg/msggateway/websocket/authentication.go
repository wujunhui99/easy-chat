package websocket

import (
	"fmt"
	"net/http"
	"time"
)

type Authentication interface {
	Auth(w http.ResponseWriter, r *http.Request) bool
	UserId(r *http.Request) string
	DeviceId(r *http.Request) string
}
type authentication struct{}

func (*authentication) Auth(w http.ResponseWriter, r *http.Request) bool {
	return true
}
func (*authentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}

	return fmt.Sprintf("%v", time.Now().UnixMilli())
}

func (*authentication) DeviceId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["deviceId"] != nil {
		return fmt.Sprintf("%v", query["deviceId"])
	}
	return ""
}
