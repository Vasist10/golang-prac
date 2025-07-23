package main 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	PhotoApi = "https://api.pexels.com/v1/search"
	VideoApi = "https://api.pexels.com/videos/search"
)

type Client struct {
	Token string
	hc   	http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token:token,hc:c}
}

type SearchResult struct {
	Page int `json:"page"`
	PerPage int `json:"per_page"`
	TotalResults int `json:"total_results"`
	NextPage string `json:"next_page"`
	Photos []Photo `json:"photos"`
}

type Photo struct {
	ID int `json:"id"`
	Width int `json:"width"`
	Height int `json:"height"`
	URL string `json:"url"`
	Photographer string `json:"photographer"`
	PhotographerUrl string `json:"photographer_url"`
	Src PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original string `json:"original"`
	Large2x string `json:"large2x"`
	Large string `json:"large"`
	Medium string `json:"medium"`
	Small string `json:"small"`
	Portrait string `json:"portrait"`
	Landscape string `json:"landscape"`
	Square string `json:"square"`
	Tiny string `json:"tiny"`
}

func (c *Client) SearchPhotos(query string, perPage int, page int) (*SearchResult, error) {

	fmt.Sprintf("Searching photos for query: %s, page: %d, perPage: %d\n", query, page, perPage)
	resp ,err := c.requestDoWithAuth("GET",url)
	defer resp.Body.Close()
	data,err :=io.util.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var result SearchResult
	err = json.Unmarshal(data, &result); 
	return &result, nil
}

func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", c.Token)
	resp, err := c.hc.Do(req)
	if err != nil {
		return resp, err 
	}
	times , err := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	if err != nil {
		return resp,nil
	}else{
		c.RemainingTimes = int32(times)
	}
	return resp, nil
}
func main(){
	os.Setenv("PEXELS_API_KEY")
	var TOKEN  = os.Getenv("PEXELS_API_KEY")

	var c = NewClient(TOKEN)

	result, err := c.SearchPhotos("nature", 1, 10)
	if err != nil {
		fmt.Println("Error searching photos:", err)
		return
	}
	if result.Page == 0 {
		fmt.Println("No photos found")
		return
	}
	fmt.Println(result)


}