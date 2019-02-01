// Code generated by go-bindata. DO NOT EDIT.
// sources:
// assets/defaults/cluster-ingress.yaml (570B)
// assets/router/cluster-role-binding.yaml (305B)
// assets/router/cluster-role.yaml (713B)
// assets/router/deployment.yaml (1.724kB)
// assets/router/namespace.yaml (243B)
// assets/router/service-account.yaml (226B)
// assets/router/service-cloud.yaml (628B)
// assets/router/service-internal.yaml (512B)

package manifests

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _assetsDefaultsClusterIngressYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x91\xb1\x6e\xdc\x30\x10\x44\x7b\x7e\xc5\xe0\x5c\x4b\xc1\xb5\x6a\x93\x14\x06\x52\x04\x08\x90\x7e\x8f\x1a\x9d\x36\x47\x2d\x09\x72\x65\x23\xf9\xfa\x40\x96\x94\xc2\x71\xc9\xc1\xf2\xf1\xed\xf0\x09\x9f\xd3\xda\x9c\xf5\xd9\xee\x95\xad\xe1\x55\x7d\xc6\xc8\x49\xd6\xe4\x78\x91\xb4\xb2\x85\x27\x3c\x5b\x73\x49\x89\x15\x31\xdb\xa4\x77\xb4\xc2\xa8\x93\xc6\x63\x04\x52\x09\x29\x25\x29\x47\x88\xa3\xae\xe6\xba\xb0\x0f\x0f\xb5\x71\x78\xf7\x46\x90\xa2\x3f\x59\x9b\x66\x1b\xa0\x7b\xd6\xe7\x42\x6b\xb3\x4e\xde\x6b\xfe\xf4\x72\x95\x54\x66\xb9\x86\x85\x2e\xa3\xb8\x0c\x01\x30\x59\x38\x9c\x6a\xc7\xb9\x15\x89\x1c\xf0\xef\x72\x77\xe0\xba\x5c\x58\xc5\x73\x0d\xc0\xa4\x26\x49\xff\xb0\xb6\x8d\xf2\x84\xaf\xd6\xd6\x4a\xf8\x2c\x8e\x6c\xe9\x37\x7c\x26\xce\x79\x44\x31\x8c\x4c\x74\xbe\xe5\x67\x13\x71\xdf\xe0\xd4\x45\xbe\xfd\x62\xf4\xfe\x3f\x7c\xf7\xf1\x42\x07\xa6\x3b\x30\xa7\x66\xd8\x6a\xdc\xad\x8e\x6e\xbe\xe4\x45\xd4\xa0\xed\x83\x36\xb1\x36\xb5\xfb\x9b\x96\xbe\xfb\x8f\x4d\xc4\xf2\xc8\xef\x49\x22\x17\x9a\x6f\xd0\x3d\xfa\xc1\xc4\xe8\xb9\xee\x09\xb0\x88\xc7\xf9\x9b\xdc\x98\xda\x19\x01\x97\x6d\xb2\xab\x39\xb1\x7f\xac\x37\x56\xa3\xb3\x6d\xe2\xaf\xb9\x3e\x58\x2f\x03\x2e\x97\x00\x54\x96\xa4\x51\xda\x80\x6b\xf8\x1b\x00\x00\xff\xff\xf1\x4e\x69\x66\x3a\x02\x00\x00")

func assetsDefaultsClusterIngressYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsDefaultsClusterIngressYaml,
		"assets/defaults/cluster-ingress.yaml",
	)
}

func assetsDefaultsClusterIngressYaml() (*asset, error) {
	bytes, err := assetsDefaultsClusterIngressYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/defaults/cluster-ingress.yaml", size: 570, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xc6, 0xc4, 0x45, 0xf8, 0xf3, 0xeb, 0x45, 0xb, 0x43, 0x73, 0xb5, 0xee, 0xed, 0x17, 0x6b, 0x3d, 0xb2, 0x5c, 0x5c, 0xa, 0x4b, 0xde, 0x77, 0xff, 0x90, 0xd0, 0x6, 0xad, 0xc9, 0x63, 0xfc, 0x46}}
	return a, nil
}

var _assetsRouterClusterRoleBindingYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x90\xb1\x4e\x03\x41\x0c\x44\xfb\xfd\x8a\x91\x52\xe7\x10\x1d\xda\x0e\xf8\x83\x20\xd1\x3b\x7b\x4e\xce\xe4\xce\x3e\xd9\xde\x14\x7c\x3d\x3a\x12\x2a\x10\xf5\xe8\x69\xe6\xcd\x0e\x2f\xa2\x63\x20\x27\x86\x5b\x4f\x76\xb8\xcd\x8c\x34\x48\x06\xde\xd8\xaf\xd2\x18\xcf\xad\x59\xd7\x1c\xca\x45\x74\xac\x78\x9d\x7b\x24\xfb\xc1\x66\xde\x70\xd1\x73\xa1\x55\xde\xd9\x43\x4c\x2b\xfc\x48\x6d\xa0\x9e\x93\xb9\x7c\x52\x8a\xe9\x70\x79\x8a\x41\xec\xe1\xfa\x58\x16\x4e\x1a\x29\xa9\x16\x60\x07\xa5\x85\x21\x81\xe0\x04\x25\xbc\x6b\xca\xc2\x25\xfa\xf1\x83\x5b\x46\x2d\x7b\xdc\x2a\xef\x4b\xee\x43\xfe\x63\x7f\xa2\x58\xa9\xfd\x95\x6f\x7e\x07\x3e\x6d\xfd\xbf\x6c\x0a\xbe\xd1\x0a\x5b\x59\x63\x92\x53\xee\x45\xcf\xce\x11\xfb\xdb\x3b\xe5\x2b\x00\x00\xff\xff\x01\x73\xe0\xaa\x31\x01\x00\x00")

func assetsRouterClusterRoleBindingYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsRouterClusterRoleBindingYaml,
		"assets/router/cluster-role-binding.yaml",
	)
}

func assetsRouterClusterRoleBindingYaml() (*asset, error) {
	bytes, err := assetsRouterClusterRoleBindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/router/cluster-role-binding.yaml", size: 305, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x1f, 0x67, 0xba, 0x2a, 0x20, 0x7a, 0x13, 0xa8, 0xfb, 0xbb, 0x73, 0xd0, 0x5e, 0xe6, 0x87, 0xda, 0xf5, 0x81, 0xed, 0x3c, 0x25, 0xb0, 0xcb, 0xef, 0xb4, 0xa0, 0xa9, 0x1c, 0x11, 0x5a, 0xfd, 0x8d}}
	return a, nil
}

var _assetsRouterClusterRoleYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x91\x31\x4f\xc3\x40\x0c\x85\xf7\xfb\x15\x56\x99\x93\x8a\x0d\x65\x65\x60\xaf\x10\xbb\x73\x71\x89\xe9\xf5\x7c\xb2\x7d\xa9\xc4\xaf\x47\x49\x5a\xa8\x40\x0c\x45\x6c\xf7\x74\x7a\xf6\xf7\x9e\xef\xe0\x31\x55\x73\x52\xb0\x28\x85\x06\x50\x49\x04\x7b\x51\x50\xa9\x4e\x6a\x2d\x3c\x8f\x6c\x60\xa3\xd4\x34\x40\x4f\x80\x06\x4a\xe6\xca\xd1\x79\x5a\x64\x11\x33\xee\x13\xb5\xe1\xc0\x79\xe8\x2e\x13\x77\x92\x28\x60\xe1\x17\x52\x63\xc9\x1d\x68\x8f\xb1\xc5\xea\xa3\x28\xbf\xa3\xb3\xe4\xf6\xf0\x60\x2d\xcb\x76\xba\x0f\x47\x72\x1c\xd0\xb1\x0b\x00\x19\x8f\xd4\x81\x14\xca\x36\xf2\xde\x1b\xce\xaf\x4a\x66\xcd\x8a\x14\xb4\x26\xb2\x2e\x34\x80\x85\x9f\x54\x6a\xb1\xd9\xd4\xc0\x66\x13\x00\xd0\x5d\xb9\xaf\x4e\xbb\x0b\xa4\x64\xeb\x20\xd7\x94\x02\xcc\xe4\x52\x35\xd2\xd9\x41\x79\x28\xc2\xd9\x6d\x51\xf3\x5a\x2b\x18\x69\x95\x46\x3a\xf1\x2a\x26\xd2\xfe\x6c\x49\x6c\xbe\x3c\x4e\xe8\x71\xfc\x09\x31\xe7\xa3\xec\x1c\xaf\x03\xde\xca\xe5\x72\xa0\xac\x34\x31\x9d\xce\x2c\xb5\x7f\xa3\xe8\x18\x23\x99\x7d\x7d\x5c\x71\x45\x25\x74\xfa\xa5\x94\x66\xbd\x66\xfb\x59\xe9\x1f\x98\x96\x09\x37\x96\xf1\xcf\xcb\xb7\xe6\xe8\xf5\x1b\x43\x2d\xc3\x1c\xfc\x23\x00\x00\xff\xff\x13\x49\xf9\x92\xc9\x02\x00\x00")

func assetsRouterClusterRoleYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsRouterClusterRoleYaml,
		"assets/router/cluster-role.yaml",
	)
}

func assetsRouterClusterRoleYaml() (*asset, error) {
	bytes, err := assetsRouterClusterRoleYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/router/cluster-role.yaml", size: 713, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x17, 0x6f, 0x7f, 0x61, 0x6d, 0xdb, 0x39, 0xe1, 0x3a, 0xf9, 0xf8, 0x54, 0x21, 0x25, 0xf8, 0x6f, 0x6c, 0xd2, 0x2, 0x5, 0x16, 0xf7, 0xdd, 0x81, 0x1f, 0xde, 0xf5, 0x6, 0x66, 0xa, 0x63, 0xd}}
	return a, nil
}

var _assetsRouterDeploymentYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x54\x41\x6f\xe2\x3c\x10\xbd\xf3\x2b\x46\xe5\xd0\x53\x4a\xfb\xb5\xfa\xb4\xeb\x1b\x82\x74\x85\x54\xda\x08\xd2\x5e\x91\x71\x86\x62\xe1\xd8\x96\x67\x92\x2e\xfb\xeb\x57\x86\xd0\x26\x2c\x45\xdd\xdb\xce\xd1\xf3\xde\x9b\x79\x1e\x7b\xfa\x30\x46\x6f\xdc\xb6\x44\xcb\xf0\xa6\x79\x0d\x05\xae\x64\x65\x18\x6a\x69\x2a\xa4\x5e\x1f\x46\xa6\x22\xc6\x00\x13\xfb\x1a\x90\x08\xc8\xa3\xd2\x2b\xad\x1a\x04\xc8\x80\x20\xbd\x37\x1a\x0b\x90\x0c\xa1\xb2\xac\x4b\xbc\xea\x6d\xb4\x2d\x44\x4b\xbe\x27\xbd\x7e\xc1\x40\xda\x59\x11\x09\x34\xa8\x6f\x7a\x25\xb2\x2c\x24\x4b\xd1\x03\xe8\x83\x95\x25\x82\x26\x20\xe4\x8e\xd4\x21\x47\x5e\xaa\x4f\x00\x46\x2e\xd1\x50\x94\x81\x28\x2e\x20\xb8\x8a\x31\xf4\x62\xb7\xf1\x94\xd0\xa0\x62\x17\xf6\x88\x52\xb2\x5a\x3f\xb4\x28\x5d\x12\x00\x63\xe9\x8d\x64\x6c\xe0\xad\x2e\x63\x98\x0e\xf3\x98\x0b\x70\x28\x1a\xa3\x0f\x84\xa1\xd6\x0a\x87\x4a\xb9\xca\xf2\xe3\xa7\x1e\x63\x58\x57\xe0\xbc\xd3\x6a\x8c\x25\xb2\xbc\xda\x54\x4b\x0c\x16\x19\xe9\x4a\xbb\x81\x23\x01\x46\xdb\xea\xe7\x3b\x28\x52\x93\xe0\x0c\x1e\x21\xdf\x5c\xd8\x60\x10\x70\x79\xd9\x40\x7d\xd0\x2e\x68\xde\x8e\x8c\x24\x8a\xed\x08\xa0\x2d\x31\x96\x89\xda\x8f\x3a\x51\x41\xb3\x56\xd2\x34\x04\xe5\x2c\x4b\x6d\x31\xb4\x2c\x27\xbb\x81\x74\x5c\x1f\xfc\xea\x52\xbe\x9e\xb1\x18\x63\x07\xc9\x2a\x63\x32\x67\xb4\xda\x0a\x98\xac\x1e\x1d\x67\x01\x29\x3e\x94\x0f\x9c\x77\x81\x5b\x45\x3f\xca\xae\x99\x7d\xeb\xb8\xd5\x63\xe6\x02\x0b\xf8\x76\xdd\xc9\xfa\xe0\xd8\x29\x67\x04\xe4\xa3\xec\x13\x39\x3a\xa7\x77\x77\x77\xfb\x57\x82\xc4\x92\xcf\x0a\xde\x7c\xbf\xfd\xff\x4b\x8a\x7d\x98\x62\x78\x3d\xfa\x5b\x1f\x69\xb4\xf5\xa9\xfb\x99\xe7\xc3\x7c\xbe\xc8\x9e\x66\x79\xa7\xc8\xee\xcb\x0a\xb8\x88\xd5\x2f\x4e\xd0\x66\x4f\xcf\x79\x3a\x5b\xcc\xd3\xd9\xcb\x64\x94\x2e\x1e\x87\xd3\x74\x9e\x0d\x47\xe9\x29\x11\xe7\xd1\xd2\x5a\xaf\x38\xd1\xfb\xcd\x70\x42\x6f\x9c\xde\x0f\x9f\x1f\xf2\xc5\x28\x9d\xe5\x93\xfb\xc9\x68\x98\xa7\x8b\xf1\x64\x76\x4a\x6e\x80\xac\x06\x7e\xa3\x07\x6c\x68\xe0\x83\xae\x25\x63\x0b\x67\x74\x8d\x16\x89\xb2\xe0\x96\x28\x3a\x02\xda\x6a\xd6\xd2\x8c\xd1\xc8\xed\x1c\x95\xb3\x05\x09\xb8\xe9\x3e\x80\x38\xe0\x1f\xc8\x5d\x22\x80\x97\xbc\x16\x30\x58\xa3\x34\xbc\xfe\x75\x9c\x3c\x35\xa8\x80\xb2\xd0\xff\x42\x23\xb5\x33\x55\x89\xd3\xb8\x52\x8e\x7e\x48\x19\xcf\xb2\xbd\xe0\xf9\x4b\x85\x66\x4c\xcd\xce\x4f\x14\x06\x8e\xab\xfd\x18\x15\x4d\x3f\x59\xb3\x15\xc0\xa1\x3a\xa4\xf6\x0d\xbc\xd7\x4e\xbe\xa0\x45\xa8\x42\xd7\x7b\x83\x9e\xba\x02\x05\xdc\xfd\x77\xdd\x79\xf8\xf3\x1d\xfc\xcf\x85\x99\xec\x3f\xc1\xef\x00\x00\x00\xff\xff\x1d\x49\xbe\x52\xbc\x06\x00\x00")

func assetsRouterDeploymentYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsRouterDeploymentYaml,
		"assets/router/deployment.yaml",
	)
}

func assetsRouterDeploymentYaml() (*asset, error) {
	bytes, err := assetsRouterDeploymentYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/router/deployment.yaml", size: 1724, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd6, 0xcd, 0x57, 0xa9, 0xea, 0xec, 0x7b, 0xbf, 0x49, 0x6c, 0x79, 0x95, 0x11, 0xb1, 0x71, 0x4d, 0x8c, 0xcc, 0x3f, 0x97, 0x5, 0xd6, 0xb, 0xbe, 0x5a, 0x62, 0x58, 0xa0, 0x10, 0xfc, 0x36, 0x41}}
	return a, nil
}

var _assetsRouterNamespaceYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\xce\xbd\x6e\xc3\x30\x0c\xc4\xf1\x5d\x4f\x71\x70\x67\xf7\x63\xd5\x43\x74\xec\xce\x44\xd7\x84\x88\x4c\x0a\x26\xad\xe7\x2f\x0c\x04\xed\xd0\xf9\x7e\xc0\xfd\x1f\x6a\xad\xe2\x53\x36\xc6\x90\x2b\x8b\x0c\xfd\xe2\x1e\xea\x56\x31\x3f\xca\xc6\x94\x26\x29\xb5\x00\x26\x1b\x2b\x7c\xd0\xe2\xae\xdf\xb9\xaa\xdd\x76\x46\x14\x40\xcc\x3c\x25\xd5\x2d\x4e\x88\x3f\xf4\xaa\xfe\x66\xde\xb8\x06\x3b\xaf\xe9\x7b\xc5\xb2\x14\xa0\xcb\x85\xfd\x89\x5f\x10\x4c\x4c\xe9\x07\x91\x0e\x99\xae\x0d\x8d\x83\xd6\xd4\x6e\x70\xc3\xe3\xb8\x10\xd2\x36\x8d\x33\x0c\x79\x97\x7c\x82\x38\xe7\xdf\x37\xc8\xd0\xf8\x1f\xb0\x1f\xb6\x76\x4e\xf6\x8a\xe5\x7d\x29\x3f\x01\x00\x00\xff\xff\x55\xe9\xbf\x8f\xf3\x00\x00\x00")

func assetsRouterNamespaceYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsRouterNamespaceYaml,
		"assets/router/namespace.yaml",
	)
}

func assetsRouterNamespaceYaml() (*asset, error) {
	bytes, err := assetsRouterNamespaceYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/router/namespace.yaml", size: 243, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xa7, 0xd4, 0xa1, 0x93, 0xf1, 0x1d, 0xdc, 0xa4, 0xc, 0x91, 0xe5, 0x54, 0xec, 0xa1, 0xd2, 0xb1, 0x34, 0x19, 0x90, 0xd7, 0x34, 0xdc, 0xce, 0xea, 0x67, 0xc7, 0x8a, 0x6c, 0xe6, 0x4b, 0xd1, 0x92}}
	return a, nil
}

var _assetsRouterServiceAccountYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\xce\x31\x4e\xc4\x30\x10\x85\xe1\xde\xa7\x78\x52\x6a\x22\xd1\xa6\xa3\xa4\xa1\x00\x89\x7e\x70\xde\xee\x8e\x88\x3d\x66\x66\xbc\x88\xdb\xa3\x20\x0e\xf0\xff\xfa\x16\x3c\xd5\x6a\xb3\x27\x2e\xe6\x70\x9b\x49\x0f\x54\xa7\x24\x77\x7c\xfc\x20\x6f\x84\x0d\xba\xa4\xf9\x8a\xe7\xc4\xb7\x1e\x07\x9c\x5f\x53\x9d\xa8\xc7\x8c\xa4\x23\xaa\x0d\xee\x65\xc1\xa0\x37\x8d\x50\xeb\x01\xe7\xf1\x77\x49\xc3\xeb\x39\xc6\x70\xab\x8c\xd0\x7e\x5d\xcb\xa7\xf6\x7d\xc3\x1b\xfd\xae\x95\xff\x86\x22\x43\xdf\xe9\x67\xbd\xe1\xfe\x58\x1a\x53\x76\x49\xd9\x0a\xb0\xe0\x45\x1a\xa1\x81\x60\x42\x12\x3e\x7b\x6a\xe3\x5a\x80\x2e\x8d\x31\xa4\x72\x3b\xa9\x3d\x6e\x7a\xc9\x07\xed\x57\x67\x44\xf9\x0d\x00\x00\xff\xff\xa5\xb8\xbc\x42\xe2\x00\x00\x00")

func assetsRouterServiceAccountYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsRouterServiceAccountYaml,
		"assets/router/service-account.yaml",
	)
}

func assetsRouterServiceAccountYaml() (*asset, error) {
	bytes, err := assetsRouterServiceAccountYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/router/service-account.yaml", size: 226, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x5, 0xfa, 0x29, 0xdd, 0x42, 0x7f, 0xe2, 0x84, 0x4a, 0x75, 0xc, 0x5f, 0xb7, 0x45, 0x41, 0x9e, 0x12, 0x68, 0xcc, 0x70, 0x87, 0xb4, 0xbe, 0xfb, 0xa4, 0x53, 0x51, 0xd2, 0x5c, 0xa7, 0x87, 0x72}}
	return a, nil
}

var _assetsRouterServiceCloudYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x90\x41\x6b\x14\x41\x10\x85\xef\xfd\x2b\x1e\xe4\x9c\xa0\x98\x83\xcc\x31\x39\x09\x41\x16\x5c\xbc\x57\x7a\x6a\x76\x9a\xf4\x54\x35\x55\x35\xab\xfb\xef\xa5\x7b\x37\xa0\x28\x1e\xe7\x31\xfd\xd5\x7b\xdf\x1d\x5e\x94\x66\x3c\x51\x25\xc9\x6c\xf8\xc6\x76\x2e\x99\x11\x8a\x56\x29\x33\x8a\x60\x31\x95\x80\x2e\x88\x95\x61\xba\x07\x5b\x8f\x73\xd5\x7d\x06\xcb\xb9\x98\xca\xc6\x12\xfe\x90\xee\xf0\x5c\x77\xef\x3f\x7c\x91\x93\xb1\x3b\xbc\x71\x2e\x4b\xc9\x38\x53\xdd\xd9\x41\xc6\xa0\xd6\x6a\xe1\x19\x14\xb0\x5d\xa2\x6c\xfc\x90\xde\x8a\xcc\xd3\xfb\xf9\x44\xad\x7c\x67\xf3\xa2\x32\xe1\xfc\x31\x6d\x1c\x34\x53\xd0\x94\x80\x3b\x7c\xa5\x8d\x51\x1c\xce\xf1\x07\x02\x10\xda\xd8\x1b\x65\x9e\xa0\x8d\xc5\xd7\xb2\xc4\x7d\xb9\x36\x49\x40\xa5\x57\xae\xde\x21\xe8\x1d\xa6\xdb\x98\xd4\x3b\xf6\x34\x2e\x8d\xa7\x21\xe4\xdd\x47\x02\x9c\x2b\xe7\x50\xfb\xfb\x59\xef\x72\x5c\x8b\x83\xaa\x2b\x56\xf2\x21\x88\x97\x85\xf3\xd0\xb5\x91\xbd\x15\x39\xe1\xe5\x09\x4d\xb5\x22\xc8\x4e\x1c\x0e\x72\xec\xb2\x32\xd5\x58\x2f\xf8\xb1\xb2\x40\x74\xc0\x6e\x6e\x9b\xce\x57\x4f\xcd\xd8\xb9\xab\x17\x10\x44\x67\xc6\x2b\xaf\x45\xe6\x71\xc7\xaf\xaa\xfa\x6c\xfe\x19\x6c\x42\xf5\x68\xb4\x2c\x25\x1f\xb4\x96\x7c\xe9\x43\x32\xd5\x04\x34\xb5\x18\xab\xef\x87\xa0\x09\x6b\x44\x1b\x6b\x9a\x69\x68\xd6\x3a\xe1\xf8\x7c\xb8\x26\x6a\x31\xe1\xf3\x87\xf1\x71\x2d\x7c\x18\xd1\xed\xcd\xef\x08\xff\x2f\xe3\xf1\xf1\xd3\x3f\x21\x9e\x7e\x05\x00\x00\xff\xff\x14\xac\xd6\xf5\x74\x02\x00\x00")

func assetsRouterServiceCloudYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsRouterServiceCloudYaml,
		"assets/router/service-cloud.yaml",
	)
}

func assetsRouterServiceCloudYaml() (*asset, error) {
	bytes, err := assetsRouterServiceCloudYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/router/service-cloud.yaml", size: 628, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x8c, 0x35, 0x3d, 0xd7, 0x8, 0xf9, 0xba, 0x0, 0x54, 0xbd, 0x2a, 0xeb, 0x98, 0x83, 0x6f, 0x28, 0x5e, 0xda, 0xd8, 0xa9, 0x45, 0x65, 0xe2, 0x35, 0x98, 0xd0, 0x6, 0x64, 0xc4, 0x82, 0x36, 0x14}}
	return a, nil
}

var _assetsRouterServiceInternalYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\xcf\x31\x4f\x23\x41\x0c\x05\xe0\x7e\x7e\xc5\x93\x52\xe7\x74\x51\xa2\xd3\x31\x1d\x4a\x95\x06\x45\x02\xd1\x9b\x59\x27\xb1\x98\x9d\x19\xd9\xde\x20\xfe\x3d\xda\x0d\x11\x0b\x34\x29\xd7\xfb\xfc\xf9\xcd\x02\xdb\x3c\x98\xb3\xe2\x91\xf5\x2c\x89\xf1\x26\x7e\x42\xc7\x07\x1a\xb2\xe3\x4c\x79\x60\x0b\x5f\xa9\x5d\x39\x2a\x9b\xc1\x1a\x27\x39\x48\x02\x95\x52\x9d\x5c\x6a\x31\x90\x32\xa8\xb5\x2c\xdc\x81\x1c\x3a\x14\x97\x9e\xff\x84\x57\x29\x5d\xbc\x1e\x08\xd4\xe4\x99\xd5\xa4\x96\x88\xf3\x2a\xf4\xec\xd4\x91\x53\x0c\xc0\x02\x0f\xd4\x33\xa8\x74\xb8\xff\xe1\x1a\xfb\x37\x13\x28\xd4\xb3\x35\x4a\x1c\x51\x1b\x17\x3b\xc9\xc1\x97\x72\xe9\x17\x80\x4c\x2f\x9c\x6d\x54\x31\x96\x8a\xd0\x3a\x38\x6b\x18\x9b\x8f\x53\x7f\x6f\x1c\xaf\xef\xda\xed\x03\x60\x9c\x39\x79\xd5\xdf\x3b\x40\xab\xea\x13\xb6\x9c\xee\x46\x9c\xdc\xdb\x94\x1b\xff\x44\xfc\xff\x7b\xf9\xd0\xea\x35\xd5\x1c\xf1\xb4\xdd\x4f\x13\x27\x3d\xb2\xef\xa7\xd0\xe7\xce\x9c\xb0\x99\xb1\xd9\xac\x6f\x44\x6c\xa6\xf4\xec\x2a\x69\xee\xac\xee\xd6\xff\x6e\x80\xa6\xd8\x47\x00\x00\x00\xff\xff\x30\x30\x02\xf6\x00\x02\x00\x00")

func assetsRouterServiceInternalYamlBytes() ([]byte, error) {
	return bindataRead(
		_assetsRouterServiceInternalYaml,
		"assets/router/service-internal.yaml",
	)
}

func assetsRouterServiceInternalYaml() (*asset, error) {
	bytes, err := assetsRouterServiceInternalYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/router/service-internal.yaml", size: 512, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x88, 0x3, 0x6d, 0x72, 0xb, 0x1e, 0xd8, 0x4, 0x82, 0xd2, 0xb1, 0x70, 0xa9, 0x3f, 0xf, 0x83, 0x3a, 0x2b, 0xeb, 0x18, 0x2b, 0x1e, 0xd2, 0xd6, 0xc3, 0xbe, 0x58, 0x72, 0xaa, 0xee, 0x2f, 0xe3}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"assets/defaults/cluster-ingress.yaml": assetsDefaultsClusterIngressYaml,

	"assets/router/cluster-role-binding.yaml": assetsRouterClusterRoleBindingYaml,

	"assets/router/cluster-role.yaml": assetsRouterClusterRoleYaml,

	"assets/router/deployment.yaml": assetsRouterDeploymentYaml,

	"assets/router/namespace.yaml": assetsRouterNamespaceYaml,

	"assets/router/service-account.yaml": assetsRouterServiceAccountYaml,

	"assets/router/service-cloud.yaml": assetsRouterServiceCloudYaml,

	"assets/router/service-internal.yaml": assetsRouterServiceInternalYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"assets": {nil, map[string]*bintree{
		"defaults": {nil, map[string]*bintree{
			"cluster-ingress.yaml": {assetsDefaultsClusterIngressYaml, map[string]*bintree{}},
		}},
		"router": {nil, map[string]*bintree{
			"cluster-role-binding.yaml": {assetsRouterClusterRoleBindingYaml, map[string]*bintree{}},
			"cluster-role.yaml":         {assetsRouterClusterRoleYaml, map[string]*bintree{}},
			"deployment.yaml":           {assetsRouterDeploymentYaml, map[string]*bintree{}},
			"namespace.yaml":            {assetsRouterNamespaceYaml, map[string]*bintree{}},
			"service-account.yaml":      {assetsRouterServiceAccountYaml, map[string]*bintree{}},
			"service-cloud.yaml":        {assetsRouterServiceCloudYaml, map[string]*bintree{}},
			"service-internal.yaml":     {assetsRouterServiceInternalYaml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
