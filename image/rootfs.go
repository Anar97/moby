package image // import "github.com/docker/docker/image"

import (
	"runtime"

	"github.com/docker/docker/layer"
	"github.com/sirupsen/logrus"
)

// TypeLayers is used for RootFS.Type for filesystems organized into layers.
const TypeLayers = "layers"

// typeLayersWithBase is an older format used by Windows up to v1.12. We		老版本的，目前不用了
// explicitly handle this as an error case to ensure that a daemon which still		当错误处理
// has an older image like this on disk can still start, even though the
// image itself is not usable. See https://github.com/docker/docker/pull/25806.
const typeLayersWithBase = "layers+base"

// RootFS describes images root filesystem						RooFS描述了镜像的根文件系统
// This is currently a placeholder that only supports layers. In the future		这目前只是一个只支持layer的占位符，在将来可能被
// this can be made into an interface that supports different implementations.		到interface里支持不同的接口
type RootFS struct {
	Type    string         `json:"type"`
	DiffIDs []layer.DiffID `json:"diff_ids,omitempty"`
}

// NewRootFS returns empty RootFS struct
func NewRootFS() *RootFS {
	return &RootFS{Type: TypeLayers}
}

// Append appends a new diffID to rootfs
func (r *RootFS) Append(id layer.DiffID) {
	r.DiffIDs = append(r.DiffIDs, id)
}

// Clone returns a copy of the RootFS
func (r *RootFS) Clone() *RootFS {
	newRoot := NewRootFS()
	newRoot.Type = r.Type
	newRoot.DiffIDs = make([]layer.DiffID, len(r.DiffIDs))
	copy(newRoot.DiffIDs, r.DiffIDs)
	return newRoot
}

// ChainID returns the ChainID for the top layer in RootFS.
func (r *RootFS) ChainID() layer.ChainID {
	if runtime.GOOS == "windows" && r.Type == typeLayersWithBase {//系统是windows且
		logrus.Warnf("Layer type is unsupported on this platform. DiffIDs: '%v'", r.DiffIDs)
		return ""
	}
	return layer.CreateChainID(r.DiffIDs)	//返回ChainID
}
