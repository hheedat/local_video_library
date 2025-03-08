package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// VideoInfo 存储视频信息
type VideoInfo struct {
	Title      string `json:"title"`
	VideoPath  string `json:"video_path"`  // URL路径
	PosterPath string `json:"poster_path"` // URL路径
	FilePath   string `json:"file_path"`   // 视频文件的实际路径
	ImagePath  string `json:"-"`           // 海报图片的实际路径，不输出到JSON
}

// 全局变量
var (
	videoData = make(map[string]*VideoInfo)
	dataLock  sync.RWMutex

	// 未匹配的文件列表
	unmatchedFiles []string
	unmatchedLock  sync.Mutex

	// 支持的文件扩展名
	videoExtensions = map[string]bool{
		".mp4": true, ".mkv": true, ".avi": true,
		".mov": true, ".wmv": true, ".flv": true,
		".webm": true,
	}
	imageExtensions = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true,
	}
)

// 加载目录配置
func loadDirectories() ([]string, error) {
	file, err := os.Open("dir.conf")
	if err != nil {
		return nil, fmt.Errorf("error opening dir.conf: %v", err)
	}
	defer file.Close()

	var dirs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			dirs = append(dirs, line)
		}
	}
	return dirs, scanner.Err()
}

// 扫描目录
func scanDirectory(dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Scanning directory: %s\n", dir)

	var videoCount, imageCount int
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return nil
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(info.Name()))
			if videoExtensions[ext] {
				videoCount++
				log.Printf("Found video: %s\n", path)
			} else if imageExtensions[ext] {
				imageCount++
				log.Printf("Found image: %s\n", path)
			}
			processFile(path, info.Name())
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking directory %q: %v\n", dir, err)
	}
	log.Printf("Directory scan completed: %s (Videos: %d, Images: %d)\n", dir, videoCount, imageCount)
}

// 处理单个文件
func processFile(path, name string) {
	ext := strings.ToLower(filepath.Ext(name))
	if !videoExtensions[ext] && !imageExtensions[ext] {
		return
	}

	baseName := getBaseName(name)
	if baseName == "" {
		// 记录未匹配的文件
		unmatchedLock.Lock()
		unmatchedFiles = append(unmatchedFiles, path)
		unmatchedLock.Unlock()
		return
	}

	dataLock.Lock()
	defer dataLock.Unlock()

	info, exists := videoData[baseName]
	if !exists {
		info = &VideoInfo{Title: baseName}
		videoData[baseName] = info
	}

	if videoExtensions[ext] {
		info.VideoPath = fmt.Sprintf("/video/%s", baseName)
		info.FilePath = path
	} else if imageExtensions[ext] {
		info.PosterPath = fmt.Sprintf("/poster/%s", baseName)
		info.ImagePath = path
	}
}

// 获取基础文件名
func getBaseName(filename string) string {
	// 移除扩展名
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 移除 pl 后缀
	baseName = regexp.MustCompile(`(?i)pl$`).ReplaceAllString(baseName, "")

	// 移除 FHD/HD 后缀
	baseName = regexp.MustCompile(`(?i)\.(?:fhd|hd)$`).ReplaceAllString(baseName, "")
	baseName = regexp.MustCompile(`(?i)-(?:fhd|hd)$`).ReplaceAllString(baseName, "")

	// 匹配番号格式
	re := regexp.MustCompile(`(?i)^([a-z]{2,})[-]?(\d+)$`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 3 {
		return ""
	}

	// 格式化文件名
	return fmt.Sprintf("%s-%s", strings.ToUpper(matches[1]), matches[2])
}

// API 处理函数
func handleAPI(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))

	dataLock.RLock()
	var videos []*VideoInfo
	for _, info := range videoData {
		if info.VideoPath != "" && info.PosterPath != "" &&
			(query == "" || strings.Contains(strings.ToLower(info.Title), query)) {
			videos = append(videos, info)
		}
	}
	dataLock.RUnlock()

	// 按标题排序
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Title < videos[j].Title
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// 请求日志中间件
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %s\n", r.Method, r.RequestURI, time.Since(start))
	}
}

// 视频文件处理函数
func handleVideo(w http.ResponseWriter, r *http.Request) {
	baseName := filepath.Base(r.URL.Path)
	dataLock.RLock()
	info, exists := videoData[baseName]
	dataLock.RUnlock()

	if !exists || info.FilePath == "" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, info.FilePath)
}

// 海报图片处理函数
func handlePoster(w http.ResponseWriter, r *http.Request) {
	baseName := filepath.Base(r.URL.Path)
	dataLock.RLock()
	info, exists := videoData[baseName]
	dataLock.RUnlock()

	if !exists || info.ImagePath == "" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, info.ImagePath)
}

//go:embed templates
var templateFS embed.FS

func main() {
	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime)

	// 解析命令行参数
	addr := flag.String("addr", ":5000", "HTTP service address")
	flag.Parse()

	// 加载目录配置
	dirs, err := loadDirectories()
	if err != nil {
		log.Fatalf("Error loading directories: %v", err)
	}

	if len(dirs) == 0 {
		log.Fatal("No directories configured in dir.conf")
	}

	// 并行扫描所有目录
	var wg sync.WaitGroup
	for _, dir := range dirs {
		wg.Add(1)
		go scanDirectory(dir, &wg)
	}
	wg.Wait()

	// 统计视频数量
	var totalVideos int
	dataLock.RLock()
	for _, info := range videoData {
		if info.VideoPath != "" && info.PosterPath != "" {
			totalVideos++
		}
	}
	dataLock.RUnlock()

	log.Printf("Total videos found: %d\n", totalVideos)

	// 输出未匹配的文件
	if len(unmatchedFiles) > 0 {
		log.Printf("\nUnmatched files (%d):\n", len(unmatchedFiles))
		sort.Strings(unmatchedFiles) // 按字母顺序排序
		for _, path := range unmatchedFiles {
			log.Printf("  %s\n", path)
		}
		log.Println() // 添加空行
	}

	// 设置路由（添加日志中间件）
	http.HandleFunc("/api/videos", logRequest(handleAPI))
	http.HandleFunc("/video/", logRequest(handleVideo))
	http.HandleFunc("/poster/", logRequest(handlePoster))

	// 设置首页模板
	tmpl, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	http.HandleFunc("/", logRequest(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmpl.Execute(w, map[string]interface{}{
			"total_videos": totalVideos,
		})
	}))

	// 启动服务器
	log.Printf("Starting server on %s\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
