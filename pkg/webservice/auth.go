/*
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

package webservice

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
)

// Create the JWT key used to create the signature
// TODO: use a private key for this instead of a phrase
var jwtKey = []byte("airshipUI_JWT_key")

const (
	username   = "username"
	password   = "password"
	expiration = "exp"
)

// The UI will either request authentication or validation, handle those situations here
func handleAuth(_ *string, request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:      configs.UI,
		Component: configs.Auth,
	}

	var err error
	switch request.SubComponent {
	case configs.Authenticate:
		if request.Authentication != nil {
			var token *string
			authRequest := request.Authentication
			token, err = createToken(authRequest.ID, authRequest.Password)
			if token != nil {
				sessions[request.SessionID].jwt = *token
				response.SubComponent = configs.Approved
				response.Token = token
			}
		} else {
			err = errors.New("No AuthRequest found in the request")
		}
	case configs.Validate:
		if request.Token != nil {
			_, err = validateToken(request)
			response.SubComponent = configs.Approved
			response.Token = request.Token
		} else {
			err = errors.New("No token found in the request")
		}
	default:
		err = errors.New("Invalid authentication request")
	}

	if err != nil {
		log.Error(err)
		e := err.Error()
		response.Error = &e
		response.SubComponent = configs.Denied
	}
	return response
}

// validate JWT (JSON Web Token)
func validateToken(request configs.WsMessage) (*string, error) {
	// update the token string to be the refresh token if it's present
	// otherwise just use the default token string
	// TODO(aschiefe): determine if we need to compare the original token claims to the refresh
	tokenString := request.Token
	if request.RefreshToken != nil {
		tokenString = request.RefreshToken
	}

	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		log.Error(err)
		return nil, err
	}

	// extract the claim from the token
	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// extract the user from the claim
		if user, ok := claim[username].(string); ok {
			// test to see if we need to sent a refresh token
			go testForRefresh(claim, request)
			return &user, nil
		}
		err = errors.New("Invalid JWT User")
		log.Error(err)
		return nil, err
	}

	err = errors.New("Invalid JWT Token")
	log.Error(err)
	return nil, err
}

// create a JWT (JSON Web Token)
func createToken(id string, passwd string) (*string, error) {
	origPasswdHash, ok := configs.UIConfig.Users[id]
	if !ok {
		return nil, errors.New("Not authenticated")
	}

	// test the password to make sure it's valid
	hash := sha512.New()
	_, err := hash.Write([]byte(passwd))
	if err != nil {
		return nil, errors.New("Error authenticating")
	}
	if origPasswdHash != hex.EncodeToString(hash.Sum(nil)) {
		return nil, errors.New("Not authenticated")
	}

	// set some claims
	claims := make(jwt.MapClaims)
	claims[username] = id
	claims[password] = passwd
	claims[expiration] = time.Now().Add(time.Hour * 1).Unix()

	// create the token
	jwtClaim := jwt.New(jwt.SigningMethodHS256)
	jwtClaim.Claims = claims

	// Sign and get the complete encoded token as string
	token, err := jwtClaim.SignedString(jwtKey)
	return &token, err
}

// from time to time we might want to send a refresh token to the UI.  The UI should not be in charge of requesting it
func testForRefresh(claim jwt.MapClaims, request configs.WsMessage) {
	// for some reason the exp is stored as an float and not an int in the claim conversion
	// so we do a little dance and cast some floats to ints and everyone goes on with their lives
	if exp, ok := claim[expiration].(float64); ok {
		if int64(exp) < time.Now().Add(time.Minute*15).Unix() {
			createRefreshToken(claim, request)
		}
	}
}

// createRefreshToken will create an oauth2 refresh token based on the timeout on the UI
func createRefreshToken(claim jwt.MapClaims, request configs.WsMessage) {
	// add the new expiration to the claim
	claim[expiration] = time.Now().Add(time.Hour * 1).Unix()

	// create the token
	jwtClaim := jwt.New(jwt.SigningMethodHS256)
	jwtClaim.Claims = claim

	// Sign and get the complete encoded token as string
	refreshToken, err := jwtClaim.SignedString(jwtKey)
	if err != nil {
		log.Error(err)
		return
	}

	// test to see if the session is still in existence before firing off a message
	if session, ok := sessions[request.SessionID]; ok {
		if err = session.webSocketSend(configs.WsMessage{
			Type:         configs.UI,
			Component:    configs.Auth,
			SubComponent: configs.Refresh,
			RefreshToken: &refreshToken,
			SessionID:    request.SessionID,
		}); err != nil {
			session.onError(err)
		}
	}
}
