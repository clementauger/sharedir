package main

import (
	"io"
	"net/http"
	"net/url"
	"sort"

	humanize "github.com/dustin/go-humanize"
)

type TplExecer interface {
	Execute(io.Writer, interface{}) error
}

func DirList(tpl TplExecer) DirListHandler {
	return func(w http.ResponseWriter, r *http.Request, f http.File) {
		dirs, err := readdir(f)
		if err != nil {
			http.Error(w, "Error reading directory", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{
			"Path":      r.URL.Path,
			"Items":     dirs,
			"PageTitle": "sharedir",
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type dirItem struct {
	Name    string
	URL     string
	Size    string
	ModTime string
}

func readdir(f http.File) ([]dirItem, error) {
	out := []dirItem{}
	dirs, err := f.Readdir(-1)
	if err != nil {
		return out, err
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: name}
		out = append(out, dirItem{
			Name: name,
			URL:  url.String(),
			// Size: fmt.Sprint(d.Size()),
			Size: humanize.Bytes(uint64(d.Size())),
			// ModTime: fmt.Sprint(d.ModTime().Format(time.Kitchen)),
			ModTime: humanize.Time(d.ModTime()),
		})
	}
	return out, nil
}
