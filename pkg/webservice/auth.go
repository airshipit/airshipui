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

// The UI will either request authentication or validation, handle those situations here
func handleAuth(request configs.WsMessage) configs.WsMessage {
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
			sessions[request.SessionID].jwt = *token
			response.SubComponent = configs.Approved
			response.Token = token
		} else {
			err = errors.New("No AuthRequest found in the request")
		}
	case configs.Validate:
		if request.Token != nil {
			err = validateToken(*request.Token)
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
		response.Error = err.Error()
		response.SubComponent = configs.Denied
	}

	return response
}

// validate JWT (JSON Web Token)
func validateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	return errors.New("Invalid JWT Token")
}

// create a JWT (JSON Web Token)
// TODO (aschiefe): for demo purposes, this is not to be used in production
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
	claims["username"] = id
	claims["password"] = passwd
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// create the token
	jwtClaim := jwt.New(jwt.SigningMethodHS256)
	jwtClaim.Claims = claims

	// Sign and get the complete encoded token as string
	token, err := jwtClaim.SignedString(jwtKey)
	return &token, err
}
