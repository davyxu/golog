package golog

import "sync"

const (

	// 缓冲队列长度
	asyncBufferSize = 100

	// 开启内存池范围的写入大小
	maxTextBytes = 1024
)

var (
	asyncWriteBuff   chan []byte
	asyncContextPool *sync.Pool
)

// 开启异步写入模式
func EnableASyncWrite() {

	asyncWriteBuff = make(chan []byte, asyncBufferSize)
	asyncContextPool = new(sync.Pool)

	// 拷贝数据的池
	asyncContextPool.New = func() interface{} {
		return make([]byte, maxTextBytes)
	}

	go func() {

		for {

			// 从队列中获取一个要写入的日志
			b := <-asyncWriteBuff

			// 写入目标
			globalWriter.Write(b)

			// 必须是由pool分配的，才能用池释放
			if cap(b) < maxTextBytes {
				asyncContextPool.Put(b)
			}

		}

	}()
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
