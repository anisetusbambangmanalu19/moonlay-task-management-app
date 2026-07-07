package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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

// geminiAPIError wraps the HTTP status code so callers can decide whether to retry.
type geminiAPIError struct {
	StatusCode int
	Message    string
}

func (e *geminiAPIError) Error() string {
	return fmt.Sprintf("Gemini API error %d: %s", e.StatusCode, e.Message)
}

// isRetryable reports whether an error is transient and worth retrying
// (overload / rate limit), as opposed to a permanent error (bad request, auth, etc).
func isRetryable(err error) bool {
	apiErr, ok := err.(*geminiAPIError)
	if !ok {
		return false
	}
	return apiErr.StatusCode == http.StatusServiceUnavailable || // 503, model overloaded
		apiErr.StatusCode == http.StatusTooManyRequests || // 429, rate limited
		apiErr.StatusCode >= 500 // any other transient server-side error
}

// Chat handles the chatbot request using RAG (Retrieval-Augmented Generation).
// Flow: user question -> query DB for task context -> build prompt -> Gemini API -> return answer
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
	prompt := fmt.Sprintf(`Kamu adalah asisten manajemen task (chatbot) yang cerdas dan ramah di aplikasi Moonlay Task Management.
Tugas utamamu adalah membantu pengguna mengelola, memantau, dan menjawab pertanyaan terkait data task mereka.

Panduan menjawab:
1. Jika pengguna menyapa (misal: "hai", "halo", "selamat pagi"), balaslah dengan ramah dan tawarkan bantuan terkait task.
2. Jika pengguna bertanya di luar konteks manajemen task atau di luar data yang diberikan, tolak dengan sopan dan ingatkan bahwa kamu adalah asisten khusus task management.
3. Jawab pertanyaan terkait task secara akurat HANYA berdasarkan DATA TASK SAAT INI di bawah.
4. Gunakan Bahasa Indonesia yang jelas, ringkas, dan profesional namun santai.

Info Status task: "todo" = belum dikerjakan, "in_progress" = sedang dikerjakan, "done" = selesai.
Format deadline menggunakan UTC timezone.

=== DATA TASK SAAT INI ===
%s
=========================

Pertanyaan pengguna: %s

Jawabanmu:`, string(tasksJSON), req.Question)

	// --- STEP 4: Call Gemini REST API directly, with retry on transient errors ---
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Konfigurasi API chatbot tidak lengkap. Periksa GEMINI_API_KEY di .env"})
		return
	}

	answer, err := callGeminiAPIWithRetry(c.Request.Context(), apiKey, prompt, 3)
	if err != nil {
		log.Printf("Gemini API Error: %v", err)
		if isRetryable(err) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Server AI sedang sibuk. Silakan coba lagi dalam beberapa saat.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan respons dari AI: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}

// callGeminiAPIWithRetry wraps callGeminiAPI with retry + exponential backoff,
// but only retries on transient errors (503 overload, 429 rate limit, 5xx).
// Permanent errors (bad request, invalid key, etc) fail immediately.
func callGeminiAPIWithRetry(ctx context.Context, apiKey, prompt string, maxAttempts int) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		answer, err := callGeminiAPI(ctx, apiKey, prompt)
		if err == nil {
			return answer, nil
		}

		lastErr = err

		if !isRetryable(err) || attempt == maxAttempts {
			return "", err
		}

		// Exponential backoff with jitter: ~1s, ~2s, ~4s + random 0-300ms
		backoff := time.Duration(1<<uint(attempt-1))*time.Second + time.Duration(rand.Intn(300))*time.Millisecond
		log.Printf("Gemini API attempt %d/%d failed (%v), retrying in %v", attempt, maxAttempts, err, backoff)

		select {
		case <-time.After(backoff):
			// continue to next attempt
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	return "", lastErr
}

// callGeminiAPI calls the Gemini API via direct HTTP REST (no SDK needed).
// Auth is sent via the x-goog-api-key header, which works for all standard
// Gemini API keys (format: "AIza...", generated at aistudio.google.com).
func callGeminiAPI(ctx context.Context, apiKey, prompt string) (string, error) {
	const url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"

	// Give each individual attempt its own timeout so one slow call
	// doesn't block the whole request indefinitely.
	reqCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	body := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", apiKey)

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

	// Prefer the structured error from Gemini's JSON body when present,
	// otherwise fall back to the HTTP status code.
	if geminiResp.Error != nil {
		return "", &geminiAPIError{StatusCode: geminiResp.Error.Code, Message: geminiResp.Error.Message}
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return "", &geminiAPIError{StatusCode: resp.StatusCode, Message: string(respBytes)}
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "Maaf, saya tidak dapat menghasilkan jawaban saat ini. Coba ulangi pertanyaan Anda.", nil
}
