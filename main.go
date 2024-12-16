package main

import (
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}
var store = sessions.NewCookieStore([]byte("my-key"))
var statisticService = &service.StatisticService{}
var filterService = &service.FilterService{}

// mengambil session baru
func getSession(r *http.Request) *sessions.Session {
	session, err := store.Get(r, "chat-session")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		session, _ = store.New(r, "chat-session")
	}
	return session
}

// mengkonversi map ke format csv
func mapToCSV(data map[string][]string) (string, error) {
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// Write header (keys of the map)
	header := make([]string, 0, len(data))
	for key := range data {
		header = append(header, key)
	}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Write rows (values of the map)
	numRows := len(data[header[0]]) // Assuming all columns have the same number of rows
	for i := 0; i < numRows; i++ {
		row := make([]string, len(header))
		for j, key := range header {
			row[j] = data[key][i]
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	// Set up the router
	router := mux.NewRouter()

	// File upload endpoint
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Ambil file dari request
		file, _, err := r.FormFile("file")
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Gagal mengambil file: " + err.Error(),
			})
			return
		}
		defer file.Close()

		// Baca file menjadi string
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Gagal membaca file: " + err.Error(),
			})
			return
		}

		// Proses file
		data, err := fileService.ProcessFile(string(fileBytes))
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Gagal memproses file: " + err.Error(),
			})
			return
		}

		// Analisa data menggunakan Tapas untuk mencari appliance dengan konsumsi energi tertinggi dan terendah
		token := os.Getenv("HUGGINGFACE_TOKEN")
		//Pemisahan query menjadi query_highest dan query_lowest (karena huggingface tidak bisa mereturn 2 nilai untuk highest dan lowest)
		query_highest := "what's the appliance with the most energy consumption?"
		result_highest, err := aiService.AnalyzeData(data, query_highest, token)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Gagal menganalisa data: " + err.Error(),
			})
			return
		}

		query_lowest := "what's the appliance with the least energy consumption?"
		result_lowest, err := aiService.AnalyzeData(data, query_lowest, token)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Gagal menganalisa data: " + err.Error(),
			})
			return
		}

		// Simpan data ke session dengan format yang benar
		session := getSession(r)
		session.Values["data"] = data
		err = session.Save(r, w)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Gagal menyimpan session: " + err.Error(),
			})
			return
		}

		//mengembalikan hasil analisis dari tapas
		annotated_result := fmt.Sprintf("From the provided data, the appliance with the least electricity usage is the %s, and the one with the most electricity usage is the %s.", result_lowest, result_highest)

		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"answer": annotated_result,
		})
	}).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//mengambil data dari session
		// Ambil session
		session := getSession(r)

		// Cek apakah data ada di session
		data, exists := session.Values["data"]
		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Harap upload file terlebih dahulu",
			})
			return
		}

		// Type assertion dengan pengecekan yang lebih baik
		tableData, ok := data.(map[string][]string)
		if !ok {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Format data tidak valid, silakan upload ulang file",
			})
			return
		}

		//menerima query dari user
		// Decode request body
		var chatRequest struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&chatRequest); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Invalid request format: " + err.Error(),
			})
			return
		}

		// Proses chat dengan AI, mengkonversi data ke format csv
		token := os.Getenv("HUGGINGFACE_TOKEN")
		context, err := mapToCSV(tableData)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Gagal memproses data: " + err.Error(),
			})
			return
		}
		//mengirim query ke AI service
		context = "Energy usage CSV:\n" + context
		response, err := aiService.ChatWithAI(context, chatRequest.Query, token)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"answer": "Gagal memproses chat: " + err.Error(),
			})
			return
		}

		//mengembalikan respons AI
		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"answer": response.GeneratedText,
		})
	}).Methods("POST")

	// Endpoint statistik
	router.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Ambil session
		session := getSession(r)

		// Cek apakah data ada di session
		data, exists := session.Values["data"]
		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Harap upload file terlebih dahulu",
			})
			return
		}

		// Type assertion
		tableData, ok := data.(map[string][]string)
		if !ok {
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Format data tidak valid",
			})
			return
		}

		// Hitung statistik
		stats, err := statisticService.CalculateStatistics(tableData)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Gagal menghitung statistik: " + err.Error(),
			})
			return
		}

		// Format response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   stats,
		})
	}).Methods("GET")

	// Tambahkan endpoint filter
	router.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Ambil session
		session := getSession(r)

		// Ambil data dari session
		data, exists := session.Values["data"]
		if !exists {
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Data belum diupload",
			})
			return
		}

		// Type assertion untuk data
		tableData, ok := data.(map[string][]string)
		if !ok {
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Format data tidak valid",
			})
			return
		}

		// Buat filter request dari query parameters
		filter := model.FilterRequest{
			Room:      r.URL.Query().Get("room"),
			Appliance: r.URL.Query().Get("appliance"),
			Date:      r.URL.Query().Get("date"),
		}

		// Filter data menggunakan filter service
		filteredData, err := filterService.FilterData(tableData, filter)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Gagal memfilter data: " + err.Error(),
			})
			return
		}

		// Kirim response
		json.NewEncoder(w).Encode(model.FilterResponse{
			Status: "success",
			Data:   filteredData,
		})
	}).Methods("GET")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}

// menggunakan gob register agar sessions dapat berjalan
func init() {
	//Mendaftarkan tipe data map[string][]string ke package gob
	gob.Register(map[string][]string{})
}
