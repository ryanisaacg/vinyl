package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	mediaRoot := os.Args[1]
	refreshToken := os.Args[2]
	refreshUrl := fmt.Sprintf("http://127.0.0.1:32400/library/sections/29/refresh?X-Plex-Token=%s", refreshToken)

	// Upload media
	http.HandleFunc("/media", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1024 * 1024 * 16)
		multipartFormData := r.MultipartForm
		files := 0
		existing := 0
		failed := 0
		for _, v := range multipartFormData.File["music"] {
			filePath := fmt.Sprintf("%s/%s", mediaRoot, v.Filename)
			if _, err := os.Stat(filePath); err == nil {
				existing += 1
				continue
			}
			uploadedFile, err := v.Open()
			if err != nil {
				failed += 1
				fmt.Println(err)
				continue
			}
			fileContents, err := io.ReadAll(uploadedFile)
			if err != nil {
				failed += 1
				fmt.Println(err)
				continue
			}
			err = uploadedFile.Close()
			if err != nil {
				failed += 1
				fmt.Println(err)
				continue
			}
			err = os.WriteFile(fmt.Sprintf("%s/%s", mediaRoot, v.Filename), fileContents, 0644)
			if err != nil {
				failed += 1
				fmt.Println(err)
				continue
			}
			files += 1
		}
		http.Redirect(w, r, fmt.Sprintf("?files=%d&failed=%d&existing=%d", files, failed, existing), http.StatusSeeOther)
		w.Write([]byte("Redirecting - if you aren't redirected, <a href=\"/\"> click here </a>"))
	})
	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		http.Get(refreshUrl)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		w.Write([]byte("Redirecting - if you aren't redirected, <a href=\"/\"> click here </a>"))
	})

	// Serve website
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)

	fmt.Println("listening on 8085")
	http.ListenAndServe(":8085", nil)
}
