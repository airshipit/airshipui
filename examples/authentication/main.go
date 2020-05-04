/*
 Copyright (c) 2020 AT&T. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// page struct is used for templated HTML
type page struct {
	Title string
}

// id and password passed from the test page
type authRequest struct {
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

func main() {
	// we're not picky, so we'll take everything and sort it out later
	http.HandleFunc("/", handler)

	log.Println("Example Auth Server listening on :12321")
	err := http.ListenAndServe(":12321", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// URI check for /basic-auth, /cookie and /oauth, everything else gets a 404
// Also a switch for GET and POST, everything else gets a 415
func handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	uri := r.RequestURI
	if uri == "/basic-auth" || uri == "/cookie" || uri == "/oauth" {
		switch method {
		case http.MethodGet:
			get(uri, w)
		case http.MethodPost:
			post(uri, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			log.Printf("Method %s for %s being rejected, not implemented", method, uri)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("URI %s being rejected, not found", uri)
	}
}

// handle the GET function and return a templated page
func get(uri string, w http.ResponseWriter) {
	var p page

	switch uri {
	case "/basic-auth":
		p = page{
			Title: "Basic Auth",
		}
	case "/cookie":
		p = page{
			Title: "Cookie",
		}
	case "/oauth":
		p = page{
			Title: "OAuth",
		}
	}

	if p != (page{}) {
		// parse and merge the template
		err := template.Must(template.ParseFiles("./examples/authentication/templates/index.html")).Execute(w, p)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			log.Printf("Error getting the templated html: %v", err)
			http.Error(w, "Error getting the templated html", http.StatusInternalServerError)
		}
	}
}

// handle the POST function and return a mock authentication
func post(uri string, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var authAttempt authRequest
	err = json.Unmarshal(body, &authAttempt)

	if err == nil {
		// TODO: make the id and password part of a conf file somewhere
		id := authAttempt.ID
		passwd := authAttempt.Password
		if id == "airshipui" && passwd == "Open Sesame!" {
			w.WriteHeader(http.StatusCreated)

			response := map[string]interface{}{
				"id":         id,
				"name":       "Some Name",
				"expiration": time.Now().Add(time.Hour * 24).Unix(),
			}

			switch uri {
			case "/basic-auth":
				response["X-Auth-Token"] = base64.StdEncoding.EncodeToString([]byte(id + ":" + passwd))
				response["type"] = "basic-auth"
				postHelper(response, w)
			case "/cookie":
				response["type"] = "cookie"
				cookieHandler(response, w)
			case "/oauth":
				response["type"] = "oauth"
				jwtHandler(id, passwd, response, w)
			}
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			http.Error(w, "Bad id or password", http.StatusUnauthorized)
		}
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		log.Printf("Error unmarshalling the request: %v", err)
		http.Error(w, "Error unmarshalling the request", http.StatusBadRequest)
	}
}

// potentially more complex logic happens here with cookie data
func cookieHandler(response map[string]interface{}, w http.ResponseWriter) {
	cookie, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling cookie response: %v", err)
	}
	b, err := encrypt(cookie)
	if err != nil {
		log.Printf("Error encrypting cookie response: %v", err)
		postHelper(nil, w)
	} else {
		response["cookie"] = b
		postHelper(response, w)
	}
}

// potentially more complex logic happens here with JWT data
func jwtHandler(id string, passwd string, response map[string]interface{}, w http.ResponseWriter) {
	token, err := createToken(id, passwd)
	if err != nil {
		log.Printf("Error creating JWT token: %v", err)
		postHelper(nil, w)
	} else {
		response["jwt"] = token
		postHelper(response, w)
	}
}

// Helper function to reduce the number of error checks that have to happen in other functions
func postHelper(returnData map[string]interface{}, w http.ResponseWriter) {
	if returnData == nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	} else {
		log.Printf("Auth data %s\n", returnData)
		b, err := json.Marshal(returnData)
		if err != nil {
			log.Printf("Error marshaling the response: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
		} else {
			_, err := w.Write(b)
			if err != nil {
				log.Printf("Error sending POST response to client: %v", err)
			} else {
				go notifyElectron(b)
			}
		}
	}
}

// This is intended to send an auth completed message to the system so that it knows there was a successful login
func notifyElectron(data []byte) {
	// TODO: probably need to pull the electron url out into its own
	resp, err := http.Post("http://localhost:8080/auth", "application/json; charset=UTF-8", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error sending auth complete to electron.  The response is %v, the error is %v\n", resp, err)
	}
}

// aes requires a 32 byte key, this is random for demo purposes
func randBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

// this creates a random ciphertext for demo purposes
// this is not intended to be reverseable or to be used in production
func encrypt(data []byte) ([]byte, error) {
	b, err := randBytes(256 / 8)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(b)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// create a JWT (JSON Web Token) for demo purposes, this is not to be used in production
func createToken(id string, passwd string) (string, error) {
	// create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// set some claims
	claims := make(jwt.MapClaims)
	claims["username"] = id
	claims["password"] = passwd
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token.Claims = claims

	//Sign and get the complete encoded token as string
	return (token.SignedString([]byte("airshipui")))
}
