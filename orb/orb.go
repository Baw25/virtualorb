package orb

import (
	"fmt"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Baw25/virtualorb/status_report"
	"github.com/gin-gonic/gin"
)

const (
	MOCK_VORB_BACKEND = "https://virtualorb-mock-backend.free.beeceptor.com"
)

func PostStatusReport(c *gin.Context) {
	newStatusReport, reportError := status_report.GenerateSingleStatusReport()
	if reportError != nil {
		log.Fatal(reportError)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Problem generating status report!"})
	}

	// Submit the status report to backend for review
	httpStatus, submitErr := postVorbBackend(newStatusReport)
	if submitErr != nil {
		log.Fatal(submitErr)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Problem submitting status report!"})
	}

	c.IndentedJSON(httpStatus, newStatusReport)
}

func PostSignup(c *gin.Context) {
	fmt.Println("POSTING THE SIGNUPS")
}

func PostSignupBatch(c *gin.Context) {
	fmt.Println("POSTING THE SIGNUPS IN BATCHES")
}

func postVorbBackend(report status_report.StatusReport) (int, error) {
	jsonData, err := json.Marshal(report)

	req, err := http.NewRequest(http.MethodPost, MOCK_VORB_BACKEND+"/report", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return http.StatusBadRequest, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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

	log.Println("Submitted Report Sucessfully!", respBody)
	return resp.StatusCode, nil
}
