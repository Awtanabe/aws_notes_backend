package migrations

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text"`
	ImageURL  string    `json:"image_url" gorm:"type:varchar(500)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SeedData(db *gorm.DB) error {
	// Check if data already exists
	var count int64
	db.Model(&Post{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Sample posts
	posts := []Post{
		{
			Title:    "Goでのバックエンド開発入門",
			Content:  "Goは高速で効率的なバックエンド開発に適したプログラミング言語です。シンプルな構文と強力な並行処理機能を持ち、Webアプリケーションの構築に最適です。",
			ImageURL: "https://images.unsplash.com/photo-1587620962725-abab7fe55159?w=800",
		},
		{
			Title:    "Next.jsでモダンなフロントエンド",
			Content:  "Next.jsはReactベースのフレームワークで、サーバーサイドレンダリングや静的サイト生成をサポートしています。SEOに優れた高速なWebアプリケーションを構築できます。",
			ImageURL: "https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=800",
		},
		{
			Title:    "Dockerで開発環境を統一",
			Content:  "Dockerを使用することで、開発環境を簡単に構築・共有できます。コンテナ技術により、どの環境でも同じように動作するアプリケーションを作成できます。",
			ImageURL: "https://images.unsplash.com/photo-1605745341112-85968b19335b?w=800",
		},
		{
			Title:    "MySQLデータベースの基礎",
			Content:  "MySQLは世界中で使われているオープンソースのリレーショナルデータベースです。高い信頼性とパフォーマンスを持ち、Webアプリケーションのデータ管理に適しています。",
			ImageURL: "https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=800",
		},
		{
			Title:    "REST APIの設計原則",
			Content:  "RESTful APIは、HTTPメソッドを適切に使用し、リソース指向の設計を行います。GET、POST、PUT、DELETEを使い分けることで、直感的で使いやすいAPIを構築できます。",
			ImageURL: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800",
		},
	}

	for _, post := range posts {
		if err := db.Create(&post).Error; err != nil {
			return err
		}
	}

	return nil
}
