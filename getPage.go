package Hornetgo

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"strings"
)

func GetPageByDiy(name string) ([]byte, error) {

	result := localCache.Get(name)

	if by, ok := result.([]byte); ok {
		return by, nil
	}

	fileName, err := getFileName(name)
	if err != nil {
		return []byte{}, err
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		Error(err)
	} else {
		localCache.Put(name, data, time.Hour*1)
	}

	return data, err

}

// GetPageByTemplate GetPageByTemplate
func GetPageByTemplate(name string, data interface{}) ([]byte, error) {

	pageData, err := getPageFile(name)
	if err != nil {
		return []byte{}, err
	}

	nameList := strings.Split(name, "/")

	t := template.New(nameList[len(nameList)-1]).Delims("<%", "%>")

	t, err = t.Parse(string(pageData))
	if err != nil {
		return []byte{}, err
	}

	w := &bytes.Buffer{}

	err = t.Execute(w, data)
	if err != nil {
		return []byte{}, err
	}

	return w.Bytes(), nil

}

func getPageFile(name string) ([]byte, error) {

	result := localCache.Get(name)

	if by, ok := result.([]byte); ok {
		return by, nil
	}

	fileName, err := getFileName(name)
	if err != nil {
		return []byte{}, err
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		Error(err)
	} else {
		localCache.Put(name, data, time.Hour*1)
	}

	return data, err

}

func getFileName(name string) (string, error) {

	fileDri, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(fileDri + HornetInfo.AppConfig.WebConfig.ViewsPath)
	if err != nil {
		return "", err
	}

	return filepath.Join(abs, name), nil

}

func render(name string, data interface{}, ctx *HornetContent) {

	body, err := GetPageByTemplate(name, data)
	if err != nil {
		Error(err)
		body = []byte("not find")
	}

	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	ctx.Write(body)

	return
}
