package controller

import (
	"errors"
	"net/http"
	"strconv"

	helper "app/helpers"
	middleware "app/http/middlewares"
	postrequests "app/http/requests/post"
	"app/services"

	"github.com/go-chi/chi/v5"
)

type PostController struct {
	postService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{postService: postService}
}

func (c *PostController) Create(w http.ResponseWriter, r *http.Request) {
	authorID, ok := middleware.UserID(r)
	if !ok {
		helper.WriteResponse(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req postrequests.CreatePostRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		helper.WriteResponse(w, http.StatusUnprocessableEntity, fieldErrors[0].Message, fieldErrors)
		return
	}

	post, err := c.postService.Create(r.Context(), authorID, req.Title, req.Content)
	if err != nil {
		helper.WriteResponse(w, http.StatusInternalServerError, "failed to create post", nil)
		return
	}

	helper.WriteResponse(w, http.StatusCreated, "post created", post)
}

func (c *PostController) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := c.postService.GetAll(r.Context())
	if err != nil {
		helper.WriteResponse(w, http.StatusInternalServerError, "failed to fetch posts", nil)
		return
	}

	helper.WriteResponse(w, http.StatusOK, "posts fetched", posts)
}

func (c *PostController) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parsePostID(r)
	if err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, "invalid post id", nil)
		return
	}

	post, err := c.postService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrPostNotFound) {
			helper.WriteResponse(w, http.StatusNotFound, "post not found", nil)
			return
		}
		helper.WriteResponse(w, http.StatusInternalServerError, "failed to fetch post", nil)
		return
	}

	helper.WriteResponse(w, http.StatusOK, "post fetched", post)
}

func (c *PostController) Update(w http.ResponseWriter, r *http.Request) {
	authorID, ok := middleware.UserID(r)
	if !ok {
		helper.WriteResponse(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := parsePostID(r)
	if err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, "invalid post id", nil)
		return
	}

	var req postrequests.UpdatePostRequest
	if err := helper.DecodeJSON(r, &req); err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if fieldErrors := helper.ValidateStruct(req); fieldErrors != nil {
		helper.WriteResponse(w, http.StatusUnprocessableEntity, fieldErrors[0].Message, fieldErrors)
		return
	}

	post, err := c.postService.Update(r.Context(), id, authorID, req.Title, req.Content)
	if err != nil {
		writePostError(w, err)
		return
	}

	helper.WriteResponse(w, http.StatusOK, "post updated", post)
}

func (c *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	authorID, ok := middleware.UserID(r)
	if !ok {
		helper.WriteResponse(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := parsePostID(r)
	if err != nil {
		helper.WriteResponse(w, http.StatusBadRequest, "invalid post id", nil)
		return
	}

	if err := c.postService.Delete(r.Context(), id, authorID); err != nil {
		writePostError(w, err)
		return
	}

	helper.WriteResponse(w, http.StatusOK, "post deleted", nil)
}

func parsePostID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func writePostError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrPostNotFound):
		helper.WriteResponse(w, http.StatusNotFound, "post not found", nil)
	case errors.Is(err, services.ErrForbidden):
		helper.WriteResponse(w, http.StatusForbidden, err.Error(), nil)
	default:
		helper.WriteResponse(w, http.StatusInternalServerError, "something went wrong", nil)
	}
}
