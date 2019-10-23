package PowerShellScripts

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type staticFilesFile struct {
	data  string
	mime  string
	mtime time.Time
	// size is the size before compression. If 0, it means the data is uncompressed
	size int
	// hash is a sha256 hash of the file contents. Used for the Etag, and useful for caching
	hash string
}

var staticFiles = map[string]*staticFilesFile{
	"GetADComputers.ps1": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\\\x90\xcbj\xc30\x10E\xf7\x02\xfdÅda\x17\xe4\x0f\bd\x11\x1cj\xdaE\v}\xacJ\x17\xae\x19\xa8C<j5\xa3v\xd1\xe4ߋ\U000b056c\xc4̜s\xa5\xd1\f\r)Vk\xd4~\xf8\x8aJ\x01\x0f\xed@bͼ;5\x04K\x14\r\xa9[\xadG\xc6\xdd\xf6\xdbtޔU±ó\x0f\xea\x1e?6\xd4)\xdc+\xf7ߑ\xac\xb1f\x86'\xd2\x18\x18\xfe0\xb2\xa6۶\"\x98\xa2\x04\u007f\xd6\x00\xc0|\xec\x8cudM\x19\xa9ʄb\"˳|\x10\xf4\xb3\x97\xaa\xce^}\x1d\x99S\x915\x11\x97Nu\xba3\x81{k\xf6\xc7\r\xeeX\xb4e\xed[%\x84\xcbm\xe6\xf9\"K\xbce\xe5\xfbb\xc1\xf4[L\xdfX\x1e\xd3j\xcf?\x14\x14\xea\xb1\x11\xcf\x10\x8fƣk9M$\x0e\x84\xfe:xw\x96^\xbc\xbbO\x8eK\xb3@\"\xf8\x0f\x00\x00\xff\xffɜ\xe3\xe9\xbf\x01\x00\x00",
		hash:  "d48d818a3b2bb3be9d7deb905313ef23447424803dfb74acecfa8fd7cc51ea63",
		mime:  "",
		mtime: time.Unix(1556326957, 0),
		size:  447,
	},
}

// NotFound is called when no asset is found.
// It defaults to http.NotFound but can be overwritten
var NotFound = http.NotFound

// ServeHTTP serves a request, attempting to reply with an embedded file.
func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/")
	f, ok := staticFiles[path]
	if !ok {
		if path != "" && !strings.HasSuffix(path, "/") {
			NotFound(rw, req)
			return
		}
		f, ok = staticFiles[path+"index.html"]
		if !ok {
			NotFound(rw, req)
			return
		}
	}
	header := rw.Header()
	if f.hash != "" {
		if hash := req.Header.Get("If-None-Match"); hash == f.hash {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("ETag", f.hash)
	}
	if !f.mtime.IsZero() {
		if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && f.mtime.Before(t.Add(1*time.Second)) {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("Last-Modified", f.mtime.UTC().Format(http.TimeFormat))
	}
	header.Set("Content-Type", f.mime)

	// Check if the asset is compressed in the binary
	if f.size == 0 {
		header.Set("Content-Length", strconv.Itoa(len(f.data)))
		io.WriteString(rw, f.data)
	} else {
		if header.Get("Content-Encoding") == "" && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			header.Set("Content-Encoding", "gzip")
			header.Set("Content-Length", strconv.Itoa(len(f.data)))
			io.WriteString(rw, f.data)
		} else {
			header.Set("Content-Length", strconv.Itoa(f.size))
			reader, _ := gzip.NewReader(strings.NewReader(f.data))
			io.Copy(rw, reader)
			reader.Close()
		}
	}
}

// Server is simply ServeHTTP but wrapped in http.HandlerFunc so it can be passed into net/http functions directly.
var Server http.Handler = http.HandlerFunc(ServeHTTP)

// Open allows you to read an embedded file directly. It will return a decompressing Reader if the file is embedded in compressed format.
// You should close the Reader after you're done with it.
func Open(name string) (io.ReadCloser, error) {
	f, ok := staticFiles[name]
	if !ok {
		return nil, fmt.Errorf("Asset %s not found", name)
	}

	if f.size == 0 {
		return ioutil.NopCloser(strings.NewReader(f.data)), nil
	}
	return gzip.NewReader(strings.NewReader(f.data))
}

// ModTime returns the modification time of the original file.
// Useful for caching purposes
// Returns zero time if the file is not in the bundle
func ModTime(file string) (t time.Time) {
	if f, ok := staticFiles[file]; ok {
		t = f.mtime
	}
	return
}

// Hash returns the hex-encoded SHA256 hash of the original file
// Used for the Etag, and useful for caching
// Returns an empty string if the file is not in the bundle
func Hash(file string) (s string) {
	if f, ok := staticFiles[file]; ok {
		s = f.hash
	}
	return
}
