package snowflake

import (
	"testing"
)

func TestSnowFlake(t *testing.T) {
	worker, err := NewWorker(1)
	if err != nil {
		t.Fatalf("new worker err:%v", err)
		return
	}
	ch := make(chan int64)
	count := 10000
	for i := 0; i < count; i++ {
		go func() {
			id := worker.Next()
			ch <- id
		}()
	}
	defer close(ch)
	m := make(map[int64]int)
	for i := 0; i < count; i++ {
		id := <-ch
		// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
		_, ok := m[id]
		if ok {
			t.Error("ID is not unique!\n")
			return
		}
		// 将 id 作为 key 存入 map
		m[id] = i
	}
	// 成功生成 snowflake ID
	t.Log("snowflake ID Get successed!")
}

func TestMultiSnowFlake(t *testing.T) {
	worker, err := NewWorker(1)
	if err != nil {
		t.Fatalf("new worker err:%v", err)
		return
	}
	t.Log("---worker 1---")
	id := worker.Next()
	t.Logf("id:%d", id)
	id = worker.Next()
	t.Logf("id:%d", id)
	id = worker.Next()
	t.Logf("id:%d", id)

	worker, err = NewWorker(2)
	if err != nil {
		t.Fatalf("new worker err:%v", err)
		return
	}
	t.Log("---worker 2---")
	id = worker.Next()
	t.Logf("id:%d", id)
	id = worker.Next()
	t.Logf("id:%d", id)
	id = worker.Next()
	t.Logf("id:%d", id)
}
