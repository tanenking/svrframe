//go:build linux
// +build linux

package application

import (
	"fmt"
	"github.com/tanenking/svrframe/constants"
	"os"
	"os/exec"
	"path/filepath"
)

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writePid() bool {
	var f *os.File
	var err1 error
	path := "./pid/"
	is, _exists := isDir(path)
	if _exists && !is {
		fmt.Println("pid不是文件夹")
		return false
	}
	if !_exists {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("%v\n", err)
			return false
		}
	}
	fpath := fmt.Sprintf("%s%s-%s-%d", path, constants.ProjectName, constants.Service_Type, config.GetServiceInfo().ServiceID)
	if checkFileIsExist(fpath) { //如果文件存在
		os.Remove(fpath)
	}
	f, err1 = os.Create(fpath) //创建文件
	check(err1)
	pidinfo := fmt.Sprintf("%d", os.Getpid())
	f.WriteString(pidinfo)

	return true
}

func doFork() bool {
	if os.Getppid() == 1 {
		return false
	}

	filePath, _ := filepath.Abs(os.Args[0])
	cmd := exec.Command(filePath, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	return true
}
