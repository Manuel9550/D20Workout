package health

// Functions just for testing that the service is up

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type HealthCheck struct {
	TestBool   bool   `json:"testBool"`
	TestInt    int    `json:"testInt"`
	TestString string `json:"testString"`
}

func GetTest(w http.ResponseWriter, r *http.Request) {

	randInt := rand.Intn(100)
	randBool := false
	randString := ""
	if randInt <= 50 {
		randBool = true
	} else {
		randBool = false
	}

	randString = generate(15)

	returnObject := HealthCheck{
		TestBool:   randBool,
		TestInt:    randInt,
		TestString: randString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnObject)
	return
}

// Taken from https://linuxhint.com/golang-generate-random-string/

func generate(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}
