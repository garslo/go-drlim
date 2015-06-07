package limiters

import (
	"encoding/json"
	"log"

	"github.com/garyburd/redigo/redis"
)

type redisCounter struct {
	conn         redis.Conn
	keyPrefix    string
	host         string
	localCounts  *memoryCounter
	globalCounts *memoryCounter
}

func NewRedisCounter(conn redis.Conn, keyPrefix, host string) Counter {
	return &redisCounter{
		conn:         conn,
		keyPrefix:    keyPrefix,
		host:         host,
		localCounts:  newMemoryCounter(make(map[string]int)),
		globalCounts: newMemoryCounter(make(map[string]int)),
	}
}

func (me *redisCounter) Zero() {
	me.localCounts.Zero()
}

func (me *redisCounter) Update() {
	// TODO: or protobuf, etc.
	data, err := json.Marshal(me.localCounts.counts)
	if err != nil {
		log.Printf("error serializing local counts, bailing: %v", err)
		return
	}
	_, err = me.conn.Do("SET", me.keyPrefix+"_"+me.host, data)
	if err != nil {
		log.Printf("error setting local counts in redis, bailing: %v", err)
		return
	}
	rawKeys, err := me.conn.Do("KEYS", me.keyPrefix+"*")
	if err != nil {
		log.Printf("error fetching keys in redis, bailing: %v", err)
		return
	}
	_, ok := rawKeys.([]string)
	if !ok {
		log.Printf("error converting raw redis keys to strings, bailing: %v", err)
		return
	}
	//for _, key := range keys {
	//
	//}
}

func (me *redisCounter) Count(host string) int {
	return me.globalCounts.Count(host)
}

func (me *redisCounter) Inc(host string) {
	me.localCounts.Inc(host)
	me.globalCounts.Inc(host)
}
