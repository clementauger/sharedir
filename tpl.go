package main

import (
	"html/template"
	"io"
	"path/filepath"
)

func assetTemplate(fpath string) (*template.Template, error) {
	return byteTemplate(MustAsset(fpath))
}

func byteTemplate(data []byte) (*template.Template, error) {
	return template.New("b").Parse(string(data))
}

func fileTemplate(fpath ...string) (*fsTemplate, error) {
	t := &fsTemplate{
		fpath: fpath,
	}
	return t, nil
}

type fsTemplate struct {
	fpath []string
}

func (t *fsTemplate) parse() (*template.Template, error) {
	var n string
	if len(t.fpath) > 0 {
		n = filepath.Base(t.fpath[0])
	}
	return template.New(n).ParseFiles(t.fpath...)
}

func (t *fsTemplate) Execute(wr io.Writer, data interface{}) error {
	tpl, err := t.parse()
	if err != nil {
		return err
	}
	return tpl.Execute(wr, data)
}
