package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	fetch "./src"
)

type Location struct {
	Lat   float64 `json:"Lat"`
	Lng   float64 `json:"Lng"`
	Name  string  `json:"Name"`
	Dates string  `json:"Dates"`
}

var (
	loca []Location
	Data []fetch.Artist
	loc  []string
	errr error
)

const port = ":8080"

func main() {
	Initiate()
	if errr != nil {
		log.Fatal("Error fetching artist data:", errr)
	}
	mux := http.NewServeMux()

	mux.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style"))))
	mux.HandleFunc("/", Index)
	mux.HandleFunc("/artist/", ArtistP)
	mux.HandleFunc("/search", Search)
	mux.HandleFunc("/filter", Filter)
	http.HandleFunc("/getdata", func(w http.ResponseWriter, r *http.Request) {
		// Marshal the struct into JSON
		jsonData, err := json.Marshal(loca)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response content type and send the JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonData)
	})

	fmt.Printf("(http://localhost%v): Server started on port\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// func Geo(w http.ResponseWriter, r *http.Request) {
// 	// Replace "YourApiKeyHere" with your actual Google Maps API key
// 	apiKey := "AIzaSyAbgRjHe8P_pbCsi7dJkaTJqTll3R01ogo"

// 	// Simulated location names, replace with your actual location names
// 	var locationsToGeocode []string = Data[id-1].Locations.Locations

// 	// Convert location names to coordinates using the Google Maps Geocoding API
// 	locations := []Location{}
// 	for _, locationName := range locationsToGeocode {
// 		lat, lng, err := GeocodeLocation(locationName, apiKey)
// 		if err != nil {
// 			log.Printf("Error geocoding location %s: %v", locationName, err)
// 			continue
// 		}
// 		locations = append(locations, Location{Lat: lat, Lng: lng, Name: locationName})
// 	}
// 	fmt.Println(locations)

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(locations)
// }

func Initiate() {
	Data, errr = fetch.GetArtist()
	var zebi []string
	for _, art := range Data {
		for _, memb := range art.Locations.Locations {
			zeb := strings.Replace(memb, "_", " ", -1)
			klwa := strings.Split(zeb, "-")
			zeboklwa := strings.Title(klwa[0]) + ", " + strings.ToUpper(klwa[1])
			zebi = append(zebi, zeboklwa)

		}
	}
	for i := range zebi {
		duplicate := false
		for j := range loc {
			if zebi[i] == loc[j] {
				duplicate = true
				break
			}
		}
		if !duplicate {
			loc = append(loc, zebi[i])
		}
	}
	sort.Strings(loc)
}

func Index(w http.ResponseWriter, r *http.Request) {
	k := template.Must(template.ParseFiles("./templates/404.html"))
	if r.URL.Path != "/" {
		k.Execute(w, struct {
			Value []fetch.Artist
			Zebi  []string
		}{Data, loc})
		return
	}
	t := template.Must(template.ParseFiles("./templates/index.html"))

	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Println(loca)
	t.Execute(w, struct {
		Succes bool
		Value  []fetch.Artist
		Zebi   []string
	}{true, Data, loc})
}

func ArtistP(w http.ResponseWriter, r *http.Request) {
	// Extract the dynamic part of the URL
	t := template.Must(template.ParseFiles("./templates/htmllll.html"))
	k := template.Must(template.ParseFiles("./templates/404.html"))
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) == 3 {
		id, err := strconv.Atoi(pathParts[2])

		if id > len(Data) || err != nil {
			k.Execute(w, struct {
				Value []fetch.Artist
				Zebi  []string
			}{Data, loc})
			return
		}
		apiKey := "AIzaSyAbgRjHe8P_pbCsi7dJkaTJqTll3R01ogo"

		// Simulated location names, replace with your actual location names
		var locationsToGeocode []string = Data[id-1].Locations.Locations

		// Convert location names to coordinates using the Google Maps Geocoding API
		locations := []Location{}
		for _, locationName := range locationsToGeocode {
			lat, lng, err := GeocodeLocation(locationName, apiKey)
			if err != nil {
				log.Printf("Error geocoding location %s: %v", locationName, err)
				continue
			}
			locations = append(locations, Location{Lat: lat, Lng: lng, Name: locationName, Dates: strings.Join(Data[id-1].Relations.DatesLocations[locationName], ", ")})
		}
		stringify, _ := json.Marshal(locations)

		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(locations)
		fmt.Println(string(stringify))
		t.Execute(w, struct {
			Succes  bool
			Avalues []fetch.Artist
			Value   fetch.Artist
			Zebi    []string
			Lase9   string
		}{true, Data, Data[id-1], loc, string(stringify)})

	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	// Extract the search query from the URL query parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	found := false
	// Perform a search operation based on the query
	var searchResults []fetch.Artist
	for _, artist := range Data {

		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			searchResults = append(searchResults, artist)
			found = true
			continue
		}
		i, errrr := strconv.Atoi(query)
		if i == artist.CreationDate && errrr == nil {
			searchResults = append(searchResults, artist)
			found = true
			continue
		}
		for _, loc := range artist.Locations.Locations {
			if strings.Contains(strings.ToLower(loc), strings.ToLower(query)) {
				searchResults = append(searchResults, artist)
				found = true
				continue
			}
		}

		if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(query)) {
			searchResults = append(searchResults, artist)
			found = true
			continue
		}

		for _, members := range artist.Members {
			if strings.Contains(strings.ToLower(members), strings.ToLower(query)) {
				searchResults = append(searchResults, artist)
				found = true
				continue
			}
		}
	}

	// Render the search results in your HTML template
	t := template.Must(template.ParseFiles("./templates/search_results.html"))
	t.Execute(w, struct {
		Success bool
		Query   string
		Results []fetch.Artist
		Zebi    []string
		Value   []fetch.Artist
	}{found, query, searchResults, loc, Data})
}

