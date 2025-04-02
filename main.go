package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	// 数据库连接参数
	dbHost := "caojingkangdb.rwlb.rds.aliyuncs.com" // PolarDB集群连接地址
	dbPort := "3306"                                // 默认端口3306
	dbUser := "workuser"                            // 数据库账号
	dbPass := "123456@cjk"                          // 数据库账号的密码
	dbName := "testsql"                             // 需要连接的数据库名

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println(dsn)
	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 创建一个操作游标
	var result string
	err = db.QueryRow("SELECT VERSION()").Scan(&result)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	// 输出数据库版本
	fmt.Printf("Connected to database, version: %s\n", result)

	// 执行 SQL 查询
	rows, err := db.Query("SELECT * FROM `<YOUR_TABLE_NAME>`") // 需要检索的数据表名
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	// 处理查询结果
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	// 检查是否有错误
	if err := rows.Err(); err != nil {
		log.Fatalf("Error during iteration: %v", err)
	}
}
