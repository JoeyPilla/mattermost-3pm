package main

import (
	"fmt"
	"github.com/joho/godotenv"
	mattermost "github.com/mattermost/mattermost-server/model"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	Client := initMattermost()

	userId := goDotEnvVariable("USERID")
	channelId := "mfje975pe386jbkawdye6pu76r"

	waitForNextWholeHour()

	called := false
	for {
		end := getEndTime()
		timeUntil := (time.Until(end) + -1*time.Hour).Round(time.Minute)
		str := fmtDuration(timeUntil)
		if time.Now().Weekday() == 5 && !called {
			str = "https://pbs.twimg.com/ext_tw_video_thumb/1246049815802212357/pu/img/HyOu5VcjHNXrk4rI.jpg"
			called = true
		}

		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		images := []string{
			"https://static.wikia.nocookie.net/spongebob/images/5/54/Wet_Painters_066.png/revision/latest?cb=20191215185631",
			"https://static.wikia.nocookie.net/spongebob/images/8/8d/Patrick%27s_Staycation_110.png/revision/latest?cb=20150717073046",
			"https://static.wikia.nocookie.net/spongebob/images/8/81/Mall_Girl_Pearl_179.png/revision/latest?cb=20191014034123",
			"https://static.wikia.nocookie.net/spongebob/images/7/76/Mermaid_Pants_156.png/revision/latest?cb=20200114023022",
			"https://static.wikia.nocookie.net/spongebob/images/8/83/The_Goofy_Newbie_074.png/revision/latest?cb=20190926203104",
			"https://static.wikia.nocookie.net/spongebob/images/e/e8/Pat_the_Dog_077.png/revision/latest?cb=20210713193621",
		}

		rng := r1.Intn(len(images))
		img := images[rng]

		if timeUntil.Hours() > -1 && timeUntil.Hours() < 8.0 {
			postMessage(Client, userId, channelId, img)
			postMessage(Client, userId, channelId, str)
		}

		waitForNextWholeHour()
	}
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func postMessage(client *mattermost.Client4, userId string, channelId string, message string) {
	newPost := mattermost.Post{
		UserId:    userId,
		ChannelId: channelId,
		Message:   message,
	}
	p, r := client.CreatePost(&newPost)

	if r.Error != nil {
		log.Fatal("Couldn't make post: ", r.Error.Message)
	}

	log.Print("Post created: ", p.Message)

}

func fmtDuration(d time.Duration) string {

	d = d.Round(time.Second)
	h, d := handleTime(d, time.Hour)
	m, d := handleTime(d, time.Minute)
	s, d := handleTime(d, time.Second)

	return fmt.Sprintf("There are %02d hours, %02d minutes, and %02d seconds until 3:00PM.", h, m, s)
}

func handleTime(d time.Duration, timeType time.Duration) (time.Duration, time.Duration) {
	t := d / timeType
	d -= t * timeType

	return t, d
}

func initMattermost() *mattermost.Client4 {
	Client := mattermost.NewAPIv4Client("https://mattermost.t3")
	Client.Login(goDotEnvVariable("USERNAME"), goDotEnvVariable("PASSWORD"))

	return Client
}

func getEndTime() time.Time {
	now := time.Now()
	date := now.Format("2006-01-02")

	endTimeString := fmt.Sprintf("%s 3:00pm (EST)", date)

	endTime, _ := parseTime(endTimeString)

	return endTime
}

func parseTime(t string) (time.Time, error) {
	parse := "2006-01-02 3:04pm (MST)"

	return time.Parse(parse, t)
}

func waitForNextWholeHour() {
	_, m, s := getCurrentHoursMinutesSeconds()
	waitTime := getWaitTime(0, m, s)

	time.Sleep(waitTime)
}

func getWaitTime(h, m, s int) time.Duration {
	return getTimeUntil(h, m, s)
}

func getTimeUntil(h int, m int, s int) time.Duration {
	hUntil := getHoursUntil(h)
	mUntil := getMinutesUntil(m)
	sUntil := getSecondsUntil(s)

	return time.Duration(hUntil + mUntil + sUntil)
}

func getHoursUntil(h int) int {
	hoursInDay := 24
	hoursUntil := (hoursInDay - 1) - (h%hoursInDay)*int(time.Hour)
	return hoursUntil
}

func getMinutesUntil(m int) int {
	minutesInHours := 60
	minutesUntil := (minutesInHours - m%minutesInHours - 1) * int(time.Minute)
	return minutesUntil
}

func getSecondsUntil(s int) int {
	secondsInMinute := 60
	secondsUntil := (secondsInMinute - s) * int(time.Second)
	return secondsUntil
}

func getCurrentHoursMinutesSeconds() (int, int, int) {
	now := time.Now()

	hours := now.Hour()
	minutes := now.Minute()
	seconds := now.Second()

	return hours, minutes, seconds
}
