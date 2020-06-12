package blog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Blog struct {
	Title string
	Pages []*Page
}

func BlogFromDir(title string, dir string) (*Blog, error) {
	blog := &Blog{
		Title: title,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".txt" {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file: %w", err)
		}

		relpath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("unable to get relative path %q : %q: %w", dir, path, err)
		}

		slug := SlugFromPath(relpath)
		page, err := ParsePageString(slug, string(content))
		if err != nil {
			return fmt.Errorf("unable to parse content %q: %w", path, err)
		}

		blog.Pages = append(blog.Pages, page)

		return nil
	})

	return blog, err
}
