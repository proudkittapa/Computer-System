package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	h := CustomerHandler{}
	h.Initialize()

	r.GET("/customers", h.GetAllCustomer)
	r.GET("/customers/:id", h.GetCustomer)
	r.POST("/customers", h.SaveCustomer)
	r.PUT("/customers/:id", h.UpdateCustomer)
	r.DELETE("/customers/:id", h.DeleteCustomer)

	return r
}

type CustomerHandler struct {
	DB *gorm.DB
}
type Customer struct {
	Id        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

func (h *CustomerHandler) Initialize() {
	db, err := gorm.Open("mysql", "webservice:P@ssw0rd@tcp(127.0.0.1:3306)/db_webservice?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Customer{})

	h.DB = db
}

func (h *CustomerHandler) GetAllCustomer(c *gin.Context) {
	customers := []Customer{}

	h.DB.Find(&customers)

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c *gin.Context) {
	customer := Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := h.DB.Delete(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
func TestGetAllCustomer(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/customers", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)
}

func TestCreateCustomer(t *testing.T) {
	r := setupRouter()

	body := []byte(`{
		"firstName": "John",
		"lastName": "Doe",
		"age": 25,
		"email": "john.doe@mail.com"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)

	var resp map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &resp)

	id = int(resp["id"].(float64))
	// assert field of response body
	assert.Nil(t, err)
	assert.Equal(t, "John", resp["firstName"])
	assert.Equal(t, "Doe", resp["lastName"])
	assert.Equal(t, float64(25), resp["age"])
	assert.Equal(t, "john.doe@mail.com", resp["email"])
}

func TestGetCustomer(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/customers/%d", id), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)

	var resp map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &resp)

	// assert field of response body
	assert.Nil(t, err)
	assert.Equal(t, float64(id), resp["id"])
	assert.Equal(t, "John", resp["firstName"])
	assert.Equal(t, "Doe", resp["lastName"])
	assert.Equal(t, float64(25), resp["age"])
	assert.Equal(t, "john.doe@mail.com", resp["email"])
}

func TestUpdateCustomer(t *testing.T) {
	r := setupRouter()

	body := []byte(`{
		"firstName": "John",
		"lastName": "Doe",
		"age": 26,
		"email": "john.doe@mail.com"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/customers/%d", id), bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)

	var resp map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &resp)

	// assert field of response body
	assert.Nil(t, err)
	assert.Equal(t, float64(id), resp["id"])
	assert.Equal(t, "John", resp["firstName"])
	assert.Equal(t, "Doe", resp["lastName"])
	assert.Equal(t, float64(26), resp["age"])
	assert.Equal(t, "john.doe@mail.com", resp["email"])
}

func TestDeleteCustomer(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/customers/%d", id), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
func main() {
	r := setupRouter()
	r.Run()
}
