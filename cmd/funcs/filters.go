package funcs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func CreationDateMinmax(artists AllArtists) {
	min := artists[0].CreationDate
	max := artists[0].CreationDate
	for i := range artists {
		if min > artists[i].CreationDate {
			min = artists[i].CreationDate
		}
		if max < artists[i].CreationDate {
			max = artists[i].CreationDate
		}
	}
	for i := range artists {
		artists[i].RangeCreationDate.Min = min
		artists[i].RangeCreationDate.Max = max
	}
}

func FirstAlbumMinMax(artists AllArtists) {
	min := artists[0].CreationDate
	max := artists[0].CreationDate
	for i := range artists {
		dates := strings.Split(artists[i].FirstAlbum, "-")
		year, err := strconv.Atoi(dates[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		if min > year {
			min = year
		}
		if max < year {
			max = year
		}
	}
	for i := range artists {
		artists[i].RangeFirstAlbumDate.Min = min
		artists[i].RangeFirstAlbumDate.Max = max
	}
}

func ChooseMembers(artists AllArtists, members []string) (AllArtists, error) {
	if members == nil {
		return artists, nil
	}

	var result AllArtists
	val, err := checkValuesMembers(members)
	if err != nil {
		return result, err
	}

	for _, number := range val {
		if number < 1 {
			return result, ErrorBad
		}
		for i, w := range artists {
			if number == len(artists[i].Members) {
				result = append(result, w)
			}
		}
	}
	return result, nil
}

func checkValuesMembers(values []string) ([]int, error) {
	if len(values) < 1 {
		return nil, ErrorPageNotFound
	}

	result := []int{}
	for _, w := range values {
		number, err := strconv.Atoi(w)
		if err != nil {
			return nil, ErrorBad
		}
		result = append(result, number)
	}
	return result, nil
}

func ChooseCreationDate(artists AllArtists, creationDate []string) (AllArtists, error) {
	var result AllArtists

	val, err := checkValues(creationDate)
	if err != nil {
		return result, err
	}

	for i := range artists {
		if compareDates(artists[i].CreationDate, val[0], val[1]) {
			result = append(result, artists[i])
		}
	}

	return result, nil
}

func ChooseFirstAlbum(artists AllArtists, furstAlbumDate []string) (AllArtists, error) {
	var result AllArtists

	val, err := checkValues(furstAlbumDate)
	if err != nil {
		return result, err
	}

	for i := range artists {
		dates := strings.Split(artists[i].FirstAlbum, "-")
		year, err := strconv.Atoi(dates[2])
		if err != nil {
			return result, err
		}
		if compareDates(year, val[0], val[1]) {
			result = append(result, artists[i])
		}
	}

	return result, nil
}

func checkValues(values []string) ([]int, error) {
	if values == nil {
		return nil, ErrorPageNotFound
	}

	result := []int{}

	for _, w := range values {

		number := strings.Split(w, " - ")

		if len(number) != 2 {
			return nil, ErrorBad
		}

		number1, err := strconv.Atoi(number[0])
		if err != nil {
			return nil, ErrorBad
		}
		number2, _ := strconv.Atoi(number[1])
		if err != nil {
			return nil, ErrorBad
		}
		result = append(result, number1)
		result = append(result, number2)
	}

	if result[1] < result[0] {
		return nil, ErrorBad
	}

	return result, nil
}

func compareDates(inStuctDate int, valueDateBegin int, valueDateEnd int) bool {
	if inStuctDate >= valueDateBegin && inStuctDate <= valueDateEnd {
		return true
	}
	return false
}

func ChooseLocations(artists AllArtists, location []string) (AllArtists, error) {
	var saveword OnlyRelations
	var ids []int

	for _, w := range location {
		location = strings.Split(w, ",")
	}

	relations, err := MakeOnlyRelation()
	if err != nil {
		return artists, ErrorInternal
	}

	for _, word := range location {
		if saveword.Index != nil {
			relations.Index = saveword.Index
		}
		for i := range relations.Index {
			for key := range relations.Index[i].DatesLocations {
				if strings.Contains(strings.ToLower(key), strings.ToLower(word)) {
					saveword.Index = append(saveword.Index, relations.Index[i])
					if compareID(ids, relations.Index[i].Id) {
						ids = append(ids, relations.Index[i].Id)
					}
				}
			}
		}
	}

	artists = compareInputIDWithArtistsID(artists, ids)

	return artists, nil
}

func MakeOnlyRelation() (OnlyRelations, error) {
	var relations OnlyRelations

	res, err := http.Get(pageRelation)
	if exitError(err) {
		return relations, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if exitError(err) {
		return relations, err
	}

	err = json.Unmarshal(body, &relations)
	if exitError(err) {
		return relations, err
	}
	return relations, nil
}

func compareID(ids []int, id int) bool {
	for _, w := range ids {
		if w == id {
			return false
		}
	}
	return true
}

func compareInputIDWithArtistsID(artists AllArtists, ids []int) AllArtists {
	var result AllArtists

	for _, input_id := range ids {
		for i := range artists {
			if input_id == artists[i].Id {
				result = append(result, artists[i])
			}
		}
	}
	return result
}
