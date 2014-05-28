package util

import (
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io"
	"log"
	"os"
	"reflect"
)

var (
	debug = flag.Bool("debug", false, "Display debug information")
)

func main() {
	flag.Parse()
}

func PrintRed(messages ...interface{}) {
	PrintWithColor(ct.Red, messages)
}

func PrintGreen(messages ...interface{}) {
	PrintWithColor(ct.Green, messages)
}

func PrintOrange(messages ...interface{}) {
	PrintWithColor(ct.Yellow, messages)
}

func PrintWithColor(color ct.Color, messages []interface{}) {
	ct.ChangeColor(color, false, ct.None, false)

	args := make([]string, 0)
	for _, message := range messages {
		args = append(args, message.(string))
	}

	printFunc := reflect.ValueOf(fmt.Println)
	printFunc.Call(BuildReflectArguments(args))

	ct.ChangeColor(ct.Black, false, ct.None, false)
}

func LogError(err interface{}) {
	if *debug {
		panic(err)
	} else if reflect.TypeOf(err).String() == "string" {
		PrintRed(err)
		os.Exit(1)
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

func BuildReflectArguments(rawArgs []string) []reflect.Value {
	args := make([]reflect.Value, 0)

	for _, arg := range rawArgs {
		args = append(args, reflect.ValueOf(arg))
	}

	return args
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
