package main

import (
	"errors"
	"fmt"
	"groupie-tracker/cmd/funcs"
	"html/template"
	"log"
	"net/http"
)

const (
	page400Html = "./templates/html/400.html"
	page404Html = "./templates/html/404.html"
	page500Html = "./templates/html/500.html"
	indexHtml   = "./templates/html/index.html"
	artistHtml  = "./templates/html/artist.html"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		pageNotFound404(w, r)
		return
	}

	ts, err := template.ParseFiles(indexHtml)
	if err != nil {
		log.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}

	artists, err := funcs.MakeAllArtists()
	if err != nil {
		log.Println(err)
		pageInternalServerError500(w, r)
		return
	}

	rel, err := funcs.MakeOnlyRelation()
	if err != nil {
		log.Println(err)
		pageInternalServerError500(w, r)
		return

	}

	for i := range artists {
		for j := range rel.Index {
			artists[i].RelationsData.DatesLocations = rel.Index[j].DatesLocations
			i++
			continue
		}
		break
	}

	var all funcs.ForTemplate
	all.SearchArtist = artists
	all.AllArtist = artists

	err = ts.Execute(w, all)
	if err != nil {
		log.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}
}

func artistByNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ts, err := template.ParseFiles(artistHtml)
	if err != nil {
		log.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}

	id := r.FormValue("id")
	log.Println("id:", id)

	result, err := funcs.MakeOneArtist(id)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, funcs.ErrorPageNotFound) {
			pageNotFound404(w, r)
			return
		} else {
			pageInternalServerError500(w, r)
			return
		}
	}

	err = ts.Execute(w, result)
	if err != nil {
		log.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ts, err := template.ParseFiles(indexHtml)
	if err != nil {
		log.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}

	artists, err := funcs.MakeAllArtists()
	if err != nil {
		log.Println(err)
		pageInternalServerError500(w, r)
		return
	}

	rel, err := funcs.MakeOnlyRelation()
	if err != nil {
		log.Println(err)
		pageInternalServerError500(w, r)
		return

	}

	for i := range artists {
		for j := range rel.Index {
			artists[i].RelationsData.DatesLocations = rel.Index[j].DatesLocations
			i++
			continue
		}
		break
	}

	var all funcs.ForTemplate
	all.AllArtist = artists

	r.ParseForm()
	searchBar := r.Form["search"]
	fmt.Println("search bar:", searchBar)

	if searchBar != nil {
		all.SearchArtist, err = funcs.Search(artists, searchBar)
		if checkErrors(err, w, r) {
			return
		}
	}

	err = ts.Execute(w, all)
	if err != nil {
		log.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}
}

func filters(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ts, err := template.ParseFiles(indexHtml)
	if err != nil {
		log.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}

	artists, err := funcs.MakeAllArtists()
	if err != nil {
		log.Println(err)
		pageInternalServerError500(w, r)
		return
	}

	rel, err := funcs.MakeOnlyRelation()
	if err != nil {
		log.Println(err)
		pageInternalServerError500(w, r)
		return

	}

	for i := range artists {
		for j := range rel.Index {
			artists[i].RelationsData.DatesLocations = rel.Index[j].DatesLocations
			i++
			continue
		}
		break
	}

	var all funcs.ForTemplate
	all.AllArtist = artists

	// var artistsWithFilters funcs.AllArtists

	funcs.CreationDateMinmax(artists)
	funcs.FirstAlbumMinMax(artists)

	r.ParseForm()
	members := r.Form["members"]
	creationDate := r.Form["CreationDate"]
	firstAlbumDate := r.Form["FirstAlbumData"]
	relation := r.Form["LocationsOfConcerts"]
	fmt.Println("members", members)
	fmt.Println("Creation Date:", creationDate)
	fmt.Println("First Album Data:", firstAlbumDate)
	fmt.Println("Locations of Concerts:", relation)

	all.SearchArtist, err = funcs.ChooseCreationDate(artists, creationDate)
	if checkErrors(err, w, r) {
		return
	}

	all.SearchArtist, err = funcs.ChooseFirstAlbum(all.SearchArtist, firstAlbumDate)
	if checkErrors(err, w, r) {
		return
	}

	if members != nil {
		all.SearchArtist, err = funcs.ChooseMembers(all.SearchArtist, members)
		if checkErrors(err, w, r) {
			return
		}
	}

	if relation != nil {
		all.SearchArtist, err = funcs.ChooseLocations(all.SearchArtist, relation)
		if checkErrors(err, w, r) {
			return
		}
	}

	err = ts.Execute(w, all)
	if err != nil {
		log.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}
}

func pageNotFound404(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles(page404Html)
	if err != nil {
		fmt.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	ts.Execute(w, nil)
}

func pageInternalServerError500(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles(page500Html)
	if err != nil {
		fmt.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	ts.Execute(w, nil)
}

func pageBadRequest(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles(page400Html)
	if err != nil {
		fmt.Println(err.Error())
		pageInternalServerError500(w, r)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	ts.Execute(w, nil)
}

func checkErrors(err error, w http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, funcs.ErrorBad) {
			pageBadRequest(w, r)
			return true
		} else if errors.Is(err, funcs.ErrorPageNotFound) {
			pageNotFound404(w, r)
			return true
		} else {
			pageInternalServerError500(w, r)
			return true
		}
	}
	return false
}
