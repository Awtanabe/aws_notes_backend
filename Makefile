.PHONY: migrate seed run

# マイグレーションとシードデータを実行
migrate:
	go run cmd/migrate/main.go

# アプリケーションを起動
run:
	go run main.go

# シードデータのみ実行（migrateと同じ）
seed: migrate
