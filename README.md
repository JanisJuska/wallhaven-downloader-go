# wallhaven-downloader-go

A command-line wallpaper downloader from [Wallhaven API](https://wallhaven.cc), written in Go.

This tool allows you to search, filter, and download wallpapers directly from the Wallhaven API with a flexible and scriptable CLI interface.

---

## ✨ Features

* 🔍 Query-based wallpaper search
* 🎯 Advanced filtering (categories, purity, resolution, ratios, etc.)
* 📄 Automatic pagination handling (24 results per page)
* 📦 Batch downloading with exact count control
* 🗂 Preserves original filenames and extensions
* ⚡ Lightweight and fast
* 📖 Built-in `--help` flag for CLI usage guidance
* 📊 **Live download progress bars (per-file) with speed & ETA**
* 📈 **Final summary (files downloaded, skipped, total size, average speed)**
* 🔑 **API key support (required for NSFW and user-specific content)**

---

## 📦 Installation

### 1. Clone the repository

```bash
git clone https://github.com/JanisJuska/wallhaven-downloader-go.git
cd wallhaven-downloader-go
```

### 2. Build the binary

```bash
go build -o wallhaven
```

### 3. (Optional) Move to PATH

```bash
mv wallhaven /usr/local/bin/
```

---

## 🚀 Usage

```
Usage: wallhaven [OPTIONS] --query <QUERY>
```

### ⚙️ Options

| Flag           | Short | Description                                              | Default         |
| -------------- | ----- | -------------------------------------------------------- | --------------- |
| `--query`      | `-q`  | Query to search for                                      | none            |
| `--count`      | `-n`  | Total wallpapers to download (auto pagination)           | `24`            |
| `--category`   | `-c`  | 100/010/001 or combined (general/anime/people)           | `111`           |
| `--purity`     | `-p`  | 100/110/111 (sfw/sketchy/nsfw)                           | `110`           |
| `--sort`       | `-s`  | date_added, relevance, random, views, favorites, toplist | `date_added`    |
| `--order`      | `-o`  | desc, asc                                                | `desc`          |
| `--colors`     | `-C`  | Dominant color filter (hex without #)                    | none            |
| `--resolution` | `-r`  | Minimum resolution (e.g. 1920x1080)                      | none            |
| `--ratios`     | `-R`  | Aspect ratio filter, comma-separated                     | none            |
| `--directory`  | `-d`  | Output directory                                         | `./wallpapers/` |
| `--help`       | `-h`  | Print help and usage information                         | —               |

---

## 🔑 API Key Support

Wallhaven requires an API key for:

* 🔞 NSFW content (`purity=111`)
* 👤 User-specific collections

### Set your API key

#### Bash / Zsh

```bash
export WALLHAVEN_API_KEY=your_key_here
```

#### Fish

```fish
set -x WALLHAVEN_API_KEY your_key_here
```

You can generate your API key here:  
https://wallhaven.cc/settings/account

API documentation:  
https://wallhaven.cc/help/api

---

### ⚠️ Important

To download NSFW wallpapers, you must:

1. Set your API key
2. Use `-p 111`

---

### Example (NSFW)

```bash
export WALLHAVEN_API_KEY=your_key_here

wallhaven -q "anime" -p 111 -n 24
```

---

## 📌 Examples

### Show help

```bash
wallhaven --help
```

---

### Download from first page

```bash
wallhaven -q "Cyberpunk 2077"
```

---

### Download 72 wallpapers (auto pagination)

```bash
wallhaven -q "nature" -n 72
```

---

### Advanced filtering

```bash
wallhaven \
  -q "space" \
  -n 48 \
  -c 100 \
  -p 110 \
  -s views \
  -o desc \
  -r 1920x1080 \
  -R 16x9 \
  -d ./downloads/
```

---

## 📊 Download Output

During downloads, each file displays a live progress bar including:

* Percentage completed
* Estimated time remaining (ETA)
* Current download speed

After completion, a summary is shown:

```text
Done: 50/50 files downloaded (3 skipped) (123.47 MB) in 21.2s (5.82 MB/s)
```

---

## 📄 Pagination

Wallhaven returns **24 results per page**.

* `--count` determines how many wallpapers to download
* The app automatically:

  * Calculates how many pages are needed
  * Fetches them
  * Stops exactly at the requested count

---

## 📁 File Naming

* Files are saved using their original names from Wallhaven

  Example:

  ```
  wallhaven-abc123.jpg
  wallhaven-xyz789.png
  ```

* Original file extensions are preserved (`.jpg`, `.png`, etc.)

* Existing files are automatically **skipped**

---

## 🙌 Acknowledgements

This project was heavily inspired by:

* https://github.com/Moskas/whdl

Special thanks for the CLI design and UX ideas.

---

## 🛠 Future Improvements

* ~~`--help` flag with improved CLI output~~ ✅
* ~~Download progress feedback (current file, progress, etc.)~~ ✅
* ~~Concurrent downloads for better performance~~ ✅
* ~~API key support (for NSFW and user-specific content)~~ ✅
* [ ] `--dry-run` mode (preview downloads without saving)
* [ ] Retry logic for failed downloads
* [ ] Resume partial downloads
* [ ] Optional GUI application (experimental idea)

---

## 📜 License

MIT License

---

## 🚧 Status

**v1.2 — Stable with API key support, progress UI, and concurrent downloads**

---

Enjoy downloading wallpapers! 🎉
