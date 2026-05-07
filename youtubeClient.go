package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
	"log"
)

type Video struct {
	Title string
	ID    string
	URL   string
}

var searchAmount = "10"

var showErrorOutput = false

func search(query string) ([]Video, error) {
	cmd := exec.Command("yt-dlp", "-J", "ytsearch"+searchAmount+":"+query)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("The program has encountered some kind of problem that I am too high to deal with.\nProbably with yt-dlp, fat chance it either can't find the video you're searching for or can't find js run time on the page. \nIt could be alot really, rate limiting, bad network conncection, youtube changed something, you're using a bad or old version of yt-dlp, there's an issue with the program.\n It could also be bad formating. I suggest you reword your search and try again. Fuck I'm tired.")
		if showErrorOutput == true {
			fmt.Println("Details:", string(out))
		}
		return nil, err
	}

	fmt.Println("YOU ARE USING YT-DLP FOR THIS VIDEO SEARCH!")

	var data map[string]interface{}
	json.Unmarshal(out, &data)

	entries := data["entries"].([]interface{})
	videos := []Video{}

	for _, e := range entries {
		v := e.(map[string]interface{})
		id := v["id"].(string)

		videos = append(videos, Video{
			Title: v["title"].(string),
			ID:    id,
			URL:   "https://www.youtube.com/watch?v=" + id,
		})
	}

	return videos, nil
}

