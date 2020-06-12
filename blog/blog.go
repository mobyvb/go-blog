package blog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Blog struct {
	Pages []*Page
}

func BlogFromDir(dir string) (*Blog, error) {
	blog := &Blog{}

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

func (blog *Blog) WriteToDir(output string) error {
	for _, page := range blog.Pages {
		html := page.HTML()

		path := filepath.Join(output, string(page.Slug)+".html")
		err := ioutil.WriteFile(path, []byte(html), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %q: %w", path, err)
		}
	}

	return nil
}
