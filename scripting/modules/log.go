package modules

import (
	"github.com/Sirupsen/logrus"
	"github.com/bbuck/dragon-mud/logger"
	"github.com/bbuck/dragon-mud/scripting/engine"
	"github.com/bbuck/dragon-mud/scripting/keys"
)

// Log is the definition of the Lua logging module.
var Log = map[string]interface{}{
	"error": func(eng *engine.Lua) int {
		performLog(eng, func(l *logrus.Entry, msg string) {
			l.Error(msg)
		})

		return 0
	},
	"warn": func(eng *engine.Lua) int {
		performLog(eng, func(l *logrus.Entry, msg string) {
			l.Warn(msg)
		})

		return 0
	},
	"info": func(eng *engine.Lua) int {
		performLog(eng, func(l *logrus.Entry, msg string) {
			l.Info(msg)
		})

		return 0
	},
	"debug": func(eng *engine.Lua) int {
		performLog(eng, func(l *logrus.Entry, msg string) {
			l.Debug(msg)
		})

		return 0
	},
}

func loggerForEngine(eng *engine.Lua) *logrus.Entry {
	llog := eng.GetGlobal(keys.Logger)
	if log, ok := llog.Interface().(*logrus.Entry); ok {
		return log
	}

	name := eng.GetGlobal(keys.EngineID).AsString()
	if name == "" {
		name = "Unknown Engine"
	}

	log := logger.LogWithSource(name)
	eng.SetGlobal(keys.Logger, log)
	eng.WhitelistFor(log)

	return log
}

func performLog(eng *engine.Lua, fn func(*logrus.Entry, string)) {
	data := eng.Nil()
	if eng.StackSize() == 2 {
		data = eng.PopTable()
	}
	msg := eng.PopString()

	log := loggerForEngine(eng)

	if !data.IsNil() && data.IsTable() {
		m := data.AsMapStringInterface()
		log = log.WithFields(logrus.Fields(m))
	}

	fn(log, msg)
}
