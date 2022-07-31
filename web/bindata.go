// Code generated by go-bindata.
// sources:
// .build/assets/about.html
// .build/assets/close.svg
// .build/assets/edit.css
// .build/assets/edit.html
// .build/assets/edit.js
// .build/assets/index.js
// .build/assets/links.css
// .build/assets/links.html
// DO NOT EDIT!

package web

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
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
	bytes []byte
	info  os.FileInfo
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

var _aboutHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x93\x4f\x6f\xd4\x30\x10\xc5\xef\xfd\x14\x43\x2e\xbd\x74\xd7\x5b\xc1\x01\x2d\xce\x22\x54\xfe\x08\xa9\x12\xa8\x2a\x07\x8e\xb3\xf6\x24\xb6\xea\x78\x82\x3d\xe9\x36\xdf\x1e\xc5\x6e\x97\xd2\x45\xec\x65\x63\xbf\xf8\xbd\xe7\x9f\x1d\xfd\xea\xe3\xb7\xab\xdb\x9f\xdf\x3f\x81\x93\x21\xec\xce\xf4\xd3\x1f\xa1\xdd\x9d\x01\x00\x68\xf1\x12\x68\xf7\x85\x61\xbb\x85\x0f\x7b\x9e\x44\xab\x3a\x55\xe5\x81\x04\xc1\x89\x8c\x2b\xfa\x35\xf9\xfb\xb6\xb9\xe2\x28\x14\x65\x75\x3b\x8f\xd4\x80\xa9\xa3\xb6\x11\x7a\x10\xb5\xb8\xbf\x03\xe3\x30\x65\x92\x76\x92\x6e\xf5\xb6\x51\x8f\x46\xc1\xc7\x3b\x70\x89\xba\xb6\x51\x59\xf9\x68\xe9\x61\x6d\x72\x6e\x8a\xba\xfc\x12\x85\xb6\xc9\x32\x07\xca\x8e\x48\x9a\xd3\x75\x4b\x8d\xbc\x55\xaa\xe3\x28\x79\xdd\x33\xf7\x81\x70\xf4\x79\x6d\x78\x50\x26\xe7\xf7\x1d\x0e\x3e\xcc\xed\x0d\x06\x3a\xe0\xbc\x7d\xb3\xd9\x5c\xbc\xde\x6c\xfe\x17\xa1\x55\x45\xa1\xf7\x6c\xe7\xc7\x44\xeb\xef\xc1\x04\xcc\xb9\x6d\x4a\xcd\xc7\x26\x45\x73\x97\x0b\xaa\x6b\x1f\xef\xb2\x56\xee\xf2\x99\x32\xee\x6e\x19\x4c\x22\x14\x02\x84\xd2\x5a\xe6\x91\xa0\xe9\x59\x35\x80\xd1\x82\x38\x8a\x30\xf3\x94\x60\x4c\xd4\x51\x22\x5b\x5f\xcb\x8e\x93\x98\x49\x20\xe2\x40\x6b\xf8\xcc\x09\xe8\x01\x87\x31\xd0\x05\xf4\xac\x3a\x1f\x31\x1a\x5a\x59\xcc\x6e\xcf\x98\xec\x1a\xbe\x76\x8b\x59\x5d\xed\x33\x60\x48\x84\x76\x06\xc1\x3b\x8a\xc7\x98\xf3\x10\x60\x4f\x80\x93\xf0\x80\xe2\x0d\x86\x30\x43\x22\xeb\x13\x19\x21\x0b\xc2\xc5\xe4\x40\xfb\xec\x85\x8a\xa9\x97\xf3\x0c\x91\xe5\x85\x13\x18\x8c\x30\x62\x16\x02\x5f\x66\xe1\xc7\xcd\x35\x88\x43\x29\x83\xd2\xbf\x74\x39\xf8\x10\x8e\x11\x20\xbc\xd6\x6a\x7c\x89\x28\x53\xe5\x93\x05\xb8\x03\x0c\xe1\xb8\x95\x0c\x3d\x2f\xad\x34\xbe\x38\xef\x9e\x55\xd1\x9b\xdd\xd3\x93\x56\xb8\xfb\x87\xb7\x71\x18\x7b\xaa\x86\x6c\x50\x3c\xc7\x12\xf2\xac\x62\x8d\xe8\x59\x91\xf5\xa2\x8e\xf3\xa7\xd8\x8b\x7e\xc2\xfe\x02\x0e\xce\x1b\x57\x77\x1a\x48\x2a\x9e\x3f\xb1\xa7\x60\xea\x29\x3d\x41\xc9\x7f\x51\xd1\xca\xfa\xfb\xe5\x12\xd6\xdb\xa7\x55\xfd\x3c\xcf\x7e\x07\x00\x00\xff\xff\x48\x3b\xa9\xc0\xb7\x03\x00\x00"

func aboutHtmlBytes() ([]byte, error) {
	return bindataRead(
		_aboutHtml,
		"about.html",
	)
}

