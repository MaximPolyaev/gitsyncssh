package main

import (
	"errors"
	"fmt"
	"github.com/melbahja/goph"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := sync(); err != nil {
		log.Fatal(err)
	}
}

func sync() error {
	localPath, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println("you localPath:", localPath)

	idRsa := os.Getenv("HOME") + "/.ssh/id_rsa"

	fmt.Println("you id_rsa:", idRsa)

	fmt.Print("Enter user: ")
	var user string
	_, err = fmt.Scanln(&user)
	if err != nil {
		return err
	}

	fmt.Print("Enter addr: ")
	var addr string
	_, err = fmt.Scanln(&addr)
	if err != nil {
		return err
	}

	fmt.Print("Enter password: ")
	var passwd string
	_, err = fmt.Scanln(&passwd)
	if err != nil {
		return err
	}

	client, err := goph.New(user, addr, goph.Password(passwd))
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	out, err := client.Run("cd projects/gitclone && git ls-files -m && git ls-files -o")
	if err != nil {
		log.Fatal(err)
	}

	if out == nil {
		fmt.Println("not changes")
		return nil
	}

	strOut := string(out)
	if strOut == "" {
		fmt.Println("not changes")
		return nil
	}

	files := strings.FieldsFunc(strOut, func(r rune) bool {
		return r == '\n'
	})

	fmt.Println("changed:")
	fmt.Println(files)

	for _, file := range files {
		targetFilePath := localPath + "/" + file
		targetPath := filepath.Dir(targetFilePath)
		err := os.MkdirAll(targetPath, os.ModeDir|os.ModePerm)
		if err != nil {
			return errors.New("mkdir " + targetPath + ": " + err.Error())
		}
		err = client.Download("projects/gitclone/"+file, targetFilePath)

		if errors.Is(err, os.ErrNotExist) {
			err := os.Remove(localPath + "/" + file)
			if err != nil {
				return errors.New("remove " + file + ": " + err.Error())
			}
			fmt.Println("removed:", file)
			continue
		}

		if err != nil {
			return errors.New("download " + file + ": " + err.Error())
		}
		fmt.Println("sync: ", file)
	}
	return nil
}
