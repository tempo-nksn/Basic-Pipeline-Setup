package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var _ = Describe("Server", func() {
	var (
		db_test  *gorm.DB
		router   *gin.Engine
		response *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		d, err := gorm.Open("postgres", "postgres://fpfujvlpxoelcm:93ec17a8d9323c29e05b569ea2bd77fa2d7dc96564d2f5b6eaa521807bf8b787@ec2-54-243-128-95.compute-1.amazonaws.com:5432/d34p830c249n99")
		db_test = d
		if err != nil {
			panic(err)
		}
		router = CreateRouter(db_test)
	})

	Describe("Version 1 API at /api/v1", func() {
		Describe("The / endpoint", func() {
			BeforeEach(func() {
				response = performRequest(router, "GET", "/api/v1/")
			})

			It("Returns with Status 200", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Returns the String 'Hello User, your taxi is booked'", func() {
				Expect(response.Body.String()).To(Equal("Hello User, your taxi is booked"))
			})
		})
	})
})

func TestHello(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	fmt.Println(req)
	fmt.Println(rec)
}
