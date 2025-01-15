package fileutil

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexmullins/zip"
)

const (
	// DefaultFileMode is the default permission for created files
	DefaultFileMode = 0644
	// MaxTreeSize is the maximum size of the tree string before truncation
	MaxTreeSize = 4090
	// TreeTruncatedMessage is displayed when the tree is too large
	TreeTruncatedMessage = "Too many files to display"
)

// AppendFile appends a line to a file, creating it if it doesn't exist
// Returns error if the operation fails
func AppendFile(path string, line string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, DefaultFileMode)
	if err != nil {
		return fmt.Errorf("failed to open file for append: %w", err)
	}
	defer file.Close()
	
	if _, err := file.WriteString(line + "\n"); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

// Tree generates a tree-like string representation of a directory structure
// prefix is used for indentation, isFirstDir is used for the root directory
func Tree(path string, prefix string, isFirstDir ...bool) string {
	var sb strings.Builder

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return ""
	}

	for i, file := range files {
		isLast := i == len(files)-1
		var pointer string
		if isLast {
			pointer = prefix + "└── "
		} else {
			pointer = prefix + "├── "
		}
		if isFirstDir == nil {
			pointer = prefix
		}
		if file.IsDir() {
			fmt.Fprintf(&sb, "%s📂 - %s\n", pointer, file.Name())
			if isLast {
				sb.WriteString(Tree(filepath.Join(path, file.Name()), prefix+"    ", false))
			} else {
				sb.WriteString(Tree(filepath.Join(path, file.Name()), prefix+"│   ", false))
			}
		} else {
			fmt.Fprintf(&sb, "%s📄 - %s (%.2f kb)\n", pointer, file.Name(), float64(file.Size())/1024)
		}
	}

	tree := sb.String()
	if len(tree) > MaxTreeSize {
		tree = TreeTruncatedMessage
	}
	return tree
}

// Zip compresses a directory into a zip file
// Returns error if any operation fails during compression
func Zip(dirPath string, zipName string) error {
	zipFile, err := os.Create(zipName)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %q: %w", filePath, err)
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return fmt.Errorf("failed to create zip entry: %w", err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %q: %w", filePath, err)
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return fmt.Errorf("failed to write file to zip: %w", err)
		}

		return nil
	})

	return err
}

// ZipWithPassword compresses a directory into a password-protected zip file
// Returns error if any operation fails during compression
func ZipWithPassword(dirPath string, zipName string, password string) error {
	zipFile, err := os.Create(zipName)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %q: %w", filePath, err)
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		header := &zip.FileHeader{
			Name:     relPath,
			Method:   zip.Deflate,
			Modified: info.ModTime(),
		}
		header.SetPassword(password)

		zipEntry, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip entry: %w", err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %q: %w", filePath, err)
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return fmt.Errorf("failed to write file to zip: %w", err)
		}

		return nil
	})

	return err
}

// Copy copies a file or directory to a destination
// Returns error if the operation fails
func Copy(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if srcInfo.IsDir() {
		return CopyDir(src, dst)
	}
	return CopyFile(src, dst)
}

// CopyFile copies a single file from src to dst
// Returns error if the operation fails
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// CopyDir recursively copies a directory tree from src to dst
// Returns error if the operation fails
func CopyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IsDir checks if a path is a directory
// Returns true if the path exists and is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Exists checks if a path exists
// Returns true if the path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads the entire contents of a file
// Returns the contents as a string and any error encountered
func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// ReadLines reads a file line by line
// Returns a slice of strings containing each line and any error encountered
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	return lines, nil
}
