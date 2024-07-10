package main

import (
	// "log"

	"github.com/zalando/go-keyring"
	// "github.com/99designs/keyring"
	"encoding/base64"
	// "fmt"
	"strings"
)

const (
	service = "Tempus"
	user    = "login"
)

// func debugKeyring() {
//     service := "Tempus"
//     user := "login"
//     password := "secr123123et"
// 	encoded := base64.StdEncoding.EncodeToString([]byte(password))
//     // set password
//     err := keyring.Set(service, user, encoded)
//     if err != nil {
//         log.Fatal(err)
//     }
//
//     // get password
//     secret, err := keyring.Get(service, user)
//     if err != nil {
//         log.Fatal(err)
//     }
//
//     log.Println(secret)
//
//     decodedByte, err := base64.StdEncoding.DecodeString(secret)
//     decoded := string(decodedByte)
//     if err != nil {
//     		fmt.Println("decode error:", err)
//     		return
//     	}
//     log.Println(decoded)
// }

//TODO we can make custom type e.g. SaveData - with []string and put everything saved there. OR we can make the thing like in inputs[variable] - to make it convinient

//TODO inconsistend global funcs (or is it called 'exported funcs?')

func storeCredentialsToKeyring(url, login, password, calendar string) error {
	url = base64.StdEncoding.EncodeToString([]byte(url))
	login = base64.StdEncoding.EncodeToString([]byte(login))
	password = base64.StdEncoding.EncodeToString([]byte(password))
	calendar = base64.StdEncoding.EncodeToString([]byte(calendar))

	credentials := url + " " + login + " " + password + " " + calendar

	err := keyring.Set(service, user, credentials)
	if err != nil {
		return err
	}
	return nil
}

func getCredentialsFromKeyring_wrapper() (url, login, password, calendar string, err error) {
	secret, err := keyring.Get(service, user)
	if err != nil {
		return
	}

	v := strings.Split(secret, " ")

	urlByte, err := base64.StdEncoding.DecodeString(v[0])
	if err != nil {
		return
	}
	loginByte, err := base64.StdEncoding.DecodeString(v[1])
	if err != nil {
		return
	}
	passwordByte, err := base64.StdEncoding.DecodeString(v[2])
	if err != nil {
		return
	}
	calendarByte, err := base64.StdEncoding.DecodeString(v[3])
	if err != nil {
		return
	}

	url = string(urlByte)
	login = string(loginByte)
	password = string(passwordByte)
	calendar = string(calendarByte)

	return
}

func getCredentialsFromKeyring() (Credentials, error) {
	//TODO inconsistent approach compared to caldav.go
	url, username, password, calendar, err := getCredentialsFromKeyring_wrapper()
	if err != nil {
		return Credentials{}, err
	}
	return Credentials{
		URL:          url,
		Username:     username,
		Password:     password,
		CalendarPath: calendar,
	}, nil

}

// func debugKeyring() {
//
// 	kr, err := keyring.Open(keyring.Config{
// 		AllowedBackends: []keyring.BackendType{
// 			keyring.SecretServiceBackend,
// 		},
// 		LibSecretCollectionName: "Defaultkeyring",
// 		ServiceName:             "myapp",
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	err = kr.Set(keyring.Item{
// 		Key: "foo",
// 		// Data: []byte("secret-bar"),
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	v, err := kr.Get("llamas")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	log.Printf("llamas was %v", v)
//
//
//
// }
