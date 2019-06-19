package metadata // import "github.com/docker/docker/distribution/metadata"

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/docker/pkg/ioutils"
)

// Store implements a K/V store for mapping distribution-related IDs		将分布相关的ID映射到磁盘上layer的ID和镜像的ID
// to on-disk layer IDs and image IDs. The namespace identifies the type of	命名空间确定了映射的类型
// mapping (i.e. "v1ids" or "artifacts"). MetadataStore is goroutine-safe.
type Store interface {
	// Get retrieves data by namespace and key.
	Get(namespace string, key string) ([]byte, error)
	// Set writes data indexed by namespace and key.
	Set(namespace, key string, value []byte) error
	// Delete removes data indexed by namespace and key.
	Delete(namespace, key string) error
}

// FSMetadataStore uses the filesystem to associate metadata with layer and
// image IDs.
type FSMetadataStore struct {
	sync.RWMutex
	basePath string
}

// NewFSMetadataStore creates a new filesystem-based metadata store.
func NewFSMetadataStore(basePath string) (*FSMetadataStore, error) {
	if err := os.MkdirAll(basePath, 0700); err != nil {
		return nil, err
	}
	return &FSMetadataStore{
		basePath: basePath,
	}, nil
}

func (store *FSMetadataStore) path(namespace, key string) string {		//basePath/namespace/key
	return filepath.Join(store.basePath, namespace, key)
}

// Get retrieves data by namespace and key. The data is read from a file named	通过命名空间和key返回数据
// after the key, stored in the namespace's directory.
func (store *FSMetadataStore) Get(namespace string, key string) ([]byte, error) {
	store.RLock()
	defer store.RUnlock()

	return ioutil.ReadFile(store.path(namespace, key))
}

// Set writes data indexed by namespace and key. The data is written to a file	写入数据
// named after the key, stored in the namespace's directory.
func (store *FSMetadataStore) Set(namespace, key string, value []byte) error {
	store.Lock()
	defer store.Unlock()

	path := store.path(namespace, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutils.AtomicWriteFile(path, value, 0644)
}

// Delete removes data indexed by namespace and key. The data file named after	直接删除目录
// the key, stored in the namespace's directory is deleted.
func (store *FSMetadataStore) Delete(namespace, key string) error {
	store.Lock()
	defer store.Unlock()

	path := store.path(namespace, key)
	return os.Remove(path)
}
