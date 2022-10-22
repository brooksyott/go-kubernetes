package logger

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	// Use this approach so that it can be read from a configuration file
	// encoding: console can be changed to json
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "console",
		"outputPaths": ["stdout", "/temp/zap.log"],
		"errorOutputPaths": ["stderr"],
		"initialFields": {"service": "go-hellokube"},
		"encoderConfig": {
          "timeKey": "datetime",
          "callerKey": "caller",
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// You can add additional values here, ie. the kubernetes node
	cfg.InitialFields["node"] = "testNode"

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}
