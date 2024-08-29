package motog

import (
	"encoding/json"
	"fmt"
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/generic"
	iotbike "futureEVChronJobs/pkg/repo/iot_bike"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

/*
	{
	    "device_token": "782a421b-fa6c-49ec-9d62-0a4768075116",
	    "refresh_token": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE3MjI0MjYzNDUsIm5iZiI6MTcyMjQyNjM0NSwianRpIjoiNDIxZTlhMjYtYzcwZC00YzFhLWFlYzItMjYyODQyZGJkOTA3IiwiaWRlbnRpdHkiOnsiaWQiOjcwNzQ0LCJkYiI6MCwiY28iOjQ3LCJuYW1lIjoiRnV0dXJlIEVWIC0gQWdyYSIsInR5cGUiOiJ1c2VyIiwicmVhZF9vbmx5IjowLCJ0eiI6LTMzMCwidHpfcyI6IkFzaWEvS29sa2F0YSIsInNzbyI6MCwiZGV2aWNlIjoid2ViIiwiYWxpYXMiOiIifSwidHlwZSI6InJlZnJlc2gifQ.AaQ4xxdB0D7IxL7QaGVHZXgZX2G_cMQcqWaI3A5USR8",
	    "token": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE3MjI0MjYzNDUsIm5iZiI6MTcyMjQyNjM0NSwianRpIjoiNTI0NmMxZTYtNWM4ZS00YmE2LWE5MWUtODE0MmIwZDhlZDQ3IiwiZXhwIjoxNzIyNjg1NTQ1LCJpZGVudGl0eSI6eyJpZCI6NzA3NDQsImRiIjowLCJjbyI6NDcsIm5hbWUiOiJGdXR1cmUgRVYgLSBBZ3JhIiwidHlwZSI6InVzZXIiLCJyZWFkX29ubHkiOjAsInR6IjotMzMwLCJ0el9zIjoiQXNpYS9Lb2xrYXRhIiwic3NvIjowLCJkZXZpY2UiOiJ3ZWIiLCJhbGlhcyI6IiJ9LCJmcmVzaCI6ZmFsc2UsInR5cGUiOiJhY2Nlc3MifQ.D8kKvNZJz92yh2xKmcbaAoW_XIF2WSYW2wzFEbgg9yk"
	}
*/
type token struct {
	DeviceToken  string `json:"device_token"`
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
}

func GetDataFromPullAPI() {
	resp := roadCastSignin()
	if resp == nil {
		return
	}
	token := strings.Split(resp.Token, " ")
	if len(token) < 2 {
		return
	}
	baseUrlAuth := viper.GetString("pullapi.baseurlauth")
	url := baseUrlAuth + "token/pull_api?token=" + token[1]
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(body))
	var data response
	err = json.Unmarshal(body, &data)

	if err != nil {
		fmt.Println(err)

	}
	var long, lat float64
	repo := iotbike.NewProfileRepository("iotBike")
	for i, d := range data.Data {
		if data.Data[i].TotalDistance == "" && data.Data[i].TotalDistanceFloat == 0 {
			continue
		} else if data.Data[i].TotalDistanceFloat != 0 {
			data.Data[i].TotalDistance = strconv.FormatFloat(data.Data[i].TotalDistanceFloat, 'f', -1, 64)
		}
		if data.Data[i].ExternalPower != 0 {
			data.Data[i].BatteryLevel = (data.Data[i].ExternalPower / 40.5) * 100
			if data.Data[i].BatteryLevel > 100 {
				data.Data[i].BatteryLevel = 100
			}
		}
		data.Data[i].Location.Type = "Point"
		long, _ = strconv.ParseFloat(d.Longitude, 64)
		lat, _ = strconv.ParseFloat(d.Latitude, 64)
		data.Data[i].Location.Coordinates = []float64{long, lat}
		var updateFields bson.M
		conv, _ := bson.Marshal(data.Data[i])
		bson.Unmarshal(conv, &updateFields)
		filter := bson.M{"deviceId": d.DeviceId}
		repo.UpdateOne(filter, bson.M{"$set": updateFields})
		bikeLog := BikeLog{
			DeviceID:            data.Data[i].DeviceId,
			DeviceName:          data.Data[i].Name,
			Location:            data.Data[i].Location,
			DeviceTotalDistance: data.Data[i].TotalDistanceFloat,
			DeviceTime:          data.Data[i].LastUpdate,
			Type:                data.Data[i].Type,
		}
		repoDev := generic.NewRepository("bikeLog")
		repoDev.InsertOne(bikeLog)

	}
}

func roadCastSignin() *token {
	username := viper.GetString("pullapi.username")
	password := viper.GetString("pullapi.password")
	baseUrlAuth := viper.GetString("pullapi.baseurlauth")
	url := baseUrlAuth + "login"
	method := "POST"

	payload := strings.NewReader(`{
    "username":"` + username + `",
    "password":"` + password + `"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var data token
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
	}
	return &data
}

type response struct {
	Data []entity.IotBikeDB `json:"data"`
}

func ImmoblizeDeviceRoadcast(deviceID int, action string) {
	token := roadCastSignin()
	if token == nil {
		return
	}
	baseUrlAuth := viper.GetString("pullapi.baseurlauth")
	url := baseUrlAuth + "set_owl_mode"
	method := "POST"

	payload := strings.NewReader(`{
    "device_id": "` + strconv.Itoa(deviceID) + `",
    "type": "` + action + `"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token.Token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(body))
}
