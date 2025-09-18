package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool" // 导入 pgx 连接池
)

const (
	// pgx 的连接字符串格式与标准库兼容，无需修改
	dbSource = "postgres://root:secret@localhost:5432/simpleBank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	// 1. 使用 pgxpool 建立连接池（推荐生产环境使用）
	pool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("无法创建数据库连接池:", err)
	}
	defer pool.Close() // 测试结束后关闭连接池

	// 2. pgxpool.Pool 实现了 DBTX 接口，可直接传给 New
	testQueries = New(pool)

	// 3. 运行测试
	os.Exit(m.Run())
}
