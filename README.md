<<<<<<< HEAD
# 🛠️ Gomon - A Nodemon Clone Go File Watcher 🚀

Gomon is a Go-based file watcher that automatically reloads your application when file changes are detected. Inspired by nodemon, it helps streamline development by automatically restarting your app.

## Build Stats 
[![Go](https://github.com/Itzhep/gomon/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/Itzhep/gomon/actions/workflows/go.yml)

## Project Stats
![GitHub repo size](https://img.shields.io/github/repo-size/Itzhep/gomon)
![GitHub Release](https://img.shields.io/github/v/release/Itzhep/gomon)
![GitHub Repo stars](https://img.shields.io/github/stars/Itzhep/gomon)
![GitHub License](https://img.shields.io/github/license/Itzhep/gomon)
![GitHub Issues](https://img.shields.io/github/issues/Itzhep/gomon)
![GitHub Forks](https://img.shields.io/github/forks/Itzhep/gomon)

## 📦 Features

- 🔄 **Automatic file watching and reloading**: Detects file changes and restarts your application automatically.
- 🎨 **CLI with color support**: Enhanced visibility with color-coded output.
- 📝 **Simple and clean configuration**: Minimal setup required to get started.
- 🐳 **Docker support**: Run your application inside Docker containers.
- 🌐 **Live reload server**: Supports live reloading for browser-based applications.

## 🏗️ Installation

### Via Go

To install Gomon, use the following command:

```bash
go install github.com/Itzhep/gomon@latest
```

### Manual Build

1. Clone the repository:

    ```bash
    git clone https://github.com/Itzhep/gomon.git
    cd gomon
    ```

2. Build the project:

    ```bash
    go build -o gomon
    ```

3. Install the Project:

    ```bash
    go install
    ```

4. Move to bin (optional):

    ```bash
    move gomon.exe C:\path\to\your\bin
    ```

## 🚀 Usage

1. Start Gomon with the path to your main application file:

    ```bash
    gomon start --app path/to/your/app.go
    ```

2. **Press `rs`** in the CLI to manually restart the application.

## 🛠️ Configuration

Gomon supports a variety of configurations directly from the CLI. You can specify the file to watch and other options like color settings for better CLI appearance.

### CLI Options

- `--app, -a`: Path to the Go application to run (required)
- `--debounce, -d`: Debounce duration for file changes (default: 1s)
- `--docker`: Use Docker for restarting the app
- `--exclude, -e`: Directories to exclude from watching (default: .git, vendor, node_modules)
- `--verbose, -v`: Enable verbose logging

## 📝 Example

Here's a basic example of how to use Gomon:

```bash
gomon start --app path/to/your/app.go
```

This command starts Gomon, watches for file changes, and restarts your application automatically.

## 🐳 Docker Support

Gomon can be run inside a Docker container. Use the provided `Dockerfile` and `docker-compose.yml` for easy setup.

## 🗂️ Contributing

Contributions are welcome! Please submit a pull request or open an issue if you find any bugs or have suggestions.

## 🌟 Star the Project

If you find Gomon useful, please give it a star on [GitHub](https://github.com/Itzhep/gomon) to support development and stay updated with new features.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

🚀 Happy coding with Gomon!
=======
# 🛠️ Gomon - A Nodemon clone Go File Watcher 🚀

Gomon is a Go-based file watcher that automatically reloads your application when file changes are detected. Inspired by nodemon, it helps streamline development by automatically restarting your app.
## Build Stats 
[![Go](https://github.com/Itzhep/gomon/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/Itzhep/gomon/actions/workflows/go.yml)
## Project Stats
![GitHub repo size](https://img.shields.io/github/repo-size/Itzhep/gomon)
![GitHub Release](https://img.shields.io/github/v/release/Itzhep/gomon)
![GitHub Repo stars](https://img.shields.io/github/stars/Itzhep/gomon)

## 📦 Features

- 🔄 Automatic file watching and reloading
- 🎨 CLI with color support for better visibility
- 🔑 GitHub integration for releases
- 📝 Simple and clean configuration
- Docker Support
- .git support
## 🏗️ Installation

### Via Go

To install Gomon, use the following command:

```bash
go install github.com/Itzhep/gomon@latest
```

### Manual Build

1. Clone the repository:

    ```bash
    git clone https://github.com/Itzhep/gomon.git
    cd gomon
    ```

2. Build the project:

    ```bash
    go build -o gomon
    ```
3. Install the Project:

    ```bash
    go install
    ```
4. Move to bin (optinal):
```bash
move gomon.exe C:\path\to\your\bin
```
## 🚀 Usage

1. Start Gomon with the path to your main application file:

    ```bash
    gomon --app path/to/your/app.go
    ```

2. **Press `rs`** in the CLI to manually restart the application.

## 🛠️ Configuration

Gomon supports a variety of configurations directly from the CLI. You can specify the file to watch and other options like color settings for better CLI appearance.

## 📝 Example

Here's a basic example of how to use Gomon:

```bash
gomon -app path/to/your/app.go
```

This command starts Gomon, watches for file changes, and restarts your application automatically.
and in New update V2 it will build your app and run it 
## 🗂️ Contributing

Contributions are welcome! Please submit a pull request or open an issue if you find any bugs or have suggestions.

## 🌟 Star the Project

If you find Gomon useful, please give it a star on [GitHub](https://github.com/Itzhep/gomon) to support development and stay updated with new features.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
![GitHub License](https://img.shields.io/github/license/Itzhep/gomon)

---

🚀 Happy coding with Gomon!
>>>>>>> 4a08f121f7a71f3eb64a4ac6bbf605c6dedd6bfd
