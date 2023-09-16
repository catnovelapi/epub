package newepub

import (
	"github.com/bmaupin/go-epub"
	"github.com/catnovelapi/BuilderHttpClient"
	"log"
)

type Cover struct {
	epub         *epub.Epub
	path         string
	header       string
	url          string
	relativePath string
}

type CoverOption interface {
	AddOption(*Cover)
}

type CoverOptionFunc func(*Cover)

func (e CoverOptionFunc) AddOption(cover *Cover) {
	e(cover)
}

func NewCover(epub *epub.Epub, coverOption ...CoverOption) *Cover {
	cover := &Cover{
		epub:   epub,
		header: "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)"}
	for _, option := range coverOption {
		option.AddOption(cover)
	}
	return cover
}

func WithHeader(header string) CoverOption {
	return CoverOptionFunc(func(cover *Cover) {
		cover.header = header
	})
}

func WithCoverSavePath(path string) CoverOption {
	return CoverOptionFunc(func(cover *Cover) {
		cover.path = path
	})
}
func WithCoverUrl(url string) CoverOption {
	return CoverOptionFunc(func(cover *Cover) {
		cover.url = url
	})
}

func (c *Cover) SetEpubCover() *Cover {
	if c.relativePath == "" {
		log.Printf("SetEpubCover error: c.relativePath == \"\"\n")
		return c
	}
	c.epub.SetCover(c.relativePath, "")
	return c
}
func (c *Cover) GetEpubCoverPath() string {
	return c.relativePath
}

func (c *Cover) downloadCover() *Cover {
	if Exists(c.path) {
		c.relativePath = c.addImage()
	}
	image := BuilderHttpClient.Get(c.url, BuilderHttpClient.Header(map[string]any{"User-Agent": c.header})).Byte()
	if len(image) == 0 {
		log.Println("downloadCover error: len(image) == 0")
		return c
	}
	NewOpen(c.path, "w").Write(image)
	c.relativePath = c.addImage()
	return c
}

func (c *Cover) addImage() string {
	image, err := c.epub.AddImage(c.path, "")
	if err != nil {
		log.Printf("addImage error: %v\n", err)
		return ""
	}
	return image
}
