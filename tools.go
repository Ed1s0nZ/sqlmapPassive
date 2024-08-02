package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var processedFiles = make(map[string]bool) // 记录已处理文件的map
var mutex sync.Mutex

func calculateAndSaveMD5(input string) error {
	hash := md5.New()
	hash.Write([]byte(input))
	md5Hash := hex.EncodeToString(hash.Sum(nil))
	dir := "./sqlmapReq/"
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(files) >= numTxt {
		return errors.New("文件夹中已有" + fmt.Sprint(numTxt) + "个文件，不再写入")
	}
	err = os.WriteFile(filepath.Join(dir, md5Hash+".txt"), []byte(input), 0644)
	if err != nil {
		return err
	}
	return nil
}

func sqlmap() {
	fileQueue := make(chan string, 100)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(fileQueue, &wg)
	}

	go watchDirectory(fileQueue)

	wg.Wait()
}

func watchDirectory(fileQueue chan<- string) {
	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalf("Failed to read directory: %s", err)
		}

		mutex.Lock()
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".txt" {
				filePath := filepath.Join(dir, file.Name())
				if !processedFiles[filePath] {
					fileQueue <- filePath
					processedFiles[filePath] = true
				}
			}
		}
		mutex.Unlock()

		// 清理已处理文件的map
		cleanProcessedFiles()

		time.Sleep(5 * time.Second)
	}
}

func worker(fileQueue <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case filePath, ok := <-fileQueue:
			if !ok {
				return // 通道已关闭，退出goroutine
			}
			executeCommand(filePath)
			err := os.Remove(filePath)
			if err != nil {
				log.Printf("Failed to delete file: %s", err)
			} else {
				fmt.Printf("Deleted file: %s\n", filePath)
			}
			mutex.Lock()
			delete(processedFiles, filePath)
			mutex.Unlock()
		case <-time.After(5 * time.Second):
			// 在等待一段时间后检查是否还有新任务或文件
			files, err := os.ReadDir(dir)
			if err != nil {
				log.Fatalf("Failed to read directory: %s", err)
			}

			if len(files) == 0 {
				continue
			}
		}
	}
}

func executeCommand(filePath string) {
	cmd := exec.Command("python3", "./sqlmap/sqlmap.py", "-r", filePath, "--batch", "--output-dir=sqlmapResult")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Executing command for file: %s\n", filePath)

	if err := cmd.Start(); err != nil {
		log.Printf("Command execution failed to start: %s", err)
		// return
	}
	// 等待命令执行完毕
	if err := cmd.Wait(); err != nil {
		log.Printf("Command execution failed: %s", err)
	} else {
		fmt.Printf("Command executed successfully for file: %s\n", filePath)

	}
}

// 定期清理已处理文件的map，删除一些旧的文件路径
func cleanProcessedFiles() {
	mutex.Lock()
	defer mutex.Unlock()

	for path := range processedFiles {
		if fileIsOld(path) {
			delete(processedFiles, path)
		}
	}
}

// 示例函数：判断文件是否过旧
func fileIsOld(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error getting file info: %s", err)
		return false // 返回false以避免误删除文件
	}
	// 判断文件创建时间，假设超过一定时间的文件认为过旧
	threshold := time.Now().Add(-24 * time.Hour)
	return fileInfo.ModTime().Before(threshold)
}
