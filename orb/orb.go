package orb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Baw25/virtualorb/signup"
	"github.com/Baw25/virtualorb/status_report"
	"github.com/gin-gonic/gin"
)

const (
	MOCK_VORB_BACKEND = "https://virtualorb-mock-backend-2.free.beeceptor.com"
)

type SignupRequestBody struct {
	Images []string `json:"images"`
	Name   string   `json:"name"`
}

type SignupResponseBody struct {
	ActionID string `json:"action_id"`
	Message  string `json:"message"`
}

func PostStatusReport(c *gin.Context) {
	newStatusReport, reportError := status_report.GenerateSingleStatusReport()
	if reportError != nil {
		log.Fatal(reportError)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Problem generating status report!"})
	}

	httpRequest, requestErr := generateReportRequest(newStatusReport, MOCK_VORB_BACKEND+"/report")
	if requestErr != nil {
		log.Fatal(requestErr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Problem submitting status report!"})
	}

	// Submit the status report to backend for review
	httpStatus, postErr := postVorbBackend(&httpRequest)
	if postErr != nil {
		log.Fatal(postErr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Problem submitting status report!"})
	}

	c.IndentedJSON(httpStatus, newStatusReport)
}

func PostSignup(c *gin.Context) {
	var requestBody SignupRequestBody
	var backendRequestBody signup.Signup
	body, requestErr := ioutil.ReadAll(c.Request.Body)
	if requestErr != nil {
		log.Fatal(requestErr)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad signup request!"})
	}

	unmarshallErr := json.Unmarshal(body, &requestBody)
	if unmarshallErr != nil {
		log.Fatal(unmarshallErr)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad signup request!"})
	}

	// Generate random action_id which serves as the id for the images
	actionID := signup.GenerateRandomId()
	backendRequestBody.ActionID = actionID
	// Generate encryption key, but usually would use one from storage
	encryptionKey, keyErr := signup.GenerateEncryptKey()
	if keyErr != nil {
		log.Fatal(keyErr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error on encryption!"})
	}
	// Encrypt the name of person or app
	encryptedName, encryptErr := signup.EncryptStringValue(requestBody.Name, encryptionKey)
	if encryptErr != nil {
		log.Fatal(encryptErr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error on encryption!"})
	}
	backendRequestBody.Name = encryptedName

	// Encrypt the iris images even though they would be encrypted on the backend
	encryptedImages := []string{}
	for i := range requestBody.Images {
		encryptedImage, encryptImageErr := signup.EncryptStringValue(requestBody.Images[i], encryptionKey)
		if encryptImageErr != nil {
			log.Fatal(encryptImageErr)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error on encryption!"})
		}
		encryptedImages = append(encryptedImages, encryptedImage)
	}
	backendRequestBody.Signals = encryptedImages

	httpRequest, requestErr := generateSignupRequest(backendRequestBody, MOCK_VORB_BACKEND+"/signup")
	if requestErr != nil {
		log.Fatal(requestErr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Problem submitting signup!"})
	}

	// Submit the signup to mock backend
	var respBody SignupResponseBody
	httpStatus, postErr := postVorbBackend(&httpRequest)
	if postErr != nil {
		log.Fatal(postErr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Problem submitting signup!"})
	}
	respBody.ActionID = actionID
	respBody.Message = "Signup successfully submitted!"

	c.IndentedJSON(httpStatus, respBody)
}

func generateReportRequest(payload status_report.StatusReport, url string) (http.Request, error) {
	jsonData, err := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return *req, err
	}

	return *req, nil
}

func generateSignupRequest(payload signup.Signup, url string) (http.Request, error) {
	jsonData, err := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return *req, err
	}

	return *req, nil
}

// Function to submit to mock virtualorb backend
func postVorbBackend(request *http.Request) (int, error) {
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return http.StatusInternalServerError, err
	}

	log.Println("Mock Request Successful!", respBody)
	return resp.StatusCode, nil
}
