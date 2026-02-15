# Day 4：数据层

- 使用 `database/sql` + 驱动；**本系列先用 SQLite**（无需安装数据库），Postgres 可选。
- 连接池、Query/Exec、预编译语句
- 简单 CRUD：插入（Create）、查一条/查列表（Read）、更新（Update）、删除（Delete）

**说明**：本日学习**只需跑 SQLite**，无需安装 PostgreSQL；`day4/postgres` 为可选示例（已装 Postgres 时可体验另一驱动），没装也不影响学习。

本目录 SQLite 示例使用**纯 Go 驱动**（modernc.org/sqlite），**无需 CGo、无需安装 gcc**，直接 `go run ./day4/sqlite` 即可。

## 运行

```bash
# 先跑 SQLite（无需安装数据库，推荐）
go run ./day4/sqlite

# 可选：PostgreSQL（需本地已安装并启动 Postgres；示例含 CRUD，占位符为 $1,$2，可与 sqlite 对照）
go run ./day4/postgres
```
