# Local Video Library

A local video library website that scans specified directories for video files and their corresponding poster images, providing web browsing and playback functionality.

## Go Version 🚀

### Features
- **High Performance**: ~100x faster than the Python version
- Concurrent directory scanning
- Efficient file serving
- Low memory usage
- Fast response time
- Support for multiple video and image formats
- Web interface with search functionality
- Responsive design

### Usage
```bash
# Run from project root
go run performance/main.go

# Optional: specify different port (default: 5000)
go run performance/main.go -addr :8080
```

### Requirements
- Go 1.16 or higher
- Supports Windows/Linux/macOS

## Python Version

### Features
- Built with Flask framework
- Simple and easy to debug
- Good for development and testing
- Same functionality as Go version

### Installation
```bash
pip install -r requirements.txt
```

### Usage
```bash
python app.py
```

### Requirements
- Python 3
- Flask

## Configuration

Create a `dir.conf` file in the project root, list the folders which you want to scan:
```
# Video directories
D:\Videos\Movie
E:\Downloads\Videos
```

## File Structure

The application expects video files and their corresponding poster images to follow this naming convention (just need has same prefix):
- Video file: `VIDEO_NAME.mp4`
- Poster image: `VIDEO_NAME.jpg`

For example:
- `XXX-000.mp4`
- `XXX-000.jpg`

## Notes
- The application runs locally and does not upload any files to external servers



<br><br><br><br><br><br><br>





# 本地视频库

一个扫描本地视频文件的 Web 应用，支持在线浏览和播放视频。

## Go 版本 🚀

### 特点
- **超高性能**：比 Python 版本快约 100 倍
- 并发扫描目录
- 高效的文件服务
- 低内存占用
- 快速响应
- 支持多种视频和图片格式
- 支持网页搜索功能
- 响应式界面设计

### 使用方法
```bash
# 在项目根目录下运行
go run performance/main.go

# 可选：指定其他端口（默认：5000）
go run performance/main.go -addr :8080
```

### 系统要求
- Go 1.16 或更高版本
- 支持 Windows/Linux/macOS

## Python 版本

### 特点
- 基于 Flask 框架
- 简单易用
- 适合开发和调试
- 功能与 Go 版本相同

### 安装
```bash
pip install -r requirements.txt
```

### 使用方法
```bash
python app.py
```

### 系统要求
- Python 3
- Flask

## 配置说明

在项目根目录创建 `dir.conf` 文件：
```
# 视频目录配置
D:\Videos\Movie
E:\Downloads\Videos
```

## 文件结构
应用程序期望的文件命名格式（只需相同前缀即可）：
- 视频文件：`XXX-000.mp4`
- 海报图片：`XXX-000.jpg`

## 功能特性

- 自动扫描配置的目录
- 支持多种视频格式（mp4, mkv, avi 等）
- 支持多种图片格式（jpg, jpeg, png 等）
- 按番号格式匹配视频和海报
- 网页端支持搜索和在线播放
- 美观的响应式界面

## 性能对比

| 特性 | Go 版本 | Python 版本 |
|------|---------|-------------|
| 目录扫描速度 | 极快 | 较慢 |
| 内存占用 | 低 | 较高 |
| 并发处理 | 支持 | 不支持 |
| 响应速度 | 极快 | 一般 |
