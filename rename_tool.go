package main

/*
重命名工具使用说明：

功能：
  - 将视频和图片文件名标准化为 "系列代号-编号" 的格式
  - 例如：abc123pl.jpg -> ABC-123.jpg
  - 支持处理 .mp4, .mkv, .avi, .mov, .wmv, .flv, .webm 等视频文件
  - 支持处理 .jpg, .jpeg, .png, .gif, .webp 等图片文件

使用方法：
  1. 预览模式（不实际重命名）:
     go run rename_tool.go -path="E:\Videos" -readonly=true

  2. 执行重命名:
     go run rename_tool.go -path="E:\Videos" -readonly=false

  3. 递归处理子目录:
     go run rename_tool.go -path="E:\Videos" -readonly=false -recursive=true

  4. 简易输出模式（只显示重命名信息）:
     go run rename_tool.go -path="E:\Videos" -readonly=true -simple=true

参数说明：
  -path: 指定要处理的目录路径（必需）
  -readonly: 是否为预览模式（默认为 true）
  -recursive: 是否递归处理子目录（默认为 false）
  -simple: 是否使用简易输出模式（默认为 false）

文件名处理规则：
  - 移除文件名中的 "pl" 后缀
  - 移除文件名中的 "FHD" 或 "HD" 标记
  - 将系列代号转为大写
  - 在系列代号和编号之间添加连字符
  - 保持原始文件扩展名不变
*/

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Stats 记录重命名操作的统计信息
type Stats struct {
	total   int
	renamed int
	skipped int
	errors  int
}

func main() {
	// 命令行参数
	path := flag.String("path", "", "Path to the video files")
	readOnly := flag.Bool("readonly", true, "Read-only mode (true: only preview changes)")
	recursive := flag.Bool("recursive", false, "Recursively process subdirectories")
	simple := flag.Bool("simple", false, "Simple output mode (only show rename operations)")
	flag.Parse()

	if *path == "" {
		fmt.Println("Please provide the path to the video files using the -path flag")
		fmt.Println("Usage example:")
		fmt.Println("  go run rename_tool.go -path='F:\\Videos' -readonly=true")
		fmt.Println("  go run rename_tool.go -path='F:\\Videos' -readonly=false -recursive=true")
		fmt.Println("  go run rename_tool.go -path='F:\\Videos' -readonly=true -simple=true")
		return
	}

	// 打印运行模式
	if !*simple {
		mode := "PREVIEW MODE"
		if !*readOnly {
			mode = "RENAME MODE"
		}
		fmt.Printf("\nRunning in %s\n", mode)
		fmt.Printf("Processing directory: %s\n", *path)
		if *recursive {
			fmt.Println("Recursive mode: enabled")
		}
		fmt.Println(strings.Repeat("-", 50))
	}

	// 执行重命名操作
	stats := Stats{}
	err := renameVideos(*path, *readOnly, *recursive, *simple, &stats)
	if err != nil {
		fmt.Printf("\nError during execution: %v\n", err)
	}

	// 打印统计信息
	if !*simple {
		printStats(stats)
	}
}

