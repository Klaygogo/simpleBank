package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool" // 导入 pgx 连接池
	"github.com/techschool/simplebank/util"
)

var testQueries *Queries
var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("无法加载配置:", err)
	}
	// 1. 使用 pgxpool 建立连接池（推荐生产环境使用）
	pool, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("无法创建数据库连接池:", err)
	}
	defer pool.Close() // 测试结束后关闭连接池

	// 2. pgxpool.Pool 实现了 DBTX 接口，可直接传给 New
	testQueries = New(pool)

	// 3. 运行测试
	os.Exit(m.Run())
}
