package mysqlx

import (
	"fmt"
	"time"

	"github.com/tanenking/svrframe/config"
	"github.com/tanenking/svrframe/logx"

	_ "github.com/go-sql-driver/mysql"
)

const (
	//最大连接数
	mysql_max_open_conns = 100
	//闲置连接数
	mysql_max_idle_conns = 20
	//最大连接周期
	mysql_max_lifttime = 60 * time.Second
)

var (
	mysqls map[string]*mysqlClient
)

func init() {
	mysqls = make(map[string]*mysqlClient)
}

func InitMysqlHelper() error {
	configs := config.GetMysqlConfigs()
	if configs == nil || len(configs) <= 0 {
		return nil
	}

	for _, cfg := range configs {
		cli, err := makeMysql(cfg)
		if err != nil {
			logx.ErrorF("makeMysql err -> %v", err)
			return err
		}
		mysqls[cfg.Name] = cli
		logx.InfoF("mysql helper [ %s ] create success", cfg.Name)
	}

	logx.InfoF("InitMysqlHelper success")
	return nil
}

func GetMysqlClient(dbname string) MysqlClient {
	cli, ok := mysqls[dbname]
	if !ok {
		return nil
	}
	return cli
}

func Select(dbname string, dest interface{}, query string, args ...interface{}) (err error) {
	m := GetMysqlClient(dbname)
	if m == nil {
		err = fmt.Errorf("not found mysql client -> name %s", dbname)
		logx.ErrorF("Select err -> %v", err)
		return
	}
	return m.Select(dest, query, args...)
}

func Delete(dbname string, query string, args ...interface{}) (rowsAffected int64, err error) {
	m := GetMysqlClient(dbname)
	if m == nil {
		err = fmt.Errorf("not found mysql client -> name %s", dbname)
		logx.ErrorF("Delete err -> %v", err)
		return
	}
	return m.Delete(query, args...)
}
func Update(dbname string, query string, args ...interface{}) (rowsAffected int64, err error) {
	m := GetMysqlClient(dbname)
	if m == nil {
		err = fmt.Errorf("not found mysql client -> name %s", dbname)
		logx.ErrorF("Update err -> %v", err)
		return
	}
	return m.Update(query, args...)
}

func Insert(dbname string, query string, args ...interface{}) (lastInsertID int64, err error) {
	m := GetMysqlClient(dbname)
	if m == nil {
		err = fmt.Errorf("not found mysql client -> name %s", dbname)
		logx.ErrorF("Insert err -> %v", err)
		return
	}
	return m.Insert(query, args...)
}

func Query(dbname string, query string, args ...interface{}) (results []*result, err error) {
	m := GetMysqlClient(dbname)
	if m == nil {
		err = fmt.Errorf("not found mysql client -> name %s", dbname)
		logx.ErrorF("Query err -> %v", err)
		return
	}
	return m.Query(query, args...)
}