// vomit code works?
func searchWithColly(search string) ([]Video, error) {
	//	fmt.Println("YOU ARE USING GOCOLLY FOR THIS SEARCH")
	query := "site:youtube.com " + search
	searchURL := "https://duckduckgo.com/html/?q=" + url.QueryEscape(query)

	c := colly.NewCollector(
		colly.AllowedDomains("duckduckgo.com", "html.duckduckgo.com"),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*duckduckgo.com*",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Referer", "https://duckduckgo.com/")
	})

	var videos []Video

	// Each result block
	c.OnHTML(".result", func(e *colly.HTMLElement) {
		title := e.ChildText(".result__a")
		rawLink := e.Request.AbsoluteURL(e.ChildAttr(".result__a", "href"))

		u, err := url.Parse(rawLink)
		if err != nil {
			return
		}

		realURL := u.Query().Get("uddg")
		if realURL == "" {
			realURL = rawLink
		}

		if title != "" && strings.Contains(realURL, "youtube.com") {
			videos = append(videos, Video{
				Title: title,
				URL:   realURL,
			})
		}
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("STATUS:", r.StatusCode)
		fmt.Println("BODY LENGTH:", len(r.Body))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	err := c.Visit(searchURL)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return videos, nil
}

func play(url string) {
	cmdmpv := exec.Command("mpv", url)
	//cmdmpv.Stdout = os.Stdout
	//cmdmpv.Stderr = os.Stderr
	cmdmpv.Start()
	cmdmpv.Wait()
}

func clear() {

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//fmt.Println(cmd.Start())
	//fmt.Println(cmd.Stdout)
	cmd.Run()

}

func logo() {
	fmt.Println("⠀⠀⠀⠀⠀⣀⣠⣤⣤⣤⣶⢶⡶⣶⣲⣶⢶⡶⣶⢶⣶⡶⣶⡶⣶⢶⣶⢶⡶⣶⡶⣶⢶⡶⣶⢶⡶⣶⢶⡶⣶⢶⣦⣶⢦⣤⣄⡀⠀⠀⠀⠀⠀\n" +
		"⠀⠀⠀⣴⣾⣟⣯⣟⣷⣻⡾⣿⡽⣷⣟⣾⢿⣽⣻⣟⣾⣽⡷⣟⣯⣿⢾⣻⣽⡷⣟⣯⣿⣻⣽⢿⣽⣻⣯⣟⣯⣿⢾⣽⣯⣟⣾⢿⡷⣆⠀⠀⠀\n" +
		"⠀⠀⣼⣟⣾⣽⢾⣻⣞⣯⢿⣳⣿⣻⣞⣯⣿⣞⡿⣞⡿⣾⡽⣿⡽⣾⣟⣯⡷⣿⣻⣽⣞⣯⣟⡿⣾⢷⣻⣾⣻⢾⣯⡷⣟⣾⢯⣿⡽⣟⣧⠀⠀\n" +
		"⠀⢀⣿⣼⢿⣼⡿⣟⣿⣻⢿⡿⣼⣧⢿⣻⣼⣻⢿⣿⣻⣧⡿⣟⣿⢧⣟⣧⣿⣟⣧⡿⣼⡿⣼⣿⣻⢿⡿⣼⣻⣿⣼⣻⣿⣻⢿⣼⢿⣟⣿⡀⠀\n" +
		"⠀⢸⣟⣾⢯⣷⢿⣯⡷⣿⢯⡿⣷⣻⢿⣽⣳⣿⣻⡾⣽⣳⡿⣯⣟⣯⡿⣽⡾⣽⣳⡿⣯⣟⣷⣯⣟⣯⣿⣻⢷⣻⣾⣽⢾⣻⣯⣟⡿⣾⢯⣇⠀\n" +
		"⠀⣼⡿⣽⣻⣽⣟⣾⣽⣟⣯⡿⣷⣻⣯⣟⣷⢯⣷⢿⣻⡽⠟⣷⣻⣽⣻⣽⢿⡽⣯⣟⣷⢿⣳⣟⣾⣿⣾⡽⣟⡿⣾⣽⣻⢷⣻⡾⣿⣽⣻⢿⠀\n" +
		"⠆⣿⡿⣽⣻⡾⣽⣾⣳⣯⢿⣽⡷⣟⣷⣻⡾⣿⣽⣻⢯⡇⠀⠈⠙⠳⢿⣿⣯⡿⣯⣟⣾⢿⣽⡾⣷⣻⡿⣟⡿⣽⡷⣯⣟⡿⣽⣻⢷⣯⣟⣿⠀\n" +
		"⠀⣿⣽⣟⡷⣿⣻⢾⣽⢯⣿⢾⣽⢿⣽⣳⣿⣳⣯⣟⡿⡇⠀⠀⠀⠀⠀⠉⠻⠽⣷⣻⣽⢿⡾⣽⡷⣿⣽⣻⣽⢿⡽⣟⣾⣟⡿⣽⣟⣾⣽⣾⠀\n" +
		"⡅⣿⢾⣽⣻⢷⣟⣯⣿⣻⢾⡿⣽⣻⡾⣟⣾⣽⣳⡿⣽⡇⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠹⢯⣿⢯⣟⣷⣯⣟⣾⣟⡿⣯⣷⣻⣽⣟⣾⣻⢾⣽⠀\n" +
		"⠁⣿⣏⣷⣿⣏⡿⣷⣏⡿⣏⣿⢿⣷⡿⣏⣷⣿⣹⢿⣹⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣹⡿⣿⡾⣷⣏⣷⣿⣹⣷⢿⣹⣾⣏⣷⣿⣏⣿⠀\n" +
		"⠀⣿⣞⡿⣾⣽⣻⣽⡾⣟⡿⣽⣻⣾⡽⣿⣽⣞⣯⡿⣯⡇⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⡾⣿⣻⣽⡷⣟⣯⣟⣷⢯⣷⣟⡿⣽⣳⣯⣟⣾⣽⣾⠂\n" +
		"⡀⣿⣾⣻⢷⣯⣟⣾⣽⢿⣽⣟⡷⣯⢿⣳⡿⣞⣯⡿⣷⡇⠀⠀⠀⠀⢀⣠⣶⡿⣿⣽⣻⢷⣟⡷⣿⣻⣽⡾⣯⣿⣳⣯⢿⣻⣽⡾⣯⡷⣟⣾⠀\n" +
		"⠃⣿⣳⣿⣻⡾⣽⣳⣯⣿⢾⣽⣻⣟⡿⣽⣻⢯⣷⢿⣳⡇⠀⣀⣴⣾⣻⣟⡷⣿⣳⣯⣟⡿⣾⣻⢷⣟⡷⣿⣽⣞⣯⣟⡿⣽⣳⣿⣳⣿⣻⣽⠀\n" +
		"⠀⢿⣻⣞⣷⢿⣻⣽⣳⣯⣿⣳⡿⣞⣿⢯⣟⡿⣽⣻⢯⣷⢾⣟⡿⣞⣷⣟⣿⣳⣿⣳⡿⣽⡷⣿⣻⡾⣿⣽⣞⣯⢿⣞⣿⣻⣽⣞⡿⣾⡽⣿⠀\n" +
		"⠀⢸⡿⣽⡾⣟⣯⣷⣟⡷⣯⣷⢿⣻⣽⣻⢯⣿⣻⣽⢿⣽⣻⡾⣟⣯⣷⢿⡾⣽⡾⣷⣟⣯⢿⣳⡿⣽⣷⣻⢾⣯⣟⡿⣾⣽⣳⡿⣽⡷⣿⡇⠀\n" +
		"⠀⠀⣿⣯⡟⣿⢳⡟⣾⣿⣽⣾⡟⣯⣷⣿⢻⣷⣯⡟⣿⡞⣯⣿⢻⣽⣾⢻⣽⣿⣽⣷⣯⡟⣿⣽⢻⣷⣯⡟⣿⣾⣽⢻⣷⣯⣿⣽⣯⡟⣷⠀⠀\n" +
		"⠀⠀⠙⣯⣿⡽⣿⡽⣟⣾⣽⣳⡿⣯⣟⣾⣟⣾⣽⣻⢷⣿⣻⣽⣻⣽⡾⣟⡿⣞⣷⣟⣾⢿⡽⣯⣿⢾⣽⣻⢷⣯⣟⡿⣾⣽⢾⣻⡾⣟⠏⠀⠀\n" +
		"⠀⠀⠀⠉⠺⣟⣷⢿⣻⣽⢾⣯⣟⣷⢿⣳⣯⣟⣾⢯⣿⢾⣽⣳⣿⣳⢿⣻⣽⣟⡷⣯⣟⣯⣿⣻⢾⣻⣽⢯⣿⢾⣽⣻⢷⣯⡿⣯⡿⠋⠀⠀⠀\n" +
		"⠀⠀⠀⠀⠀⠀⠈⠉⠉⠛⢛⣚⡙⢞⠛⠛⠚⠙⠫⠻⠽⠯⠿⠗⠯⠛⠿⠹⠷⠯⠟⠯⠟⠳⠏⠿⠿⠽⠚⢛⠙⠛⠓⠋⠋⠁⠉⠁⠀⠀⠀⠀⠀")
}
func searchDisplayPlayVideos(gocolly bool) {

	//var query string
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Search YouTube: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)
	var videos []Video
	var err error
	if gocolly == false {
		videos, err = search(query)
	} else {
		videos, err = searchWithColly(query)
	}
	//fmt.Println(videos)
	if err != nil {
		var answer string
		//fmt.Println("Error:", err)
		fmt.Println("Basically there was an error Searching for videos!")
		fmt.Println("Press whatever to return")
		fmt.Scanln(&answer)
		clear()
		return
		//return
	}

	fmt.Println("\nResults:\n")

	for i := 0; i <= len(videos)-1; i++ {

		fmt.Printf("[%d] %s\n", (i + 1), videos[i].Title)
	}
	var num int
	fmt.Println("Press any Number to select a video, or 0 to back out")
	fmt.Scanln(&num)
	if num == 0 {
		return
	} else {
		go play(videos[num-1].URL)
	}
	clear()

}
//if audio is true then do audio, else do video,
//command for downloading playlist of vids as mp3 files
//yt-dlp -t mp3 -o "%(playlist_index)s - %(title)s.%(ext)s" https://www.youtube.com/playlist?list=PLVEBmKtb8S7YRnpzrDO-s93E22xzjkg7G
func downloadVideo(audio, downloadPlaylist bool) {
	var fileType, vidurl string
	
	fmt.Println("Enter your video url")
	if downloadPlaylist == true {fmt.Println("example link: https://www.youtube.com/playlist?list=PLVEBmKtb8S7YRnpzrDO-s93E22xzjkg7G")}
	fmt.Scanln(&vidurl)
	if audio == false {
	fmt.Println("Choose a file type to download your video (or playlist) in \nSome file types to try:\nmp4\nmkv\nwebm\nmov\navi\nflv\nts\nm2ts\nmpeg\nogv")
	}
	if audio == true {
		fmt.Println("Choose a file type to download your video (or playlist) in \nSome file types to try:\nm4a\nmp3\nopus\nogg\naac\nwav\nflac\nwma")
	}	
	fmt.Scanln(&fileType)
	if downloadPlaylist == false {
		cmd := exec.Command("yt-dlp", "-t", fileType, vidurl)
		err := cmd.Run()
		if err != nil {log.Fatal(err)}
	}

	if downloadPlaylist == true {
		cmd := exec.Command("yt-dlp", "-t", fileType, `"%(playlist_index)s - %(title)s.%(ext)s"`, vidurl)
                err := cmd.Run()
                if err != nil {log.Fatal(err)}

	}

}

