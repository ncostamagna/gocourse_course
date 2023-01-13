package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/ncostamagna/go_lib_response/response"
	"github.com/ncostamagna/gocourse_course/internal/course"
)

func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {

	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeStoreCourse,
		encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return r

}

func decodeStoreCourse(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var req course.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeGetCourse(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}
	p := mux.Vars(r)
	req := course.GetReq{
		ID: p["id"],
	}

	return req, nil
}

func decodeGetAllCourse(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := course.GetAllReq{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil
}

func decodeUpdateCourse(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var req course.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	path := mux.Vars(r)
	req.ID = path["id"]

	return req, nil
}

func decodeDeleteCourse(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	path := mux.Vars(r)
	req := course.DeleteReq{
		ID: path["id"],
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(resp)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)

	w.WriteHeader(resp.StatusCode())

	_ = json.NewEncoder(w).Encode(resp)
}

func authorization(token string) error {
	if token != os.Getenv("TOKEN") {
		return errors.New("invalid token")
	}
	return nil
}
