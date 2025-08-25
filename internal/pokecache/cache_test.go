package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	clearInterval := 5 * time.Second
	cache := NewCache(clearInterval)

	cases := []struct {
		key string
		val []byte
	}{
        {
            key: "https://foo.bar",
            val: []byte("testy data"),
        }, {
            key: "file://funny.cache/root/system32",
            val: []byte("{more test data}"),
        },
    }

    for i, c := range cases {
        t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
            cache.Add(c.key, c.val)
            val, ok := cache.Get(c.key)
            if !ok {
                t.Errorf("expected to find key")
                return
            }
            if string(val) != string(c.val) {
                t.Errorf("expected to find value")
                return
            }
        })
    }
}

func TestReapLoop(t *testing.T) {
    const clearInterval = 5 * time.Millisecond
    const waitTime = clearInterval + 5*time.Millisecond
    cache := NewCache(clearInterval)
    cache.Add("https://foo.bar", []byte("testy"))
    _, ok := cache.Get("https://foo.bar")
    if !ok {
        t.Errorf("expected to find key")
        return
    }
    time.Sleep(waitTime)

    _, ok = cache.Get("https://foo.bar")
    if ok {
        t.Errorf("expected not to find key")
        return
    }

}
