// Day4 示例：PostgreSQL + database/sql，简单 CRUD（与 sqlite 对照，占位符用 $1,$2）
// 需要本地已安装并启动 Postgres；可选设置环境变量 DB_DSN，否则用下方默认。
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // 匿名导入：init() 里 sql.Register("postgres", ...)，Open 时才能用
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("无法连接数据库:", err)
	}
	fmt.Println("已连接 PostgreSQL")

	// 建表（Postgres 用 SERIAL 自增）
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// 插入；Postgres 占位符为 $1, $2（与 SQLite 的 ? 不同）
	var id int
	err = db.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id", "小李").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("插入 ID:", id)

	// 查一条
	var name string
	err = db.QueryRow("SELECT name FROM users WHERE id = $1", id).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("查询到:", name)

	// 查多行
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var rid int
		var n string
		rows.Scan(&rid, &n)
		fmt.Printf("  id=%d name=%s\n", rid, n)
	}

	// 更新
	_, err = db.Exec("UPDATE users SET name = $1 WHERE id = $2", "小李（已改）", id)
	if err != nil {
		log.Fatal(err)
	}
	_ = db.QueryRow("SELECT name FROM users WHERE id = $1", id).Scan(&name)
	fmt.Println("更新后 name =", name)

	// 删除
	res, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	n, _ := res.RowsAffected()
	fmt.Printf("删除影响行数: %d\n", n)
}
