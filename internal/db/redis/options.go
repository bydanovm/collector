package redis

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisUserOptions func(*redisUserExtSetting)

func WithLock() redisUserOptions {
	return func(set *redisUserExtSetting) { set.lock = true }
}

func WithUnLock() redisUserOptions {
	return func(set *redisUserExtSetting) { set.unlock = true }
}

func WithRevertKey() redisUserOptions {
	return func(set *redisUserExtSetting) { set.revertKey = true }
}

func WithRevertZSet() redisUserOptions {
	return func(set *redisUserExtSetting) { set.revertZSet = true }
}

func WithWithoutKey() redisUserOptions {
	return func(set *redisUserExtSetting) { set.withoutKey = true }
}

type redisUserExtSetting struct {
	pattern    string
	key        string
	fullKey    string // Сформированный ключ
	value      interface{}
	lock       bool          // Заблокировать запись
	unlock     bool          // Разблокировать запись
	revertKey  bool          // Перевернуть ключ и паттерн местами
	revertZSet bool          // Записать в массив ZSet наоборот
	withoutKey bool          // Игнорировать поле ключ
	timeoutOp  time.Duration // Таймаут выполнения операции
	exprBlock  time.Duration // Время блокировки записи
	exprRec    time.Duration // Время жизни записи
	err        []error       // Ошибки
}

func (r *redisUserExtSetting) UnionKeys() string {
	return r.pattern + ":" + r.key
}
func getExtSet(pattern string, key int64, opts ...redisUserOptions) redisUserExtSetting {
	var dBExtSet = &redisUserExtSetting{}
	for _, opt := range opts {
		opt(dBExtSet)
	}
	if dBExtSet.withoutKey {
		dBExtSet.key = pattern
	} else {
		if dBExtSet.revertKey {
			dBExtSet.key = fmt.Sprintf("%s:%d", pattern, key)
		} else {
			dBExtSet.key = fmt.Sprintf("%d:%s", key, pattern)
		}
	}
	return *dBExtSet
}

type rdsOptions func(*redis.Options)

func WithRdsAddr(a string) rdsOptions {
	return func(o *redis.Options) { o.Addr = a }
}

func WithRdsPass(a string) rdsOptions {
	return func(o *redis.Options) { o.Password = a }
}

func WithRdsDBNum(a int) rdsOptions {
	return func(o *redis.Options) { o.DB = a }
}
