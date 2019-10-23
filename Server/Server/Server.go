package Server

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"
    //. "../FireEye/Endpoints"
	// "./PowerShell/Utils"
	//"./Utils"
)

func StartServer() {

    // NewServeMux you can have multiples of these passing them to listen and serve
    h := http.NewServeMux()

    // Get current dir Handler
    fs := http.FileServer(http.Dir(getCurrentDir()))

    // Set Handles 
    h.Handle("/files/", http.StripPrefix("/files/", fs), )
    h.Handle("/success", logSuccessHandler("C:\\Temp\\ThisIsAtest.txt"))
    h.Handle("/FE/hello", helloHandler())
    
    // Logger returns handler that is passed to listen and serve
    hl := logger(h)

    // Start server to listen on Logger Handle
    log.Println("\nListening.......")
    log.Fatal(http.ListenAndServe("192.161.1.1:7777", hl))
}

func helloHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %q", "Benjamin")
    })
}

func logSuccessHandler(pathToFile string) http.Handler {
	/*
	Create Request Body in PowerShell
	$body = @{
		Hostname = $(hostname)
	}
	$contentType = "application/x-www-form-urlencoded"
	# Make call to server Success route
	Invoke-WebRequest -Uri "http://192.168.1.1:7777/success" -Method Post -Body $body -ContentType $contentType -UseBasicParsing | out-null
	*/
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        hostname := fmt.Sprintf("%s", r.Form["Hostname"][0])
        logSuccess(GetFile(pathToFile), hostname)
    })
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func logSuccess(f *os.File, host string)  {
    formatedHost := fmt.Sprintf("\r\n%v", host)
    f.WriteString(formatedHost)
    defer f.Close()
}

func GetFile(pathToFile string) *os.File  {
    
    if _, err := os.Stat(pathToFile); err != nil {
        os.Mkdir(pathToFile, os.FileMode(0777))
    }

    if _, err := os.Stat(pathToFile); err != nil{
        file, err := os.Create(pathToFile)
        check(err)
        return file
    }

   file, err := os.OpenFile(pathToFile, os.O_APPEND, 0777)
   check(err)
   return file
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ignorefavicon := fmt.Sprintf("%s", r.URL)
        if(ignorefavicon == "/favicon.ico"){
            return
        }
		log.Printf("Method: %s %s requested %s",r.Method, r.RemoteAddr, r.URL)
		defer h.ServeHTTP(w, r)
	})
}

func getCurrentDir() string  {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    check(err)
    return fmt.Sprintf("%v\\", dir)
}