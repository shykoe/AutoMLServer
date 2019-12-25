package main

import (
	"archive/tar"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func createTar(path string, target string) error {
	fw, err := os.Create(target)
	if err != nil {
		return err
	}
	defer fw.Close()
	tw := tar.NewWriter(fw)
	defer func() {
		if err := tw.Close(); err != nil {
			log.Error(err)
		}
	}()
	return filepath.Walk(path, func(fileName string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		// 这里需要处理下 hdr 中的 Name，因为默认文件的名字是不带路径的，
		// 打包之后所有文件就会堆在一起，这样就破坏了原本的目录结果
		// 例如： 将原本 hdr.Name 的 syslog 替换程 log/syslog
		// 这个其实也很简单，回调函数的 fileName 字段给我们返回来的就是完整路径的 log/syslog
		// strings.TrimPrefix 将 fileName 的最左侧的 / 去掉，
		hdr.Name = strings.TrimPrefix(strings.TrimPrefix(fileName, path), string(filepath.Separator))
		log.Info(hdr.Name)
		if hdr.Name == ""{
			return nil
		}
		// 写入文件信息
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// 判断下文件是否是标准文件，如果不是就不处理了，
		// 如： 目录，这里就只记录了文件信息，不会执行下面的 copy
		if !fi.Mode().IsRegular() {
			return nil
		}

		// 打开文件
		fr, err := os.Open(fileName)
		defer fr.Close()
		if err != nil {
			return err
		}

		// copy 文件数据到 tw
		n, err := io.Copy(tw, fr)
		if err != nil {
			return err
		}

		// 记录下过程，这个可以不记录，这个看需要，这样可以看到打包的过程
		log.Printf("成功打包 %s ，共写入了 %d 字节的数据\n", fileName, n)

		return nil
	})
}
