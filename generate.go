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
	jsDir := filepath.Join(outDir, "js")
	if err := ensureDir(cssDir); err != nil {
		return err
	}
	if err := ensureDir(jsDir); err != nil {
		return err
	}

	css := `:root{
  --max-width:1100px;
  --accent:#2b8a6b;
  --muted:#666666;
  --bg:#fffaf7;
  --page-pad:16px;
  --radius:10px;
}
*{box-sizing:border-box}
html,body{height:100%}
body{
  margin:0;
  font-family:Inter,system-ui,-apple-system,Segoe UI,Roboto,Helvetica,Arial,sans-serif;
  color:#222;
  background:#fff;
  line-height:1.45;
  padding:0 var(--page-pad);
}
.header{
  background:var(--bg);
  border-bottom:1px solid rgba(0,0,0,0.04);
}
.header-inner{
  max-width:var(--max-width);
  margin:0 auto;
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap:16px;
  padding:18px var(--page-pad);
}
.brand{display:flex;flex-direction:column}
.brand strong{font-size:1.05rem}
.brand .tag{font-size:0.9rem;color:var(--muted)}
.nav{display:flex;gap:10px;align-items:center;flex-wrap:wrap}
.nav a{color:#222;text-decoration:none;padding:8px 6px;border-radius:8px}
.nav a:hover{background:rgba(0,0,0,0.03)}
.booking{
  background:var(--accent);
  color:#fff;
  padding:10px 14px;
  border-radius:8px;
  text-decoration:none;
  white-space:nowrap;
  box-shadow:0 2px 6px rgba(43,138,107,0.12);
}
.container{
  max-width:var(--max-width);
  margin:28px auto;
  padding:0 var(--page-pad);
}
h1{margin-top:0}
section{margin-top:18px}
ul{padding-left:1.15rem}
footer{
  color:var(--muted);
  text-align:center;
  padding:28px 0;
  border-top:1px solid rgba(0,0,0,0.04);
}

/* responsive / mobile nav */
.nav-toggle{
  display:none;
  background:transparent;
  border:0;
  padding:8px;
  font-size:1.2rem;
  cursor:pointer;
}
.nav-mobile{
  display:flex;
  flex-direction:column;
  gap:8px;
}
.nav-mobile a{display:block;padding:10px;border-radius:8px}
@media (max-width:800px){
  .header-inner{flex-direction:row;align-items:center;padding:14px var(--page-pad)}
  .nav{display:none}
  .nav-toggle{display:inline-flex}
  .nav.open{display:flex;position:absolute;top:72px;left:0;right:0;background:white;padding:12px 16px;border-bottom:1px solid rgba(0,0,0,0.04);box-shadow:0 6px 18px rgba(0,0,0,0.06);z-index:40;flex-direction:column}
  .booking{width:100%;text-align:center;padding:12px}
  .brand{width:auto}
  .container{margin:18px auto;padding:0 var(--page-pad)}
}
@media (max-width:420px){
  :root{--page-pad:12px}
}`
	if err := os.WriteFile(filepath.Join(cssDir, "main.css"), []byte(css), 0644); err != nil {
		return err
	}

	js := `(function(){
  const btn = document.querySelector('.nav-toggle');
  const nav = document.querySelector('nav.nav');
  if(!btn || !nav) return;
  btn.addEventListener('click', function(){
    nav.classList.toggle('open');
    const expanded = nav.classList.contains('open');
    btn.setAttribute('aria-expanded', expanded ? 'true' : 'false');
  });
  // close when clicking outside
  document.addEventListener('click', function(e){
    if(!nav.classList.contains('open')) return;
    if(nav.contains(e.target) || btn.contains(e.target)) return;
    nav.classList.remove('open');
    btn.setAttribute('aria-expanded','false');
  });
})();`
	if err := os.WriteFile(filepath.Join(jsDir, "main.js"), []byte(js), 0644); err != nil {
		return err
	}
	return nil
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
