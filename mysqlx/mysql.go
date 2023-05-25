package mysqlx

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tanenking/svrframe/config"
	"github.com/tanenking/svrframe/logx"

	"github.com/jmoiron/sqlx"
)

type MysqlClient interface {
	Select(dest interface{}, query string, args ...interface{}) (err error)
	Delete(query string, args ...interface{}) (rowsAffected int64, err error)
	Update(query string, args ...interface{}) (rowsAffected int64, err error)
	Insert(query string, args ...interface{}) (lastInsertID int64, err error)
	Query(query string, args ...interface{}) (results []*result, err error)
}

type mysqlClient struct {
	*config.MysqlConfig
	mysqlDB *sqlx.DB
}

func makeMysql(cfg *config.MysqlConfig) (cli *mysqlClient, err error) {
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)

	cli = &mysqlClient{
		MysqlConfig: cfg,
	}
	cli.mysqlDB, err = sqlx.Open("mysql", dbDSN)
	if err != nil {
		return
	}
	cli.mysqlDB.SetMaxOpenConns(mysql_max_open_conns)
	cli.mysqlDB.SetMaxIdleConns(mysql_max_idle_conns)
	cli.mysqlDB.SetConnMaxLifetime(mysql_max_lifttime)

	if err = cli.mysqlDB.Ping(); err != nil {
		return
	}

	logx.InfoF("mysql conn success -> %s", dbDSN)
	return
}

func (m *mysqlClient) _update_or_delete(query string, args ...interface{}) (rowsAffected int64, err error) {
	var ret sql.Result
	ret, err = m.mysqlDB.Exec(query, args...)
	if err != nil {
		logx.ErrorF("%v", err)
		return
	}
	rowsAffected, err = ret.RowsAffected()
	if err != nil {
		logx.ErrorF("%v", err)
		return
	}
	return
}

func (m *mysqlClient) Select(dest interface{}, query string, args ...interface{}) (err error) {
	err = m.mysqlDB.Select(dest, query, args...)
	if err != nil {
		logx.ErrorF("%v", err)
		return
	}
	return
}

func (m *mysqlClient) Delete(query string, args ...interface{}) (rowsAffected int64, err error) {
	s := strings.ToLower(query[0:6])
	if strings.Compare(s, "delete") != 0 {
		err = fmt.Errorf("this function just run delete")
		logx.ErrorF("%v", err)
		return
	}
	return m._update_or_delete(query, args...)
}
func (m *mysqlClient) Update(query string, args ...interface{}) (rowsAffected int64, err error) {
	s := strings.ToLower(query[0:6])
	if strings.Compare(s, "update") != 0 {
		err = fmt.Errorf("this function just run update")
		logx.ErrorF("%v", err)
		return
	}
	return m._update_or_delete(query, args...)
}

func (m *mysqlClient) Insert(query string, args ...interface{}) (lastInsertID int64, err error) {
	s := strings.ToLower(query[0:6])
	if strings.Compare(s, "insert") != 0 {
		err = fmt.Errorf("this function just run insert")
		logx.ErrorF("%v", err)
		return
	}
	var ret sql.Result
	ret, err = m.mysqlDB.Exec(query, args...)
	if err != nil {
		logx.ErrorF("%v", err)
		return
	}
	lastInsertID, err = ret.LastInsertId()
	if err != nil {
		logx.ErrorF("%v", err)
		return
	}
	return
}

func (m *mysqlClient) Query(query string, args ...interface{}) (results []*result, err error) {
	s := strings.ToLower(query[0:6])
	if strings.Compare(s, "select") != 0 {
		s = strings.ToLower(query[0:4])
		if strings.Compare(s, "call") != 0 {
			err = fmt.Errorf("this function just run select or call")
			logx.ErrorF("%v", err)
			return
		}
	}

	results = []*result{}
	var rows *sql.Rows
	//查询数据，取所有字段
	rows, err = m.mysqlDB.Query(query, args...)
	if err != nil {
		return
	}

	for {
		res := &result{}
		//返回所有列
		res.Fileds, err = rows.Columns()
		if err != nil {
			return
		}
		//这里表示一行所有列的值，用[]byte表示
		vals := make([][]byte, len(res.Fileds))
		//这里表示一行填充数据
		scans := make([]interface{}, len(res.Fileds))
		//这里scans引用vals，把数据填充到[]byte里
		for k := range vals {
			scans[k] = &vals[k]
		}
		i := 0
		for rows.Next() {
			//填充数据
			rows.Scan(scans...)
			//每行数据
			r := row{
				Values: map[string]string{},
			}
			//把vals中的数据复制到row中
			for k, v := range vals {
				key := res.Fileds[k]
				//这里把[]byte数据转成string
				r.Values[key] = string(v)
			}
			//放入结果集
			res.Rows = append(res.Rows, r)
			i++
		}

		results = append(results, res)
		if !rows.NextResultSet() {
			break
		}
	}

	if err != nil {
		logx.ErrorF("%v", err)
	}
	return
}
