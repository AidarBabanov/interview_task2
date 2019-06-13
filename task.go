package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	now := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
	secret := []byte(os.Getenv("secret"))
	query := fmt.Sprintf("timestamp=%d", now)
	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(query))
	signedQuery := hex.EncodeToString(hash.Sum(nil))
	query = fmt.Sprintf("%s&signature=%s", query, signedQuery)
	u := url.URL{
		Scheme:     "https",
		Host:       "api.binance.com",
		Path:       "/api/v3/account",
		ForceQuery: true,
		RawQuery:   query,
	}
	fmt.Println(u.String())

	request, _ := http.NewRequest("GET", u.String(), strings.NewReader(""))
	request.Header.Set("X-MBX-APIKEY", os.Getenv("apikey"))

	client := http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()

	if err != nil {
		fmt.Println("Error: " + err.Error())
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
	}
}