func renameVideos(path string, readOnly, recursive, simple bool, stats *Stats) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", path, err)
	}

	for _, entry := range entries {
		if recursive && entry.IsDir() {
			subdir := filepath.Join(path, entry.Name())
			if err := renameVideos(subdir, readOnly, recursive, simple, stats); err != nil {
				if !simple {
					fmt.Printf("Error processing subdirectory %s: %v\n", subdir, err)
				}
			}
			continue
		}

		if !entry.IsDir() && (isVideoFile(entry.Name()) || isImageFile(entry.Name())) {
			stats.total++
			oldName := entry.Name()
			newName := getNewFileName(oldName)

			if newName == "" {
				stats.skipped++
				if !simple {
					fmt.Printf("SKIP: %s (no match found)\n", oldName)
				}
				continue
			}

			if newName == oldName {
				stats.skipped++
				if !simple {
					fmt.Printf("SKIP: %s (already in correct format)\n", oldName)
				}
				continue
			}

			oldPath := filepath.Join(path, oldName)
			newPath := filepath.Join(path, newName)

			// 检查目标文件是否已存在
			if _, err := os.Stat(newPath); err == nil {
				stats.errors++
				fmt.Printf("ERROR: Target file already exists: %s\n", newPath)
				continue
			}

			if !simple {
				fmt.Printf("Found: %s", oldName)
				fmt.Printf("  -> %s\n", newName)
			}

			if !readOnly {
				if err := os.Rename(oldPath, newPath); err != nil {
					stats.errors++
					fmt.Printf("ERROR: Failed to rename: %v\n", err)
				} else {
					stats.renamed++
					fmt.Printf("RENAMED: %s -> %s\n", oldName, newName)
				}
			} else {
				stats.renamed++
				fmt.Printf("PREVIEW: Would rename %s -> %s\n", oldName, newName)
			}
		}
	}

	return nil
}

func isVideoFile(filename string) bool {
	extension := strings.ToLower(filepath.Ext(filename))
	videoExtensions := map[string]bool{
		".mp4":  true,
		".avi":  true,
		".mkv":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".webm": true,
		".srt":  true,
	}
	return videoExtensions[extension]
}

func isImageFile(filename string) bool {
	extension := strings.ToLower(filepath.Ext(filename))
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return imageExtensions[extension]
}

func getNewFileName(oldName string) string {
	// 保存原始扩展名
	ext := filepath.Ext(oldName)
	baseName := strings.TrimSuffix(oldName, ext)

	// 移除 pl 后缀（不区分大小写）
	baseName = regexp.MustCompile(`(?i)pl$`).ReplaceAllString(baseName, "")

	// 移除 .FHD 或 .HD 后缀（不区分大小写）
	baseName = regexp.MustCompile(`(?i)\.(?:fhd|hd)$`).ReplaceAllString(baseName, "")

	// 移除 -FHD 或 -HD 后缀（不区分大小写）
	baseName = regexp.MustCompile(`(?i)-(?:fhd|hd)$`).ReplaceAllString(baseName, "")

	// 匹配番号格式（支持带和不带连字符的格式）
	re := regexp.MustCompile(`(?i)^([a-z]{2,})[-]?(\d+)$`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 3 {
		return "" // 未找到匹配的格式
	}

	// 格式化新文件名
	seriesCode := strings.ToUpper(matches[1])
	number := matches[2]

	// 返回新文件名（保持原始扩展名）
	return fmt.Sprintf("%s-%s%s", seriesCode, number, ext)
}

func printStats(stats Stats) {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("Summary:")
	fmt.Printf("Total files processed: %d\n", stats.total)
	fmt.Printf("Files to be renamed: %d\n", stats.renamed)
	fmt.Printf("Files skipped: %d\n", stats.skipped)
	fmt.Printf("Errors encountered: %d\n", stats.errors)
	fmt.Println(strings.Repeat("-", 50))
}

// 添加测试函数
func init() {
	// 测试用例
	testCases := []struct {
		input    string
		expected string
	}{
		{"ABC-123.FHD.mkv", "ABC-123.mkv"},
		{"def456pl.jpg", "DEF-456.jpg"},
		{"GHI-789-HD.mp4", "GHI-789.mp4"},
		{"MNO-345.HD.avi", "MNO-345.avi"},
		{"pqr-678-fhd.wmv", "PQR-678.wmv"},
	}

	// 运行测试
	fmt.Println("\nRunning self-test...")
	fmt.Println(strings.Repeat("-", 50))
	for _, tc := range testCases {
		result := getNewFileName(tc.input)
		if result == tc.expected {
			fmt.Printf("PASS: %s -> %s\n", tc.input, result)
		} else {
			fmt.Printf("FAIL: %s -> %s (expected: %s)\n", tc.input, result, tc.expected)
		}
	}
	fmt.Println(strings.Repeat("-", 50))
}
