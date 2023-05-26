package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"golang.org/x/exp/slices"
)

const (
	BASE_TITLE      = "Breaking News English Lite"
	BASE_URL        = "https://breakingnewsenglish.com/"
	DIST_DIR        = "pages"
	NUMBER_OF_ITEMS = 30
)

type Content struct {
	MainPage      string
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
var americanAccent = []string{
	"level3", "level6",
}

func main() {
	for level, file := range levels {
		generate(level, file, NUMBER_OF_ITEMS)
	}
}

func generate(l string, f string, n int) {
	now := time.Now()

	feed := &feeds.Feed{
		Title:       fmt.Sprintf("%s %s", BASE_TITLE, strings.Title(l)),
		Link:        &feeds.Link{Href: fmt.Sprintf("%s%s", BASE_URL, f)},
		Description: "",
		Author:      &feeds.Author{Name: "longkey1", Email: "longkey1@gmail.com"},
		Created:     now,
	}

	resp, err := http.Get(feed.Link.Href)
	if err != nil {
		log.Fatalf("failed to get html: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("failed to load html: %s", err)
	}

	feed.Items = []*feeds.Item{}
	doc.Find("#primary li").Each(func(i int, s *goquery.Selection) {
		if len(feed.Items) == n {
			return
		}

		// date
		d := strings.Replace(strings.TrimSpace(s.Find("tt").Text()), ":", "", 1)
		if len(d) == 0 {
			return
		}
		date, err := time.Parse("2006-01-02", d)
		if err != nil {
			log.Fatalf("failed to parse date: %s", err)
		}

		// title
		title := strings.TrimSpace(s.Find("a").Text())
		log.Printf("%s %s", date.Format("2006-01-02"), title)

		// page
		href, _ := s.Find("a").Attr("href")
		mainPage := fmt.Sprintf("https://breakingnewsenglish.com/%s", strings.TrimSpace(href))
		content, err := getContent(l, mainPage)
		if err != nil {
			log.Printf("Skipped %s %s for failed to getContent: %s\n", date.Format("2006-01-02"), title, err)
			return
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: content.ListeningPage},
			Description: content.Text,
			Created:     date,
			Enclosure:   &feeds.Enclosure{Url: content.Audio, Type: "audio/mpeg", Length: content.AudioLength},
		})

		// sleep
		time.Sleep(500 * time.Millisecond)
	})

	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(DIST_DIR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	fp, err := os.Create(path.Join(DIST_DIR, fmt.Sprintf("%s.xml", l)))
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

	log.Printf("generated: %s", path.Join(DIST_DIR, fmt.Sprintf("%s.xml", l)))
}

func getContent(level string, page string) (Content, error) {
	res, err := http.Get(page)
	if err != nil {
		return Content{}, fmt.Errorf("failed to get html: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return Content{}, fmt.Errorf("failed to fetch data: %d, %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("failed to load html: %s", err)
	}

	listeningPage := strings.Replace(page, ".html", "-l.html", 1)
	text := strings.TrimSpace(doc.Find("article").Text())
	audio := strings.Replace(page, ".html", ".mp3", 1)
	if slices.Contains(americanAccent, level) {
		audio = strings.Replace(page, ".html", "-a.mp3", 1)
	}
	audioLength, err := getAudioLength(audio)
	if err != nil {
		return Content{}, fmt.Errorf("failed to get audio length: %w", err)
	}

	return Content{
		MainPage:      page,
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
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch audio file: %s, status code: %d, status: %s", file, res.StatusCode, res.Status)
	}

	length := res.Header.Get("Content-Length")

	return length, nil
}
