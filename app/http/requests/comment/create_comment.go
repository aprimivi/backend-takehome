package comment

type CreateCommentRequest struct {
	AuthorName string `json:"author_name" validate:"required,max=100"`
	Content    string `json:"content" validate:"required"`
}
