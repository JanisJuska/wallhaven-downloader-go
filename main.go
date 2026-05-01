package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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
	a := app.New()
	w := a.NewWindow("Wallhaven Downloader")

	// --- Inputs ---
	queryEntry := widget.NewEntry()
	queryEntry.SetPlaceHolder("Search...")

	resolutionEntry := widget.NewEntry()
	resolutionEntry.SetPlaceHolder("resolution: 1920x1080")

	ratioEntry := widget.NewEntry()
	ratioEntry.SetPlaceHolder("aspect ratio: 16x9")

	colorEntry := widget.NewEntry()
	colorEntry.SetPlaceHolder("color hex (without #)")

	// --- Category ---
	category := "111"
	catGeneral := widget.NewCheck("General", func(b bool) { category = toggleBit(category, 0, b) })
	catAnime := widget.NewCheck("Anime", func(b bool) { category = toggleBit(category, 1, b) })
	catPeople := widget.NewCheck("People", func(b bool) { category = toggleBit(category, 2, b) })

	// --- Purity ---
	purity := "100"
	purSFW := widget.NewCheck("SFW", func(b bool) { purity = toggleBit(purity, 0, b) })
	purSketchy := widget.NewCheck("Sketchy", func(b bool) { purity = toggleBit(purity, 1, b) })
	purNSFW := widget.NewCheck("NSFW", func(b bool) { purity = toggleBit(purity, 2, b) })

	// --- Directory ---
	dirPath := "./wallpapers/"
	dirLabel := widget.NewLabel(dirPath)

	selectDirBtn := widget.NewButton("Select Directory", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				dirPath = uri.Path() + string(os.PathSeparator)
				dirLabel.SetText(dirPath)
			}
		}, w)
	})

	// --- Terminal Output ---
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Logs...")
	output.Disable()

	logLine := func(s string) {
		output.SetText(output.Text + s + "\n")
	}

	// --- Count entry ---
	countEntry := widget.NewEntry()
	countEntry.SetPlaceHolder("Count (e.g. 50)")
	// countEntry.SetText("") // default

	// --- Sorting ---
	sortSelect := widget.NewSelect(
		[]string{"date_added", "relevance", "random", "views", "favorites", "toplist"},
		nil,
	)
	sortSelect.SetSelected("date_added")

	orderSelect := widget.NewSelect(
		[]string{"desc", "asc"},
		nil,
	)
	orderSelect.SetSelected("desc")

	sortSelect.OnChanged = func(val string) {
		if val == "random" {
			orderSelect.Disable()
		} else {
			orderSelect.Enable()
		}
	}

	// --- Download Button ---
	startBtn := widget.NewButton("Download", func() {
		go func() {
			logLine("Starting download...")

			// --- Parse count ---
			count := 24
			if countEntry.Text != "" {
				_, err := fmt.Sscanf(countEntry.Text, "%d", &count)
				if err != nil || count <= 0 {
					logLine("Invalid count value")
					return
				}
			}

			if count > 500 {
				logLine("Count is too large (max 500 for safety)")
				return
			}

			query := "q=" + strings.Join(strings.Fields(queryEntry.Text), "+")

			searchFlags := []string{
				query,
				fmt.Sprintf("categories=%s", category),
				fmt.Sprintf("purity=%s", purity),
				fmt.Sprintf("sorting=%s", sortSelect.Selected),
				fmt.Sprintf("order=%s", orderSelect.Selected),
				fmt.Sprintf("atleast=%s", resolutionEntry.Text),
				fmt.Sprintf("ratios=%s", ratioEntry.Text),
				fmt.Sprintf("colors=%s", colorEntry.Text),
			}

			// --- First request to get lastPage ---
			data, err := getRequestData(query)
			if err != nil {
				logLine("Error: " + err.Error())
				return
			}

			pages, remainder := getPages(count, data.Meta.LastPage)

			var wallpapers []Wallpaper

			for i := 1; i <= pages; i++ {
				search := fmt.Sprintf("%s&page=%d", strings.Join(searchFlags, "&"), i)

				resp, err := getRequestData(search)
				if err != nil {
					logLine("Fetch error: " + err.Error())
					return
				}

				limit := len(resp.Data)

				if i == pages {
					limit -= remainder
				}

				wallpapers = append(wallpapers, resp.Data[:limit]...)
			}

			os.MkdirAll(dirPath, os.ModePerm)

			var downloaded int64
			var skipped int64
			var wg sync.WaitGroup

			for _, wp := range wallpapers {
				dst := filepath.Join(dirPath, filepath.Base(wp.Path))

				// skip existing
				if _, err := os.Stat(dst); err == nil {
					atomic.AddInt64(&skipped, 1)
					continue
				}

				wg.Add(1)
				go func(url, dst string) {
					defer wg.Done()

					err := downloadFileSimple(url, dst)
					if err != nil {
						logLine("Failed: " + url)
						return
					}

					atomic.AddInt64(&downloaded, 1)
					logLine("Downloaded: " + filepath.Base(dst))
				}(wp.Path, dst)
			}

			wg.Wait()

			logLine(fmt.Sprintf(
				"Done: %d downloaded, %d skipped (total requested: %d)",
				downloaded, skipped, count,
			))
		}()
	})

	// --- Layout ---
	topBar := container.NewBorder(
		nil,
		nil,
		nil,
		startBtn,
		container.NewGridWithColumns(2,
			queryEntry,
			countEntry,
		),
	)

	left := container.NewVBox(
		widget.NewLabel("Categories"),
		catGeneral, catAnime, catPeople,
	)

	middle := container.NewVBox(
		widget.NewLabel("Purity"),
		purSFW, purSketchy, purNSFW,
	)

	sorting := widget.NewCard("", "", container.NewVBox(
		widget.NewLabel("Sort by"),
		sortSelect,
		orderSelect,
	))

	right := container.NewVBox(
		widget.NewLabel("Others: "),
		resolutionEntry,
		ratioEntry,
		colorEntry,
		selectDirBtn,
		dirLabel,
	)

	controls := container.NewGridWithColumns(4, left, middle, sorting, right)

	content := container.NewBorder(
		topBar,
		output,
		nil,
		nil,
		controls,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

// --- Helpers ---

func toggleBit(s string, index int, value bool) string {
	b := []rune(s)
	if value {
		b[index] = '1'
	} else {
		b[index] = '0'
	}
	return string(b)
}

// simplified downloader (no mpb)
func downloadFileSimple(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
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
		url = fmt.Sprintf("%s&apikey=%s", url, apiKey) // request.Header.Set("X-API-Key", apiKey)
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
	json.Unmarshal(body, &data) // if err != nil { // return Response{}, err // }

	return data, nil
}
