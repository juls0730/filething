package routes

import (
	"filething/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type UploadResponse struct {
	Usage int64 `json:"usage"`
	File  File  `json:"file"`
}

func UploadFile(c echo.Context) error {
	user := c.Get("user").(*models.User)

	fullPath := strings.Trim(c.Param("*"), "/")
	basePath := fmt.Sprintf("%s/%s/%s/", os.Getenv("STORAGE_PATH"), user.ID, fullPath)

	currentUsage, err := calculateStorageUsage(basePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	reader, err := c.Request().MultipartReader()
	if err != nil {
		fmt.Println(err)
		return err
	}

	part, err := reader.NextPart()
	if err != nil {
		fmt.Println(err)
		return err
	}

	filepath := filepath.Join(basePath, part.FileName())

	if _, err = os.Stat(filepath); err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "File with that name already exists"})
	}

	dst, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	// Read the file manually because otherwise we are limited by the arbitrarily small size of /tmp
	buffer := make([]byte, 4096)
	totalSize := int64(0)

	for {
		n, readErr := part.Read(buffer)

		if readErr != nil && readErr == io.ErrUnexpectedEOF {
			dst.Close()
			os.Remove(filepath)
			return c.JSON(http.StatusRequestTimeout, map[string]string{"message": "Upload canceled"})
		}

		if readErr != nil && readErr != io.EOF {
			fmt.Println(err)
			return readErr
		}

		totalSize += int64(n)

		if currentUsage+totalSize > user.Plan.MaxStorage {
			dst.Close()
			os.Remove(filepath)
			return c.JSON(http.StatusInsufficientStorage, map[string]string{"message": "Insufficient storage space"})
		}

		if _, err := dst.Write(buffer[:n]); err != nil {
			fmt.Println(err)
			return err
		}

		if n == 0 || readErr == io.EOF {
			entry, err := os.Stat(filepath)

			if err != nil {
				fmt.Println(err)
				return err
			}

			uploadFile := &UploadResponse{
				Usage: currentUsage + totalSize,
				File: File{
					Name:         entry.Name(),
					IsDir:        entry.IsDir(),
					Size:         entry.Size(),
					LastModified: entry.ModTime().Format("1/2/2006"),
				},
			}

			return c.JSON(http.StatusOK, uploadFile)
		}
	}
}

func calculateStorageUsage(basePath string) (int64, error) {
	var totalSize int64

	// Read the directory
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return 0, err
	}

	// Iterate over directory entries
	for _, entry := range entries {
		if entry.IsDir() {
			// Recursively calculate size of directories
			dirPath := filepath.Join(basePath, entry.Name())
			dirSize, err := calculateStorageUsage(dirPath)
			if err != nil {
				return 0, err
			}
			totalSize += dirSize
		} else {
			// Calculate size of file
			_ = filepath.Join(basePath, entry.Name())
			fileInfo, err := entry.Info()
			if err != nil {
				return 0, err
			}
			totalSize += fileInfo.Size()
		}
	}

	return totalSize, nil
}

type File struct {
	Name         string `json:"name"`
	IsDir        bool   `json:"is_dir"`
	Size         int64  `json:"size"`
	LastModified string `json:"last_modified"`
}

func GetFiles(c echo.Context) error {
	user := c.Get("user").(*models.User)

	fullPath := strings.Trim(c.Param("*"), "/")
	basePath := fmt.Sprintf("%s/%s/%s/", os.Getenv("STORAGE_PATH"), user.ID, fullPath)

	f, err := os.Open(basePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	jsonFiles := make([]File, 0)

	for _, f := range files {
		jsonFiles = append(jsonFiles, File{
			Name:         f.Name(),
			IsDir:        f.IsDir(),
			Size:         f.Size(),
			LastModified: f.ModTime().Format("1/2/2006"),
		})
	}

	return c.JSON(http.StatusOK, jsonFiles)
}

func GetUsage(c echo.Context) error {
	user := c.Get("user").(*models.User)

	fullPath := strings.Trim(c.Param("*"), "/")
	basePath := fmt.Sprintf("%s/%s/%s/", os.Getenv("STORAGE_PATH"), user.ID, fullPath)
	storageUsage, err := calculateStorageUsage(basePath)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]int64{"usage": storageUsage})
}