func aboutHtml() (*asset, error) {
	bytes, err := aboutHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "about.html", size: 951, mode: os.FileMode(420), modTime: time.Unix(1656918125, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _closeSvg = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2c\x8e\xc1\x8a\x84\x30\x0c\x86\x5f\x25\xe4\xbe\x69\x9b\x62\xcd\x2e\xd6\xc3\x9e\xdd\x87\x58\x50\x5a\xc1\x19\x65\x2c\x56\xe6\xe9\x87\x56\x2f\x3f\xe1\xfb\xbf\x84\x74\xfb\x11\xe0\x7c\x2c\xcf\xdd\x63\x4c\x69\xfb\x51\x2a\xe7\x4c\xd9\xd2\xfa\x0a\x8a\xb5\xd6\x6a\x3f\x02\x42\x9e\xc7\x14\x3d\x5a\x87\x10\xa7\x39\xc4\x74\xcd\xc7\x3c\xe5\xdf\xf5\xf4\xa8\x41\x83\x75\x60\x1d\xf6\xdd\xf6\x9f\x22\x8c\x1e\xff\x58\xa8\x81\x6f\x72\x3c\xb0\x23\x2b\xd0\x52\x03\x46\xc0\x34\x24\x52\x79\x25\xed\x2d\xdd\xdc\x5c\x5e\xdd\x58\x98\x0c\x43\x89\xc1\x08\xb0\x26\xc3\x8b\x94\x4b\x35\x0a\xff\xaa\x65\x6d\xc0\xc8\x1b\x55\xdf\x95\x87\xfb\x4f\x00\x00\x00\xff\xff\xe5\x79\xbd\x91\xd8\x00\x00\x00"

func closeSvgBytes() ([]byte, error) {
	return bindataRead(
		_closeSvg,
		"close.svg",
	)
}

func closeSvg() (*asset, error) {
	bytes, err := closeSvgBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "close.svg", size: 216, mode: os.FileMode(420), modTime: time.Unix(1656918125, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _editCss = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x93\xdf\x8e\xda\x3c\x10\xc5\xef\xbf\xa7\xb0\xe0\x06\x24\x9c\x35\xec\xc2\x7e\xd8\x52\x1f\xa2\x77\xbd\x74\xec\x49\x62\xad\xe3\xb1\x6c\x07\xc2\x46\xbc\x7b\x95\x3f\x40\x5a\xb6\xaa\xaa\x08\x14\xec\x1f\x9e\x99\x73\x8e\x73\xd4\x97\x2e\x97\xea\xa3\x0c\xd8\x38\xcd\x97\x45\x51\x88\x02\x5d\xa2\x85\xac\x8d\xbd\xf0\xc5\x77\x69\xe1\x2c\x2f\x8b\x4d\x94\x2e\xd2\x08\xc1\x4c\x40\x34\x9f\xc0\xdf\x76\xbe\x1d\x7f\x9e\xc1\x94\x55\xe2\xaf\x8c\x5d\x0b\x0c\x75\x97\xa0\x4d\x54\x5a\x53\x3a\xae\xc0\x25\x08\xd7\x65\x2e\x43\x77\x36\x3a\x55\xfc\xb0\xef\xff\x57\xcb\x50\x1a\xc7\x19\x91\x4d\x42\xe1\x31\x9a\x64\xd0\xf1\x00\x56\x26\x73\x82\xeb\xb2\x09\xb6\xfb\x87\x66\x5e\xbf\x68\x46\x4c\x05\x19\xf3\xad\xf0\x52\x6b\xe3\x4a\xbe\xdb\xfb\x56\x28\xb4\x18\xf8\xf2\x78\x3c\x8a\x1c\x83\x86\x40\x83\xd4\xa6\x89\xfc\xcd\xb7\xd3\x0a\xdf\xfa\x96\x44\xb4\x46\x93\xa5\x52\x4a\x60\x93\xac\x71\xc0\x1d\x3a\x10\x39\xb6\x34\x56\x52\xe3\x99\x33\xb2\xf3\x2d\x39\xf8\x96\x84\x32\x97\x2b\xb6\xe9\x9f\x6c\xb7\x1e\x26\xe0\x05\xaa\x26\x76\xcf\x27\xb2\x63\x31\x02\x9c\x9e\x21\xff\x30\x89\x1a\xe7\x9b\x44\xbd\x95\x0a\x2a\xb4\x1a\x42\x37\x35\xa9\xb5\xbe\xa1\x35\x7e\xfe\x91\x50\xb5\xef\x7e\x99\x71\x1c\x7e\x7f\x60\x4f\x6a\xcf\xec\xb1\x50\xa4\x99\x8a\xbb\xad\x6f\x45\x0a\xd2\x4d\x76\x0c\xaf\xbd\xa5\x64\xc7\x58\x1d\x09\xc8\x08\xd4\x38\x8a\x4d\x12\xf7\x3d\x1e\x95\xb4\xf0\x63\xc5\xd6\x8f\x35\x8a\xc1\xf4\x15\x13\x7a\x32\x46\xe0\x26\x74\x8e\x29\x61\x4d\xfb\xc2\xcf\xa2\xdf\x76\x43\x6f\xe1\x5f\x3d\x19\x86\xce\xac\x71\x1f\xb3\x18\xd3\x49\x94\xe2\xd0\x3f\x23\x52\x34\xea\x4b\x04\x60\x96\x84\x01\xfd\x26\x6f\xa2\xb2\x63\x31\x0a\xa5\x41\x61\x90\x83\x1c\xbd\xf7\x23\x96\x55\x2e\xcd\xe4\x1f\xc9\x29\x11\x7d\x93\xfd\x87\x91\xfe\x46\xcd\xf8\xc2\xa2\x4c\x7c\x98\xed\xba\x54\x36\x76\xf7\xd8\xcb\x3c\xa2\x6d\x12\x88\x84\x9e\x33\x31\x20\x9c\x89\x51\x0d\x7e\x0f\xf2\xa0\xc4\x63\x0c\x53\xcb\x12\x78\x13\xec\x6a\xf1\x12\x5f\x94\xc5\x08\x59\x3c\x95\x8b\xf5\x1c\xba\xd7\x18\x6d\xb8\xbb\xf1\x20\x02\x78\x90\x89\x3b\x9c\xde\xe6\x7b\xe3\x45\xff\xdf\xb7\xa4\xff\x12\xaa\x09\x11\x03\xf7\x68\x86\x43\x66\x49\x41\x2f\x95\x49\x97\x29\x27\xaa\xc9\x8d\xa2\x39\x7c\x1a\x08\x2b\x96\xbd\xbf\x6f\x08\xdb\x10\x96\x6d\xdf\xf7\x1b\xb2\x5d\x8b\x89\xe6\x4c\x4c\x47\x51\x38\x81\x4b\xf1\xa6\xb0\x8d\xd9\xc9\xc4\xee\x86\x65\xaf\xbf\x73\xc6\x55\x10\xcc\x28\x23\xaf\xf0\x04\xe1\x01\x1f\xae\xff\xfd\x0c\x00\x00\xff\xff\xd0\xd6\xde\x37\xdb\x04\x00\x00"

func editCssBytes() ([]byte, error) {
	return bindataRead(
		_editCss,
		"edit.css",
	)
}

func editCss() (*asset, error) {
	bytes, err := editCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "edit.css", size: 1243, mode: os.FileMode(420), modTime: time.Unix(1659240567, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _editHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x52\xb1\x72\xe3\x20\x10\xed\xfd\x15\x7b\xd4\x67\xa3\x99\xbb\xe2\xc6\x07\x4a\xe1\x78\x52\x26\x93\x71\x93\x12\xa3\x95\x21\x41\x82\xc0\xca\x89\xfe\x3e\x83\x50\x6c\x65\xd2\x85\x86\xdd\x7d\xec\x7b\x8f\x05\xf1\xeb\xf6\x7e\x77\x78\x7a\xd8\x83\xa1\xce\xd5\x2b\x51\x36\x00\x61\x50\x35\x39\x00\x10\x64\xc9\x61\x7d\xe7\x05\x2f\x51\xa9\x76\x48\x0a\x0c\x51\x58\xe3\xeb\x60\xcf\x92\xed\x7c\x4f\xd8\xd3\xfa\x30\x06\x64\xa0\x4b\x26\x19\xe1\x3b\xf1\x4c\xfb\x1f\xb4\x51\x31\x21\xc9\x81\xda\xf5\x3f\xc6\x67\x22\x67\xfb\x17\x30\x11\x5b\xc9\x78\xe2\xd8\x58\xda\xe8\x94\xd8\x04\xe6\x15\xd1\x49\x96\x68\x74\x98\x0c\x22\x5d\x01\x1a\x03\xce\xfc\xb9\xe1\x3b\x5d\x76\x97\xb6\x9c\xb7\xbe\xa7\xb4\x39\x79\x7f\x72\xa8\x82\x4d\x1b\xed\xbb\xdc\x72\xd3\xaa\xce\xba\x51\x3e\x2a\x87\x6f\x6a\xdc\xfe\xad\xaa\xdf\x7f\xaa\xea\x67\xd2\x82\x7f\x8e\x4c\x1c\x7d\x33\xce\x6e\x5a\x1f\x3b\x50\x03\x79\xed\xbb\xe0\x90\x50\x32\xdf\xb6\xb3\x57\x00\xd1\xd8\x33\xd8\x46\xb2\xa3\x8a\x97\xe2\xa2\xac\x5d\x62\xb5\xe0\x8d\x3d\x2f\x40\xdb\x87\x81\x16\x16\xd8\x74\x74\x88\x8e\x41\x70\x4a\xa3\xf1\xae\xc1\x28\xd9\xbe\x27\x8c\x40\x06\x61\x88\x0e\xc8\x43\x32\x3e\x12\xf6\x99\x71\xe2\xb8\xb8\x58\x0a\x5c\xb5\xbb\xf0\x45\x5b\xf0\x7c\x99\x7a\x55\x92\xa4\xa3\x0d\x04\x29\xea\xeb\xb3\x3d\x4f\x66\x0b\x52\x46\x52\x26\x21\x78\xf9\x56\x1f\x01\x00\x00\xff\xff\x1e\x7c\x18\x18\x6e\x02\x00\x00"

func editHtmlBytes() ([]byte, error) {
	return bindataRead(
		_editHtml,
		"edit.html",
	)
}

func editHtml() (*asset, error) {
	bytes, err := editHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "edit.html", size: 622, mode: os.FileMode(420), modTime: time.Unix(1656918125, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _editJs = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x56\x6f\x6f\xdb\xbc\x11\x7f\x9f\x4f\x91\x70\x40\x40\x22\x0c\x9d\xf4\x79\xf3\x20\x9a\x50\xa4\xa9\xb7\x62\x68\x97\x22\xc9\x8b\x0d\x5d\x31\xd0\xe4\x49\xe2\x4c\x91\x34\x49\x25\xf6\x1c\x7d\xf7\x81\x92\x6c\xcb\x89\xdb\x74\xcf\x0b\xdb\xbc\xe3\xf1\x77\x7f\xf8\xbb\xa3\x1f\xb9\x3f\x96\xb6\xce\x70\xd1\x18\x11\x95\x35\xb8\x22\xeb\x8a\x2d\xf2\xad\x5c\x90\xb5\x87\xd8\x78\x73\x2c\xad\x68\x6a\x30\x91\x2d\x1a\xf0\xab\x7b\xd0\x20\xa2\xf5\xb8\x20\x6d\x56\xb1\x05\xff\xf5\x23\xd7\x5a\x0f\xa7\xc4\x4f\x0f\x09\x0f\x3c\xc2\x54\x43\x92\x36\x27\x42\x18\x9d\xa1\x33\xca\xc9\xba\x60\x21\xae\x34\xb0\x00\xf1\xab\xb7\x0e\x7c\x5c\xe1\x19\xe5\x14\x21\xd2\xb6\x04\x4b\x5b\x3f\x3f\xe7\xeb\x96\x64\x29\xdb\x65\xe5\xb3\xa3\xbd\x74\x93\xb6\xd8\x81\x92\xf5\x66\x79\x3c\xc3\xbc\xdf\xae\xf3\x58\xa9\x90\xa5\x2f\xb6\xac\x7c\xce\xfb\xa5\xb4\x06\xfe\x62\x42\xfe\xed\x7b\x2f\x83\xf7\xd6\x0f\x0a\xce\xac\xd1\x96\xcb\x31\x70\x82\x32\x39\x67\x1e\x82\xb3\x26\xc0\x03\x2c\x23\x5d\xe4\x9c\x85\xc8\x63\x13\xb2\x7a\x83\xc8\x0a\xeb\xa7\x5c\x54\xbb\x38\x1d\x59\x3b\x6c\xe8\x82\xb4\xa4\xed\xb0\x3b\x5f\x63\xf0\x7a\xeb\xfe\xf5\x69\xb3\x2d\xae\xc1\x09\xa1\x9d\x31\xe7\x6d\xb4\x71\xe5\x80\x59\xf3\xd1\x1a\xd8\x41\x71\xb2\x1e\x67\xc7\x5c\x13\x2a\xcc\x49\x36\x20\xa4\xbd\x36\xdb\x07\x98\xee\x07\xb3\x41\xd8\x06\xf4\x26\xc4\x93\x8a\xd5\x27\xe0\x12\xc6\x28\xb4\x1e\x70\x96\x95\x4f\x97\x7b\x07\x8b\x06\x42\xec\xed\xba\xed\x1f\x03\x06\x30\xf2\x6f\xf7\xb7\x7f\x7f\x1d\xd4\xce\x15\x46\x37\xd6\x44\x30\xf1\xfc\x61\xe5\x00\xd1\x23\xc4\x9d\xd3\x4a\xf0\x64\x3f\xf9\x4f\xb0\x26\x13\x15\xf7\x01\x62\xde\xc4\xe2\x77\x44\xb2\x51\x34\x46\xe2\x84\xcf\x42\xf4\xca\x94\xaa\x58\x61\x4e\xde\x88\xe7\x75\x2c\x5b\xa8\x97\xb5\x19\x84\x59\x8b\x49\x56\xb1\x3b\x58\xe4\x45\x22\x7f\xd7\x10\x3b\x98\x8e\xfd\x3d\x3f\x0d\x3c\x1d\xff\xe3\xcb\xe7\x4f\x31\xba\xa1\x4c\x59\xcd\xac\x83\xce\x88\x9e\x5c\x6c\xf1\x93\x61\x81\xeb\xae\x99\x4a\x88\x23\xb0\x2d\x47\x36\x8e\x30\xfa\xeb\xf4\x01\xd1\x59\x67\xeb\x6c\x78\xc3\xf8\xeb\xed\x7d\x6f\xdd\x12\xbc\xac\xfc\xa8\xe7\x4a\xfb\x56\xcb\x25\x85\xc8\x5d\x8a\xe8\x83\x6d\x8c\x54\xa6\xbc\xd1\x0a\x4c\xbc\x03\x11\x31\xc9\xa4\xad\x53\xeb\x63\x47\x51\xcd\x7d\xa9\xcc\x79\xb4\x0e\xd1\x27\x65\xa4\x7d\x62\xca\x18\xf0\x9f\x40\x95\x55\x9c\xfc\x76\x2e\x58\xd5\x2f\xdf\x9d\x21\xb7\x44\xa4\xa5\xb3\xd7\xae\xb0\x66\x8f\x5c\x37\xf0\xfc\x8c\x10\x61\xd1\xab\x1a\x93\x4c\x9c\xe4\xf1\xf4\x14\xc7\x5c\x50\x83\x09\x15\xef\x3d\x13\x9a\x87\xf0\x59\x85\xc8\xb8\x94\x18\x3d\xaa\x80\xc8\xd5\x58\xed\xa1\xb6\x8f\x30\xec\x90\x96\x8e\xe6\xa0\x20\x6b\xc1\x9c\x87\x47\x30\xf1\x23\x14\xbc\xd1\x29\x15\x91\x6b\xdb\x53\x8c\x39\x1e\x2b\xc3\x6b\x60\xa1\x99\xf5\x3c\xc2\x97\x84\x05\xa7\x55\xc4\x68\x82\xc8\xb7\xcb\xef\x7d\xfd\x0e\x87\x9b\xd8\x93\xae\x05\xa3\x09\x77\x6a\xd2\x78\x3d\x41\x67\x82\x6c\xb9\x8f\xd7\x8d\xd7\x57\x65\x4b\x86\x16\xdf\x5d\x81\xa4\x40\xd6\x32\xef\x08\xec\x12\xc3\xb1\x24\x99\x2a\xb0\x64\x76\x4e\x54\x81\x21\x97\xcc\xdb\x26\x42\x5f\xae\x26\x07\x96\xe2\x4c\xde\x33\x99\x03\x0b\xb6\xf1\x02\xfe\x5d\xd9\x10\x3b\x1d\xb0\xc6\xeb\xd3\x53\x5c\xa9\x10\xad\x5f\x31\x0f\x4e\x73\x01\xf7\x31\x11\x63\xdd\xd2\x23\xd3\x68\x4d\xd1\x04\xa4\x8a\x13\x74\xd6\x10\x0a\x39\x4a\x8b\xf4\x8b\x4e\x72\xf9\x5e\x9e\xc1\xd5\xb6\x2e\xd6\xab\x52\x99\x33\xa0\x73\x16\x61\x19\x87\x2e\xcd\x11\xa2\xf3\x03\x95\x2f\x1a\x31\x47\x64\x6f\xab\xbb\x2b\xad\x4c\xd2\xcb\xbc\x23\x0f\x46\x3c\x09\x69\x8e\x5c\xc7\xe8\xd5\xac\x49\x9c\xad\x3c\x14\x88\x42\xda\x18\x7b\x4a\x9e\xb9\x73\x60\xe4\x4d\xa5\xb4\xc4\x32\x05\x3c\xc0\x04\xc7\x0d\x22\x14\x5e\xba\xab\x4c\xec\xd4\x63\x9c\x8b\x3f\xe7\x86\x3f\xaa\x92\x47\xeb\x59\x13\xc0\x5f\x97\xe9\x61\x53\x46\xc2\xf2\xb6\xc0\xe8\x0b\x17\xca\x44\x1b\x2a\x44\xde\xa3\x7f\x35\xef\x7e\xbb\xfc\xfd\xfc\x06\x5d\xa1\x9b\xe8\xf5\xf9\x0d\x7a\x11\x45\x0a\x73\xe8\x83\x39\x45\xd1\x73\x13\x0a\xeb\x6b\x44\x51\x10\x5c\xc3\x3f\xf1\x25\x41\x84\x96\x10\xfb\x87\xb6\xa3\x7b\xca\xf7\x03\x0f\x70\x6d\xe4\x74\x99\x62\xc2\x92\x5e\x50\x49\x2f\x09\x69\x41\x07\x48\x0f\x42\xd6\x2d\x16\x58\xf6\x03\x3b\x3d\x10\xb4\x7e\xdd\x32\xff\x17\x6d\x77\x84\x3d\x1a\x31\x76\x50\xe6\x08\x65\x33\xdc\x31\x4e\x90\xc4\xe2\xcd\x08\xf9\x38\xfd\x3c\x7d\x98\x22\x7a\x80\xd0\xf8\x35\x8b\x4b\x2a\xc9\xba\x1c\xb3\xb8\x24\x59\xc9\xec\xfc\xf9\x79\x81\xcb\x51\x32\x66\x9c\xcc\xcf\x4b\x78\x41\xd2\xc0\x58\xec\x35\xf1\x4b\x16\x66\x87\x58\xd8\xb3\x2d\x7b\xc5\xc2\x9e\x9d\x43\x1b\xef\x51\x28\x2b\xf7\x71\xa7\x77\x77\xb7\x77\x57\xc7\xe8\x4c\x64\xfb\xf7\x5e\xee\xe6\xdf\x8f\xef\xbd\xa5\xae\x83\x5f\x60\xd4\x6d\x12\x3a\xdf\xc8\x7f\x12\xb5\x43\x84\xfa\x9d\xac\x03\x22\x54\x6f\xe5\xc6\x6b\x44\x68\x1c\xfd\x0b\x24\xeb\x02\x93\x6c\x18\xaf\x5c\xca\x69\x9a\x61\x29\x29\x30\xe9\xdd\xf4\x10\xd4\x7f\x01\xd1\x82\x9e\x5c\x92\xcc\x1d\xb0\x08\xcd\xac\x56\x11\xd1\x23\xde\x99\xe8\x03\x26\x73\x58\x35\x0e\xd1\xd9\x0f\x0d\x1c\x0f\x11\x7e\x66\x20\x2a\x6e\xca\xad\x85\x3f\x64\xa1\x95\x98\x23\x5a\x77\x06\x7f\x80\xc6\x5b\x86\x96\x70\x70\xcc\xfe\x2a\x2b\xdf\x5d\x5c\xe4\xb9\x3c\x3d\xdd\xb4\x45\x5e\xf6\xf3\x35\xcd\xcd\x34\x41\xa9\x66\x85\x15\x4d\xc0\x84\xce\x30\x21\xed\xd0\x93\x5b\x6d\x4b\xba\x4f\x69\x87\x57\xf5\xe8\x7f\x01\x00\x00\xff\xff\xf0\x90\xb4\x5b\xb9\x0b\x00\x00"

func editJsBytes() ([]byte, error) {
	return bindataRead(
		_editJs,
		"edit.js",
	)
}

func editJs() (*asset, error) {
	bytes, err := editJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "edit.js", size: 3001, mode: os.FileMode(420), modTime: time.Unix(1659240567, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _indexJs = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x56\xdb\x6e\x1b\x37\x10\x7d\xd7\x57\x4c\x50\x23\xe4\x22\x6b\x4a\x4e\xdf\x2c\xa8\x45\x22\xbb\xe8\x25\xad\x0b\xdb\x7d\x28\x8a\xa2\x60\x76\x67\x25\x46\x5c\x72\x41\x72\x2d\xab\x8e\xfe\xa0\x1f\xd0\xef\xeb\x97\x14\xbc\xec\x45\x8a\xed\x16\xe8\x4b\xfd\x60\x2d\x87\xc3\x99\x73\x38\x37\xd2\xaa\x55\x85\x13\x5a\xd1\x0c\x1e\x26\x93\x3b\x6e\xe0\xa4\x32\x35\x2c\xe0\x84\x92\x4a\x9b\x9a\x64\xf9\x04\x00\xe0\xa4\xa8\x9b\x28\xfd\xac\xa8\x9b\x41\x2a\x6d\x27\x95\xb6\x97\xb6\x46\x26\x69\x6b\x64\x27\x2d\x74\xb3\xfb\x0e\x77\xb0\x00\xc5\xef\xc4\x8a\x3b\x6d\x58\x6b\xd1\xbc\x59\xa1\x72\x4c\xa8\x12\xef\xaf\x2a\x4a\xbe\xe7\x85\x50\x4e\xdb\x35\xc9\xe0\x8b\x05\xcc\xe0\x4b\x20\x7f\xfd\xf1\xe7\xe9\x92\xc0\x39\x90\xa5\x33\xf2\x74\x49\xa2\x45\xc9\xad\xfb\xc9\xc8\x79\xc2\x8d\xb0\x80\x9e\x8e\xe2\x35\x7a\x4a\x00\x06\x5d\x6b\x14\x9c\xd0\x52\x17\x6d\xed\x7d\x15\x06\xb9\xc3\x4b\x89\x7e\x15\x35\xb3\xf9\x64\x9f\xec\x18\xb4\xe2\xf7\x03\x5b\xd1\x4e\xdc\x2b\x9c\x67\x56\x99\x9a\xad\xd0\xd1\x59\xe6\x7f\xde\xea\x56\x95\x42\xad\x96\x52\xa0\x72\xd7\x58\x38\x9a\xcd\x27\x10\xd5\x0a\x6b\x29\xa9\xb9\x59\x09\x75\xea\x74\x43\x72\xd8\x0a\x55\xea\x2d\x13\x4a\xa1\xf9\x1a\xc5\x6a\xed\xa6\x9f\xc3\x69\xb0\xcd\xd6\x71\xfd\x7a\x84\xc7\xe3\xfb\xca\xe8\x7a\x8c\xa8\x35\x62\x00\xd5\x70\xe3\x7c\x14\x5a\x23\x98\x6d\xdf\x5b\x67\x84\x5a\xd1\xb3\x8c\xd9\x46\x0a\x47\xc9\x94\x04\x34\xe9\x1e\x82\xf6\x2f\x67\xbf\x0e\x0e\xa4\xe6\xe5\xe3\x74\xbd\xeb\x10\xaf\x88\x80\x4a\x5d\x70\xaf\xc2\x1a\xee\xd6\xe1\xde\x02\x4d\xc6\x3f\xf0\x7b\xfa\x10\x62\xd2\x1a\x79\x0e\x64\xca\x1b\x31\x6d\x8d\x9c\x12\x78\x15\x8e\xc7\x80\x95\xdc\xf1\xdb\x5d\x83\xe7\x40\x3e\x58\xad\xc8\x04\x60\x9f\x31\x2e\xb7\x7c\x67\x87\x4c\xf4\x5a\x11\x03\x80\xa8\x80\xbe\xf0\x02\xa6\x37\x9d\x0c\xc0\xae\xf5\xf6\xd2\x18\x6d\x82\x2e\x43\xff\x19\xa0\xf8\xbf\xc8\x33\xae\xf6\x93\xf0\x13\x42\xa7\x5b\xe7\xc9\x84\x13\x61\x91\xa7\x03\x01\x34\x2c\x92\xc6\xcb\x97\xf1\x83\x79\xe1\xc7\x8f\x40\xc8\xbc\x4f\x6a\x76\xc7\x25\x6d\x8d\xcc\x58\xa5\x8b\xd6\xd2\xe4\xb3\x35\xf2\x42\x94\xcb\x35\x57\x2b\x8c\xb2\xbd\x8f\x5f\xbc\x5d\x8f\xf5\x9d\x50\x9b\xe3\xe4\xcc\xc1\x9a\x62\xb8\x6a\x19\x34\x48\x7f\x61\x73\x8f\xdc\xb3\xb7\xa6\x80\x17\x0b\x20\xa4\xa3\x1f\x35\xbd\xf8\x95\xff\x0e\xee\x00\xa5\xc5\x83\xfd\x3e\x54\xda\x88\x95\x50\x23\x5d\x6f\xd8\x57\x33\x73\x78\xef\x28\x21\x59\x38\xc5\x0c\xd6\xfa\x0e\x97\x92\xfb\x7c\xad\xda\x62\xd3\x6d\xf0\xb2\x4c\x52\x29\xd4\xc6\xe7\x52\x82\xcc\x43\x8d\x3f\x51\x55\x84\x93\xac\x33\xe0\x9c\xa1\x64\x6d\xb0\x22\xb9\x47\x91\xc4\xc1\xfd\xb0\xe4\x4d\x83\xaa\xbc\xd5\xd4\x63\x1b\x9c\xac\x9f\x75\x62\x1b\xae\x06\x3f\x3d\xd0\xb5\x72\x64\xec\x25\xb5\x9e\xa7\x3d\x85\xfb\x08\x95\xea\x0c\x57\x36\x34\xbe\x1c\x88\x2d\xb8\xc4\x9f\xe9\x59\x96\x58\xaf\xd0\xdd\xa0\xc4\x54\x26\xcc\xa2\x7b\xcb\x2d\xbe\x51\xe5\xe5\xbd\xf3\x78\x78\x6a\x0a\x39\xcc\x72\x18\x16\x67\xa3\x6a\xee\x53\x77\x9c\x0f\x35\x5a\xcb\x57\xa9\x5f\xfd\x53\x6c\x62\x14\x8e\x29\xc7\x88\x45\x32\xd8\x5d\x4c\x32\x73\x79\x7d\x7d\x75\x7d\x0e\x3e\xb5\x3a\x4f\xff\xe9\x26\x3a\x2e\x6b\x51\xe2\x85\xe1\x5b\x34\x9f\xb6\x8f\xe7\xed\xcc\x0e\xec\x8c\xeb\xe7\xf1\x46\x94\x06\x4a\x57\x82\x19\x73\x46\xd4\xb1\xd4\x7c\x8d\x84\xed\x45\x37\x0f\xba\x42\x19\x1a\x41\xc8\xf9\xb4\x1b\x1a\xa5\x0c\x5c\x07\xfc\x07\xa6\xba\xf3\x7e\xb8\xb1\x8a\x97\xf8\x8d\xa2\xaf\x67\xb3\xec\x93\x4a\xeb\x15\xae\x5a\x37\x68\x04\x5a\xa1\xf5\x6b\x45\x89\x6d\xdf\xd7\xc2\x91\x7c\x60\x95\xc2\x8c\xac\x31\x78\x87\xca\x5d\x60\xc5\x5b\x99\x46\xc6\xbf\xeb\xba\x5d\xe7\x7a\xf2\x5a\x8e\xbb\xb2\x8b\x4d\xf7\xc7\xab\x9b\xdb\x34\x3b\xfd\xd1\xe7\x1b\x35\x9c\xc3\xb7\x37\x57\x3f\xb0\x38\x4f\x44\xb5\xa3\x0f\xe9\x94\xff\xbf\xcf\x0e\x5b\x3a\xfc\x3f\x7a\xfa\x7c\xf0\x10\xd6\x83\xfd\xe3\x60\x3f\x63\x76\x34\x0e\x86\x29\x30\x4c\x8b\x14\x9f\xb8\x1f\x16\xa3\x31\x11\xa8\x98\xa2\xdf\xb7\xba\x35\x05\xfe\xb6\xd6\xd6\x8d\xd5\x8e\x52\xcd\xc3\xb3\x4e\x9b\x1d\x33\xd8\x48\x5e\xe0\x8d\xe3\x0e\xe9\xc3\x3e\x07\xd5\x4a\x99\x03\x99\x62\x29\x5c\x1f\xa5\x9e\x41\x37\x5d\x46\x23\xa5\xa3\x93\xa6\x90\xcf\x86\x90\x21\x3e\x1b\x37\xb8\x2b\xf5\x56\x91\xfc\xa0\xe6\x52\x37\xf0\x0a\x0d\xb7\x0e\x9f\xde\x2e\x82\xe0\x78\xdf\x7b\xf0\xa5\x10\x34\xa4\x28\x36\x8f\xa4\x7b\x9f\xa4\x84\x0c\x6f\xa3\x58\x1c\x31\x1e\xc7\x43\x34\x00\x4f\x4f\x25\x5e\x96\x97\xbe\x54\xde\x09\xeb\x50\xa1\xa1\x24\x3e\xd2\x48\x9e\x5e\x6b\xd9\x7c\x12\x3f\xfc\xc9\x63\x4b\xfe\x7d\x13\x8a\x62\x9f\xf9\x9f\xbf\x03\x00\x00\xff\xff\xd1\x8f\x2c\x94\xf3\x0a\x00\x00"

func indexJsBytes() ([]byte, error) {
	return bindataRead(
		_indexJs,
		"index.js",
	)
}

func indexJs() (*asset, error) {
	bytes, err := indexJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "index.js", size: 2803, mode: os.FileMode(420), modTime: time.Unix(1656918125, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _linksCss = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\x8f\x41\x6e\xe3\x30\x0c\x45\xf7\x73\x0a\x23\xd9\x8e\x02\x25\x0e\x8a\x96\xbe\x45\x6f\xc0\x58\x92\x4d\x84\x11\x05\x89\xae\xed\x0a\xbd\x7b\x11\x24\x2e\x90\x05\x17\x04\x1f\x1f\xfe\xbf\x88\x5b\xeb\x05\xfb\xeb\x90\x65\x8a\x0e\xf6\x21\x84\x2e\x48\x54\x13\xf0\x46\xbc\xc2\xee\x13\xd9\xcf\xb8\xee\xfe\x17\x8c\xc5\x14\x9f\xe9\x09\x14\xfa\xf6\x70\x3e\xa5\xe5\xb1\xce\x9e\x86\x51\xa1\xb5\xf6\xe7\xc0\x14\xaf\xa5\xce\xe4\x74\x84\x77\x6b\xd3\xd2\xdd\x30\x0f\x14\xc1\x36\x38\xa9\x3c\x81\x06\x6b\x2f\x2c\x19\xf6\xf6\x23\x74\xea\x17\x35\xce\xf7\x92\x51\x49\x22\x44\x89\xfe\x0f\x84\x51\xbe\x7c\xae\x92\xb0\x27\x5d\xe1\xf0\xb6\x5d\x0e\x61\x62\x36\x53\xe6\x4d\xe5\x9c\x7b\xa8\xca\x88\x4e\x66\x38\xa6\xa5\xb9\x8f\x6d\xee\xd5\xb6\xb7\xf1\xb8\xf1\x6d\xdb\x3e\xc3\x99\x8b\xa8\xca\x0d\xce\x36\x2d\x1b\x37\x71\x4d\xe8\x1c\xc5\x01\x6c\xc7\x54\xd4\x14\x5d\xd9\x1b\x5d\x93\x7f\x49\x38\x71\xc3\x54\x5f\x45\xa7\xbb\xe8\xdf\x6f\x00\x00\x00\xff\xff\xb7\x76\xa3\x2a\x62\x01\x00\x00"

func linksCssBytes() ([]byte, error) {
	return bindataRead(
		_linksCss,
		"links.css",
	)
}

func linksCss() (*asset, error) {
	bytes, err := linksCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "links.css", size: 354, mode: os.FileMode(420), modTime: time.Unix(1659240024, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _linksHtml = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x52\x4d\x8f\x94\x40\x10\xbd\xef\xaf\x28\x3b\x7b\x5c\x28\x36\x7a\x30\xd8\x8d\x31\xab\xf1\xb2\x89\x66\x33\x1e\x3c\xf6\x40\x01\x1d\x1a\x7a\xa4\x0b\x94\x10\xfe\xbb\x69\x60\xe2\xcc\x98\x98\x58\x97\x4e\xd5\xab\x8f\xf7\x5e\x5a\xbe\xfa\xf8\xe5\xe9\xf0\xfd\xeb\x27\xa8\xb9\xb5\xd9\x9d\x3c\x3f\xa4\x8b\xec\x0e\x00\x40\xb2\x61\x4b\xd9\x67\x07\x69\x0a\x1f\x72\x36\x23\xc1\xb3\xe9\x1a\x2f\x71\x43\xb6\xae\x96\x58\x43\xcd\x7c\x8a\xe8\xc7\x60\x46\x25\x9e\x5c\xc7\xd4\x71\x74\x98\x4e\x24\x20\xdf\x32\x25\x98\x7e\x31\x86\x23\xef\x20\xaf\x75\xef\x89\xd5\xc0\x65\xf4\x56\xe0\xbe\xc8\x9a\xae\x81\xba\xa7\x52\x09\xf4\x18\x32\x1f\xe7\xde\x8b\x15\x0d\xd1\x93\x55\xc2\xf3\x64\xc9\xd7\x44\x2c\xfe\x9e\x0b\x34\x7c\x8a\x58\xba\x8e\x7d\x5c\x39\x57\x59\xd2\x27\xe3\xe3\xdc\xb5\x98\x7b\xff\xbe\xd4\xad\xb1\x93\x7a\xd1\x96\x7e\xea\x29\x7d\x93\x24\x0f\xaf\x93\xe4\x5f\x27\x24\x6e\x8e\xc8\xa3\x2b\xa6\xfd\x62\x61\x46\xc8\xad\xf6\x5e\x89\x95\xe6\xce\x64\xc5\xea\xc7\x6c\xf7\xca\x6e\x5e\xd5\x8f\x17\xe8\x60\xff\x24\x21\xe6\x19\x7a\xdd\x55\x04\xf7\x0d\x4d\x0f\x70\xdf\xbb\x81\x09\x52\x05\x31\x2c\xcb\x55\xa7\xb4\xe6\x7a\x74\x2d\xea\x5d\xf8\x3c\xef\xb3\xf1\xb7\x97\x67\x58\x16\x91\x55\x0e\x43\xb1\xa1\x09\x96\x45\xa2\xce\xe4\xb1\x07\xfc\x9f\x15\x67\x89\xe5\x60\x6d\x34\xf4\x56\x64\xb7\x2d\x61\xed\x35\x49\xbc\x65\x39\xcf\x40\x5d\x71\x29\x46\xe2\xd9\x04\x89\x85\x19\x83\xc3\x9b\xb5\x12\xb7\x2f\xf8\x3b\x00\x00\xff\xff\x59\x5f\xf5\xba\x9a\x02\x00\x00"

func linksHtmlBytes() ([]byte, error) {
	return bindataRead(
		_linksHtml,
		"links.html",
	)
}

func linksHtml() (*asset, error) {
	bytes, err := linksHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "links.html", size: 666, mode: os.FileMode(420), modTime: time.Unix(1659240567, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
	"about.html": aboutHtml,
	"close.svg": closeSvg,
	"edit.css": editCss,
	"edit.html": editHtml,
	"edit.js": editJs,
	"index.js": indexJs,
	"links.css": linksCss,
	"links.html": linksHtml,
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
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
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
	"about.html": &bintree{aboutHtml, map[string]*bintree{}},
	"close.svg": &bintree{closeSvg, map[string]*bintree{}},
	"edit.css": &bintree{editCss, map[string]*bintree{}},
	"edit.html": &bintree{editHtml, map[string]*bintree{}},
	"edit.js": &bintree{editJs, map[string]*bintree{}},
	"index.js": &bintree{indexJs, map[string]*bintree{}},
	"links.css": &bintree{linksCss, map[string]*bintree{}},
	"links.html": &bintree{linksHtml, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
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
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
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
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

