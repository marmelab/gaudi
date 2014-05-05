package util

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	debug = flag.Bool("debug", false, "Display debug information")
)

func main() {
	flag.Parse()
}

func LogError(err interface{}) {
	if *debug {
		panic(err)
	} else {
		log.Fatal(err)
	}
}

func Debug(notice ...interface{}) {
	if *debug {
		log.Println(notice)
	}
}

func IsDir(path string) bool {
	stat := getFileStat(path)
	if stat == nil {
		return false
	}

	return stat.Mode().IsDir()
}

func IsFile(path string) bool {
	stat := getFileStat(path)
	if stat == nil {
		return false
	}

	return stat.Mode().IsRegular()
}

func getFileStat(path string) os.FileInfo {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil
	}

	return stat
}

func GetFileCheckSum(filePath string) string {
	fileBuffer, err := os.Open(filePath)
	if err != nil {
		return ""
	}

	hash := md5.New()
	io.Copy(hash, fileBuffer)

	return fmt.Sprintf("%x", hash.Sum(nil))
}

/**
 * @see: https://gist.github.com/elazarl/5507969
 */
func Copy(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
