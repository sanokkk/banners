package suite

import (
	"banner-service/internal/app/http/requests"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/lib/pq"
	"io"
	"net/http"
	"strconv"
	"testing"
)

// with running server
func Test_GetUserBanner_WithoutToken_ReturnsUnauthorized(t *testing.T) {
	client := http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:8080/user_banner?tag_id=2&feature_id=1&use_last_revision=true", nil)
	//req.Header.Set("token", "user_token")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	assert.Equal(t, response.StatusCode, http.StatusUnauthorized)
}

// with running server
func TestGetUserBanner_WithoutExistingBanner_ReturnsNotFound(t *testing.T) {
	client := http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:8080/user_banner?tag_id=2&feature_id=99999&use_last_revision=false", nil)
	req.Header.Set("token", "user_token")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	assert.Equal(t, response.StatusCode, http.StatusNotFound)
}

// with running server
func TestGetUserBanner_WithInvalidParameters_ReturnsBadRequest(t *testing.T) {
	client := http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:8080/user_banner?tag_id=2&feature_id=-1&use_last_revision=true", nil)
	req.Header.Set("token", "user_token")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	assert.Equal(t, response.StatusCode, http.StatusBadRequest)
}

// with running server
func TestGetUserBanner_WithCreatedBanner_ReturnsOk(t *testing.T) {
	client := http.Client{}

	var buffer bytes.Buffer
	body := requests.CreateBannerRequest{TagIds: pq.Int32Array{1, 2, 3}, FeatureId: 1, Content: "content", IsActive: true}
	if err := json.NewEncoder(&buffer).Encode(body); err != nil {
		t.Error(err)
	}
	createReq, _ := http.NewRequest("POST", "http://localhost:8080/banner", &buffer)
	createReq.Header.Set("token", "admin_token")

	response, err := client.Do(createReq)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	id, err := strconv.Atoi(string(bodyBytes))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	req, _ := http.NewRequest("GET", "http://localhost:8080/user_banner?tag_id=2&feature_id=1&use_last_revision=true", nil)
	req.Header.Set("token", "user_token")

	getResponse, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	assert.Equal(t, getResponse.StatusCode, http.StatusOK)

	deleteReq, _ := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/banner/%d", id), &buffer)
	deleteReq.Header.Set("token", "admin_token")

	deleteResponse, err := client.Do(deleteReq)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	assert.Equal(t, deleteResponse.StatusCode, http.StatusOK)
}
