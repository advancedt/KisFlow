package test

import (
	"KisFlow/log"
	"context"
	"testing"
)

func TestKisLogger(t *testing.T) {
	// 初始化了一个空的上下文
	ctx := context.Background()

	log.Logger().InfoFX(ctx, "TestKisLogger InfoFX\n")
	log.Logger().ErrorFX(ctx, "TestKisLogger ErrorFX\n")
	log.Logger().DebugFX(ctx, "TestKisLogger DebugFX\n")

	log.Logger().InfoF("TestKisLogger InfoF\n")
	log.Logger().ErrorF("TestKisLogger ErrorF\n")
	log.Logger().DebugF("TestKisLogger DebugF\n")
}
