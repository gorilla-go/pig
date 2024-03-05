package param

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	ContentType string
	Path        string
	Basename    string
	Ext         string
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
		Path:        filePath,
		ContentType: contentType,
		Basename:    filepath.Base(filePath),
		Ext:         filepath.Ext(filePath),
	}
}

func (f *File) Move(dest string) *File {
	err := os.Rename(f.Path, dest)
	if err != nil {
		panic(err)
	}

	f.Path = dest
	return f
}

func (f *File) ArchiveMove(dir string) (*File, string) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		panic(err)
	}

	if !fileInfo.IsDir() {
		panic(fmt.Sprintf("dir %s is not directory", dir))
	}

	// move
	now := time.Now()
	relativePath := fmt.Sprintf("%s/%s", now.Format("200601"), now.Format("02"))
	destDir := fmt.Sprintf(
		"%s/%s",
		dir,
		relativePath,
	)

	err = os.MkdirAll(destDir, 0777)
	if err != nil {
		panic(err)
	}

	basename := filepath.Base(f.Path)
	dest := fmt.Sprintf("%s/%s", destDir, basename)
	f.Move(dest)
	f.Path = dest
	return f, fmt.Sprintf("%s/%s", relativePath, basename)
}
