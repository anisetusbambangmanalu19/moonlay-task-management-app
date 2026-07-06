package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// ChatbotHandler handles RAG-based chatbot requests
type ChatbotHandler struct {
	taskRepo *repository.TaskRepository
}

// NewChatbotHandler creates a new ChatbotHandler
func NewChatbotHandler(taskRepo *repository.TaskRepository) *ChatbotHandler {
	return &ChatbotHandler{taskRepo: taskRepo}
}

// ChatRequest defines the expected JSON body for POST /api/chatbot
type ChatRequest struct {
	Question string `json:"question" binding:"required,min=1"`
}

// Gemini REST API request/response structs
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}
type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}
type geminiPart struct {
	Text string `json:"text"`
}
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// Chat handles the chatbot request using RAG (Retrieval-Augmented Generation).
// Flow: user question → query DB for task context → build prompt → Gemini API → return answer
// POST /api/chatbot
func (h *ChatbotHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- STEP 1: Retrieve task context from database (RAG retrieval step) ---
	tasks, err := h.taskRepo.GetTasksForChatbot()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data task untuk chatbot"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"answer": "Belum ada data task yang tersimpan di sistem."})
		return
	}

	// --- STEP 2: Serialize task context to JSON ---
	tasksJSON, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data task"})
		return
	}

	// --- STEP 3: Build RAG prompt ---
	prompt := fmt.Sprintf(`Kamu adalah asisten manajemen task yang cerdas. 
Kamu HANYA boleh menjawab berdasarkan data task yang diberikan di bawah ini.
Jangan menjawab pertanyaan yang tidak berkaitan dengan data task ini.
Jika pertanyaan tidak bisa dijawab dari data yang tersedia, katakan dengan jujur bahwa kamu tidak memiliki informasi tersebut.

Status task: "todo" = belum dikerjakan, "in_progress" = sedang dikerjakan, "done" = selesai.
Format deadline menggunakan UTC timezone.

=== DATA TASK SAAT INI ===
%s
=========================

Pertanyaan: %s

Jawab dalam Bahasa Indonesia dengan jelas, ringkas, dan informatif.`, string(tasksJSON), req.Question)

	// --- STEP 4: Call Gemini REST API directly ---
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Konfigurasi API chatbot tidak lengkap. Periksa GEMINI_API_KEY di .env"})
		return
	}

	answer, err := callGeminiAPI(c.Request.Context(), apiKey, prompt)
	if err != nil {
		log.Printf("Gemini API Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan respons dari AI: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}

// callGeminiAPI calls the Gemini API via direct HTTP REST (no SDK).
// Supports two key formats:
// - "AIza..." (classic API Key) → sent as ?key= query parameter
// - "AQ...."  (new Google AI Studio key) → sent as x-goog-api-key header
func callGeminiAPI(ctx context.Context, apiKey, prompt string) (string, error) {
	// Build URL
	var url string
	isNewFormat := len(apiKey) > 3 && apiKey[:3] == "AQ."
	if isNewFormat {
		// New format keys: send without ?key= param, use header instead
		url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"
	} else {
		url = fmt.Sprintf(
			"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s",
			apiKey,
		)
	}

	body := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// New format keys use x-goog-api-key header
	if isNewFormat {
		httpReq.Header.Set("x-goog-api-key", apiKey)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	// Check for API-level error
	if geminiResp.Error != nil {
		return "", fmt.Errorf("Gemini API error %d: %s", geminiResp.Error.Code, geminiResp.Error.Message)
	}

	// Extract answer
	if len(geminiResp.Candidates) > 0 &&
		len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "Maaf, saya tidak dapat menghasilkan jawaban saat ini. Coba ulangi pertanyaan Anda.", nil
}
