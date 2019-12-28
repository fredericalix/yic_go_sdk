package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
	yic "github.com/youritcity/go-sdk/youritcity"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type sensor struct {
	ID          uuid.UUID `json:"id"`
	Enable      bool      `json:"enable,omitempty"`
	Activity    string    `json:"activity,omitempty"`
	Status      string    `json:"status,omitempty"`
	RAM         float64   `json:"ram,omitempty"`
	Ramusage    float64   `json:"ramusage,omitempty"`
	Loadaverage float64   `json:"loadaverage,omitempty"`
}

var (
	ids = []uuid.UUID{
		uuid.FromStringOrNil("1377959e-97ce-46c1-9715-22c34bb9afbe"),
		uuid.FromStringOrNil("99f123d0-ad44-4752-9176-8f1ac547030c"),
		uuid.FromStringOrNil("d8f722c7-8345-4396-a17a-7084c9af6745"),
		uuid.FromStringOrNil("da60f79c-a1d5-4ce0-8bcc-ddc17860571f"),
		uuid.FromStringOrNil("e8598b0d-a488-467b-9fe9-52026d65ada5"),
	}
	enables      = []bool{true, false, true, true, true, true, true}
	activities   = []string{"normal", "slow", "fast"}
	status       = []string{"online", "online", "online", "online", "offline", "failure", "build", "maintenance"}
	rams         = []float64{2}
	ramusages    = []float64{0.3, 0.6, 1, 1.3, 1.6, 2}
	loardaerages = []float64{0.1, 0.2, 0.5, 0.8, 1, 1.1, 1.5, 2.2, 3.3, 4}
)

func genRandSensors() sensor {
	return sensor{
		ID: ids[rand.Intn(len(ids))],
		// Enable:   enables[rand.Intn(len(enables))],
		Activity: activities[rand.Intn(len(activities))],
		Status:   status[rand.Intn(len(status))],
		// RAM:         rams[rand.Intn(len(rams))],
		// Ramusage:    ramusages[rand.Intn(len(ramusages))],
		// Loadaverage: loardaerages[rand.Intn(len(loardaerages))],
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	sleepTime := flag.Duration("t", time.Second, "Time to wait between sending next sensor update.")
	url := flag.String("url", "", "url of the server.")
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "usage : %s app_token\n", os.Args[0])
		os.Exit(1)
	}
	token := yic.App{Token: flag.Args()[0]}

	// Establish a connection to yourITcity server and try to get a jwt token
	conn := yic.NewConnectionWithConfig(yic.ConnectionConfig{URI: *url, InsecureSSL: true})
	for {
		var err error
		_, err = conn.Renew(token)
		if err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(time.Second * 5)
	}
	client := conn.Client()

	// Random sensor generator
	sch := make(chan sensor, 1)
	go func(ch chan<- sensor) {
		for {
			ch <- genRandSensors()
		}
	}(sch)

	for sensor := range sch {
		// Try to send the sensor message
		code, err := send(client, sensor, *url)
		if err != nil {
			if code != 401 {
				log.Fatalln(code, err)
			}
			// It fail with 401 UnAuthorized try to renew
			_, err := conn.Renew(token)
			if err != nil {
				log.Fatalln(err)
			}
			// Try to resend the sensor now the jwt token has been renewed
			code, err = send(client, sensor, *url)
			if err != nil {
				log.Fatalln(code, err)
			}
		}
		b, _ := json.Marshal(sensor)
		log.Printf("%d %s\n", code, b)

		// Wait before sending the next sensor message
		time.Sleep(*sleepTime)
	}
}

func send(client *http.Client, s sensor, baseurl string) (code int, err error) {
	bstr, _ := json.Marshal(s)
	body := bytes.NewBuffer(bstr)

	resp, err := client.Post(baseurl+"/sensors", "application/json", body)
	if err != nil {
		if resp != nil {
			return resp.StatusCode, err
		}
		return 0, err
	}
	if resp.StatusCode == http.StatusOK {
		return resp.StatusCode, nil
	}

	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return resp.StatusCode, fmt.Errorf("%s", b)
}
