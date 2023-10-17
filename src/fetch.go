package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Artist struct {
	Id              int          `json:"id"`
	Image           string       `json:"image"`
	Name            string       `json:"name"`
	Members         []string     `json:"members"`
	CreationDate    int          `json:"creationDate"`
	FirstAlbum      string       `json:"firstAlbum"`
	LocationsURL    string       `json:"locations"`
	ConcertDatesURL string       `json:"concertDates"`
	RelationsURL    string       `json:"relations"`
	Locations       LocationData `json:"-"`
	ConcertDates    ConcertDates `json:"-"`
	Relations       RelationData `json:"-"`
}

type LocationData struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type ConcertDates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type RelationData struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(target)
}

func GetArtist() ([]Artist, error) {
	// Fetch the artist data
	artistResponse, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer artistResponse.Body.Close()

	// Decode the artist JSON response into a slice of Artist
	var artists []Artist
	if err := json.NewDecoder(artistResponse.Body).Decode(&artists); err != nil {
		return nil, err
	}

	// Fetch additional data for each artist
	for i, artist := range artists {
		if err := fetchData(artist.LocationsURL, &artists[i].Locations); err != nil {
			return nil, err
		}
		if err := fetchData(artist.ConcertDatesURL, &artists[i].ConcertDates); err != nil {
			return nil, err
		}
		if err := fetchData(artist.RelationsURL, &artists[i].Relations); err != nil {
			return nil, err
		}
	}

	return artists, nil
}

func Compare(start, end, input string) bool {
	layout1 := "02-01-2006"
	layout2 := "2006-01-02"
	dateInput, err1 := time.Parse(layout1, input)
	if err1 != nil {
		fmt.Println("Error parsing dateStr1:", err1)
		return false
	}
	dateStart, err2 := time.Parse(layout2, start)
	if err2 != nil {
		fmt.Println("Error parsing Start:", err2)
		return false
	}
	dateEnd, err2 := time.Parse(layout2, end)
	if err2 != nil {
		fmt.Println("Error parsing End:", err2)
		return false
	}
	if dateInput.After(dateStart) && dateInput.Before(dateEnd) {
		return true
	}
	return false
}

func Compare1(start, input string) bool {
	layout1 := "02-01-2006"
	layout2 := "2006-01-02"
	dateInput, err1 := time.Parse(layout1, input)
	if err1 != nil {
		fmt.Println("Error parsing dateStr1:", err1)
		return false
	}
	dateStart, err2 := time.Parse(layout2, start)
	if err2 != nil {
		fmt.Println("Error parsing Start:", err2)
		return false
	}
	if dateStart.Before(dateInput) {
		return true
	}
	return false
}
