package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"github.com/gocolly/colly"
	"strings"
	// "log"
)

type data struct {
	URL    string
	Images []string
	Links  []string
}

func crawl(crawl_url string) ([]string, []string) {
	proxyURL, _ := url.Parse("http://IIT2020060:Satwik..060@172.31.2.4:8080/")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
	}

	c := colly.NewCollector()
	max_depth := 2

	c.WithTransport(client.Transport)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	var imgs []string
	c.OnHTML("img", func(h *colly.HTMLElement) {
		source := h.Attr("src")
		imgs = append(imgs, source)
	})

	var links []string
	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		source := h.Attr("href")
		if h.Request.Depth >= max_depth{
			return
		}
		if strings.HasPrefix(source, "http"){
			fmt.Printf("Depth: %d, Link: %s\n", h.Request.Depth, source)
			links = append(links, source)
		}
		h.Request.Visit(source)
	})

	err := c.Visit(crawl_url)
	_ = err
	return imgs, links
}
func crawl_images(crawl_url string)([]string){
	proxyURL, _ := url.Parse("http://IIT2020060:Satwik..060@172.31.2.4:8080/")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
	}
	c := colly.NewCollector()

	c.WithTransport(client.Transport)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	var imgs []string
	c.OnHTML("img", func(h *colly.HTMLElement) {
		source := h.Attr("src")
		imgs = append(imgs, source)
	})


	err := c.Visit(crawl_url)
	_ = err
	fmt.Println(len(imgs))
	return imgs
}

func crawl_links(crawl_url string)([]string){
	proxyURL, _ := url.Parse("http://IIT2020060:Satwik..060@172.31.2.4:8080/")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
	}
	c := colly.NewCollector()
	// max_depth := 1

	c.WithTransport(client.Transport)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	var links []string
	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		source := h.Attr("href")
		// if h.Request.Depth >= max_depth{
		// 	return
		// }
		if strings.HasPrefix(source, "http"){
			// fmt.Printf("Depth: %d, Link: %s\n", h.Request.Depth, source)
			links = append(links, source)
		}
		// h.Request.Visit(source)
	})

	err := c.Visit(crawl_url)
	_ = err
	fmt.Println(len(links))
	return links
}
func index(w http.ResponseWriter, r *http.Request) {
	layout := "layout.html"
	tmpl, err := template.ParseFiles(layout, "search.html")
	_ = err


	if r.Method == http.MethodPost{
		crawl_url := r.FormValue("url")
		http.Redirect(w, r, "/data?url=" + crawl_url, http.StatusSeeOther)
	}
	tmpl.ExecuteTemplate(w,layout,nil)
}

func data_func(w http.ResponseWriter, r *http.Request){
	layout := "layout.html"
	tmpl , err := template.ParseFiles(layout, "data.html")
	_= err

	if r.Method == http.MethodGet{
		query := r.URL.Query()
		crawl_url := query.Get("url")
		tmpl.ExecuteTemplate(w, layout, crawl_url)
	}
}

func links_func(w http.ResponseWriter, r *http.Request){
	layout := "layout.html"
	tmpl , err := template.ParseFiles(layout, "links.html")
	_= err

	if r.Method == http.MethodGet{
		query := r.URL.Query()
		crawl_url := query.Get("url")

		links := crawl_links(crawl_url)
		// fmt.Println(len(links))
		tmpl.ExecuteTemplate(w, layout, links)
	}
}
func images_func(w http.ResponseWriter, r *http.Request){
	layout := "layout.html"
	tmpl , err := template.ParseFiles(layout, "images.html")
	_= err

	if r.Method == http.MethodGet{
		query := r.URL.Query()
		crawl_url := query.Get("url")

		imgs := crawl_images(crawl_url)
		fmt.Println(len(imgs))
		tmpl.ExecuteTemplate(w, layout, imgs)
	}
}

func main() {
	// Set up a proxy URL
	http.HandleFunc("/", index)
	http.HandleFunc("/data", data_func)
	http.HandleFunc("/images", images_func)
	http.HandleFunc("/links", links_func)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
