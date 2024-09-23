# 🛠️ Gomon - A Nodemon clone Go File Watcher 🚀

Gomon is a Go-based file watcher that automatically reloads your application when file changes are detected. Inspired by nodemon, it helps streamline development by automatically restarting your app.

## 📦 Features

- 🔄 Automatic file watching and reloading
- 🎨 CLI with color support for better visibility
- 🔑 GitHub integration for releases
- 📝 Simple and clean configuration

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
    go build
        ```
3. Install the Project:

    ```bash
    go install
    ```
## 🚀 Usage

1. Start Gomon with the path to your main application file:

    ```bash
    gomon -app path/to/your/app.go
    ```

2. **Press `rs`** in the CLI to manually restart the application.

## 🛠️ Configuration

Gomon supports a variety of configurations directly from the CLI. You can specify the file to watch and other options like color settings for better CLI appearance.


This command starts Gomon, watches for file changes, and restarts your application automatically.

## 🗂️ Contributing

Contributions are welcome! Please submit a pull request or open an issue if you find any bugs or have suggestions.

## 🌟 Star the Project

If you find Gomon useful, please give it a star on [GitHub](https://github.com/Itzhep/gomon) to support development and stay updated with new features.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
![GitHub License](https://img.shields.io/github/license/Itzhep/gomon)

---

🚀 Happy coding with Gomon!
