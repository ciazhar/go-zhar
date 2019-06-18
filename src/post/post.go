package domain

import (
	"github.com/ciazhar/zhar/rest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Post struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	Active        bool          `json:"active,omitempty"`
	ApplicationID string        `json:"applicationId" bson:"applicationId"`
	AuthorId      bson.ObjectId `json:"authorId,omitempty" bson:"-"`
	AuthorRef     *mgo.DBRef    `json:"-" bson:"author"`
	//Author			*Profile		`json:"author,omitempty" bson:"-"`
	Content         string          `json:"content"`
	Title           string          `json:"title"`
	Slug            string          `json:"slug"`
	PostCategoryID  []bson.ObjectId `json:"postCategoryId,omitempty" bson:"-"`
	PostCategoryRef []*mgo.DBRef    `json:"-" bson:"postCategory"`
	//PostCategory	[]*PostCategory	`json:"postCategory,omitempty" bson:"-"`
	ThumbnailImageUrl string    `json:"thumbnailImageUrl" bson:"thumbnailImageUrl"`
	ContentUrl        string    `json:"contentUrl" bson:"contentUrl"`
	IsPublished       bool      `json:"isPublished" bson:"isPublished"`
	IsLiked           bool      `json:"isLiked" bson:"-"`
	ViewCount         *int      `json:"viewCount" bson:"-"`
	LikeCount         *int      `json:"likeCount" bson:"-"`
	CreatedAt         time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt" bson:"updatedAt"`
}

type PostMongodbRepository interface {
	GetCollection() string
	Insert(post *Post) error
	Find(q rest.Query) ([]Post, error)
	FindById(id string) (Post, error)
	Update(post *Post) error
	Delete(id string) error
}

type PostUseCase interface {
	Insert(post *Post) error
	Find(q rest.Query) ([]Post, error)
	FindByAuthHeader(q rest.Query, authHeader string) ([]Post, error)
	FindById(id string) (Post, error)
	Update(post *Post) error
	Patch(post *Post) error
	Delete(id string) error
}
