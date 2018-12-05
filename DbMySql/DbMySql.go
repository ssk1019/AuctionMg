package DbMySql

/*
	參考: https://duyanghao.github.io/go-mysql/
*/

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DbMySql struct {
	hDb      *sql.DB
	hDbError error
}

//create the connection pool
func (v *DbMySql) Create(ip string, port string, username string, password string, dbName string) error {
	connMaxLifetime := 30
	maxIdleConns := 20
	maxOpenConns := 20

	//connect to the database
	//par := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&writeTimeout=1s&timeout=10s", username, password, addrs, port, database)
	par := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s", username, password, ip, port, dbName)
	db, err := sql.Open("mysql", par) //第一个参数为驱动名
	if err != nil {
		v.hDbError = fmt.Errorf("Failed to connect to log mysql: %s", err)
		return v.hDbError
	}
	//ping the mysql
	err = db.Ping()
	if err != nil {
		v.hDbError = fmt.Errorf("Failed to ping mysql: %s", err)
		return v.hDbError
	}
	//set db
	v.hDb = db

	//reuse the connection forever(Expired connections may be closed lazily before reuse)
	//If d <= 0, connections are reused forever.
	fmt.Printf("connMaxLifetime:%d\n", connMaxLifetime)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Hour)
	//db.SetConnMaxLifetime(10*time.Second)

	//SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	//If n <= 0, no idle connections are retained.
	fmt.Printf("maxIdleConns:%d\n", maxIdleConns)
	db.SetMaxIdleConns(maxIdleConns)

	//SetMaxOpenConns sets the maximum number of open connections to the database.
	//If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than MaxIdleConns, then MaxIdleConns will be reduced to match the new MaxOpenConns limit
	//If n <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited).
	fmt.Printf("maxOpenConns:%d\n", maxOpenConns)
	db.SetMaxOpenConns(maxOpenConns)

	return nil
}

func (v *DbMySql) Exec(strSql string, args ...interface{}) (sql.Result, error) {
	return v.hDb.Exec(strSql, args...)
}

func (v *DbMySql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return v.hDb.Query(query, args...)
}
