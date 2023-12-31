package models

import (
	"errors"
	"html"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Post model represents a post
type Post struct {
	gorm.Model
	Title          string         `gorm:"size:255;not null;unique" json:"title"`
	PostPermalinks string         `gorm:"size:255" json:"post_permalinks"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	AuthorID       uint           `gorm:"not null" json:"author_id"`
	Author         Profile        `gorm:"foreignKey:AuthorID" json:"author"`
	Tags           pq.StringArray `gorm:"type:text[]" json:"tags"`
	Thumbnails     string         `gorm:"size:255" json:"thumbnails"`
	ReadTime       string         `json:"read_time"`
}

func (p *Post) Prepare() {
	// Sanitize and trim strings
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.PostPermalinks = html.EscapeString(strings.TrimSpace(p.PostPermalinks))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))

	if p.Tags == nil {
		p.Tags = make([]string, 0) // Initialize Tags field as an empty string slice
	}

	if p.Thumbnails == "" {
		p.Thumbnails = "" // Initialize Thumbnails field as an empty string
	}

	if p.ReadTime == "" {
		p.ReadTime = "" // Initialize ReadTime field as an empty string
	}
}

func (p *Post) Validate() map[string]string {
	var err error

	var errorMessages = make(map[string]string)
	if p.Title == "" {
		err = errors.New("Title is required")
		errorMessages["Required_title"] = err.Error()

	}
	if p.Content == "" {
		err = errors.New("Required Content")
		errorMessages["Required_content"] = err.Error()
	}

	if p.AuthorID < 1 {
		err = errors.New("Required Author")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}

func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&p).Preload("Author").Error
	if err != nil {
		return &Post{}, err
	}

	return p, nil
}

func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}

	err = db.Debug().Model(&Post{}).Limit(100).Order("created_at desc").Preload("Author").Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}

	return &posts, nil
}

func (p *Post) FindPostById(db *gorm.DB, pid uint64) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Profile{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {
	var err error

	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Title: p.Title, Content: p.Content, PostPermalinks: p.PostPermalinks, Tags: p.Tags, Thumbnails: p.Thumbnails, ReadTime: p.ReadTime}).Error

	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Profile{}).Where("user_id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}
func (p *Post) DeleteAPost(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Take(&Post{}).Delete(&Post{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
func (p *Post) FindUserPosts(db *gorm.DB, uid uint32) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&posts).Error

	if err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for _, post := range posts {
			err := db.Debug().Model(&Profile{}).Where("id = ?", post.AuthorID).Take(&post.Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

// when a user is deleted, we also delete the post that the user had
func (c *Post) DeleteUserPosts(db *gorm.DB, uid uint32) (int64, error) {
	posts := []Post{}
	db = db.Debug().Model(&Post{}).Where("author_id = ?", uid).Find(&posts).Delete(&posts)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
