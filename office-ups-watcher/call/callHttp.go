package call

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HttpPost(url string, bodyData string, requestName string, accessToken string) error {
	request, err := http.NewRequest("POST", url, bytes.NewBufferString(bodyData))
	if err != nil {
		log.Printf("Error creating request for %s", requestName)
		log.Print(err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error calling %s", requestName)
		log.Print(err)
	}
	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body from %s", requestName)
		log.Println(err)
	}

	if response.StatusCode != 200 {
		log.Printf("Non-200 code returned when calling %s", requestName)
	}
	fmt.Println(string(b))
	return err
}