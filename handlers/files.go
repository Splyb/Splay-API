package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type UploadResponse struct {
	FileURL string `json:"file_url"`
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := header.Filename
	if !isValidFile(fileName) {
		http.Error(w, "Archivo no permitido", http.StatusForbidden)
		return
	}

	uploadPath := "./uploads/" + fileName
	os.MkdirAll("./uploads", os.ModePerm)

	out, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error al escribir el archivo", http.StatusInternalServerError)
		return
	}

	// Subir a Supabase Storage
	supabaseURL := os.Getenv("SUPABASE_URL")
	storageBucket := "music"
	anonKey := os.Getenv("SUPABASE_ANON_KEY")

	fileURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, storageBucket, fileName)
	err = uploadToSupabaseStorage(supabaseURL, storageBucket, fileName, uploadPath, anonKey)
	if err != nil {
		http.Error(w, "Error al subir el archivo a Supabase", http.StatusInternalServerError)
		return
	}

	// Insertar metadatos en la tabla usando Supabase API
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	err = insertTrackToDatabase(supabaseURL, serviceKey, fileName, "Unknown Artist", fileURL)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	response := UploadResponse{FileURL: fileURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func isValidFile(fileName string) bool {
	allowedExtensions := []string{".mp3", ".wav", ".flac"}
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

func uploadToSupabaseStorage(url, bucket, fileName, filePath, anonKey string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/storage/v1/object/%s/%s", url, bucket, fileName), file)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+anonKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error al subir el archivo: %s", string(body))
	}

	return nil
}

func insertTrackToDatabase(url, serviceKey, title, artist, fileURL string) error {
	data := map[string]interface{}{
		"title":       title,
		"artist":      artist,
		"file_url":    fileURL,
		"uploaded_at": nil,
	}

	body, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/v1/tracks", url), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error al insertar en la base de datos: %s", string(body))
	}

	return nil
}
