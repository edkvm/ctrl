package main

import (
	"encoding/json"
	"fmt"
	"github.com/edkvm/ctrl/gowrapper"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)



func buildUrl(city, unitSys, apiKey string) string {
	return fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%v&units=%v&appid=%v", city, unitSys, apiKey)
}

func Api(ep string) []byte {

	c := http.Client{
		Timeout: 30*time.Second,
	}

	res, err := c.Get(ep)
	if err != nil {

	}
	defer res.Body.Close()

	if res.StatusCode != 200 {

	}

	data, err := ioutil.ReadAll(res.Body)

	return data
}

type Params struct{
	City string `json:"city"`
}

func GetWeather(params Params) (string, error) {

	apiKey := os.Getenv("API_KEY")
	unitSys := os.Getenv("DEFAULT_UNIT_SYS")

	url := buildUrl(params.City, unitSys, apiKey)

	raw := Api(url)

	var data struct {
		Name string `json:"name"`
		Main map[string]interface{} `json:"main"`
	}

	err := json.Unmarshal(raw, &data)
	if err != nil {

	}

	return fmt.Sprintf("It's %v degrees in %v", data.Main["temp"], data.Name), nil
}

func main() {
	gowrapper.Start(GetWeather)
}


