package filemeta

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	// Verbose enables verbose logging in this package
	Verbose = false
)

// Sourcefile represents the metadata for a single backed up file
type Sourcefile struct {
	// Key is the location of this file in a bucket
	Key string `yaml:"key"`
}

// Index holds all of the metadata for files backed up
type Index struct {
	// Files maps the local file location to its metadata
	Files map[string]Sourcefile `yaml:"files"`
}

// NewIndex creates an Index from Yaml
func NewIndex(buf string) (*Index, error) {
	index := &Index{}
	err := yaml.Unmarshal([]byte(buf), index)
	if err != nil {
		return nil, err
	}

	return index, nil
}

// Encode the index data as Yaml
func (i *Index) Encode() (string, error) {
	out, err := yaml.Marshal(i)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

type PathWalker func(root string, index *Index) filepath.WalkFunc

func FilePathWalker(root string, index *Index) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			doLog("Add in file to index: %s", path)
			key := path
			if root != "" {
				key = fmt.Sprintf("%s/%s", root, path)
			}
			index.Files[path] = Sourcefile{
				Key: key,
			}
		}
		return err
	}
}

// NewIndexFromRoot creates a new Index populated from a filesystem directory
func NewIndexFromRoot(bucketRoot, path string, walker PathWalker) *Index {
	i := &Index{
		Files: map[string]Sourcefile{},
	}

	filepath.Walk(path, walker(bucketRoot, i))

	return i
}

func doLog(format string, args ...interface{}) {
	if Verbose {
		log.Printf(format, args...)
	}
}