package newepub

import (
	"fmt"
	"github.com/bmaupin/go-epub"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type NewEpub struct {
	sync.Mutex
	Epub       *epub.Epub
	epubConfig *ConfigEpub
}

func NewNewEpubClient(epubConfig *ConfigEpub) *NewEpub {
	if epubConfig.BookName == "" {
		log.Fatal("epubConfig.BookName is empty, please use SetBookName(bookName string) first")
		return nil
	}

	newClient := &NewEpub{Epub: epub.NewEpub(epubConfig.BookName), epubConfig: epubConfig}
	newClient.Epub.SetAuthor(epubConfig.AuthorName)
	newClient.Epub.SetDescription(epubConfig.Intro)
	newClient.Epub.SetLang("zh-CN")
	bookIntroHtml := newClient.epubConfig.getBookIntroHtml()
	if newClient.epubConfig.Cover != "" {
		co := NewCover(
			newClient.Epub,
			WithCoverSavePath(path.Dir(newClient.epubConfig.getTextPath())+"/cover.jpg"),
			WithCoverUrl(newClient.epubConfig.Cover),
		).downloadCover().SetEpubCover()
		bookIntroHtml += fmt.Sprintf(`<br/><br/><img src="%v" />`, co.GetEpubCoverPath())
	}

	_, _ = newClient.addChapter("书籍信息", bookIntroHtml)
	NewOpen(newClient.epubConfig.getTextPath(), "w").Write(newClient.epubConfig.getBookIntro())

	return newClient
}
func (epub *ConfigEpub) getBookIntroHtml() string {
	var bookInfoHtml = fmt.Sprintf(`<h1 style="text-align: center;">%v</h1>`, epub.BookName)
	if epub.AuthorName != "" {
		bookInfoHtml += fmt.Sprintf(`<p style="text-align: center;">作者：%v</p>`, epub.AuthorName)
	}
	if epub.WordCount != "" {
		bookInfoHtml += fmt.Sprintf(`<p style="text-align: center;">字数：%v</p>`, epub.WordCount)
	}
	if epub.TagName != "" {
		bookInfoHtml += fmt.Sprintf(`<p style="text-align: center;">标签：%v</p>`, epub.TagName)
	}
	if epub.Intro != "" {
		bookInfoHtml += fmt.Sprintf(`<h3>简介：%v</h3>`, epub.Intro)
	}
	return bookInfoHtml
}

func (epub *ConfigEpub) getBookIntro() string {
	return regexp.MustCompile(`<[^>]+>`).ReplaceAllString(epub.getBookIntroHtml(), "\n")
}

func (epub *NewEpub) addChapter(title, content string) (string, error) {
	var contentHtml string
	for contentLineIndex, contentLine := range strings.Split(content, "\n") {
		if contentLine != "" {
			image := regexp.MustCompile(`\[img=\d+,\d+](.*?)\[/img]|<img .*src="(.*?)">`).FindStringSubmatch(contentLine)
			if len(image) > 0 {
				imageName := regexp.MustCompile(`[\\/:*?"<>|]`).ReplaceAllString(title, "") + "_" + strconv.Itoa(contentLineIndex)
				co := NewCover(
					epub.Epub,
					WithCoverSavePath(path.Dir(epub.epubConfig.getTextPath())+"/"+imageName+".jpg"),
					WithCoverUrl(image[1]),
				).downloadCover()
				contentLine = strings.ReplaceAll(contentLine, image[0], fmt.Sprintf(`<img src="%v" />`, co.GetEpubCoverPath()))
			}
			contentHtml += fmt.Sprintf("<p>%s</p>", contentLine)
		}
	}
	contentHtml = fmt.Sprintf(`<h1>%s</h1><br/>%s`, title, contentHtml)
	return epub.Epub.AddSection(contentHtml, title, "", "")
}

func (epub *NewEpub) SaveChapter(Index int, title, content string) {
	epub.Lock()
	defer epub.Unlock()
	epub.addChapter(title, content)
	NewOpen(epub.epubConfig.getTextPath(), "a").Write(fmt.Sprintf("\n\n第%d章: %s\n%s", Index, title, content))
}

func (epub *NewEpub) Save() error {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Save Epub error:", err)
		}
	}()
	return epub.Epub.Write(epub.epubConfig.getEpubPath())
}
