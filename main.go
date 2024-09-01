package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Upload media
	http.HandleFunc("/media", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1024 * 1024 * 16)
		multipartFormData := r.MultipartForm
		fmt.Println(multipartFormData.File)
		files := 0
		for _, v := range multipartFormData.File["music"] {
			fmt.Println(v.Filename, ":", v.Size)
			uploadedFile, _ := v.Open()
			_, err := io.ReadAll(uploadedFile)
			uploadedFile.Close()
			if err != nil {
				fmt.Println("err", err)
			} else {
				files += 1
			}
		}
		fmt.Println(files)
		http.Redirect(w, r, fmt.Sprintf("?files=%d", files), http.StatusSeeOther)
		w.Write([]byte("Redirecting - if you aren't redirected, <a href=\"/\"> click here </a>"))
	})

	// Serve website
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", fs)

	fmt.Println("listening on 8085")
	http.ListenAndServe(":8085", nil)
}
