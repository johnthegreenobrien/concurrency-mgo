package cache

import (
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"

	"github.com/go-redis/cache"
	"sync"
	"time"
)

var codec *cache.Codec
var once sync.Once

func GetInstance() *cache.Codec {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			//Password: "", // no password set
			//DB:       0,  // use default DB
		})

		codec = &cache.Codec{
			Redis: client,

			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)
			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		}
	})
	return codec
}

func Get(key string, wg *sync.WaitGroup) string {
	wanted_objs := make(chan string)
	wg.Add(1)
	// singleton is thread safe and could be used with goroutines
	go func() {
		codec := GetInstance()
		var wanted string
		if err := codec.Get(key, &wanted); err == nil {
			wanted_objs <- wanted
		}
		wg.Done()
	}()
	return <-wanted_objs
}

func Set(key string, val string, wg *sync.WaitGroup) {
	wg.Add(1)
	// singleton is thread safe and could be used with goroutines
	go func() {
		codec := GetInstance()
		codec.Set(&cache.Item{
			Key:        key,
			Object:     val,
			Expiration: time.Minute,
		})
		wg.Done()
	}()
}
