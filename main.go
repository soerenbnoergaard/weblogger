package main

import (
    "errors"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"
)


type Settings struct {
    port int;
    token string;
    path string;
    verbose bool;
    noTimestamp bool;
}

var (
    settings Settings
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
    basedir := settings.path
    filename := filepath.Join(basedir, fmt.Sprintf("%s-%s.csv", settings.token, timestamp))
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

    if settings.verbose {
        fmt.Printf("GET:  %s\n", r.URL);
    }

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

    if settings.verbose {
        fmt.Printf("POST: %s\n", r.URL);
    }

    token, ok := keys["token"]
    if !ok {
        return
    }

    if token[0] != settings.token {
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
    var s string
    if settings.noTimestamp {
        s = fmt.Sprintf("%s\n", data[0])
    } else {
        s = fmt.Sprintf("%d,%s\n", time.Now().Unix(), data[0])
    }
    _, err = f.WriteString(s)

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

func parseArgs() error {
    // Optional arguments
    flag.StringVar(&settings.path, "d", getExecutableDirectory(), "Directory where datafiles are stored")
    flag.BoolVar(&settings.verbose, "v", false, "Enable verbose output")
    flag.BoolVar(&settings.noTimestamp, "t", false, "Disable timestamp")
    flag.Parse()

    // Positional arguments
    if flag.NArg() < 2 {
        fmt.Println("Usage: ./weblogger [OPTIONS] [-h] PORT TOKEN")
        return errors.New("Not enough arguments given")
    }

    // Port number
    port, err := strconv.Atoi(flag.Arg(0))
    if err != nil {
        return errors.New(fmt.Sprintf("Port %s is not a number", flag.Arg(0)))
    }
    settings.port = port

    // Token
    settings.token = flag.Arg(1)

    if settings.verbose {
        fmt.Printf("Server settings:\n%+v\n", settings)
    }
    return nil
}

func main() {
    err := parseArgs()
    if err != nil {
        log.Fatal(err)
        return
    }

    // Start server
    if settings.verbose {
        fmt.Println("Starting server...")
    }
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", settings.port), nil))
}
