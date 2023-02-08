package test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danangkonang/cake-store-restful/config"
	"github.com/danangkonang/cake-store-restful/controller"
	"github.com/danangkonang/cake-store-restful/model"
	"github.com/danangkonang/cake-store-restful/service"
	"github.com/go-playground/assert/v2"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

func CakeServer(db *config.DB) *mux.Router {
	router := mux.NewRouter()
	c := controller.NewCakeController(
		service.NewServiceCake((*config.DB)(db)),
	)
	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/cakes", c.FindCakes).Methods("GET")
	v1.HandleFunc("/cake/{id:[0-9]+}", c.FindCakeById).Methods("GET")
	v1.HandleFunc("/cake", c.SaveCake).Methods("POST")
	v1.HandleFunc("/cake", c.UpdateCake).Methods("PUT")
	v1.HandleFunc("/cake", c.DeleteCake).Methods("DELETE")
	return router
}

type Any struct{}

func (a Any) Match(v driver.Value) bool {
	return true
}

func TestFindCakes(t *testing.T) {
	sqlDB, mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	con := &config.DB{
		Mysql: sqlDB,
	}
	defer con.Mysql.Close()
	godotenv.Load()

	sample := []struct {
		name        string
		expectation int
		mc          *sqlmock.Rows
	}{
		{
			name:        "test succees",
			expectation: 200,
			mc: sqlmock.NewRows(
				[]string{"id", "title", "description", "rating", "image", "created_at", "updated_at"},
			).AddRow(1, "foo", "foo", 1, "bar.jpg", time.Now(), time.Now()),
		},
		{
			name:        "no rows",
			expectation: 200,
			mc: sqlmock.NewRows(
				[]string{"id", "title", "description", "rating", "image", "created_at", "updated_at"},
			),
		},
	}
	for _, v := range sample {
		query := "SELECT id, title, description, rating, image, created_at, updated_at FROM products"
		mock.ExpectQuery(query).WillReturnRows(v.mc)
		request, _ := http.NewRequest("GET", "/api/v1/cakes", nil)
		response := httptest.NewRecorder()
		CakeServer(&config.DB{Mysql: con.Mysql}).ServeHTTP(response, request)
		// body, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(body))
		assert.Equal(t, v.expectation, response.Code)

	}
}

func TestFindCakeById(t *testing.T) {
	sqlDB, mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	con := &config.DB{
		Mysql: sqlDB,
	}
	defer con.Mysql.Close()
	godotenv.Load()
	sample := []struct {
		name        string
		expectation int
		mc          *sqlmock.Rows
	}{
		{
			name:        "test succees",
			expectation: 200,
			mc: sqlmock.NewRows(
				[]string{"id", "title", "description", "rating", "image", "created_at", "updated_at"},
			).AddRow(1, "foo", "foo", 1, "bar.jpg", time.Now(), time.Now()),
		},
		{
			name:        "id not found",
			expectation: 400,
			mc: sqlmock.NewRows(
				[]string{"id", "title", "description", "rating", "image", "created_at", "updated_at"},
			),
		},
	}
	for _, v := range sample {
		query := "SELECT id, title, description, rating, image, created_at, updated_at FROM products WHERE id=?"
		mock.ExpectQuery(query).WithArgs(1).WillReturnRows(v.mc)
		request, _ := http.NewRequest("GET", "/api/v1/cake/1", nil)
		response := httptest.NewRecorder()
		CakeServer(&config.DB{Mysql: con.Mysql}).ServeHTTP(response, request)
		// body, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(body))
		assert.Equal(t, v.expectation, response.Code)
	}
}

