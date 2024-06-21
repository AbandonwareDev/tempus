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
    user = "login"
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

func storeCredentialsToKeyring(url,login,password string) error {
	url = base64.StdEncoding.EncodeToString([]byte(url))
	login = base64.StdEncoding.EncodeToString([]byte(login))
	password = base64.StdEncoding.EncodeToString([]byte(password))

	credentials := url + " " + login + " " + password
	
    err := keyring.Set(service, user, credentials)
    if err != nil {
        return err
    }
    return nil	
}

func getCredentialsFromKeyring() (url,login,password string, err error) {
    secret, err := keyring.Get(service, user)
    if err != nil {
        return "","","",err
    }

    v := strings.Split(secret," ")

    urlByte, err  := base64.StdEncoding.DecodeString(v[0])
        if err != nil {
        return "","","",err
    }
    loginByte, err := base64.StdEncoding.DecodeString(v[1])
        if err != nil {
        return "","","",err
    }
    passwordByte, err  := base64.StdEncoding.DecodeString(v[2])
        if err != nil {
        return "","","",err
    }
    
    url = string(urlByte)
    login = string(loginByte)
    password = string(passwordByte)
	
   	
    return 
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
