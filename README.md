# 🖥️ wallhaven-downloader-go (Desktop GUI)

A desktop GUI version of the **Wallhaven downloader**, built with Go and Fyne.

This is an experimental branch that brings the original CLI functionality into a simple graphical interface.

---

## 🚧 Status

**Early development / experimental**

Current features:

* 🔍 Query-based search
* 🎯 Filtering (categories, purity, resolution, ratios, colors)
* 🔢 Custom download count
* 📊 Sorting (date_added, views, favorites, toplist, etc.)
* 📁 Directory selection
* ⚡ Concurrent downloads
* 📝 Live log output

Planned:

* 📊 Progress bars (per file & total)
* 🖼 Wallpaper preview grid
* ⏸ Pause / cancel downloads
* 🔁 Retry failed downloads
* 🎛 Better UI/UX

---

## 🖼️ Preview

*(UI is minimal and subject to change)*

---

## 📦 Installation

### 1. Clone the repository

```bash
git clone https://github.com/JanisJuska/wallhaven-downloader-go.git
cd wallhaven-downloader-go
git checkout desktop
```

---

## 🛠 Build Instructions

### 2. Install dependencies

Make sure you have:

* Go (1.20+ recommended)
* Fyne toolkit

Install Fyne:

```bash
go get fyne.io/fyne/v2
```

---

### 3. Build the binary

```bash
go build -o wallhaven-gui
```

This will create an executable file:

```
wallhaven-gui
```

---

## 🚀 Running the App

### Option A: Run directly

```bash
./wallhaven-gui
```

---

### Option B: Install system-wide (recommended)

Move the binary to a directory in your `$PATH`:

```bash
sudo mv wallhaven-gui /usr/local/bin/
```

Now you can launch it from anywhere:

```bash
wallhaven-gui
```

---

## 🖥️ Add to Application Launcher (Linux)

To make the app appear in your desktop environment / launcher:

### 1. Create a `.desktop` file

```bash
mkdir -p ~/.local/share/applications
nano ~/.local/share/applications/wallhaven-gui.desktop
```

---

### 2. Add this content

```ini
[Desktop Entry]
Name=Wallhaven Downloader
Comment=Download wallpapers from Wallhaven
Exec=/usr/local/bin/wallhaven-gui
Icon=utilities-terminal
Terminal=false
Type=Application
Categories=Utility;
```

---

### 3. Make it executable

```bash
chmod +x ~/.local/share/applications/wallhaven-gui.desktop
```

---

### 4. Refresh desktop database (optional)

```bash
update-desktop-database ~/.local/share/applications
```

---

Now you can:

* Launch from app launcher
* Bind it in your window manager (e.g. Hyprland)
* Search it like a normal app

---

## ⚡ Optional: Hyprland keybind

Example keybind:

```ini
bind = SUPER, W, exec, wallhaven-gui
```

---

## 🧊 Optional: App Icon

Right now the app uses a default icon.

To improve:

1. Download an icon (`.png` or `.svg`)
2. Place it somewhere like:

```bash
~/.local/share/icons/wallhaven.png
```

3. Update `.desktop` file:

```ini
Icon=/home/youruser/.local/share/icons/wallhaven.png
```

---

## 📦 Uninstall

```bash
sudo rm /usr/local/bin/wallhaven-gui
rm ~/.local/share/applications/wallhaven-gui.desktop
```

---

## 🧠 Notes

* The app is self-contained (no external runtime required)
* API key is still provided via environment variable:

```bash
export WALLHAVEN_API_KEY=your_key_here
```

You may want to add this to your shell config (`.bashrc`, `.zshrc`, etc.)


---

## 🚀 Usage

### Basic workflow

1. Enter a **search query**
2. Set **count** (number of wallpapers)
3. Adjust filters:

   * Categories (General / Anime / People)
   * Purity (SFW / Sketchy / NSFW)
   * Resolution / ratio / color
4. Choose sorting:

   * date_added, views, favorites, etc.
5. Select download directory
6. Click **Download**

---

## ⚙️ Features

### 🔍 Search

Works the same as CLI:

* Supports multi-word queries
* Automatically formatted for API

---

### 🔢 Count

* Controls exact number of wallpapers downloaded
* Automatically handles pagination (24 per page)
* Built-in safety cap to prevent excessive downloads

---

### 🎯 Filtering

Supports:

* Categories (`111`, `100`, etc.)
* Purity (`100`, `110`, `111`)
* Resolution (`1920x1080`)
* Ratios (`16x9`)
* Colors (hex without `#`)

---

### 📊 Sorting

Available options:

* `date_added`
* `relevance`
* `random`
* `views`
* `favorites`
* `toplist`

Order:

* `desc`
* `asc`

---

### 📁 File Handling

* Files are saved with original names
* Existing files are skipped automatically
* Output directory is user-selectable

---

## 🔑 API Key Support

Same as CLI version.

Required for:

* 🔞 NSFW content
* 👤 User-specific results

### Set API key

#### Bash / Zsh

```bash
export WALLHAVEN_API_KEY=your_key_here
```

#### Fish

```fish
set -x WALLHAVEN_API_KEY your_key_here
```

Get your key:
https://wallhaven.cc/settings/account

API docs:
https://wallhaven.cc/help/api

---

## 🐧 Linux (Hyprland / tiling WM note)

If you're using a tiling window manager (like Hyprland), the app may open tiled by default.

To force floating mode:

```ini
windowrulev2 = float, title:^(Wallhaven Downloader)$
```

Reload config:

```bash
hyprctl reload
```

---

## ⚠️ Limitations

* No progress bars yet (log output only)
* No preview before download
* No cancel/stop button
* Basic UI (not final design)

---

## 🔀 Branches

| Branch    | Description              |
| --------- | ------------------------ |
| `main`    | Stable CLI application   |
| `desktop` | Experimental GUI version |

---

## 🙌 Acknowledgements

Based on the original CLI project:

* https://github.com/JanisJuska/wallhaven-downloader-go
* https://github.com/Moskas/whdl

---

## 📜 License

MIT License

---

## 🚀 Future Direction

The goal is to evolve this into a full-featured desktop app with:

* Visual browsing
* Download management
* Better performance controls
* Polished UI

---

Enjoy downloading wallpapers with a GUI 🎉
