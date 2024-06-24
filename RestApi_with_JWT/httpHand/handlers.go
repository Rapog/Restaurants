package httpHand

import (
	"encoding/json"
	"ex04/DB"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageStr := query.Get("page")
	page := 0
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			if p < 0 || p > 1364 {
				ErrorPage(w, pageStr)
				return
			}
			page = p
		}
	}
	limit := 3
	var offset = page * limit
	var forHTML DB.ForHTML
	forHTML.Places, forHTML.Total, _ = forHTML.Places.GetPlaces(limit, offset)
	var api DB.ApiPlaces = DB.ApiPlaces{IndexName: "places",
		Total:  forHTML.Total,
		Places: forHTML.Places,
	}
	api.LastPage = api.Total / 10
	api.PrevPage = page - 1
	if api.PrevPage < 0 {
		api.PrevPage = 0
	}
	api.NextPage = page + 1
	if api.NextPage > api.LastPage {
		api.NextPage = api.LastPage
	}
	jsonData, err := json.Marshal(api)
	if err != nil {
		log.Printf("Encod error: %s", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func HandleRequestGeo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	latStr := query.Get("lat")
	lonStr := query.Get("lon")

	var geo DB.GeolocFind = DB.GeolocFind{
		PageName: "Recommendation",
	}
	geo.Places = geo.Places.GetPlacesGeo(latStr, lonStr)

	jsonData, err := json.Marshal(geo)
	if err != nil {
		log.Printf("Encod error: %s", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func HandleRequestToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		log.Printf("Token err: %s", err)
	}
	var auth Auth = Auth{
		Token: tokenString,
	}
	err = json.NewEncoder(w).Encode(auth)
	if err != nil {
		log.Printf("Token respond err: %s", err)
	}
}

type Auth struct {
	Token string `json:"token"`
}

// https://blog.logrocket.com/jwt-authentication-go/
func VerifyJWT(endpointHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			log.Printf("Token - %v", r.Header["Authorization"][0])
			tokenSlice := strings.Split(r.Header["Authorization"][0], " ")
			log.Printf("Token[1]: %s", tokenSlice[1])
			token, err := jwt.Parse(tokenSlice[1], func(token *jwt.Token) (interface{}, error) {
				//_, ok := token.Method.(*jwt.SigningMethodECDSA)
				//if !ok {
				//	w.WriteHeader(http.StatusUnauthorized)
				//	//_, err := w.Write([]byte("You're Unauthorized!"))
				//	//if err != nil {
				//	//	return nil, err
				//	//}
				//	//return nil, err
				//}
				return sampleSecretKey, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, err2 := w.Write([]byte("You're Unauthorized due to error parsing the JWT"))
				if err2 != nil {
					return
				}
			}
			if token.Valid {
				endpointHandler(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
	})
}

func ErrorPage(w http.ResponseWriter, page string) {
	w.Header().Set("Content-Type", "application/json")
	type errorJson struct {
		Error string `json:"error"`
	}
	var error errorJson = errorJson{
		Error: "Invalid 'page' value: " + page,
	}
	jsonData, _ := json.Marshal(error)
	http.Error(w, "Server Error", http.StatusBadRequest)
	w.Write(jsonData)
}
