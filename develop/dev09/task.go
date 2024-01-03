package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Args struct {
	O        []string
	maxDepth int
	address  []string
}

func getArgs() (*Args, error) {
	O := flag.String("O", "", "filename")
	maxDepth := flag.Int("maxDepth", 1, "set the depth for recursively loading")

	flag.Parse()

	if *maxDepth < 1 {
		return nil, errors.New("depth must be positive number")
	}

	var files []string

	if len(*O) > 0 {
		files = strings.Split(*O, " ")
	}

	args := &Args{
		O:        files,
		maxDepth: *maxDepth,
	}

	if len(flag.Args()) < 1 {
		return nil, errors.New("incorrect address")
	}

	args.address = append(args.address, flag.Args()...)

	return args, nil
}

func getLinks(address string) map[string]bool {
	links := make(map[string]bool)

	parsed, _ := url.Parse(address)
	host := parsed.Hostname()

	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(address)
	if err != nil || resp == nil {
		return nil
	}
	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	document.Find("a").Each(func(index int, element *goquery.Selection) {
		link, _ := element.Attr("href")
		parsed, err := url.Parse(link)
		if err != nil || parsed.Path == "" {
			return
		}

		linkHost := parsed.Hostname()
		if linkHost != "" && linkHost != host {
			return
		}

		scheme := "https"
		if parsed.Scheme != "" {
			scheme = parsed.Scheme
		}

		newLink := fmt.Sprintf("%s://%s%s", scheme, host, parsed.Path)

		links[newLink] = true
	})

	return links
}

func getFilename(filename, address, path string) string {
	var counter int

	if filename == "" {
		if strings.HasSuffix(address, "/") {
			filename = "index.html"
		} else {
			filename = filepath.Base(address)
			if !strings.Contains(filename, ".") {
				filename = fmt.Sprintf("%s.html", filename)
			}
		}
	}

	origName := filename

	for {
		_, err := os.Stat(fmt.Sprintf("%s/%s", path, filename))

		if errors.Is(err, os.ErrNotExist) {
			break
		}

		counter++
		suffix := fmt.Sprintf(".%d", counter)

		newFileName := strings.Builder{}
		newFileName.Grow(len(origName) + len(suffix))
		newFileName.WriteString(origName)
		newFileName.WriteString(suffix)

		filename = newFileName.String()
	}

	return filename
}

func saveFile(body []byte, filename, address string) error {
	parsed, err := url.Parse(address)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s%s", parsed.Host, parsed.Path)
	path = filepath.Dir(path)

	err = os.MkdirAll(path, 0666)
	if err != nil {
		return err
	}

	filename = getFilename(filename, address, path)

	err = os.WriteFile(fmt.Sprintf("%s/%s", path, filename), body, 0666)
	if err != nil {
		return err
	}

	fmt.Printf("file %s saved in dir %s\n\n", filename, path)

	return nil
}

func download(address, filename string, maxDepth int) error {
	if maxDepth < 1 {
		return nil
	}

	fmt.Printf("Sending http request to %s\n", address)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(address)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("response status: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = saveFile(body, filename, address)
	if err != nil {
		return err
	}

	if maxDepth-1 > 1 {
		links := getLinks(address)

		for link := range links {
			err := download(link, filename, maxDepth-1)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func wget() error {
	if len(os.Args) < 2 {
		return errors.New("you must specify a web-address")
	}

	args, err := getArgs()
	if err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGSEGV)

	go func() {
		<-sigs
		fmt.Println("Stopped by exceeded time")
		os.Exit(1)
	}()

	for i, address := range args.address {
		var filename string

		if i < len(args.O) {
			filename = args.O[i]
		}

		err := download(address, filename, args.maxDepth)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := wget()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
