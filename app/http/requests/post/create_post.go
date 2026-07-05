package post

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,max=255"`
	Content string `json:"content" validate:"required"`
}
