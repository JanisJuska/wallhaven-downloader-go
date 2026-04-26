package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/spf13/pflag"
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
	searchFlags, pageSlice, remainder, directory := getFlags()

	lastPage := len(pageSlice) - 1
	lastPageCount := 24 - remainder

	for pageCount, pageString := range pageSlice {
		searchString := fmt.Sprintf("%v&%v", strings.Join(searchFlags, "&"), pageString)
		fmt.Println(searchString)
		data, err := getRequestData(searchString)
		if err != nil {
			log.Fatal(err)
		}

		if pageCount < lastPage {
			for _, x := range data.Data {
				filename := path.Base(x.Path)
				dst := fmt.Sprintf("%s%s", directory, filename)
				err := getter.GetFile(dst, x.Path)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			for i, x := range data.Data {
				if i == lastPageCount {
					break
				}
				filename := path.Base(x.Path)
				dst := fmt.Sprintf("%s%s", directory, filename)
				err := getter.GetFile(dst, x.Path)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func getFlagys() ([]string, []string, int, string) {
	var searchFlags []string
	query := pflag.StringP("query", "q", "", "Query to search for")
	count := pflag.IntP("count", "n", 24, "Total wallpapers to download; fetches multiple pages as needed")
	category := pflag.IntP("category", "c", 111, "100/010/001 or combined (general/anime/people) [default: 111]")
	purity := pflag.IntP("purity", "p", 110, "100/110/111 (sfw/sketchy/nsfw) [default: 110]")
	sort := pflag.StringP("sort", "s", "date_added", "date_added, relevance, random, views, favorites, toplist [default: date_added]")
	order := pflag.StringP("order", "o", "desc", "desc, asc [default: desc]")
	colors := pflag.StringP("colors", "C", "", "Dominant color filter (hex without #, e.g. 660000)")
	resAtleast := pflag.StringP("resolution", "r", "", "Minimum allowed resolution. Best used together with Ratios (e.g. 1920x1080)")
	ratio := pflag.StringP("ratios", "R", "", "Aspect ratio filter, comma-separated (e.g. 16x9,16x10)")
	directory := pflag.StringP("directory", "d", "./wallpapers/", "Output directory [default: ./wallpapers/]")

	pflag.Parse()

	queryString := "q=" + strings.Join(strings.Fields(*query), "+")
	catString := fmt.Sprintf("categories=%d", *category)
	pureString := fmt.Sprintf("purity=%d", *purity)
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
	url := fmt.Sprintf("https://wallhaven.cc/api/v1/search?%s", searchString)
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

	return data, nil
}
