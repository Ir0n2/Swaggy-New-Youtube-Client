package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	//out, err := cmd.Output()
	//if err != nil {
	//fmt.Println(out)
	//	return nil, err
	//}

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("The program has encountered some kind of problem that I am too high to deal with.\nProbably with yt-dlp, fat chance it either can't find the video you're searching for or can't find js run time on the page. \nIt could be alot really, rate limiting, bad network conncection, youtube changed something, you're using a bad or old version of yt-dlp, there's an issue with the program.\n It could also be bad formating. I suggest you reword your search and try again. Fuck I'm tired.")
		if showErrorOutput == true {
			fmt.Println("Details:", string(out))
		}
		return nil, err
	}

	//fmt.Println(string(out))

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

func main() {
	clear()
	var ans string
main:
	for {
		logo()
		fmt.Println("This is a menu for my youtube client!\n0: Quit\n1: Search youtube using yt-dlp\n2: Adjust Search Amount\n3: Play youtube URL")

		fmt.Scanln(&ans)

		switch ans {

		case "0":
			break main
		case "1":
			//var query string
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Search YouTube: ")
			query, _ := reader.ReadString('\n')
			query = strings.TrimSpace(query)

			videos, err := search(query)
			if err != nil {
				//fmt.Println("Error:", err)
				fmt.Println("Basically there was an error Searching for videos!")
				fmt.Println("Press whatever to return")
				fmt.Scanln(&ans)
				clear()
				continue main
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
				continue main
			} else {
				go play(videos[num-1].URL)
			}
			clear()
		case "2":
			fmt.Println("Example search amount 10: ")
			fmt.Scanln(&searchAmount)
			clear()
		case "3":
			var youtubeUrl string
			fmt.Println("Paste youtube url here: ")
			fmt.Scanln(&youtubeUrl)
			go play(youtubeUrl)
		}
	}

}
