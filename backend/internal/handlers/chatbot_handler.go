package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/anisetusbambangmanalu19/moonlay-task-management-app/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
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
	// Context: task data from DB + user question + constraint to only answer from context
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

	// --- STEP 4: Call Gemini API ---
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Konfigurasi API chatbot tidak lengkap. Periksa GEMINI_API_KEY di .env"})
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menginisialisasi AI client"})
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.3) // Lower temp for factual, grounded answers

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan respons dari AI"})
		return
	}

	// --- STEP 5: Extract answer from Gemini response ---
	var answer string
	if len(resp.Candidates) > 0 &&
		resp.Candidates[0].Content != nil &&
		len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			answer = string(text)
		}
	}

	if answer == "" {
		answer = "Maaf, saya tidak dapat menghasilkan jawaban saat ini. Coba ulangi pertanyaan Anda."
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
