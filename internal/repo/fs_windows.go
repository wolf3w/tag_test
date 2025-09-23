//go:build windows

// Package repo НЕ работает вне винды. Такие вот дела.
package repo

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/wolf3w/tag_test/internal/domain"
)

type FileStorage struct {
	path string
}

func NewFileStorage(root string) (*FileStorage, error) {
	pathToDir := root + "pictures"
	if _, err := os.Stat(pathToDir); os.IsNotExist(err) {
		// Да, перезаписываем ошибку. Потому что ошибка была сигналом об отсутствии директории.
		err = os.Mkdir(pathToDir, 0777)
		if err != nil {
			return nil, fmt.Errorf("create pictures' directory: %w", err)
		}
	}

	return &FileStorage{path: pathToDir}, nil
}

// Write открыть файл для записи байтов в него. Горутины здесь не водятся, записываем целый собранный слайс байт.
func (fs *FileStorage) Write(filename string, collectedData []byte) error {
	file, err := os.OpenFile(fs.path+"/"+filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	// Ошибку закрытия игнорим, поскольку поймаем её на записи
	defer file.Close()

	_, err = file.Write(collectedData)

	return err
}

// Read читаем весь файл, склеиваем и получаем картинку в виде слайса байт.
func (fs *FileStorage) Read(filename string) ([]byte, error) {
	file, err := os.Open(fs.path + "/" + filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	var payload []byte
	buffer := make([]byte, 1024)
	for {
		_, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read file: %w", err)
		}
		payload = append(payload, buffer...)
	}
	return payload, nil
}

// ListPictures метод для вывода списка файлов. Внимание: работает только на винде. Делать кроссплатформенное хранение
// файлов с именем вида {file_name}_{create_time}_{update_time} мне было лень. К тому же это долгий поиск по всей папке.
func (fs *FileStorage) ListPictures() ([]domain.PictureInfo, error) {
	entries, err := os.ReadDir(fs.path)
	if err != nil {
		return nil, fmt.Errorf("cannot read dir: %w", err)
	}

	picInfo := make([]domain.PictureInfo, 0, len(entries))

	for _, file := range entries {
		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("cannot fetch info: %w", err)
		}

		winFileInfo := info.Sys().(*syscall.Win32FileAttributeData)

		picInfo = append(picInfo, domain.PictureInfo{
			Name:      info.Name(),
			CreatedAt: time.Unix(0, winFileInfo.CreationTime.Nanoseconds()),
			UpdatedAt: time.Unix(0, winFileInfo.LastWriteTime.Nanoseconds()),
		})
	}
	return picInfo, nil
}