func downloadMenu() {

        clear()
        var ans string
downloadMenu:
        for {
		fmt.Println("Menu to download videos\n\n0: Quit\n1: Download video (single)\n2: Download Audio (single)\n3: Download videos (playlist)\n4: Download songs (playlist)")

                fmt.Scanln(&ans)

                switch ans {

                case "0":
			clear()
                        break downloadMenu
		case "1":
			//single video
			downloadVideo(false, false)
		case "2":
			//single audio
			downloadVideo(true, false)
		case "3":
			//playlist of videos
			downloadVideo(false, true)
		case "4":
			//playlist of songs
			downloadVideo(true, true)
		}	
	}
}

func killPort8080() {
	cmd := exec.Command("fuser", "-k", "8080/tcp")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Nothing was using port 8080.")
	} else {
		fmt.Println("Killed process using port 8080.")
	}
}

func launchServer() {
	killPort8080()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Could not get LAN IP")
		return
	}

	lanIP := "Unknown"

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			lanIP = ipnet.IP.String()
			break
		}
	}

	fmt.Println("========================================")
	fmt.Println(" YouTube Downloader GUI")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Open this address in your browser:")
	fmt.Println("  http://localhost:8080")

	if lanIP != "Unknown" {
		fmt.Printf("  http://%s:8080\n", lanIP)
	}

	fmt.Println()
	fmt.Println("To access the downloader from another")
	fmt.Println("device on your home network (such as")
	fmt.Println("your phone), you may need to allow")
	fmt.Println("incoming connections through your")
	fmt.Println("firewall permissions.")
	fmt.Println()
	fmt.Println("Press Q then ENTER to close the GUI.")
	fmt.Println("========================================")

	cmd := exec.Command("go", "run", "server.go")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		fmt.Println("Failed to launch server.go:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')

		if text == "q\n" || text == "Q\n" {
			fmt.Println("Closing downloader GUI...")
			cmd.Process.Kill()
			//cmd.Process.Signal(os.Interrupt)
			cmd.Wait()
			killPort8080()
			clear()
			break
		}
	}
}

func main() {
	clear()
	var ans string
main:
	for {
		logo()
		fmt.Println("This is a menu for my youtube client!\n0: Quit\n1: Search youtube using yt-dlp\n2: Search youtube using gocolly\n3: Adjust Search Amount\n4: Play youtube video directly from URL\n5: Download menu\n6: Browser based youtube downloader")

		fmt.Scanln(&ans)

		switch ans {

		case "0":
			break main
		case "1":
			searchDisplayPlayVideos(false)
		case "3":
			fmt.Println("Example search amount 10: ")
			fmt.Scanln(&searchAmount)
			clear()
		case "4":
			var youtubeUrl string
			fmt.Println("Paste youtube url here: ")
			fmt.Scanln(&youtubeUrl)
			go play(youtubeUrl)
		case "2":
			searchDisplayPlayVideos(true)
		case "5":
			downloadMenu()
		case "6":
			launchServer()
		}
	}

}
