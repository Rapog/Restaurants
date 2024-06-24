package httpHand

import (
	"ex01/DB"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func HomeHandlers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the jungle")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the Tuna` Page!")
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageStr := query.Get("page")
	page := 0
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			if p < 0 || p > 1364 {
				w.WriteHeader(400)
				fmt.Fprintf(w, "No page with number: %d", p)
				return
			}
			page = p
		}
	}
	limit := 10
	var offset = page * limit
	var forHTML DB.ForHTML
	forHTML.Places, forHTML.Total, _ = forHTML.Places.GetPlaces(limit, offset)
	//fmt.Println(forHTML.Total)
	//fmt.Fprintf(w, places[0].Name)
	tmpl := template.Must(template.ParseFiles("httpHand/resiki.html"))
	//if err != nil {
	//	log.Fatalf("HTML template parsing file err: %s", err)
	//}
	//forHTML.Page = page
	forHTML.Pages = forHTML.Total / 10
	forHTML.PrevPage = page - 1
	if forHTML.PrevPage < 0 {
		forHTML.PrevPage = 0
	}
	forHTML.NextPage = page + 1
	if forHTML.NextPage > forHTML.Pages {
		forHTML.NextPage = forHTML.Pages
	}
	//forHTML.LastId
	tmpl.Execute(w, forHTML)
}

//func sourceToString(source map[string]interface{}) string {
//	var sb strings.Builder
//	for key, value := range source {
//		sb.WriteString(fmt.Sprintf("%s: %v, ", key, value))
//	}
//	return sb.String()
//}
