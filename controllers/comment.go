package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/godiscourse/godiscourse/middleware"
	"github.com/godiscourse/godiscourse/models"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/views"
)

type commentRequest struct {
	TopicID string `json:"topic_id"`
	Body    string `json:"body"`
}

type commentImpl struct{}

func registerComment(router *httptreemux.TreeMux) {
	impl := &commentImpl{}

	router.POST("/comments", impl.create)
	router.GET("/topics/:id/comments", impl.comments)
}

func (impl *commentImpl) create(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var body commentRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		views.RenderErrorResponse(w, r, session.BadRequestError(r.Context()))
		return
	}
	if comment, err := middleware.CurrentUser(r).CreateComment(r.Context(), body.TopicID, body.Body); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComment(w, r, comment)
	}
}

func (impl *commentImpl) comments(w http.ResponseWriter, r *http.Request, params map[string]string) {
	offset, _ := time.Parse(time.RFC3339Nano, r.URL.Query().Get("offset"))
	if topic, err := models.ReadTopic(r.Context(), params["id"]); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else if topic == nil {
		views.RenderErrorResponse(w, r, session.NotFoundError(r.Context()))
	} else if comments, err := topic.ReadComments(r.Context(), offset); err != nil {
		views.RenderErrorResponse(w, r, err)
	} else {
		views.RenderComments(w, r, comments)
	}
}