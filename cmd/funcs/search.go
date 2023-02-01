package funcs

import (
	"strconv"
	"strings"
)

func Search(artists AllArtists, search []string) (AllArtists, error) {
	var result AllArtists
	str := ""
	splits := []string{}
	for _, w := range search {
		str = w
		splits = strings.Split(w, " - ")
	}

	for i := range artists {
		switch {
		case len(splits) == 3:
			if strings.Contains(strings.ToLower(artists[i].Name), strings.ToLower(splits[2])) {
				result = append(result, artists[i])
			}

		case len(splits) == 2 && splits[1] == "member":
			if compareSearchMembers(splits[0], artists[i].Members) {
				result = append(result, artists[i])
			}
		case len(splits) == 2:
			if strings.Contains(strings.ToLower(artists[i].Name), strings.ToLower(splits[0])) {
				result = append(result, artists[i])
			}
		default:
			switch {
			case strings.Contains(strings.ToLower(artists[i].Name), strings.ToLower(str)):
				result = append(result, artists[i])

			case strings.Contains(strings.ToLower(artists[i].FirstAlbum), strings.ToLower(str)):
				result = append(result, artists[i])

			case compareSearchMembers(str, artists[i].Members):
				result = append(result, artists[i])

			case compareSearchCreationDate(str, artists[i].CreationDate):
				result = append(result, artists[i])
			}
		}
	}

	location, err := ChooseLocations(artists, search)
	if err != nil {
		return nil, ErrorInternal
	}

	for i := range location {
		result = append(result, location[i])
	}

	return result, nil
}

func compareSearchMembers(search string, members []string) bool {
	for _, w := range members {
		if strings.Contains(strings.ToLower(w), strings.ToLower(search)) {
			return true
		}
	}
	return false
}

func compareSearchCreationDate(search string, date int) bool {
	searchDate, err := strconv.Atoi(search)
	if err != nil {
		return false
	}
	if date == searchDate {
		return true
	}
	return false
}
