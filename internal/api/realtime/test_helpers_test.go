package realtime

import (
	"testing"

	"github.com/fedora-oss/javinizer-go/internal/api/testkit"
	ws "github.com/fedora-oss/javinizer-go/internal/websocket"
)

var wsHub *ws.Hub

func initTestWebSocket(t *testing.T) {
	testkit.InitTestWebSocket(t)
	wsHub = testkit.CurrentHub()
}
