package storage

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileInfo struct {
	Name      string
	Size      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FileStorage struct {
	dir string
	mu  sync.RWMutex
}

func NewFileStorage(dir string) *FileStorage {
	os.MkdirAll(dir, 0755)
	return &FileStorage{dir: dir}
}

func (fs *FileStorage) Save(name string, content []byte) (*FileInfo, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	filePath := filepath.Join(fs.dir, name)
	now := time.Now()

	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:      name,
		Size:      info.Size(),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (fs *FileStorage) List() ([]FileInfo, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	entries, err := os.ReadDir(fs.dir)
	if err != nil {
		return nil, err
	}

	var fileInfos []FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfos = append(fileInfos, FileInfo{
			Name:      info.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
			UpdatedAt: info.ModTime(),
		})
	}

	return fileInfos, nil
}

func (fs *FileStorage) Get(name string) ([]byte, *FileInfo, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	filePath := filepath.Join(fs.dir, name)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, nil, err
	}

	return content, &FileInfo{
		Name:      name,
		Size:      info.Size(),
		CreatedAt: info.ModTime(),
		UpdatedAt: info.ModTime(),
	}, nil
}
