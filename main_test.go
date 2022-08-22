package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string, body *strings.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPostStatusReport(t *testing.T) {
	body := "{\"send_report\":\"true\"}"
	router := gin.Default()
	testRouter := SetupServerRoutes(router)

	w := performRequest(testRouter, "POST", "/v1/virtualorb/report", strings.NewReader(body))

	assert.Equal(t, http.StatusCreated, w.Code)

	var res map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &res)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res["battery_level_percent"])
	assert.NotNil(t, res["battery_voltage"])
	assert.NotNil(t, res["device_cpu_percent"])
	assert.NotNil(t, res["device_cpu_temp"])
	assert.NotNil(t, res["device_disk_space_used_percent"])
	assert.NotNil(t, res["device_disk_space_available_percent"])
}

func TestPostSignup(t *testing.T) {
	body := "{\"name\":\"thisFakeName\",\"images\":[\"wmf1zgcswrp4qr3oivw2=kxskbcmawb3ahl2y5vv13brs22ncacdopk9h1cbrefu9kcffzw9t9pdwoucb3tg3bnwm3ekuzfj2uz56\"]}"
	router := gin.Default()
	testRouter := SetupServerRoutes(router)

	w := performRequest(testRouter, "POST", "/v1/virtualorb/signup", strings.NewReader(body))

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &res)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res["action_id"])
	assert.NotNil(t, res["message"])
}
