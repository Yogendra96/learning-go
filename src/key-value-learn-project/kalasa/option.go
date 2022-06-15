// Open Source: MIT License
// Author: Leon Ding <ding@ibyte.me>
// Date: 2022/2/26 - 10:47 PM - UTC/GMT+08:00

package bottle

import (
	"fmt"
	"os"
	"strings"
)

// Option bottle setting option
type Option struct {
	Directory       string `yaml:"Directory"`       // data directory
	DataFileMaxSize int64  `yaml:"DataFileMaxSize"` // data file max size
	Enable          bool   `yaml:"Enable"`          // data whether to enable encryption
	Secret          string `yaml:"Secret"`          // data encryption key
}

var (
	// DefaultOption default initialization option
	DefaultOption = Option{
		Directory:       "./data",
		DataFileMaxSize: 10240,
	}
)

// Validation verifying configuration Items
func (o *Option) Validation() {
	if o.Directory == "" {
		panic("The data file directory cannot be empty!!!")
	}

	// The first one does not determine whether there is a backslash
	o.Directory = pathBackslashes(o.Directory)

	// Record the location of the data file
	o.Directory = strings.TrimSpace(o.Directory)

	Root = o.Directory

	if o.DataFileMaxSize != 0 {
		defaultMaxFileSize = o.DataFileMaxSize
	}

	if o.Enable {
		if len(o.Secret) < 16 && len(o.Secret) > 16 {
			panic("The encryption key contains less than 16 characters!!!")
		}
		Secret = []byte(o.Secret)
		encoder = AES()
	}

	dataDirectory = fmt.Sprintf("%sdata/", Root)

	indexDirectory = fmt.Sprintf("%sindex/", Root)

}

// SetEncryptor Set up a custom encryption implementation
func SetEncryptor(encryptor Encryptor, secret []byte) {
	if len(secret) == 0 {
		panic("The key used by the encryptor cannot be empty!!!")
	}
	encoder.enable = true
	encoder.Encryptor = encryptor
	Secret = secret
}

// SetIndexSize set the expected index size to prevent secondary
// memory allocation and data migration during running
func SetIndexSize(size int32) {
	if size == 0 {
		return
	}
	index = make(map[uint64]*record, size)
}

// SetHashFunc sets the specified hash function
func SetHashFunc(hash Hashed) {
	HashedFunc = hash
}

// pathBackslashes Check directory ending backslashes
func pathBackslashes(path string) string {
	if !strings.HasSuffix(path, "/") {
		return fmt.Sprintf("%s/", path)
	}
	return path
}

// Checks whether the target path exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
