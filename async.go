package golog

import (
	"os"
	"sync"
	"time"
)

const (

	// 缓冲队列长度
	asyncBufferSize = 100

	// 开启内存池范围的写入大小
	maxTextBytes = 1024
)

var (
	asyncWriteBuff   chan interface{}
	asyncContextPool *sync.Pool
)

// 开启异步写入模式
func EnableASyncWrite() {

	asyncWriteBuff = make(chan interface{}, asyncBufferSize)
	asyncContextPool = new(sync.Pool)

	// 拷贝数据的池
	asyncContextPool.New = func() interface{} {
		return make([]byte, maxTextBytes)
	}

	go func() {

		for {

			// 从队列中获取一个要写入的日志
			raw := <-asyncWriteBuff

			switch d := raw.(type) {
			case func():
				d()
			case []byte:
				// 写入目标
				globalWriter.Write(d)

				// 必须是由pool分配的，才能用池释放
				if cap(d) < maxTextBytes {
					asyncContextPool.Put(d)
				}
			}

		}

	}()
}

// 等待异步写入全部完成
func FlushASyncWrite(timeout time.Duration) {

	if asyncWriteBuff == nil {
		return
	}

	ch := make(chan struct{})

	asyncWriteBuff <- func() {

		// 确保文件已经写入
		if f, ok := globalWriter.(*os.File); ok {
			f.Sync()
		}

		ch <- struct{}{}
	}

	select {
	case <-ch:
	case <-time.After(timeout):
	}
}

func queuedWrite(b []byte) {
	var newb []byte

	// 超大
	if len(b) >= maxTextBytes {
		newb = make([]byte, len(b))
	} else {
		newb = asyncContextPool.Get().([]byte)[:len(b)]
	}

	copy(newb, b)
	asyncWriteBuff <- newb
}
