package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock OpenAI API response
func mockOpenAIServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp := ChatCompletionResponse{
			Choices: []Choice{
				{Message: Message{Content: "Mocked response"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestCompletionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	openAIServer := mockOpenAIServer()
	defer openAIServer.Close()

	router.POST("/completion", CompletionHandler)

	t.Run("Successful completion request", func(t *testing.T) {
		prompt := `Write me an itinerary for 5 days to Japan in the coming April. Describe the weather that month, and also 5 things to take note about this country's culture. Keep to a maximum travel area to the size of Hokkaido, if possible, to minimize traveling time between cities.\n\nFor each day, list me the following:\n- Attractions suitable for that season\n- Hotel (prefer not to change it unless traveling to another city)\n- 2 Restaurants, one for lunch and another for dinner, with shortened Google Map links\nand give me a daily summary of the above points into a paragraph or two.`
		requestBody := `{"prompt": "` + prompt + `", "session": "mock-session"}`
		req, _ := http.NewRequest("POST", "/completion", bytes.NewBuffer([]byte(requestBody)))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		expectedResponse := `{"response":"Mocked response","status":200}`
		assert.JSONEq(t, expectedResponse, resp.Body.String())
	})

	t.Run("Invalid JSON request", func(t *testing.T) {
		requestBody := `{"prompt": }` // Invalid JSON
		req, _ := http.NewRequest("POST", "/completion", bytes.NewBuffer([]byte(requestBody)))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		expectedResponse := `{"error":"invalid request","details":"invalid character '}' looking for beginning of value"}`
		assert.JSONEq(t, expectedResponse, resp.Body.String())
	})

	t.Run("OpenAI API request failure", func(t *testing.T) {
		prompt := `Write me an itinerary for 5 days to Japan in the coming April. Describe the weather that month, and also 5 things to take note about this country's culture. Keep to a maximum travel area to the size of Hokkaido, if possible, to minimize traveling time between cities.\n\nFor each day, list me the following:\n- Attractions suitable for that season\n- Hotel (prefer not to change it unless traveling to another city)\n- 2 Restaurants, one for lunch and another for dinner, with shortened Google Map links\nand give me a daily summary of the above points into a paragraph or two.`
		requestBody := `{"prompt": "` + prompt + `", "session": "mock-session"}`
		req, _ := http.NewRequest("POST", "/completion", bytes.NewBuffer([]byte(requestBody)))
		req.Header.Set("Content-Type", "application/json")

		// Temporarily replace the OpenAI API URL with an invalid one to simulate failure
		http.DefaultClient = &http.Client{}
		http.DefaultClient.Transport = &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse("http://invalid-url")
			},
		}

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		expectedResponse := `{"error":"failed to request OpenAI API","status":500}`
		assert.JSONEq(t, expectedResponse, resp.Body.String())

		// Restore the original transport
		http.DefaultClient.Transport = nil
	})
}
