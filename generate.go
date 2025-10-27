package main

import (
    "flag"
    "fmt"
    "html/template"
    "os"
    "path/filepath"
    "time"
)

type Page struct {
    Title    string
    Filename string
    OutFile  string
}

type SiteData struct {
    SiteTitle string
    Therapist string
    Booking   string
    Year      int
    PageTitle string
    Pages     []Page
}

var pages = []Page{
    {Title: "Home", Filename: "index", OutFile: "index.html"},
    {Title: "About", Filename: "about", OutFile: "about.html"},
    {Title: "Services", Filename: "services", OutFile: "services.html"},
    {Title: "Contact", Filename: "contact", OutFile: "contact.html"},
}

func ensureDir(p string) error {
    return os.MkdirAll(p, 0755)
}

func writeCSS(outDir string) error {
    cssDir := filepath.Join(outDir, "css")
    if err := ensureDir(cssDir); err != nil {
        return err
    }
    css := `body{font-family:Inter,system-ui,-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;margin:0;color:#222}
.header{background:#f6f1ea;padding:20px 24px;display:flex;justify-content:space-between;align-items:center}
.container{max-width:900px;margin:32px auto;padding:0 16px}
nav a{margin-right:12px;color:#333;text-decoration:none}
.booking{background:#2b8a6b;color:white;padding:10px 14px;border-radius:6px;text-decoration:none}
footer{color:#777;padding:20px 0;text-align:center;margin-top:40px}`
    return os.WriteFile(filepath.Join(cssDir, "main.css"), []byte(css), 0644)
}

func main() {
    out := flag.String("out", "public", "output directory")
    name := flag.String("name", "Massage Therapist", "therapist name")
    booking := flag.String("booking", "https://booking.example.com", "booking URL")
    title := flag.String("title", "Massage Therapy", "site title")
    flag.Parse()

    if err := ensureDir(*out); err != nil {
        fmt.Println("mkdir error:", err)
        os.Exit(1)
    }

    site := SiteData{
        SiteTitle: *title,
        Therapist: *name,
        Booking:   *booking,
        Year:      time.Now().Year(),
        Pages:     pages,
    }

    basePath := filepath.Join("templates", "base.tmpl")
    for _, p := range pages {
        tplPath := filepath.Join("templates", p.Filename+".tmpl")
        tpl, err := template.ParseFiles(basePath, tplPath)
        if err != nil {
            fmt.Println("template parse error:", err)
            os.Exit(1)
        }

        outPath := filepath.Join(*out, p.OutFile)
        f, err := os.Create(outPath)
        if err != nil {
            fmt.Println("create file error:", err)
            os.Exit(1)
        }

        site.PageTitle = p.Title
        if err := tpl.ExecuteTemplate(f, "base", site); err != nil {
            fmt.Println("execute template error:", err)
            f.Close()
            os.Exit(1)
        }
        f.Close()
        fmt.Println("wrote", outPath)
    }

    if err := writeCSS(*out); err != nil {
        fmt.Println("write css error:", err)
        os.Exit(1)
    }
    fmt.Println("wrote CSS")
    fmt.Println("done")
}
