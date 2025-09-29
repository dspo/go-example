package app0

import (
	"go.uber.org/fx"

	"gitee.com/huajinet/go-example/pkg/engine"
)

type ApplicationContext struct {
	fx.In

	Engine *engine.Engine
}
