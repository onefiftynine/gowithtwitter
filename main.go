package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"net/http"
	"os"
)

type TwitterApiKeys struct {
	ConsumerKey       string `json:"consumerkey"`
	ConsumerSecret    string `json:"consumersecret"`
	AccessToken       string `json:"accesstoken"`
	AccessTokenSecret string `json:"accesstokensecret"`
}

var (
	apikeys TwitterApiKeys
	api     *anaconda.TwitterApi
)

// read config file
func configure() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&apikeys)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(apikeys)
}

func main() {
	configure()
	anaconda.SetConsumerKey(apikeys.ConsumerKey)
	anaconda.SetConsumerSecret(apikeys.ConsumerSecret)
	// searchResult, _ := api.GetSearch("golang", v)
	// for _, tweet := range searchResult.Statuses {
	// 	fmt.Println(tweet.Text)
	// }
	//	availableLocations()
	getCurrentTrends()

	// ser go cron to run every hour
	gocron.Every(1).Hour().Do(getCurrentTrends)
	<-gocron.Start()
}

//Joke data structure for joke, need to change it ehen we add more joke source websites
type Joke struct {
	JokeSentence string `json:"joke"`
}

// get joke TODO add more joke sources
func getRandomJoke() string {
	var j Joke
	resp, err := http.Get("http://tambal.azurewebsites.net/joke/random")
	if err != nil {
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &j)
	fmt.Println(j.JokeSentence)
	return j.JokeSentence + "  "
}

func init() {
	api = anaconda.NewTwitterApi(apikeys.AccessToken, apikeys.AccessTokenSecret)

}

// get current trends and post to twitter
//23424848 code for india, to get code  use availableLocations
func getCurrentTrends() {
	trendResponse, _ := api.GetTrendsByPlace(23424848, nil)
	for _, trend := range trendResponse.Trends {
		fmt.Println(trend.Name)
		ch := make(chan string)
		go post(getRandomJoke()+trend.Name, ch)
		fmt.Println("created at", <-ch)
	}
}

//post
func post(status string, ch chan string) {
	fmt.Println("posting")
	tweet, err := api.PostTweet(status, nil)
	if err != nil {
		fmt.Println("error", err)
		ch <- "error"
	}
	ch <- tweet.CreatedAt
}

func availableLocations() {
	trendLocations, _ := api.GetTrendsAvailableLocations(nil)
	for _, location := range trendLocations {
		fmt.Println("Name:  ", location.Name)
		fmt.Println("Woeid", location.Woeid)
		// fmt.PrintLn(location.PlaceType.code)
	}
}
