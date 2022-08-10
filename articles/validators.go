package articles

import (
	"first/common"
	"first/users"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type ArticleModelValidator struct {
	Article struct {
		Title string                `form:"title" json:"title"`
		Image *multipart.FileHeader `form:"image" json:"description" `
		Body  string                `form:"body" json:"body"`
		Tags  string                `form:"tagList" json:"tagList"`
		Price float32               `form:"price" json:"price"`
	} `json:"article"`
	articleModel ArticleModel `json:"-"`
}

func NewArticleModelValidator() ArticleModelValidator {
	return ArticleModelValidator{}
}

func NewArticleModelValidatorFillWith(articleModel ArticleModel) ArticleModelValidator {
	articleModelValidator := NewArticleModelValidator()
	articleModelValidator.Article.Title = articleModel.Title
	// articleModelValidator.Article.Image = articleModel.Image
	articleModelValidator.Article.Body = articleModel.Body
	// for _, tagModel := range articleModel.Tags {
	// 	articleModelValidator.Article.Tags = append(articleModelValidator.Article.Tags, tagModel.Tag)
	// }
	return articleModelValidator
}

func (s *ArticleModelValidator) Bind(c *gin.Context) error {
	myUserModel := c.MustGet("my_user_model").(users.UserModel)

	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.articleModel.Slug = slug.Make(s.Article.Title)
	s.articleModel.Title = s.Article.Title
	s.articleModel.Image = s.Article.Image.Filename
	s.articleModel.Body = s.Article.Body
	s.articleModel.Price = s.Article.Price
	s.articleModel.Author = GetArticleUserModel(myUserModel)
	s.articleModel.setTags(s.Article.Tags)
	return nil
}

type CommentModelValidator struct {
	Comment struct {
		Body string `form:"body" json:"body" binding:"max=2048"`
	} `json:"comment"`
	commentModel CommentModel `json:"-"`
}

func NewCommentModelValidator() CommentModelValidator {
	return CommentModelValidator{}
}

func (s *CommentModelValidator) Bind(c *gin.Context) error {
	myUserModel := c.MustGet("my_user_model").(users.UserModel)

	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.commentModel.Body = s.Comment.Body
	s.commentModel.Author = GetArticleUserModel(myUserModel)
	return nil
}
