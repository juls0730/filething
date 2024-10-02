package routes

import (
	"archive/zip"
	"bytes"
	"filething/models"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type UploadResponse struct {
	Usage int64 `json:"usage"`
	File  File  `json:"file"`
}

func UploadFile(c fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	fullPath := strings.Trim(c.Params("*"), "/")
	basePath := fmt.Sprintf("%s/%s/%s/", os.Getenv("STORAGE_PATH"), user.ID, fullPath)

	currentUsage, err := calculateStorageUsage(fmt.Sprintf("%s/%s", os.Getenv("STORAGE_PATH"), user.ID))
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = os.Stat(basePath)
	directoryExists := err == nil

	// Create the directories if they don't exist
	if !directoryExists {
		err = os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if string(c.Request().Header.ContentType()) == "application/json" {
		if err == http.ErrNotMultipart {
			if directoryExists {
				// Directories exist, but no file was uploaded
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "A folder with that name already exists"})
			}
			// Directories were just created, and no file was provided
			entry, err := os.Stat(basePath)

			if err != nil {
				fmt.Println(err)
				return err
			}

			uploadFile := &UploadResponse{
				Usage: currentUsage + entry.Size(),
				File: File{
					Name:         entry.Name(),
					IsDir:        entry.IsDir(),
					Size:         entry.Size(),
					LastModified: entry.ModTime().Format("2 Jan 06"),
				},
			}

			return c.Status(http.StatusOK).JSON(uploadFile)
		}
		fmt.Println(err)
		return err
	}

	_, params, err := mime.ParseMediaType(string(c.Request().Header.ContentType()))
	if err != nil {
		log.Fatal(err)
	}

	reader := multipart.NewReader(c.Request().BodyStream(), params["boundary"])

	part, err := reader.NextPart()
	if err != nil {
		fmt.Println(err)
		return err
	}

	filepath := filepath.Join(basePath, part.FileName())

	if _, err = os.Stat(filepath); err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"message": "File with that name already exists"})
	}

	dst, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	// Read the file manually because otherwise we are limited by the arbitrarily small size of /tmp
	buffer := make([]byte, 1*1024*1024)
	totalSize := int64(0)

	for {
		n, readErr := part.Read(buffer)

		if readErr != nil && readErr == io.ErrUnexpectedEOF {
			dst.Close()
			os.Remove(filepath)
			return c.Status(http.StatusRequestTimeout).JSON(fiber.Map{"message": "Upload canceled"})
		}

		if readErr != nil && readErr != io.EOF {
			fmt.Println(err)
			return readErr
		}

		totalSize += int64(n)

		if currentUsage+totalSize > user.Plan.MaxStorage {
			dst.Close()
			os.Remove(filepath)
			return c.Status(http.StatusInsufficientStorage).JSON(fiber.Map{"message": "Insufficient storage space"})
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
					LastModified: entry.ModTime().Format("2 Jan 06"),
				},
			}

			return c.Status(http.StatusOK).JSON(uploadFile)
		}
	}
}

func calculateStorageUsage(basePath string) (int64, error) {
	var totalSize int64
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return 0, err
	}

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

func GetFiles(c fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	fullPath := strings.Trim(c.Params("*"), "/")
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
			LastModified: f.ModTime().Format("2 Jan 06"),
		})
	}

	return c.Status(http.StatusOK).JSON(jsonFiles)
}

func GetFile(c fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	fullPath := strings.Trim(c.Params("*"), "/")

	fileNamesParam := c.Query("filenames")
	var fileNames []string
	if fileNamesParam != "" {
		fileNames = strings.Split(fileNamesParam, ",")
	}

	if fullPath == "" && len(fileNames) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "A file is required"})
	}

	basePath := fmt.Sprintf("%s/%s", os.Getenv("STORAGE_PATH"), user.ID)
	if fullPath != "" {
		basePath = filepath.Join(basePath, fullPath)
	}

	fileInfo, err := os.Stat(basePath)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "No file found!"})
	}

	var buf bytes.Buffer
	if fileInfo.IsDir() {
		c.Type("application/zip")

		if len(fileNames) != 0 {
			err := zipFiles(&buf, filepath.Join(basePath, fullPath), fileNames)
			if err != nil {
				fmt.Println(err)
				return err
			}

			_, err = buf.WriteTo(c.Response().BodyWriter())
			return err
		}

		err := zipFiles(&buf, basePath, []string{""})
		if err != nil {
			fmt.Println(err)
			return err
		}

		_, err = buf.WriteTo(c.Response().BodyWriter())
		return err
	} else {
		return c.SendFile(basePath)
	}
}

func zipFiles(buf *bytes.Buffer, basePath string, files []string) error {
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	for _, filePath := range files {
		unescapedFilePath, err := url.PathUnescape(filePath)
		if err != nil {
			return err
		}
		err = processFile(zipWriter, basePath, unescapedFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func processFile(zipWriter *zip.Writer, basePath string, filePath string) error {
	fullFilePath := filepath.Join(basePath, filePath)

	fileInfo, err := os.Stat(fullFilePath)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate

	header.Name = filepath.ToSlash(filePath)

	if fileInfo.IsDir() {
		header.Name += "/"
	}

	headerWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		files, err := os.ReadDir(fullFilePath)
		if err != nil {
			return err
		}

		for _, file := range files {
			err := processFile(zipWriter, basePath, filepath.Join(filePath, file.Name()))
			if err != nil {
				return err
			}
		}
		return nil
	}

	file, err := os.Open(fullFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(headerWriter, file)
	return err
}

type DeleteRequest struct {
	Files []File `json:"files"`
}

func DeleteFiles(c fiber.Ctx) error {
	deleteData := new(DeleteRequest)

	if err := c.Bind().JSON(deleteData); err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "An unknown error occoured!"})
	}

	if len(deleteData.Files) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Files are required!"})
	}

	user := c.Locals("user").(*models.User)

	fullPath := strings.Trim(c.Params("*"), "/")
	basePath := fmt.Sprintf("%s/%s/%s", os.Getenv("STORAGE_PATH"), user.ID, fullPath)

	for _, file := range deleteData.Files {
		path := filepath.Join(basePath, file.Name)
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "An unknown error occoured!"})
		}
	}

	word := "file"
	fileLen := len(deleteData.Files)

	if fileLen != 1 {
		word = word + "s"
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf("Successfully deleted %d %s", fileLen, word)})
}
