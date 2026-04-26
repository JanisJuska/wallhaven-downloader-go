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
| `--ratios`     | `-R`  | Aspect ratio filter (e.g. 16x9,16x10)                    | none            |
| `--directory`  | `-d`  | Output directory                                         | `./wallpapers/` |
| `--help`       | `-h`  | Print help and usage information                         | —               |

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

## 📄 Pagination

Wallhaven returns **24 results per page**.

* `--count` determines how many wallpapers to download
* The app automatically:

  * Calculates how many pages are needed
  * Fetches them sequentially
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

---

## 🙌 Acknowledgements

This project was heavily inspired by:

* https://github.com/Moskas/whdl

Special thanks for the CLI design and UX ideas.

---

## 🛠 Future Improvements

* ~~`--help` flag with improved CLI output~~ ✅
* [ ] API key support (for NSFW and user-specific content)
* [ ] Download progress feedback (current file, progress, etc.)
* [ ] Concurrent downloads for better performance
* [ ] `--dry-run` mode (preview downloads without saving)
* [ ] Optional GUI application (experimental idea)

---

## 🔑 API Key (Planned)

Wallhaven requires an API key for:

* NSFW content
* User-specific collections

Future versions might support:

### Using Bash / Zsh

```bash
export WALLHAVEN_API_KEY=your_key_here
```

### Using Fish

```fish
set -x WALLHAVEN_API_KEY your_key_here
```

You can generate a key at:
https://wallhaven.cc/settings/account

API docs:
https://wallhaven.cc/help/api

---

## 📜 License

MIT License

---

## 🚧 Status

**v1.0 — Functional and stable**

More features and improvements are planned.

---

Enjoy downloading wallpapers! 🎉
