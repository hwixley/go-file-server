package main

import (
    "net/http"
    "log"
	"os"
	"strconv"
	"io"
)

func main() {
    // set the handler function for the default route
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // open the file
        file, err := os.Open("public" + r.URL.Path)
        if err != nil {
            // if the file doesn't exist, return a 404 response
            w.WriteHeader(http.StatusNotFound)
            return
        }
        defer file.Close()

        // get the file's stats
        stats, err := file.Stat()
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        // if the file is a directory, return a 404 response
        if stats.IsDir() {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        // set the headers for the response
        w.Header().Set("Content-Disposition", "attachment; filename="+stats.Name())
        w.Header().Set("Content-Type", "application/octet-stream")
        w.Header().Set("Content-Length", strconv.FormatInt(stats.Size(), 10))

        // copy the file to the response
        _, err = io.Copy(w, file)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    })

    // start the server
    log.Fatal(http.ListenAndServe(":8080", nil))
}
