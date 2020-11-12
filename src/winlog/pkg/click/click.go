package click

import (
	"database/sql"
	"fmt"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/util"
)

func config() string {
	clickhouseHost := util.GetEnv("CLICKHOUSE_HOST", "127.0.0.1:9000")
	return fmt.Sprintf("tcp://%s?username=&compress=true&debug=false", clickhouseHost)
}

func Connect() (err error, connect *sql.DB) {
	connect, err = sql.Open("clickhouse", config())
	connect.SetConnMaxLifetime(-1)
	connect.SetMaxOpenConns(-1)
	if err != nil {
		fmt.Println("DB Connection ErrorL", err)
		return err, nil
	}
	if err = connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println("Error Connect Ping:", err)
			return err, nil
		}
	}

	return err, connect
}
