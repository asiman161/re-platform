package replatform

import "net/http"

func (i *Implementation) Ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong"))
}
