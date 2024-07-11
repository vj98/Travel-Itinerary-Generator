package api

import (
	"backend/internal/config"
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int              `json:"index"`
	Message      Message          `json:"message"`
	LogProbs     *json.RawMessage `json:"logprobs"` // Use RawMessage if the structure is complex or not needed directly
	FinishReason string           `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Open AI api integration
func CompletionHandler(c *gin.Context) {
	var reqData struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	log.Println("prompt ", reqData.Prompt)

	requestData := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: reqData.Prompt,
			},
		},
		Temperature: 0.8,
		MaxTokens:   1500,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode request data"})
		return
	}

	// Prepare and send the HTTP request to OpenAI
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request", "status": http.StatusInternalServerError})
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error", "status": http.StatusInternalServerError})
		c.Abort()
		return
	}

	token := "Bearer " + cfg.OpenAIKey
	log.Println("token ", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to request OpenAI API", "status": http.StatusInternalServerError})
		return
	}
	defer resp.Body.Close()

	// Decode the response from OpenAI
	var compResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&compResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode OpenAI response", "status": http.StatusInternalServerError})
		return
	}

	log.Println("choices ", compResp.Choices[0].Message.Content)

	if len(compResp.Choices) > 0 {
		c.JSON(http.StatusOK, gin.H{"response": compResp.Choices[0].Message.Content, "status": http.StatusOK})
	} else {
		c.JSON(http.StatusOK, gin.H{"response": "No completion found."})
	}
}

// dummy api for testing
func DummyCall(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "Day 1: Arrival in Bangkok\n- Explore the Grand Palace and Wat Pho\n- Stay at Chatrium Hotel Riverside Bangkok\n- Weather: April is the hottest month in Thailand, with average temperatures around 30-35°C. Be prepared for high humidity.\n\nDay 2: Bangkok\n- Visit Wat Arun and take a boat ride along the Chao Phraya River\n- Explore the bustling markets and street food stalls\n- Stay at Chatrium Hotel Riverside Bangkok\n\nDay 3: Ayutthaya\n- Day trip to the ancient city of Ayutthaya, a UNESCO World Heritage Site\n- Visit the historical temples and ruins\n- Stay at Sala Ayutthaya\n\nDay 4: Chiang Mai\n- Fly to Chiang Mai\n- Explore the Old City and visit temples such as Wat Chedi Luang and Wat Phra Singh\n- Stay at Tamarind Village Chiang Mai\n- Weather: Chiang Mai is cooler than Bangkok in April, with temperatures around 25-30°C.\n\nDay 5: Chiang Mai\n- Visit the Doi Suthep Temple and enjoy the panoramic views of the city\n- Explore the local markets and try traditional Northern Thai dishes\n- Stay at Tamarind Village Chiang Mai\n\nDay 6: Chiang Rai\n- Day trip to the White Temple (Wat Rong Khun) and the Golden Triangle\n- Explore the unique art and architecture of the White Temple\n- Stay at Le Meridien Chiang Rai Resort\n\nDay 7: Departure\n- Fly back to Bangkok for departure\n- Last-minute shopping at the local markets\n- Weather: April marks the beginning of the rainy season in Thailand, so be prepared for occasional showers.\n\n5 things to note about Thai culture:\n1. Respect for the monarchy: Thai people hold their King and royal family in high regard, and any form of disrespect is considered offensive.\n2. Buddhist customs: Thailand is a predominantly Buddhist country, so it's important to respect religious sites and practices.\n3. Etiquette: Thais value politeness and courtesy, so be mindful of your behavior and gestures while interacting with locals.\n4. Dress code: When visiting temples or other religious sites, modest attire is required. Avoid wearing revealing clothing.\n5. Food customs: Thai cuisine is known for its spicy and flavorful dishes. It's common to share dishes family-style and use a spoon and fork for eating.", "status": http.StatusOK})
}
