package bootstrap

import (
	"github.com/gofiber/fiber"
	"time"
)

type Configuration func(*Bootstrap)

type Bootstrap struct {
	*fiber.App
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time
}

func New(appName, appOwner string, cfgs ...Configuration) *Bootstrap {
	b := &Bootstrap{
		App:          fiber.New(),
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
	}
	for _, cfg := range cfgs {
		cfg(b)
	}
	return b
}

func (b *Bootstrap) Configure(cfgs ...Configuration) {
	for _, cfg := range cfgs {
		cfg(b)
	}
}
