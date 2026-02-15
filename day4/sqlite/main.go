// Day4 示例：SQLite + database/sql 简单 CRUD（Create/Read/Update/Delete）
// 本文件学习：Open、Exec（建表/插入/更新/删除）、QueryRow、Query、Scan
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // 匿名导入：仅触发该包的 init()；init() 里会调 sql.Register("sqlite", ...)，把驱动登记到 database/sql，之后 Open("sqlite", dsn) 才能用
)

func main() {
	// 连接 SQLite（modernc 驱动），内存库不落盘，适合示例；无需 CGO
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // 程序退出前关闭连接

	// ---------- 建表：Exec 执行不返回结果的 SQL ----------
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// ---------- 插入：用 ? 占位符传参数，防止 SQL 注入 ----------
	res, err := db.Exec("INSERT INTO users (name) VALUES (?)", "小王")
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId() // 拿到自增的 id（database/sql 的 Result 接口）

	fmt.Println("插入 ID:", id)

	// ---------- 查询一行：QueryRow + Scan（都是 database/sql 的 API）----------
	var name string                                                             // 用来接这一行里 name 列的值
	err = db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)   // &name 传地址，Scan 把查到的值填进 name；有错误时 err 非 nil
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("查询到:", name)

	// ---------- 查询多行：Query 返回 rows，Next 逐行，Scan 把每列的值填进变量 ----------
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var rid int64
		var n string
		rows.Scan(&rid, &n)
		fmt.Printf("  id=%d name=%s\n", rid, n)
	}

	// ---------- 更新：Exec 执行 UPDATE ----------
	_, err = db.Exec("UPDATE users SET name = ? WHERE id = ?", "小王（已改）", id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("更新后:", id)
	err = db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("   name =", name)

	// ---------- 删除：Exec 执行 DELETE ----------
	res2, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	n, _ := res2.RowsAffected()
	fmt.Printf("删除 id=%d，影响行数: %d\n", id, n)
	rows2, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()
	fmt.Println("删除后列表:")
	for rows2.Next() {
		var rid int64
		var n string
		rows2.Scan(&rid, &n)
		fmt.Printf("  id=%d name=%s\n", rid, n)
	}
}
