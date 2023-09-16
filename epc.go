package newepub

import (
	"fmt"
	"path"
	"strings"
)

type ConfigEpub struct {
	BookName   string
	Intro      string
	WordCount  string
	AuthorName string
	TagName    string
	Cover      string
	dirPath    string
}

type Option interface {
	AddOption(*ConfigEpub)
}
type OptionFunc func(*ConfigEpub)

func (e OptionFunc) AddOption(epub *ConfigEpub) {
	e(epub)
}

func WithBookName(bookName string) Option {
	return OptionFunc(func(epub *ConfigEpub) {
		epub.BookName = bookName
	})
}

func WithAuthorName(authorName string) Option {
	return OptionFunc(func(epub *ConfigEpub) {
		epub.AuthorName = authorName
	})
}
func WithWordCount(wordCount any) Option {
	return OptionFunc(func(epub *ConfigEpub) {
		epub.WordCount = fmt.Sprintf("%v", wordCount)
	})
}
func WithIntro(intro string) Option {
	return OptionFunc(func(epub *ConfigEpub) {
		epub.Intro = intro
	})
}
func WithCover(cover string) Option {
	return OptionFunc(func(epub *ConfigEpub) {
		epub.Cover = cover
	})
}
func WithTagName(tagNameList []string) Option {
	return OptionFunc(func(epub *ConfigEpub) {
		epub.TagName = strings.Join(tagNameList, ",")
	})
}

func WithPath(dirPath string) Option {
	return OptionFunc(func(epub *ConfigEpub) {
		dir := path.Join(dirPath, epub.BookName)
		if !Exists(dir) {
			MkdirAll(dir)
		}
		epub.dirPath = dirPath
	})
}

func NewConfigEpub(opts ...Option) *ConfigEpub {
	c := &ConfigEpub{}
	for _, opt := range opts {
		opt.AddOption(c)
	}
	return c
}

func (epub *ConfigEpub) SetPath(dirPath string) *ConfigEpub {
	dir := path.Join(dirPath, epub.BookName)
	if !Exists(dir) {
		MkdirAll(dir)
	}
	epub.dirPath = dirPath
	return epub
}
func (epub *ConfigEpub) getTextPath() string {
	return path.Join(epub.dirPath, epub.BookName, epub.BookName+".txt")
}

func (epub *ConfigEpub) getEpubPath() string {
	return path.Join(epub.dirPath, epub.BookName, epub.BookName+".epub")
}
