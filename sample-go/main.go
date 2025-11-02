package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// コネクションプールのサイズ
const poolSize = 5

func main() {
	dsn := "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB接続エラー:", err)
	}
	defer db.Close()

	// 🔧 コネクションプールの設定
	db.SetMaxOpenConns(poolSize)            // 最大同時接続数（PostgreSQL側で3バックエンドまで）
	db.SetMaxIdleConns(poolSize)            // アイドル状態で保持する接続数
	db.SetConnMaxLifetime(30 * time.Second) // 接続の寿命（任意）

	fmt.Println("=== コネクションプール実験開始　プール数", poolSize, " ===")

	var wg sync.WaitGroup

	// 🔁 5つの同時クエリを実行してみる
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			start := time.Now()

			// 実際に時間のかかるクエリを投げる
			_, err := db.Exec("SELECT pg_sleep(3)") // 3秒間スリープ
			if err != nil {
				fmt.Printf("クエリ%d エラー: %v\n", n, err)
				return
			}

			fmt.Printf("クエリ%d 完了！（経過: %.1fs）\n", n, time.Since(start).Seconds())
		}(i)
	}

	wg.Wait()
	fmt.Println("=== 実験終了 ===")
}
