package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/danangkonang/cake-store-restful/helper"
	"github.com/danangkonang/cake-store-restful/model"
	"github.com/danangkonang/cake-store-restful/service"
	"github.com/gorilla/mux"
)

type cakeController struct {
	service service.ServiceCake
}

func NewCakeController(fclty service.ServiceCake) *cakeController {
	return &cakeController{
		service: fclty,
	}
}

func (c *cakeController) SaveCake(w http.ResponseWriter, r *http.Request) {
	var p model.ProductPostRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		helper.LoggerError("json", err.Error())
		helper.MakeRespon(w, 400, "invalid json format", nil)
		return
	}
	validationErrors, err := helper.Validation(p)
	if err != nil {
		helper.MakeRespon(w, 400, err.Error(), validationErrors)
		return
	}
	p.CreatedAt = time.Now()

	if err := c.service.SaveCake(&p); err != nil {
		helper.LoggerError("SaveCake", err.Error())
		helper.MakeRespon(w, 400, "internal server error", nil)
		return
	}
	helper.MakeRespon(w, 200, "success", nil)
}

func (c *cakeController) FindCakes(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.FindCakes()
	if err != nil {
		helper.LoggerError("FindCakes", err.Error())
		helper.MakeRespon(w, 400, "internal server error", nil)
		return
	}
	helper.MakeRespon(w, 200, "success", res)
}

func (c *cakeController) FindCakeById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	res, err := c.service.FindCakeById(id)
	switch {
	case err == sql.ErrNoRows:
		helper.MakeRespon(w, 400, "cake not found", nil)
		return
	case err != nil:
		helper.LoggerError("FindCakeById", err.Error())
		helper.MakeRespon(w, 400, "internal server error", nil)
		return
	}
	helper.MakeRespon(w, 200, "success", res)
}

func (c *cakeController) UpdateCake(w http.ResponseWriter, r *http.Request) {
	var p model.ProductUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		helper.LoggerError("json", err.Error())
		helper.MakeRespon(w, 400, "invalid json format", nil)
		return
	}
	defer r.Body.Close()
	validationErrors, err := helper.Validation(p)
	if err != nil {
		helper.MakeRespon(w, 400, err.Error(), validationErrors)
		return
	}
	p.UpdatedAt = time.Now()

	if err := c.service.UpdateCake(&p); err != nil {
		helper.LoggerError("UpdateCake", err.Error())
		helper.MakeRespon(w, 400, "internal server error", nil)
		return
	}
	helper.MakeRespon(w, 200, "success", nil)
}

func (c *cakeController) DeleteCake(w http.ResponseWriter, r *http.Request) {
	var p model.ProductDeleteRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		helper.LoggerError("json", err.Error())
		helper.MakeRespon(w, 400, "invalid json format", nil)
		return
	}
	defer r.Body.Close()
	validationErrors, err := helper.Validation(p)
	if err != nil {
		helper.MakeRespon(w, 400, err.Error(), validationErrors)
		return
	}

	if err := c.service.DeleteCake(&p); err != nil {
		helper.LoggerError("DeleteCake", err.Error())
		helper.MakeRespon(w, 400, "internal server error", nil)
		return
	}
	helper.MakeRespon(w, 200, "success", nil)
}
