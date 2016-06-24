package drops

import (
	"bytes"
	"sync"
)

var buffPool sync.Pool = sync.Pool{
	New: func() interface{} { return bytes.NewBuffer(nil) },
}

func GetBuff() *bytes.Buffer {
	return buffPool.Get().(*bytes.Buffer)
}

func PutBuff(buff *bytes.Buffer) {
	buff.Reset()
	buffPool.Put(buff)
}
