package funcs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	ErrorBad          = errors.New("error: BadRequset")
	ErrorInternal     = errors.New("error: InternalServer")
	ErrorPageNotFound = errors.New("error: Page Not Found")
)

type AllArtists []Artist

type Artist struct {
	Id                  int      `json:"id"`
	Image               string   `json:"image"`
	Name                string   `json:"name"`
	Members             []string `json:"members"`
	CreationDate        int      `json:"creationDate"`
	FirstAlbum          string   `json:"firstalbum"`
	Relations           string   `json:"relations"`
	RelationsData       Relation
	RangeCreationDate   MinMax
	RangeFirstAlbumDate MinMax
}

type MinMax struct {
	Min int
	Max int
}

type Relation struct {
	Id             int
	DatesLocations map[string][]string
}

type OnlyRelations struct {
	Index []Relation `json:"index"`
}

type ForTemplate struct {
	AllArtist    []Artist
	SearchArtist []Artist
}

const (
	pageArtists  = "https://groupietrackers.herokuapp.com/api/artists"
	pageRelation = "https://groupietrackers.herokuapp.com/api/relation"
)

func MakeAllArtists() (AllArtists, error) {
	var result AllArtists

	res, err := http.Get(pageArtists)
	if exitError(err) {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if exitError(err) {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if exitError(err) {
		return nil, err
	}

	return result, nil
}

func MakeOneArtist(id string) (Artist, error) {
	var result Artist

	err := checkId(id)
	if err != nil {
		return Artist{}, ErrorPageNotFound
	}

	res, err := http.Get(pageArtists + "/" + id)
	if exitError(err) {
		return Artist{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if exitError(err) {
		return Artist{}, err
	}

	err = json.Unmarshal(body, &result)
	if exitError(err) {
		return Artist{}, err
	}
	rel, err := http.Get(result.Relations)
	if exitError(err) {
		return Artist{}, err
	}

	bodyrel, err := ioutil.ReadAll(rel.Body)
	if exitError(err) {
		return Artist{}, err
	}

	err = json.Unmarshal(bodyrel, &result.RelationsData)

	return result, nil
}

func exitError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func checkId(id string) error {
	number, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return ErrorPageNotFound
	}

	if number < 1 || number > 52 {
		return ErrorPageNotFound
	}

	return nil
}
