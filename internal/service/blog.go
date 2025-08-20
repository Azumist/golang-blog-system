package service

import (
	"fmt"
	"strings"

	"blog-system/internal/domain"
	"blog-system/internal/repository"
)

type BlogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) *BlogService {
	return &BlogService{repo: repo}
}

// -- articles --
func (s *BlogService) CreateArticle(title, content, author string, tagNames []string) (*domain.Article, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title field cannot be empty")
	}
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content field cannot be empty")
	}
	if strings.TrimSpace(author) == "" {
		return nil, fmt.Errorf("author field cannot be empty")
	}

	article := &domain.Article{
		Title:   title,
		Content: content,
		Author:  author,
	}

	if err := (*s).repo.CreateArticle(article); err != nil {
		return nil, err
	}

	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		tag := &domain.Tag{Name: tagName}
		// will fail is exists already, don't return error
		(*s).repo.CreateTag(tag)

		tags, _ := (*s).repo.GetAllTags()
		for _, t := range tags {
			if t.Name == tagName {
				(*s).repo.AddTagToArticle(article.ID, t.ID)
				break
			}
		}
	}

	return (*s).repo.GetArticle(article.ID)
}

func (s *BlogService) GetArticle(id int) (*domain.Article, error) {
	return (*s).repo.GetArticle(id)
}

func (s *BlogService) GetAllArticles() ([]*domain.Article, error) {
	return (*s).repo.GetAllArticles()
}

func (s *BlogService) UpdateArticle(id int, title, content string) (*domain.Article, error) {
	article, err := (*s).repo.GetArticle(id)
	if err != nil {
		return nil, fmt.Errorf("article not found")
	}

	if strings.TrimSpace(title) != "" {
		article.Title = title
	}
	if strings.TrimSpace(content) != "" {
		article.Content = content
	}

	if err := (*s).repo.UpdateArticle(article); err != nil {
		return nil, err
	}

	return (*s).repo.GetArticle(id)
}

func (s *BlogService) DeleteArticle(id int) error {
	_, err := (*s).repo.GetArticle(id)
	if err != nil {
		return fmt.Errorf("article not found")
	}

	return (*s).repo.DeleteArticle(id)
}

// -- comments --
func (s *BlogService) AddComment(articleID int, author, content string) (*domain.Comment, error) {
	if strings.TrimSpace(author) == "" {
		return nil, fmt.Errorf("author field cannot be empty")
	}
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content field cannot be empty")
	}

	_, err := (*s).repo.GetArticle(articleID)
	if err != nil {
		return nil, fmt.Errorf("article not found")
	}

	comment := &domain.Comment{
		ArticleID: articleID,
		Author:    author,
		Content:   content,
	}

	if err := (*s).repo.CreateComment(comment); err != nil {
		return nil, err
	}

	return comment, nil
}
