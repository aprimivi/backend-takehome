package controller

import (
	"errors"
	"net/http"
	"strconv"

	helper "app/helpers"
	commentrequests "app/http/requests/comment"
	"app/services"

	"github.com/go-chi/chi/v5"
)

type CommentController struct {
	commentService *services.CommentService
}

func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

func (c *CommentController) Create(w http.ResponseWriter, r *http.Request) {
	postID, err := parseIDParam(r, "postID")
	if err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, "invalid post id", nil)
		return
	}

	var req commentrequests.CreateCommentRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		helper.WriteResponse(w, http.StatusUnprocessableEntity, fieldErrors[0].Message, fieldErrors)
		return
	}

	comment, err := c.commentService.Create(r.Context(), postID, req.AuthorName, req.Content)
	if err != nil {
		if errors.Is(err, services.ErrPostNotFound) {
			helper.WriteResponse(w, http.StatusNotFound, "post not found", nil)
			return
		}
		helper.WriteResponse(w, http.StatusInternalServerError, "failed to create comment", nil)
		return
	}

	helper.WriteResponse(w, http.StatusCreated, "comment created", comment)
}

func (c *CommentController) GetAllByPost(w http.ResponseWriter, r *http.Request) {
	postID, err := parseIDParam(r, "postID")
	if err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, "invalid post id", nil)
		return
	}

	comments, err := c.commentService.GetAllByPost(r.Context(), postID)
	if err != nil {
		if errors.Is(err, services.ErrPostNotFound) {
			helper.WriteResponse(w, http.StatusNotFound, "post not found", nil)
			return
		}
		helper.WriteResponse(w, http.StatusInternalServerError, "failed to fetch comments", nil)
		return
	}

	helper.WriteResponse(w, http.StatusOK, "comments fetched", comments)
}

func parseIDParam(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, name), 10, 64)
}
