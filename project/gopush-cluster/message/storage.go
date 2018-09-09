

package main

import (
	log "github.com/alecthomas/log4go"
	"encoding/json"
	"errors"
	"github.com/Terry-Mao/gopush-cluster/rpc"
)

const (
	RedisStorageType = "redis"
	MySQLStorageType = "mysql"
	ketamaBase       = 255
	saveBatchNum     = 1000
)

var (
	UseStorage     Storage
	ErrStorageType = errors.New("unknown storage type")
)

// Stored messages interface
type Storage interface {
	// GetPrivate get private msgs.
	GetPrivate(key string, mid int64) ([]*rpc.Message, error)
	// SavePrivate Save single private msg.
	SavePrivate(key string, msg json.RawMessage, mid int64, expire uint) error
	// Save private msgs return failed keys.
	SavePrivates(keys []string, msg json.RawMessage, mid int64, expire uint) ([]string, error)
	// DelPrivate delete private msgs.
	DelPrivate(key string) error
}

// InitStorage init the storage type(mysql or redis).
func InitStorage() error {
	if Conf.StorageType == RedisStorageType {
		UseStorage = NewRedisStorage()
	} else if Conf.StorageType == MySQLStorageType {
		UseStorage = NewMySQLStorage()
	} else {
		log.Error("unknown storage type: \"%s\"", Conf.StorageType)
		return ErrStorageType
	}
	return nil
}
