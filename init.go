package gows

import "github.com/Gandalf-Le-Dev/abyssnet.protocol.gows/internal"

var (
	framePadding  = frameHeader{}                         // 帧头填充物
	binaryPool    = internal.NewBufferPool(128, 256*1024) // 缓冲池
	defaultLogger = new(stdLogger)                        // 默认日志工具
)