func TestCreateCake(t *testing.T) {
	sqlDB, mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	con := &config.DB{
		Mysql: sqlDB,
	}
	defer con.Mysql.Close()
	godotenv.Load()
	sample := []struct {
		name        string
		expectation int
		body        *model.ProductPostRequest
	}{
		{
			name:        "succees",
			expectation: 200,
			body: &model.ProductPostRequest{
				Title:       "foo",
				Description: "bar",
				Rating:      9,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty title",
			expectation: 400,
			body: &model.ProductPostRequest{
				Title:       "",
				Description: "bar",
				Rating:      9,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty description",
			expectation: 400,
			body: &model.ProductPostRequest{
				Title:       "foo",
				Description: "",
				Rating:      9,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty rating",
			expectation: 400,
			body: &model.ProductPostRequest{
				Title:       "foo",
				Description: "bar",
				Rating:      0,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty image",
			expectation: 400,
			body: &model.ProductPostRequest{
				Title:       "foo",
				Description: "bar",
				Rating:      9,
				Image:       "",
			},
		},
		{
			name:        "empty all",
			expectation: 400,
			body: &model.ProductPostRequest{
				Title:       "",
				Description: "",
				Rating:      0,
				Image:       "",
			},
		},
		{
			name:        "empty all",
			expectation: 400,
			body:        &model.ProductPostRequest{},
		},
	}
	for _, v := range sample {
		query := "INSERT INTO products (title, description, rating, image, created_at) VALUES (?, ?, ?, ?, ?)"
		mock.ExpectExec(query).WithArgs("foo", "bar", v.body.Rating, Any{}, Any{}).WillReturnResult(sqlmock.NewResult(0, 0))
		b, err := json.Marshal(v.body)
		if err != nil {
			t.Fatal(err)
		}
		request, _ := http.NewRequest("POST", "/api/v1/cake", bytes.NewBuffer(b))
		response := httptest.NewRecorder()
		CakeServer(&config.DB{Mysql: con.Mysql}).ServeHTTP(response, request)
		// body, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(body))
		assert.Equal(t, v.expectation, response.Code)
	}
}

func TestUpdateCake(t *testing.T) {
	sqlDB, mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	con := &config.DB{
		Mysql: sqlDB,
	}
	defer con.Mysql.Close()
	godotenv.Load()
	sample := []struct {
		name        string
		expectation int
		body        *model.ProductUpdateRequest
	}{
		{
			name:        "succees",
			expectation: 200,
			body: &model.ProductUpdateRequest{
				Id:          1,
				Title:       "foo",
				Description: "bar",
				Rating:      9,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty id",
			expectation: 400,
			body: &model.ProductUpdateRequest{
				Title:       "foo",
				Description: "bar",
				Rating:      9,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "no id found",
			expectation: 400,
			body: &model.ProductUpdateRequest{
				Id:          0,
				Title:       "foo",
				Description: "bar",
				Rating:      9,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty title",
			expectation: 400,
			body: &model.ProductUpdateRequest{
				Id:          1,
				Title:       "",
				Description: "bar",
				Rating:      9,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty description",
			expectation: 400,
			body: &model.ProductUpdateRequest{
				Id:          1,
				Title:       "foo",
				Description: "",
				Rating:      9,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty rating",
			expectation: 400,
			body: &model.ProductUpdateRequest{
				Id:          1,
				Title:       "foo",
				Description: "bar",
				Rating:      0,
				Image:       "http://foo.bar/img/img.jpg",
			},
		},
		{
			name:        "empty image",
			expectation: 400,
			body: &model.ProductUpdateRequest{
				Id:          1,
				Title:       "foo",
				Description: "bar",
				Rating:      9,
				Image:       "",
			},
		},
		{
			name:        "empty all",
			expectation: 400,
			body: &model.ProductUpdateRequest{
				Id:          0,
				Title:       "",
				Description: "",
				Rating:      0,
				Image:       "",
			},
		},
		{
			name:        "empty all",
			expectation: 400,
			body:        &model.ProductUpdateRequest{},
		},
	}
	for _, v := range sample {
		upq := "UPDATE products SET title=?, description=?, rating=?, image=?, updated_at=? WHERE id=?"
		mock.ExpectExec(upq).WithArgs(v.body.Title, v.body.Description, v.body.Rating, v.body.Image, Any{}, v.body.Id).WillReturnResult(sqlmock.NewResult(0, 0))
		b, err := json.Marshal(v.body)
		if err != nil {
			t.Fatal(err)
		}
		request, _ := http.NewRequest("PUT", "/api/v1/cake", bytes.NewBuffer(b))
		response := httptest.NewRecorder()
		CakeServer(&config.DB{Mysql: con.Mysql}).ServeHTTP(response, request)
		// body, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(body))
		assert.Equal(t, v.expectation, response.Code)
	}
}

func TestDeleteCake(t *testing.T) {
	sqlDB, mock, err := NewMock()
	if err != nil {
		t.Fatal(err)
	}
	con := &config.DB{
		Mysql: sqlDB,
	}
	defer con.Mysql.Close()
	godotenv.Load()
	type Body struct {
		Id int `json:"id"`
	}
	sample := []struct {
		name        string
		expectation int
		body        Body
		mc          *sqlmock.Rows
	}{
		{
			name:        "succees",
			expectation: 200,
			body: Body{
				Id: 1,
			},
			mc: sqlmock.NewRows(
				[]string{"id", "title", "description", "rating", "image", "created_at", "updated_at"},
			).AddRow(1, "foo", "foo", 1, "bar.jpg", time.Now(), time.Now()),
		},
		{
			name:        "not found",
			expectation: 400,
			body:        Body{},
			mc: sqlmock.NewRows(
				[]string{"id", "title", "description", "rating", "image", "created_at", "updated_at"},
			).AddRow(1, "foo", "foo", 1, "bar.jpg", time.Now(), time.Now()),
		},
		{
			name:        "not found",
			expectation: 400,
			body: Body{
				Id: 0,
			},
			mc: sqlmock.NewRows(
				[]string{"id", "title", "description", "rating", "image", "created_at", "updated_at"},
			),
		},
	}
	for _, v := range sample {
		mock.ExpectExec("DELETE FROM products WHERE id=?").WithArgs(v.body.Id).WillReturnResult(sqlmock.NewResult(0, 0))

		b, err := json.Marshal(v.body)
		if err != nil {
			t.Fatal(err)
		}
		request, _ := http.NewRequest("DELETE", "/api/v1/cake", bytes.NewBuffer(b))
		response := httptest.NewRecorder()
		CakeServer(&config.DB{Mysql: con.Mysql}).ServeHTTP(response, request)
		// body, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(body))
		assert.Equal(t, v.expectation, response.Code)
	}
}
