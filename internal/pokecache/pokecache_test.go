package pokecache

import (
	"testing"
	"time"
)

func TestCreateCache(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)

	if cache.cacheMap == nil {
		t.Error("cache is nil")
	}
}

func TestCreateCacheAddAndGet(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)

	cases := []struct {
		inputKey string
		inputValue []byte
	} {
		{
			inputKey: "key",
        	inputValue: []byte("value"),
		},
		{
            inputKey: "another_key",
            inputValue: []byte("another_value"),
        },
		{
            inputKey: "new_key",
            inputValue: []byte("new_value"),
        },
	}

	for _, c := range cases {
		cache.Add(c.inputKey, c.inputValue)

		actual, ok := cache.Get(c.inputKey)
		if !ok {
			t.Errorf("key %s not found", c.inputKey)
		}
		if string(actual) != string(c.inputValue) {
			t.Errorf("%s value is not equal to expected value: %s", string(actual), string(c.inputValue));
		}
	}

}

func TestCreateCacheReap(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)

	key := "key"
	cache.Add(key, []byte(""))

	time.Sleep(interval + time.Millisecond)

	_, ok := cache.Get(key)
	if ok {
		t.Errorf("%s should have been reaped", key)
	}
}

func TestCreateCacheReapFail(t *testing.T) {
	interval := time.Millisecond * 10
	cache := NewCache(interval)

	key := "key"
	cache.Add(key, []byte(""))

	time.Sleep(interval / 2)

	_, ok := cache.Get(key)
	if !ok {
		t.Errorf("%s should not have been reaped", key)
	}
}