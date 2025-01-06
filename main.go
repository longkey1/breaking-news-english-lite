package main

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"golang.org/x/exp/slices"
)

const (
	BaseTitle     = "Breaking News English Lite"
	BaseUrl       = "https://breakingnewsenglish.com/"
	DistDir       = "pages"
	NumberOfItems = 30
)

type Content struct {
	MainPage      string
	Title         string
	Date          time.Time
	ListeningPage string
	Text          string
	Audio         string
	AudioLength   string
}

var levels = map[string]string{
	"level0": "easy-news-english.html",
	"level1": "simple-english-news.html",
	"level2": "easy-english-news.html",
	"level3": "graded-news-stories.html",
	"level4": "graded-news-articles.html",
	"level5": "english-news-readings.html",
	"level6": "news-for-kids.html",
}
var usLevels = []string{
	"level3", "level6",
}

func main() {
	now := time.Now()
	generateIndex(now)

	var wg sync.WaitGroup
	for level, file := range levels {
		wg.Add(1)
		go func(l string, f string) {
			generatePageAndFeed(now, l, f, NumberOfItems, false)
			generatePageAndFeed(now, l, f, NumberOfItems, true)
			defer wg.Done()
		}(level, file)
	}
	wg.Wait()
}

func generateIndex(now time.Time) {
	tpl := template.Must(template.ParseFiles("templates/index.tpl"))

	values := map[string]interface{}{
		"UpdatedAt": now,
	}

	fp, err := os.Create(path.Join(DistDir, "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	defer func(fp *os.File) {
		_ = fp.Close()
	}(fp)
	if err := tpl.ExecuteTemplate(fp, "index.tpl", values); err != nil {
		log.Fatal(err)
	}
}

func generatePageAndFeed(now time.Time, l string, f string, n int, us bool) {
	if us == true && slices.Contains(usLevels, l) == false {
		return
	}

	ttlFmt := "%s %s"
	if us {
		ttlFmt = "%s(US) %s"
	}
	feed := &feeds.Feed{
		Title:       fmt.Sprintf(ttlFmt, BaseTitle, cases.Title(language.Und, cases.NoLower).String(l)),
		Link:        &feeds.Link{Href: fmt.Sprintf("%s%s", BaseUrl, f)},
		Description: "",
		Author:      &feeds.Author{Name: "longkey1", Email: "longkey1@gmail.com"},
		Created:     now,
	}

	resp, err := http.Get(feed.Link.Href)
	if err != nil {
		log.Fatalf("failed to get html: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("failed to load html: %s", err)
	}

	var Contents []*Content
	feed.Items = []*feeds.Item{}
	doc.Find("#primary li").Each(func(i int, s *goquery.Selection) {
		if len(feed.Items) == n {
			return
		}

		// content
		content, err := newContent(l, s, us)
		if err != nil {
			log.Printf("skipped %3d, %s", i, err)
			return
		}

		// contents
		Contents = append(Contents, &content)
		lv := l
		if us == true {
			lv = fmt.Sprintf("%s-us", l)
		}
		log.Printf("%s %s %s", lv, content.Date.Format("2006-01-02"), content.Title)

		// feed
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       content.Title,
			Link:        &feeds.Link{Href: content.ListeningPage},
			Description: content.Text,
			Created:     content.Date,
			Enclosure:   &feeds.Enclosure{Url: content.Audio, Type: "audio/mpeg", Length: content.AudioLength},
		})

		// sleep
		time.Sleep(1000 * time.Millisecond)
	})

	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(DistDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	fpFmt := "%s.xml"
	if us == true {
		fpFmt = "%s-us.xml"
	}
	fp, err := os.Create(path.Join(DistDir, fmt.Sprintf(fpFmt, l)))
	if err != nil {
		log.Fatal(err)
	}
	defer func(fp *os.File) {
		_ = fp.Close()
	}(fp)

	_, err = fp.WriteString(atom)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("generated: %s", path.Join(DistDir, fmt.Sprintf(fpFmt, l)))

	fp2Fmt := "%s.html"
	if us == true {
		fp2Fmt = "%s-us.html"
	}
	fp2, err := os.Create(path.Join(DistDir, fmt.Sprintf(fp2Fmt, l)))
	if err != nil {
		log.Fatal(err)
	}
	defer func(fp2 *os.File) {
		_ = fp2.Close()
	}(fp2)

	tpl := template.Must(template.ParseFiles("templates/page.tpl"))

	fdFmt := "%s.xml"
	if us == true {
		fdFmt = "%s-us.xml"
	}
	values := map[string]interface{}{
		"Title":     fmt.Sprintf(ttlFmt, BaseTitle, cases.Title(language.Und, cases.NoLower).String(l)),
		"Feed":      fmt.Sprintf(fdFmt, l),
		"Contents":  Contents,
		"UpdatedAt": now,
	}
	if err = tpl.ExecuteTemplate(fp2, "page.tpl", values); err != nil {
		log.Fatal(err)
	}

	log.Printf("generated: %s", path.Join(DistDir, fmt.Sprintf(fp2Fmt, l)))
}

func newContent(l string, s *goquery.Selection, a bool) (Content, error) {
	// date
	d := strings.Replace(strings.TrimSpace(s.Find("tt").Text()), ":", "", 1)
	if len(d) == 0 {
		return Content{}, fmt.Errorf("not found date string from %s page", l)
	}
	date, err := time.Parse("2006-01-02", d)
	if err != nil {
		return Content{}, fmt.Errorf("failed to parse date: %w", err)
	}

	// title
	title := strings.TrimSpace(s.Find("a").Text())

	// href
	href, _ := s.Find("a").Attr("href")

	// mainPage
	mainPage := fmt.Sprintf("https://breakingnewsenglish.com/%s", strings.TrimSpace(href))

	// listingPage, text, audio, audioLength
	res, err := http.Get(mainPage)
	if err != nil {
		return Content{}, fmt.Errorf("failed to get html: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return Content{}, fmt.Errorf("failed to fetch data: %d, %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return Content{}, fmt.Errorf("failed to load html: %w", err)
	}

	listeningPage := strings.Replace(mainPage, ".html", "l.html", 1)
	if slices.Contains(usLevels, l) {
		listeningPage = strings.Replace(mainPage, ".html", "-l.html", 1)
	}
	text := strings.TrimSpace(doc.Find("article").Text())
	audio := strings.Replace(mainPage, ".html", ".mp3", 1)
	if a == true {
		audio = strings.Replace(mainPage, ".html", "-a.mp3", 1)
	}
	audioLength, err := getAudioLength(audio)
	if err != nil {
		return Content{}, fmt.Errorf("failed to get audio length: %w", err)
	}

	return Content{
		Title:         title,
		MainPage:      mainPage,
		Date:          date,
		ListeningPage: listeningPage,
		Text:          text,
		Audio:         audio,
		AudioLength:   audioLength,
	}, nil
}

func getAudioLength(file string) (string, error) {
	res, err := http.Head(file)
	if err != nil {
		return "", fmt.Errorf("failed to get audio file: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch audio file: %s, status code: %d, status: %s", file, res.StatusCode, res.Status)
	}

	length := res.Header.Get("Content-Length")

	return length, nil
}
