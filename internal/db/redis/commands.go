package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Константы для настройки таймаутов и срока жизни
const (
	timeExpiration time.Duration = 1000 // Время жизни блокировки
	timeSleep      time.Duration = 100  // Интервал повтора при блокировке
	timeOut        time.Duration = 5000 // Таймаут операции
)

var (
	errDeleteBlocking = errors.New("error delete blocking for record")
	errSetBlocking    = errors.New("error set blocking for record")
	errBlocking       = errors.New("record blocked by another process")
)

// UserRedisImpl представляет реализацию Redis для пользовательских операций
type UserRedisImpl struct {
	Redis   redis.Cmdable        // Клиент Redis
	Config  *redisUserExtSetting // Конфигурация операции
	process *Process             // Управление процессами
	clone   bool                 // Флаг клонирования
}

// Возвращает экземпляр Redis для пользовательских операций. Если экземпляр помечен как клонированный,
// это гарантирует, что перед возвратом будет правильно установлен полный ключ. В противном случае создается и
// возвращает новый клонированный экземпляр с настройками времени ожидания по умолчанию и истечения срока действия.
func (uri *UserRedisImpl) getInstance() *UserRedisImpl {
	if uri.clone {
		if len(uri.Config.fullKey) == 0 && len(uri.Config.pattern) > 0 && len(uri.Config.key) > 0 {
			uri.Config.fullKey = uri.Config.UnionKeys()
		} else if len(uri.Config.fullKey) == 0 && len(uri.Config.key) > 0 {
			uri.Config.fullKey = uri.Config.key
		}
		return uri
	} else {
		return &UserRedisImpl{
			Config: &redisUserExtSetting{
				timeoutOp: timeOut * time.Millisecond,
				exprBlock: timeExpiration * time.Millisecond,
				exprRec:   0,
			},
			Redis:   uri.Redis,
			clone:   true,
			process: uri.process,
		}
	}
}

func (uri *UserRedisImpl) resetInstance() *UserRedisImpl {
	uri.clone = false
	return uri
}

// SetPattern устанавливает паттерн для формирования ключа
func (uri *UserRedisImpl) SetPattern(pattern string) *UserRedisImpl {
	redis := uri.getInstance()
	redis.Config.pattern = pattern
	return redis
}

// SetValue устанавливает значение для записи в Redis
func (uri *UserRedisImpl) SetValue(value interface{}) *UserRedisImpl {
	redis := uri.getInstance()
	redis.Config.value = value
	return redis
}

// SetKey устанавливает ключ для операции
func (uri *UserRedisImpl) SetKey(key string) *UserRedisImpl {
	redis := uri.getInstance()
	redis.Config.key = key
	return redis
}

// SetExprRec устанавливает время жизни записи в Redis
func (uri *UserRedisImpl) SetExprRec(expr time.Duration) *UserRedisImpl {
	redis := uri.getInstance()
	redis.Config.exprRec = expr
	return redis
}

// SetTimeoutOp устанавливает таймаут операции Redis
func (uri *UserRedisImpl) SetTimeoutOp(timeout time.Duration) *UserRedisImpl {
	redis := uri.getInstance()
	redis.Config.timeoutOp = timeout
	return redis
}

// SetExprBlock устанавливает время жизни блокировки записи в Redis
func (uri *UserRedisImpl) SetExprBlock(expr time.Duration) *UserRedisImpl {
	redis := uri.getInstance()
	redis.Config.exprBlock = expr
	return redis
}

// Lock включает механизм блокировки для обеспечения консистентности операций
func (uri *UserRedisImpl) Lock() *UserRedisImpl {
	redis := uri.getInstance()
	redis.Config.lock = true
	return redis
}

// Unlock отключает механизм блокировки после завершения операций
func (uri *UserRedisImpl) Unlock() *UserRedisImpl {
	redis := uri.getInstance()
	redis.Config.unlock = true
	return redis
}

