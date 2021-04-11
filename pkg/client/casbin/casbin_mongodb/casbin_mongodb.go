package casbin_mongodb

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/mongodb-adapter/v3"
	log "github.com/sirupsen/logrus"
	"go-lab/app/go_web/go_gin/conf"
	casbin2 "go-lab/pkg/client/casbin"
	"sync"
)

var onceCasbin sync.Once // 保证全局只有1个
var casbinEnforcer *casbin.SyncedEnforcer

var casbinEnforcerInitErr error

const connStringFmtWithoutPwd = "mongodb://%s/%s?authSource=admin&replicaSet=%s"    // 没有密码时用的连接字符串
const connStringFmtWithPwd = "mongodb://%s:%s@%s/%s?authSource=admin&replicaSet=%s" // 有密码时用的连接字符串

// casbin Enforcer(单例)
func GetCabinEnforcer() *casbin.SyncedEnforcer {
	config := conf.GConfig.MongoDbConfigs.DB
	if casbinEnforcer != nil {
		return casbinEnforcer
	}
	onceCasbin.Do(createMongodbCasbinEnforcer(&config))
	if casbinEnforcerInitErr != nil {
		log.Error("GetCabinEnforcer error :" + casbinEnforcerInitErr.Error())
		return nil
	} else {
		return casbinEnforcer
	}
}

func createMongodbCasbinEnforcer(config *conf.MongoDb) func() {
	return func() {
		var connectionString string
		if config.Password != "" {
			connectionString = fmt.Sprintf(connStringFmtWithPwd, config.User, config.Password, config.Addr, config.Db, config.Rs)
		} else {
			connectionString = fmt.Sprintf(connStringFmtWithoutPwd, config.Addr, config.Db, config.Rs)
		}
		a, err := mongodbadapter.NewAdapter(connectionString)
		if err != nil {
			log.WithError(err).Error("initCasbinEnforcers error when  new adapter")
			casbinEnforcerInitErr = err
			return
		}
		enforcer, err := casbin.NewSyncedEnforcer(casbin2.Model, a)
		if err != nil {
			log.Errorf("[initCasbinEnforcers] error when new enforcer,  err [%v]", err)
			casbinEnforcerInitErr = err
			return
		}
		err = enforcer.LoadPolicy()
		if err != nil {
			log.Errorf("[initCasbinEnforcers] error when enforcer load policy, err [%v]", err)
			casbinEnforcerInitErr = err
			return
		}
		casbinEnforcer = enforcer
	}
}
