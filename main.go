package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

type userComment struct {
	By    string
	Text  string
	Karma int
}

func main() {
	var list []userComment

	comments := getPostComments()
	for _, id := range comments {
		username, text := getComment(id)
		karma := getUserKarma(username)

		if karma <= 100 {
			fmt.Println("Ignoring: ", username, "(", karma, ")")
			continue
		}

		list = append(list, userComment{username, text, karma})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Karma > list[j].Karma
	})

	createOPML(list)
}

func createOPML(list []userComment) {
	output, err := os.Create("list.opml")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = output.Close()
	}()

	_, err = io.WriteString(output, `<?xml version="1.0" encoding="UTF-8"?>
<opml version="2.0">
<head>
<title>Hacker News Personal Blogs</title>
</head>
<body>
<outline title="HN Personal Blogs" text="HN Personal Blogs">`)
	if err != nil {
		panic(err)
	}

	for _, comment := range list {
		blog, err := extractBlogURL(comment.Text)
		if err != nil {
			fmt.Println("No URL:", comment.By, "(", comment.Karma, ")", err)
			continue
		}

		feed, err := findAtomFeed(blog)
		if err != nil {
			fmt.Println("No feed: ", blog, "by", comment.By, "(", comment.Karma, ")", err)
			continue
		}

		feedEscaped := bytes.Buffer{}
		err = xml.EscapeText(&feedEscaped, []byte(feed))
		if err != nil {
			panic(err)
		}

		byEscaped := bytes.Buffer{}
		err = xml.EscapeText(&byEscaped, []byte(comment.By))
		if err != nil {
			panic(err)
		}
		_, err = io.WriteString(output, fmt.Sprintf(`
	<outline type="rss" title="%s" text="%s" xmlUrl="%s" htmlUrl="%s"/>`,
			byEscaped.String(), byEscaped.String(), feedEscaped.String(), feedEscaped.String()))
		if err != nil {
			panic(err)
		}
	}

	_, err = io.WriteString(output, `
</outline>
</body>
</opml>
`)
	if err != nil {
		panic(err)
	}
}

func extractBlogURL(text string) (string, error) {
	re := regexp.MustCompile(`href="([^"]+)"`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	re = regexp.MustCompile(`(https?://)(www\.)[^ ]+`)
	matches = re.FindStringSubmatch(text)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("no blog url found in %s", text)
}

func findAtomFeed(url string) (string, error) {
	content := getBlogContent(url, 3)
	if content == "" {
		return "", fmt.Errorf("no content for %s", url)
	}

	re := regexp.MustCompile(`<link\s+[^>]*rel="?alternate"?[^>]+>`)
	matches := re.FindAllString(content, -1)
	for _, match := range matches {
		if regexp.MustCompile(`type="?application/((rss)|(atom))\+xml"?`).MatchString(match) {
			re = regexp.MustCompile(`href="?([^"\s]+)"?`)
			matches = re.FindStringSubmatch(match)
			if len(matches) > 1 {
				feedUrl := matches[1]
				if feedUrl[0] == '/' {
					feedUrl = url + feedUrl
				} else if !strings.HasPrefix(feedUrl, "http") {
					feedUrl = url + "/" + feedUrl
				}
				return feedUrl, nil
			}
		}
	}
	return "", fmt.Errorf("no feed found for %s", url)
}

func getBlogContent(url string, retries int) string {
	for i := 0; i < retries; i++ {
		res, err := http.Get(url)
		if err != nil {
			continue
		}
		defer func() {
			_ = res.Body.Close()
		}()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			continue
		}
		return string(data)
	}
	return ""
}

func getPostComments() []int {
	post := struct {
		Kids []int
	}{}
	getAndParse("https://hacker-news.firebaseio.com/v0/item/36575081.json", &post)
	return post.Kids
}

func getComment(id int) (string, string) {
	comment := struct {
		By   string
		Text string
	}{}
	getAndParse(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id), &comment)
	return comment.By, html.UnescapeString(comment.Text)
}

func getUserKarma(username string) int {
	user := struct {
		Karma int
	}{}
	getAndParse(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/user/%s.json", username), &user)
	return user.Karma
}

func getAndParse(url string, v any) {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}