func Fix(s string) string {
	s = strings.Replace(s, ", ", "-", -1)
	s = strings.Replace(s, " ", "_", -1)
	s = strings.ToLower(s)
	return s
}

func Filter(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("first-album-start")
	end := r.FormValue("first-album-end")
	members := r.Form["Members"]
	location := r.FormValue("Country")
	creationstart, _ := strconv.Atoi(r.FormValue("minC"))
	creationend, _ := strconv.Atoi(r.FormValue("maxC"))

	found := false
	var filterresults, creationresult, membersresult, locationresult []fetch.Artist
	for _, artist := range Data {
		if start != "" {
			if fetch.Compare(start, end, artist.FirstAlbum) {
				filterresults = append(filterresults, artist)
			}
		} else {
			filterresults = Data
		}
	}

	if location == "" {
		locationresult = filterresults
	} else {
		location = Fix(location)
		for _, artist := range filterresults {
			for _, loccccc := range artist.Locations.Locations {
				if loccccc == location {
					locationresult = append(locationresult, artist)
					break
				}
			}
		}
	}
	if members == nil {
		membersresult = locationresult
	} else {
		for _, artist := range locationresult {
			for _, m := range members {
				mm, _ := strconv.Atoi(m)
				if mm == len(artist.Members) {
					membersresult = append(membersresult, artist)
					break
				}
			}
		}
	}
	for _, artist := range membersresult {
		if artist.CreationDate >= creationstart && artist.CreationDate <= creationend {
			creationresult = append(creationresult, artist)
		}
	}

	if creationresult != nil {
		found = true
	}
	t := template.Must(template.ParseFiles("./templates/search_results.html"))
	t.Execute(w, struct {
		Success bool
		Query   string
		Results []fetch.Artist
		Zebi    []string
		Value   []fetch.Artist
	}{found, "zebi", creationresult, loc, Data})
}

func GeocodeLocation(locationName, apiKey string) (float64, float64, error) {
	// Make a request to the Google Maps Geocoding API
	resp, err := http.Get("https://maps.googleapis.com/maps/api/geocode/json?address=" + locationName + "&key=" + apiKey)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// Parse the response
	var geocodeResponse struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&geocodeResponse)
	if err != nil {
		return 0, 0, err
	}

	// Check if any results were found
	if len(geocodeResponse.Results) == 0 {
		return 0, 0, nil // No results found
	}

	// Extract the coordinates from the first result
	lat := geocodeResponse.Results[0].Geometry.Location.Lat
	lng := geocodeResponse.Results[0].Geometry.Location.Lng
	fmt.Println(lat, lng)
	return lat, lng, nil
}
