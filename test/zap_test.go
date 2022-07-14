package test

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestZapDemo1(t *testing.T) {
	logger, err := zap.NewDevelopment(zap.AddCaller(), zap.AddStacktrace(zapcore.InfoLevel))
	if err != nil {
		return
	}
	defer logger.Sync()

	url := "http://example.org/api"
	logger.Info("failed to fetch URL",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	sugar := logger.Sugar()
	logger.With()
	sugar.Infow("failed to fetch URL",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}

func TestReference(t *testing.T) {
	a1 := myDemoStruct{i: 1, name: "nelson"}
	a2 := a1
	fmt.Printf("a1 : %p \n", &a1)
	fmt.Printf("a2 : %p \n", &a2)
}

type myDemoStruct struct {
	i    int
	name string
}
