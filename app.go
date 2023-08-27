package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	ServeStatic()
	StartServer("8000")

}

// Listen and serve port
func StartServer(addr string) {
	log.Printf("Starting server listening on: http://localhost:%v", addr)
	log.Fatal(http.ListenAndServe(":"+addr, nil))
}

// Gets the relative path of the requested path, making `website` the base path
func RelPath(path string) string {
	return "website" + path
}

// Provides handler function to read and respond with static files
func ServeStatic() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		url := strings.Replace(path, "/https:/", "https://", -1)
		if path != "/favicon.ico" {
			go func() {
				err := downloadVid(url)
				if err != nil {
					fmt.Println(err)
				}
			}()
			fmt.Fprint(w, "OK")
		}
	}
	http.HandleFunc("/", h1)
}

func downloadVid(url string) error {
	app := "yt-dlp"
	arg1 := url
	arg0 := "-P"
	savePath := os.Getenv("SAVE_PATH")
	fmt.Println(app, arg0, savePath, arg1)
	cmd := exec.Command(app, arg0, savePath, arg1)
	exitCode := cmd.ProcessState.ExitCode()
	fmt.Println("Exit code:", exitCode)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(stdout))
	return nil
}
