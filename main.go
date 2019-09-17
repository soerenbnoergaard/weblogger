package main

import (
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

var (
    activeToken = ""
    activePath = ""
)

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func getExecutableDirectory() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        log.Fatal(err)
        return "."
    }
    return dir
}

func getFileNameFromTimestamp(timestamp string) (string, error) {
    basedir := activePath
    filename := filepath.Join(basedir, fmt.Sprintf("%s-%s.csv", activeToken, timestamp))
    if fileExists(filename) {
        return filename, nil
    } else {
        return filename, errors.New("File not found")
    }
}

func getActiveFileName() (string, error) {
    timestamp := time.Now().Format("2006-01-02")
    return getFileNameFromTimestamp(timestamp)
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
    // Determine the file to read from based on the timestamp setting.
    // If no "date" parameter is given, the current date is used.

    keys := r.URL.Query()
    timestamp, ok := keys["date"]

    var filename string
    var err error
    if ok {
        filename, err = getFileNameFromTimestamp(timestamp[0])
    } else {
        filename, err = getActiveFileName()
    }

    if err != nil {
        fmt.Fprintf(w, "No data found on this date")
        return
    }

    buf, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(fmt.Sprintf("Could not read file: %s", filename))
        return
    }
    content := string(buf)
    fmt.Fprintf(w, "%s", content)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
    keys := r.URL.Query()

    token, ok := keys["token"]
    if !ok {
        return
    }

    if token[0] != activeToken {
        log.Printf("Incorrect token: %s\n", token[0])
        return
    }

    data, ok := keys["data"]
    if !ok {
        return
    }

    filename, _ := getActiveFileName()
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

    if err != nil {
        log.Fatal(fmt.Sprintf("Could not open file: %s\n", filename))
        return
    }

    defer f.Close()

    w.Header().Set("Access-Control-Allow-Origin", "*")
    _, err = f.WriteString(fmt.Sprintf("%s\n", data[0]))

    if err != nil {
        log.Fatal(fmt.Sprintf("Could not append to file: %s\n", filename))
        return
    }

    return
}

func handler(w http.ResponseWriter, r *http.Request) {
    switch (r.Method) {
    case "GET":
        handleGetRequest(w, r)
    case "POST":
        handlePostRequest(w, r)
    default:
        break
    }
}

func main() {
    // Parse arguments
    if len(os.Args) < 3 {
        fmt.Println("Usage: ./weblogger PORT TOKEN [PATH]")
        return
    }

    // Port number
    port, err := strconv.Atoi(os.Args[1])
    if err != nil {
        log.Fatal(fmt.Sprintf("Could not parse port number: %s\n", os.Args[1]))
        return
    }

    // Token
    activeToken = os.Args[2]

    // Data file path
    if len(os.Args) < 4 {
        activePath = getExecutableDirectory()
    } else {
        activePath = os.Args[3]
    }

    // Start server
    fmt.Println("Starting server...")
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
