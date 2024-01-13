package pig

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	ContentType string
	FilePath    string
}

func NewFile(filePath string) *File {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	// detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		panic(err)
	}
	contentType := http.DetectContentType(buffer)

	return &File{
		FilePath:    filePath,
		ContentType: contentType,
	}
}

func (f *File) Move(dest string) *File {
	err := os.Rename(f.FilePath, dest)
	if err != nil {
		panic(err)
	}

	f.FilePath = dest
	return f
}

func (f *File) ArchiveMove(dir string) *File {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		panic(err)
	}

	if !fileInfo.IsDir() {
		panic(fmt.Sprintf("dir %s is not directory", dir))
	}

	// move
	now := time.Now()
	destDir := fmt.Sprintf(
		"%s/%s/%s",
		dir,
		now.Format("200601"),
		now.Format("02"),
	)

	err = os.MkdirAll(destDir, 0777)
	if err != nil {
		panic(err)
	}

	basename := filepath.Base(f.FilePath)
	dest := fmt.Sprintf("%s/%s", destDir, basename)
	f.Move(dest)
	f.FilePath = dest
	return f
}
