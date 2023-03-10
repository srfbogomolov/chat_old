package main

import (
    "flag"
    "log"
    "net/http"
    "path/filepath"
    "sync"
    "text/template"
    "os"
    "github.com/srfbogomolov/trace"
)

// templateHandler represents a single template
type templateHandler struct {
    once     sync.Once
    filename string
    templ    *template.Template
}

// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
        t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
    })
    t.templ.Execute(w, r)
}

func main() {
    var addr = flag.String("addr", ":8080", "The addr of the application.")
    flag.Parse()
    r := newRoom()
    r.tracer = trace.New(os.Stdout)
    http.Handle("/", &templateHandler{filename: "chat.html"})
    http.Handle("/room", r)

    //get the room going
    go r.run()

    // start the web server
    log.Println("Starting web server on", *addr)
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
