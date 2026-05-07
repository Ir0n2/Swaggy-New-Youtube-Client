package main

import (
	"bytes"
	"strings"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
        "path/filepath"
	"io"
	"archive/zip"
)

func main() {

	http.HandleFunc("/", login)
	http.HandleFunc("/page", reload)
	http.HandleFunc("/temp.zip", downloadZip)
	
	//serve dir temp on /temp
	f := http.FileServer(http.Dir("./temp"))
	http.Handle("/temp/", http.StripPrefix("/temp/", f))

	//run server
	log.Println("Serving on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func title(url string) string {

	cmd := exec.Command("yt-dlp", "--get-title", url)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	title := strings.TrimSpace(out.String())
	return title
}

func login(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "POST" {

		fmt.Println("run")
               	//erase previous downloads before downloading another
		err := os.RemoveAll("temp")
		if err != nil {
			fmt.Println(err)
		}
		//get form values
		url := r.FormValue("url")
	        mode := r.FormValue("mode")
	        format := r.FormValue("format")
		
		//make dir
		tempDir := "temp"
       		os.MkdirAll(tempDir, os.ModePerm)	
	        //make filepath to put file in that temp dir
		outputTemplate := filepath.Join(tempDir, title(url)+".%(ext)s")
		//if mode download playlist from youtube to temp dir, else download 1 video
		if mode == "playlist" {
			downloadVideo(false, true, format, url, outputTemplate)
		} else {
			downloadVideo(true, false, format, url, outputTemplate)

		}

		http.ServeFile(w, r, "static/page.html")
		
		fmt.Println("fin")
	} 

	if r.Method == "GET" {
		http.ServeFile(w, r, "static/main.html")		

	}

}

func downloadZip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=temp.zip")

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	files, _ := os.ReadDir("./temp")

	for _, f := range files {
		filePath := "./temp/" + f.Name()
		file, _ := os.Open(filePath)
		defer file.Close()

		fw, _ := zipWriter.Create(f.Name())
		io.Copy(fw, file)
	}
}
//need this handler for page.html which auto redirects back to home page after the post request
func reload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static/page.html")
	}
}

func downloadVideo(audio, downloadPlaylist bool, fileType, url, title string) {

        if downloadPlaylist == false {
                cmd := exec.Command("yt-dlp", "-o", title,"-t", fileType, url)
                err := cmd.Run()
                if err != nil {log.Fatal(err)}
        }

        if downloadPlaylist == true {
                cmd := exec.Command("yt-dlp", "-t", fileType, "-o", "temp/%(playlist_index)s - %(title)s.%(ext)s", url)
                err := cmd.Run()
                if err != nil {log.Fatal(err)}

        }

}
