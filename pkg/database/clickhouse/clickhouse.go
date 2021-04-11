package clickhouse

import (
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func NewClickHouseClient(addrs []string, user, password, dbName string) (*sqlx.DB, error) {
	if len(addrs) <= 0 {
		return nil, errors.New("invalid clickhouse address")
	}

	var dsn = fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s", addrs[0], user, password, dbName)
	if len(addrs) > 1 {
		// alt_hosts-单个逗号分隔的地址主机列表，用于负载均衡
		dsn += fmt.Sprintf("&alt_hosts=%s", strings.Join(addrs[1:], ","))
	}

	sqlDB, err := sqlx.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.Errorf("clickhouse ping error: %s", err)
		}
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(time.Hour)

	return sqlDB, nil
}
