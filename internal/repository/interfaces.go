package repository

import "blog-system/internal/domain"

type BlogRepository interface {
	CreateArticle(article *domain.Article) error
	GetArticle(id int) (*domain.Article, error)
	GetAllArticles() ([]*domain.Article, error)
	UpdateArticle(*domain.Article) error
	DeleteArticle(id int) error

	CreateComment(comment *domain.Comment) error
	GetCommentsByArticleID(articleID int) ([]*domain.Comment, error)
	DeleteComment(id int) error

	CreateTag(tag *domain.Tag) error
	GetAllTags() ([]*domain.Tag, error)
	GetTagsByArticleID(articleID int) ([]*domain.Tag, error)
	AddTagToArticle(articleID int, tagID int) error
	RemoveTagFromArticle(articleID int, tagID int) error
}
