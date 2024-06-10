package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const uploadDir = "./uploads/"

func main() {
	http.HandleFunc("/store/add", handleAdd)
	http.HandleFunc("/store/ls", handleList)
	http.HandleFunc("/store/rm", handleRemove)
	http.HandleFunc("/store/update", handleUpdate)
	http.HandleFunc("/store/wc", handleWordCount)
	http.HandleFunc("/store/freq-words", handleFreqWords)

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // Max size 10MB

	files := r.MultipartForm.File["file"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//defer file.Close()

		// Check if file already exists
		fileName := fileHeader.Filename
		_, err = os.Stat(uploadDir + fileName)
		if !os.IsNotExist(err) {
			w.WriteHeader(http.StatusConflict)
			io.WriteString(w, "File already exists: "+fileName+"\n")
			return
		}

		// Create a new file
		f, err := os.Create(uploadDir + fileName)
		if err != nil {
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//defer f.Close()

		// Copy file contents to the new file
		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f.Close()

		io.WriteString(w, "File uploaded successfully: "+fileName+"\n")
		file.Close()
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	files, err := filepath.Glob(uploadDir + "*")
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		io.WriteString(w, fileInfo.Name()+"\n")
	}
}

func handleRemove(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Missing 'file' query parameter\n")
		return
	}

	err := os.Remove(uploadDir + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "File removed successfully: "+fileName+"\n")
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Missing 'file' query parameter\n")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	f, err := os.Create(uploadDir + fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	io.WriteString(w, "File updated successfully\n")
}

func handleWordCount(w http.ResponseWriter, r *http.Request) {
	files, err := filepath.Glob(uploadDir + "*")
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wordCount := 0
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			wordCount++
		}
	}

	fmt.Fprintf(w, "Word count: %d\n", wordCount)
}

func handleFreqWords(w http.ResponseWriter, r *http.Request) {
	limit := 10
	order := "asc"

	query := r.URL.Query()
	if limitStr := query.Get("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}
	if orderStr := query.Get("order"); orderStr != "" {
		order = orderStr
	}

	files, err := filepath.Glob(uploadDir + "*")
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wordFreq := make(map[string]int)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			word := scanner.Text()
			wordFreq[word]++
		}
	}

	sortedWords := make([]string, 0, len(wordFreq))
	for word := range wordFreq {
		sortedWords = append(sortedWords, word)
	}

	if order == "asc" {
		sort.Strings(sortedWords)
	} else if order == "dsc" {
		sort.Sort(sort.Reverse(sort.StringSlice(sortedWords)))
	}

	freqWords := make([]string, 0, limit)
	for i, word := range sortedWords {
		if i >= limit {
			break
		}
		freqWords = append(freqWords, fmt.Sprintf("%s: %d", word, wordFreq[word]))
	}

	fmt.Fprintf(w, "Frequent words:\n%s\n", strings.Join(freqWords, "\n"))
}
