package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	// we have to put the temperature inside the main struct as it is defined inside the main body in the openweathermap api
	Main struct {
		Celsius        float64 `json:"temp"`
		Celsiusexp     float64 `json:"-"`
		Pressure       float64 `json:"pressure"`
		Humidity       float64 `json:"humidity"`
		MinTemperature float64 `json:"-"`
		MinTemp        float64 `json:"temp_min"`
		MaxTemperature float64 `json:"-"`
		MaxTemp        float64 `json:"temp_max"`
		// now I had to do it as it gives me the value in kelvin
		// the - tells the encoder and decoder to ignore the field
		// Celsius float64 `json:"temperature"` this is also wrong json as OWM API is giving the json as temp and not temperature
		// always remeber that when we make it in json like json: "temp" there should not be any space between : and ""
	} `json:"main"`
	Rain struct {
		Rain_Volume_for_last_one_hour float64 `json:"rain.1h"`
	} `json:"rain"`
	// Kelvin float64 `json: "temp"`
	// so this is a wrong thing as temperature must go inside the main struct
}

// the below function is to get apikey from our .apiConfig
func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("getting a fucking error")
		return apiConfigData{}, err
		// if there is error in the filename then we return the error with no configdata
	}
	var c apiConfigData
	// if we can access the data we can store it in var c
	err = json.Unmarshal(bytes, &c)
	// json.Unmarshal takes two fields one the data that is to be decoded and a pointer to the variable where it will be stored
	// now we are converting it to golang and if throws error we say error
	// if we are unmarshalling a data which is not a valid json we can get an error
	// or the structure does not match with the structure of GoLang
	// json.Unmarshal can throw additional errors than marshal when the destination variable is not a pointer or if the JSON data contains more fields than the GoLang can handle

	if err != nil {
		fmt.Println("getting a fucking error")
		return apiConfigData{}, err
	}
	return c, nil

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go!\n"))
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	// first we are loading the apiconfig and storing it in the apiConfig
	if err != nil {
		// fmt.Println("getting a fucking error")
		return weatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		// fmt.Println("getting a fucking error")
		// if()
		return weatherData{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return weatherData{}, fmt.Errorf("city not available")
	}
	defer resp.Body.Close()
	// the request isn't closed itself we have to close it
	var d weatherData
	// in the json.NewDecoder(resp.Body) creates a NewDecoder that reads the response from the resp Body (as it is in JSON format) of the HTTP GET request. The body contains the response body which represents a JSON string with the weather dara
	// Now .Decode(&d) we are using the & so that the decoder can write the decoded data directly to the address of the 'd' variable
	// .Decode decodes the JSON data from response  body into Go struct represented by variable 'd'
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		// fmt.Println("Getting a fucking error")
		return weatherData{}, err
	}
	// d.Main.Celsius = d.Main.Celsius - 273.15
	// d.Main.MaxTemp = d.Main.MaxTemp - 273.15
	// d.Main.MinTemp = d.Main.MinTemp - 273.15
	// d.Main.Celsiusexp = d.Main.Celsius - 273.15
	// d.Main.MaxTemperature = d.Main.MaxTemp - 273.15
	// d.Main.MinTemperature = d.Main.MinTemp - 273.15
	// in this case the error occured as we are telling to ignore while encoding and decoding so we can directly convert the main part
	return d, nil
}

func main() {
	// json.NewEncoder(w)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		// in this we want that the person can go to /weather/ and then after that write the city name to get the weather of that city
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			// in this we are storing the third thing which comes after slash
			// for example if my site is http://localhost:8080/weather/london
			// the r.URL.Path is the above thing
			// the "/" is used to tell ki kha se split ho rha h
			// 3 implies the third one
			// we get an array like "http://localhost:8080","weather","london"
			// so by 3 we are storing london in our city variable
			data, err := query(city)
			if err != nil {
				// fmt.Println("getting a fucking error")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// w HTTP.responsewriter that is used to send the data back to the user
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			// we are telling that the response is in json format and the character encoding will be done in utf-8
			json.NewEncoder(w).Encode(data)
		})
	fmt.Println("Starting the server ...")
	http.ListenAndServe(":8080", nil)
	fmt.Println("Server started...")
}
