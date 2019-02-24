package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/joho/godotenv"
	"github.com/tempo-nksn/Tempo-Backend/models"

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

func performGETRequest(r http.Handler, path string, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performPOSTRequest(r http.Handler, url string, data string) *httptest.ResponseRecorder {
	var dataBytes = []byte(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(dataBytes))
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(req)
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
		DATABASE := os.Getenv("DB_DRIVER")
		databaseURL := os.Getenv("DATABASE_URL")
		if DATABASE == "" && databaseURL == "" {
			err := godotenv.Load("../.env")
			if err != nil {
				log.Fatal("Error loading .env file")
			}
			DATABASE = os.Getenv("DB_DRIVER")
			databaseURL = os.Getenv("DATABASE_URL")
		}
		d, err := gorm.Open(DATABASE, databaseURL)
		if err != nil {
			log.Fatal(err)
			panic("failed to establish database connection")
		}
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
	Describe("Version 1 API at /dashboard", func() {
		Describe("The / endpoint", func() {
			var user models.DashBoard
			BeforeEach(func() {
				var cred = `{
					"u_name": "test3",
					"password": "test3"
				}`
				r := performPOSTRequest(router, "/login", cred)
				type Body struct {
					Code   string
					Expire string
					Token  string
				}
				var body Body
				json.Unmarshal(r.Body.Bytes(), &body)
				// fmt.Println("response is: ", body.Token)
				var bearer = "Bearer " + body.Token
				response = performGETRequest(router, "/api/v1/dashboard/", bearer)
				fmt.Println("response is ", response.Body)
				json.Unmarshal(response.Body.Bytes(), &user)
				// claims := jwt.ExtractClaims()
				// id := claims["id"]
			})
			// var user models.DashBoard

			It("Returns with Status 200", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Returns the user's name", func() {
				Expect(user.Name).To(Equal("Bharath"))
			})

			It("Returns the user's email'", func() {
				Expect(user.Email).To(Equal("test3@gmail.com"))
			})

			It("Returns the user's Phone number'", func() {
				Expect(user.Phone).To(Equal("123456789"))
			})

			It("Returns the user's Wallet'", func() {
				Expect(user.Wallet).To(Equal(int64(13232)))
			})
		})
	})
})
