package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const query = "http://en.wikipedia.org/w/api.php?action=query&titles=Image:%s.jpg&prop=imageinfo&iiprop=url"
const regxp = `>http://upload.wikimedia.org/wikipedia/commons/.*/%s.jpg`

func grabImage(imgurl string, word string) {
	fmt.Printf("Requesting with: %s\n", imgurl)
	imgRsp, _ := http.Get(imgurl)
	defer imgRsp.Body.Close()
	imgBody, err := ioutil.ReadAll(imgRsp.Body)
	if err != nil {
		fmt.Println("Couldn't download image")
		return
	}
	ioutil.WriteFile(fmt.Sprintf("images/%s.jpg", word), imgBody, 0777)
}

func main() {
	os.Mkdir("images", 0777)

	data, _ := ioutil.ReadFile("wordlist.txt")
	list := strings.Split(string(data), "\n")

	for _, word := range list {
		resp, _ := http.Get(fmt.Sprintf(query, word))
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		imgUrl := regexp.MustCompile(fmt.Sprintf(regxp, word))
		match := imgUrl.Find(body)
		if match == nil {
			fmt.Printf("No match found for: %s\n", word)
			continue
		}

		imgurl := string(match)
		imgurl = imgurl[1:len(imgurl)]

		go grabImage(imgurl, word)
	}
}
