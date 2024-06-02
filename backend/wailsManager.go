package manager

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"optio/backend/config"
	"optio/backend/stat"
)

type NotificationType int

const (
	Info NotificationType = iota
	Success
	Error
)

func (nt NotificationType) String() string {
	return [...]string{"info", "success", "error"}[nt]
}

// Manager handles collections of Files for conversion.
type Manager struct {
	ctx context.Context
}

// NewManager creates a new Manager.
func NewManager(c *config.Config, s *stat.Stat) *Manager {
	return &Manager{}
}

// WailsInit performs setup when Wails is ready.
func (fm *Manager) Startup(ctx context.Context) error {
	print("Manager startup")
	fm.ctx = ctx
	runtime.LogInfo(ctx, "Manager initialized...")
	return nil
}

func (fm *Manager) Notify(msg string, nt NotificationType) {
	runtime.EventsEmit(fm.ctx, "notify", map[string]interface{}{
		"msg":  msg,
		"type": nt.String(),
	})
}
