package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/pflag"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Response struct {
	Data []Wallpaper `json:"data"`
	Meta Meta        `json:"meta"`
}

type Wallpaper struct {
	ID         string   `json:"id"`
	URL        string   `json:"url"`
	ShortURL   string   `json:"short_url"`
	Views      int      `json:"views"`
	Favorites  int      `json:"favorites"`
	Source     string   `json:"source"`
	Purity     string   `json:"purity"`
	Category   string   `json:"category"`
	DimensionX int      `json:"dimension_x"`
	DimensionY int      `json:"dimension_y"`
	Resolution string   `json:"resolution"`
	Ratio      string   `json:"ratio"`
	FileSize   int      `json:"file_size"`
	FileType   string   `json:"file_type"`
	CreatedAt  string   `json:"created_at"`
	Colors     []string `json:"colors"`
	Path       string   `json:"path"`
	Thumbs     Thumbs   `json:"thumbs"`
}

type Thumbs struct {
	Large    string `json:"large"`
	Original string `json:"original"`
	Small    string `json:"small"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

func main() {
	var downloaded int64
	var totalBytes int64
	var skipped int64

	var wg sync.WaitGroup
	p := mpb.New(
		mpb.WithWidth(60),
	)

	sem := make(chan struct{}, 5)

	searchFlags, pageSlice, remainder, directory := getFlags()

	os.MkdirAll(directory, os.ModePerm)

	lastPage := len(pageSlice) - 1

	var allWallpapers []Wallpaper

	for pageCount, pageString := range pageSlice {
		searchString := fmt.Sprintf("%v&%v", strings.Join(searchFlags, "&"), pageString)

		data, err := getRequestData(searchString)
		if err != nil {
			log.Fatal(err)
		}

		limit := len(data.Data)

		if pageCount == lastPage {
			limit = len(data.Data) - remainder
		}

		allWallpapers = append(allWallpapers, data.Data[:limit]...)
	}

	totalBar := p.New(
		int64(len(allWallpapers)),
		mpb.BarStyle(),
		mpb.PrependDecorators(
			decor.Name("Downloading "),
			decor.CountersNoUnit("%d / %d"),
		),
	)

	start := time.Now()

	for _, x := range allWallpapers {
		filename := path.Base(x.Path)
		dst := fmt.Sprintf("%s%s", directory, filename)

		// skip existing files
		if _, err := os.Stat(dst); err == nil {
			atomic.AddInt64(&skipped, 1)
			totalBar.Increment()
			continue
		}

		wg.Add(1)

		go func(url, dst string) {
			defer wg.Done()
			defer totalBar.Increment() // always increments

			sem <- struct{}{}

			err := downloadFile(url, dst, p)

			<-sem

			if err != nil {
				log.Printf("failed: %s\n", url)
				return
			}

			atomic.AddInt64(&downloaded, 1)

			// get file size after download
			info, err := os.Stat(dst)
			if err == nil {
				atomic.AddInt64(&totalBytes, info.Size())
			}
		}(x.Path, dst)
	}

	wg.Wait()
	p.Wait()

	elapsed := time.Since(start).Seconds()

	mb := float64(totalBytes) / (1024 * 1024)
	speed := mb / elapsed

	fmt.Println()

	totalFiles := len(allWallpapers)

	fmt.Printf(
		"Done: %d/%d files (%d skipped) (%.2f MB) in %.1fs (%.2f MB/s)\n",
		downloaded,
		totalFiles,
		skipped,
		mb,
		elapsed,
		speed,
	)
}

func getFlags() ([]string, []string, int, string) {
	var searchFlags []string
	query := pflag.StringP("query", "q", "", "Query to search for")
	count := pflag.IntP("count", "n", 24, "Total wallpapers to download; fetches multiple pages as needed")
	category := pflag.StringP("category", "c", "111", "100/010/001 or combined (general/anime/people) ")
	purity := pflag.StringP("purity", "p", "100", "100/110/111 (sfw/sketchy/nsfw) ")
	sort := pflag.StringP("sort", "s", "date_added", "date_added, relevance, random, views, favorites, toplist ")
	order := pflag.StringP("order", "o", "desc", "desc, asc ")
	colors := pflag.StringP("colors", "C", "", "Dominant color filter (hex without #, e.g. 660000)")
	resAtleast := pflag.StringP("resolution", "r", "", "Minimum allowed resolution. Best used together with Ratios (e.g. 1920x1080)")
	ratio := pflag.StringP("ratios", "R", "", "Aspect ratio filter, comma-separated (e.g. 16x9,16x10)")
	directory := pflag.StringP("directory", "d", "./wallpapers/", "Output directory ")
	help := pflag.BoolP("help", "h", false, "Print help")

	pflag.Usage = func() {
		fmt.Println("Usage: wallhaven [OPTIONS] --query <QUERY>")
		fmt.Println("\nOptions:")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if *help {
		pflag.Usage()
		os.Exit(0)
	}

	queryString := "q=" + strings.Join(strings.Fields(*query), "+")
	catString := fmt.Sprintf("categories=%s", *category)
	pureString := fmt.Sprintf("purity=%s", *purity)
	sortString := fmt.Sprintf("sorting=%s", *sort)
	ordString := fmt.Sprintf("order=%s", *order)
	colString := fmt.Sprintf("colors=%s", *colors)
	resString := fmt.Sprintf("atleast=%s", *resAtleast)
	ratioString := fmt.Sprintf("ratios=%s", *ratio)

	data, err := getRequestData(queryString)
	if err != nil {
		log.Fatal(err)
	}

	pages, remainder := getPages(*count, data.Meta.LastPage)

	var pageSlice []string
	var pagesString string

	for i := 1; i <= pages; i++ {
		pagesString = fmt.Sprintf("page=%d", i)
		pageSlice = append(pageSlice, pagesString)
	}

	searchFlags = append(searchFlags, queryString, catString, pureString, sortString, ordString, resString, ratioString, colString)

	return searchFlags, pageSlice, remainder, *directory
}

func getPages(count int, lastPage int) (int, int) {
	const PAGE int = 24
	pages := 1
	remainder := PAGE - count

	if count > PAGE {
		var i int
		var inRange bool
		for i = 1; i <= lastPage; i++ {
			if PAGE*i > count {
				inRange = true
				break
			}
		}

		if !inRange {
			log.Fatalf("Count is too large. Max --count possible for a single query is %d", PAGE*lastPage)
		}

		pages *= i
		remainder = PAGE*i - count
		return pages, remainder
	}

	return pages, remainder
}

func getRequestData(searchString string) (Response, error) {
	apiKey := os.Getenv("WALLHAVEN_API_KEY")

	url := fmt.Sprintf("https://wallhaven.cc/api/v1/search?%s", searchString)

	if apiKey != "" {
		url = fmt.Sprintf("%s&apikey=%s", url, apiKey)
		// request.Header.Set("X-API-Key", apiKey)
	} else {
		log.Println("No API key provided (WALLHAVEN_API_KEY). NSFW results may be unavailable.")
	}

	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, err
	}

	response, err := client.Do(request)
	if err != nil {
		return Response{}, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}

	var data Response

	json.Unmarshal(body, &data)
	// if err != nil {
	// 	return Response{}, err
	// }

	return data, nil
}

func downloadFile(url, dst string, p *mpb.Progress) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	size := resp.ContentLength
	if size <= 0 {
		size = -1
	}

	filename := path.Base(dst)

	bar := p.New(
		size,
		mpb.BarStyle(),
		mpb.PrependDecorators(
			decor.Name(shortName(filename)+" "),
			decor.Percentage(),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 60),
			decor.AverageSpeed(decor.SizeB1024(0), "% .2f"),
		),
	)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer out.Close()

	reader := bar.ProxyReader(resp.Body)
	defer reader.Close()

	_, err = io.Copy(out, reader)
	return err
}

func shortName(name string) string {
	if len(name) > 30 {
		return name[:27] + "..."
	}
	return name
}
