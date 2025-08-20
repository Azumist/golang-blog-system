package repository

import (
	"blog-system/internal/domain"
	"database/sql"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// -- articles --
func (r *SQLiteRepository) CreateArticle(article *domain.Article) error {
	query := `INSERT INTO articles (title, content, author) VALUES (?, ?, ?)`
	result, err := (*r).db.Exec(query, (*article).Title, (*article).Content, (*article).Author)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	(*article).ID = int(id)

	return nil
}

func (r *SQLiteRepository) GetArticle(id int) (*domain.Article, error) {
	query := `SELECT id, title, content, author, created_at, updated_at FROM articles WHERE id = ?`
	row := (*r).db.QueryRow(query, id)

	var article domain.Article
	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Author,
		&article.CreatedAt,
		&article.UpdatedAt)
	if err != nil {
		return nil, err
	}

	article.Tags, _ = (*r).GetTagsByArticleID(article.ID)
	article.Comments, _ = (*r).GetCommentsByArticleID(article.ID)

	return &article, nil
}

func (r *SQLiteRepository) GetAllArticles() ([]*domain.Article, error) {
	query := `SELECT id, title, content, author, created_at, updated_at FROM articles ORDER BY created_at DESC`
	rows, err := (*r).db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*domain.Article
	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Author,
			&article.CreatedAt,
			&article.UpdatedAt)
		if err != nil {
			return nil, err
		}

		article.Tags, _ = (*r).GetTagsByArticleID(article.ID)
		articles = append(articles, &article)
	}

	return articles, nil
}

func (r *SQLiteRepository) UpdateArticle(article *domain.Article) error {
	query := `UPDATE articles SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := (*r).db.Exec(query, (*article).Title, (*article).Content, (*article).ID)
	return err
}

func (r *SQLiteRepository) DeleteArticle(id int) error {
	query := `DELETE FROM articles WHERE id = ?`
	_, err := (*r).db.Exec(query, id)
	return err
}

// -- comments --
func (r *SQLiteRepository) CreateComment(comment *domain.Comment) error {
	query := `INSERT INTO comments (article_id, author, content) VALUES (?, ?, ?)`
	result, err := (*r).db.Exec(query, (*comment).ArticleID, (*comment).Author, (*comment).Content)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	(*comment).ID = int(id)
	return nil
}

func (r *SQLiteRepository) GetCommentsByArticleID(articleID int) ([]*domain.Comment, error) {
	query := `SELECT id, article_id, author, content, created_at FROM comments WHERE article_id = ? ORDER BY created_at ASC`
	rows, err := (*r).db.Query(query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ArticleID,
			&comment.Author,
			&comment.Content,
			&comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (r *SQLiteRepository) DeleteComment(id int) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := (*r).db.Exec(query, id)
	return err
}

// -- tags --
func (r *SQLiteRepository) CreateTag(tag *domain.Tag) error {
	query := `INSERT INTO tags (name) VALUES (?)`
	result, err := (*r).db.Exec(query, (*tag).Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	(*tag).ID = int(id)
	return nil
}

func (r *SQLiteRepository) GetAllTags() ([]*domain.Tag, error) {
	query := `SELECT id, name FROM tags ORDER BY name`
	rows, err := (*r).db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}

func (r *SQLiteRepository) GetTagsByArticleID(articleID int) ([]*domain.Tag, error) {
	query := `SELECT t.id, t.name FROM tags t 
		JOIN article_tags at ON t.id = at.tag_id 
		WHERE at.article_id = ?`
	rows, err := (*r).db.Query(query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	return tags, nil
}

func (r *SQLiteRepository) AddTagToArticle(articleID int, tagID int) error {
	query := `INSERT OR IGNORE INTO article_tags (article_id, tag_id) VALUES (?, ?)`
	_, err := (*r).db.Exec(query, articleID, tagID)
	return err
}

func (r *SQLiteRepository) RemoveTagFromArticle(articleID int, tagID int) error {
	query := `DELETE FROM article_tags WHERE article_id = ? AND tag_id = ?`
	_, err := (*r).db.Exec(query, articleID, tagID)
	return err
}
