package definition

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gitoa.ru/go-4devs/config/definition"
	"gitoa.ru/go-4devs/config/definition/group"
	"gitoa.ru/go-4devs/config/definition/option"
	"gitoa.ru/go-4devs/config/definition/proto"
	"gitoa.ru/go-4devs/log/level"
)

func Configure(ctx context.Context, def *definition.Definition) error {
	def.Add(
		option.String("test", "test description", option.Default("defult")),
		option.Int("int", "test int description", option.Default(2)),
		group.New("group", "group description", option.String("test", "test description")),
		option.Time("start", "start at", option.Default(time.Now())),
		option.Duration("timer", "timer", option.Default(time.Hour)),
		group.New("log", "logger",
			option.New("level", "log level", level.Level(0), option.Default(level.Debug)),
			option.New("logrus", "logrus level", logrus.Level(0), option.Default(logrus.DebugLevel)),
			proto.New("sevice", "cutom service log", option.New("level", "log level", level.Level(0), option.Default(level.Debug))),
		),
		option.New("erors", "skiped errors", []string{}),
		proto.New("proto_errors", "proto errors", option.New("erors", "skiped errors", []string{})),
	)

	return nil
}