// SetNX устанавливает значение ключа, если он еще не существует в Redis.
//
// Он поддерживает механизм блокировки, если он настроен, и обрабатывает время ожидания операции.
//
// Пример шаблона для выполнения: redis.SetPattern().SetKey().setValue().Lock().SetNX(ctx)
func (uri *UserRedisImpl) SetNX(ctx context.Context) *redis.BoolCmd {
	rds := uri.getInstance()
	defer rds.resetInstance()

	// Выставляется таймаут
	ctx, cancel := context.WithTimeout(ctx, rds.Config.timeoutOp)
	defer cancel()

	if rds.Config.lock {
		idProc, err := uri.block(ctx, rds.Config.fullKey)
		defer uri.unBlock(ctx, rds.Config.fullKey, idProc)
		if err != nil {
			err = fmt.Errorf("%s:%w", GetFunctionName(), err)
			rds.Config.err = append(rds.Config.err, err)
			var cmd = &redis.BoolCmd{}
			cmd.SetErr(err)
			return cmd
		}
	}

	return rds.Redis.SetNX(ctx, rds.Config.fullKey, rds.Config.value, rds.Config.exprRec)
}

// Включить блокировку по ключу для полного выполнения транзакции
//
// Блокировка по ключу реализуется для случаев практически одновременного чтения ключа
// и его изменения в процессе какого-либо бизнес процесса, что может привести к нарушении
// консистентности
//
// Включается блокировка ключа при чтении (Get) и отключается при записи (Create, Update),
// если в вышестоящей функции имеется атрибут обязательной блокировки
//
// У блокировки имеется TLS, по истечении которого блокировка с ключа сбрасывается для последующих операций
func (uri *UserRedisImpl) block(ctx context.Context, key string) (idProc uint64, err error) {
	idProc = uri.process.GetID()
	idProcStr := strconv.FormatUint(idProc, 10)
	keyLock := key + ":lock"

	setBlock := func() error {
		// Попытка установки блокировки
		set, err := uri.Redis.SetNX(ctx, keyLock, idProcStr, uri.Config.exprBlock).Result()
		if err != nil {
			return fmt.Errorf("%s:%w:%w", GetFunctionName(), errSetBlocking, err)
		}

		// Если блокировка установлена - всё ок
		// Иначе - ошибка
		if set {
			return nil // Блокировка успешно установлена
		} else {
			return fmt.Errorf("%s:%w:%w", GetFunctionName(), errBlocking, err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("%s:%w", GetFunctionName(), ctx.Err())
		default:
			if err := setBlock(); err != nil {
				// Ошибка установки блокировки
				if errors.Is(err, errSetBlocking) {
					return 0, fmt.Errorf("%s:%w", GetFunctionName(), err)
				}
				time.Sleep(timeSleep * time.Millisecond)
				continue
			}
		}
		return idProc, nil // Блокировка установлена и ключ готов к использованию
	}
}

// Разблокирует ключ Redis, сняв его блокировку.
// Если указан idProc, проверяет, принадлежит ли блокировка указанному процессу, прежде чем снимать ее.
func (uri *UserRedisImpl) unBlock(ctx context.Context, key string, idProc ...uint64) error {
	keyLock := key + ":lock"

	if len(idProc) > 0 {
		value, err := uri.Redis.Get(ctx, keyLock).Result()
		if err != nil {
			return fmt.Errorf("%s:%w:%w", GetFunctionName(), errDeleteBlocking, err)
		} else {
			value, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return fmt.Errorf("%s:%w:%w", GetFunctionName(), errDeleteBlocking, err)
			}
			if value != idProc[0] {
				return fmt.Errorf("%s:%w:%w", GetFunctionName(), errBlocking, err)
			}
		}
	}
	if err := uri.Redis.Del(ctx, keyLock).Err(); err != nil {
		return fmt.Errorf("%s:%w:%w", GetFunctionName(), errDeleteBlocking, err)
	}

	return nil // Блокировка успешно удалена
}
