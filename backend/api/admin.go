package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginSuccesResponseAdmin struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type AdminErrorRespone struct {
	Error string `json:"error"`
}

//jwtkey untuk membuat signature
var jwtKey = []byte("key")

// Struct claim sebagai object yang akan diencode oleh jwt
// jwt.StandardClaims sebagai embedded type untuk provide standart claim yang biasanya ada pada JWT
type ClaimsAdmin struct {
	Username string
	// Role     string
	jwt.StandardClaims
}

func (api *API) loginadmin(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)
	var admin Admin
	err := json.NewDecoder(req.Body).Decode(&admin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := api.adminsRepo.Loginadmin(admin.Username, admin.Password)

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(AdminErrorRespone{Error: err.Error()})
		return
	}

	// adminRoles, err := api.adminsRepo.GetAdminRole(*res)

	// Deklarasi expiry time untuk token jwt (time millisecond)
	// claim menggunakan variable yang sudah didefinisikan diatas
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: admin.Username,
		// Role:     adminRoles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	//token encoded claim dengan salah satu algoritma yang dipakai
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	//jwt string dari token yang sudah dibuat menggunakan JWT key yang telah dideklarasikan
	//return internal error ketika ada kesalahan ketika pembuatan JWT string

	tokenStr, err := token.SignedString(jwtkey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Set token string kedalam cookie response
	//Return response berupa username dan token JWT yang sudah login

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expirationTime,
	})

	encoder.Encode(LoginSuccesResponse{Username: *res, Token: tokenStr})
}

func (api *API) logoutadmin(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	token, err := req.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// return unauthorized ketika token kosong
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// return bad request ketika field token tidak ada
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if token.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	c := http.Cookie{
		Name:   "token",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout succes"))
}
