package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type WeatherAPIResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

type WeatherResponse struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
}

type IPGeoResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		if ip != "127.0.0.1" && ip != "::1" {
			return ip
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "127.0.0.1" || ip == "::1" {
		resp, err := http.Get("https://api.ipify.org?format=text")
		if err == nil {
			buf := new(strings.Builder)
			_, err = io.Copy(buf, resp.Body)
			if err == nil {
				return strings.TrimSpace(buf.String())
			}
		}
	}
	return ip
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	apikey := os.Getenv("WEATHER_API_KEY")
	if apikey == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		return
	}

	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")
	city := r.URL.Query().Get("city")

	var query string
	if city != "" {
		query = city
	} else if lat != "" && lon != "" {
		query = fmt.Sprintf("%s,%s", lat, lon)
	} else {
		ip := getIP(r)
		resp, err := http.Get("http://ip-api.com/json/" + ip)
		if err != nil {
			http.Error(w, "Unable to fetch IP-based location", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var ipData IPGeoResponse
		if err := json.NewDecoder(resp.Body).Decode(&ipData); err != nil {
			http.Error(w, "Failed to decode IP-based location", http.StatusInternalServerError)
			return
		}
		query = fmt.Sprintf("%f,%f", ipData.Lat, ipData.Lon)
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apikey, query)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var check map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&check); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}
	if errData, ok := check["error"]; ok {
		errMsg := errData.(map[string]interface{})["message"].(string)
		http.Error(w, "Weather API error: "+errMsg, http.StatusBadRequest)
		return
	}

	bodyBytes, _ := json.Marshal(check)
	var apiResp WeatherAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		http.Error(w, "Failed to decode structured weather", http.StatusInternalServerError)
		return
	}

	weather := WeatherResponse{
		Location:    apiResp.Location.Name,
		Temperature: apiResp.Current.TempC,
		Description: apiResp.Current.Condition.Text,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(weather)
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	http.HandleFunc("/weather", getWeather)

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	fmt.Println("Server started successfully")
}
