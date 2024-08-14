package motog

import (
	"bytes"
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
	"unicode"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

// signup
/*
{
    "statusCode": "0000",
    "statusMessage": "SUCCESS",
    "statusInfo": null,
    "responseData": {
        "userName": "info@futureev.in",
        "idToken": "eyJraWQiOiJTUDB4endISVk3OUJJSzNKdDN1YkpxcHhxNk1mUXRqVTJnb3YwQm9HRjlZPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkMzllODhkZi0yODMxLTQyZjktODc3Zi1iMjQ3Njg0ZTcwZTUiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiY3VzdG9tOmIyYiI6IjAiLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb21cL3VzLWVhc3QtMV9YdFN6VGJrdGMiLCJjdXN0b206Z3JvdXAiOiJVc2VyIiwiY29nbml0bzp1c2VybmFtZSI6ImQzOWU4OGRmLTI4MzEtNDJmOS04NzdmLWIyNDc2ODRlNzBlNSIsIm9yaWdpbl9qdGkiOiI3YjMyNjZkNy1mMmY4LTQ4MjMtOTIyYi1iNzE4OWViZWI5NWYiLCJhdWQiOiI0dGVpcDQxOWs1NWhqdXB0aXA4amJvYWxudiIsImN1c3RvbTpjdXN0b21lcklkIjoiMTE4IiwiZXZlbnRfaWQiOiJkODcwYTMwMy1kNGM2LTQ3OTAtYjgzNy1iMTZkOTcxMzNhNDEiLCJ0b2tlbl91c2UiOiJpZCIsImF1dGhfdGltZSI6MTcyMTg4MzEwMSwiY3VzdG9tOm9yZ0lkIjoiMSIsImV4cCI6MTcyMTg4NDkwMSwiY3VzdG9tOnJvbGUiOiJCMkIiLCJpYXQiOjE3MjE4ODMxMDEsImp0aSI6IjNlNGMxMjYwLTU1Y2EtNDExMi1hOTdjLWM0NTNjMmMwNmI0NCIsImVtYWlsIjoiaW5mb0BmdXR1cmVldi5pbiJ9.RF5PO0zmq0Zep-j_X-TZ8Zl_dT2JXWwj6yVBM24ZxGapIgWphivvSPgS0iQDFbM53ATJyFAd5-gmatNqriN695vkCpKU3kq1z1_uacTU74gHErP1pVj2oGDYcPjkd-NPnyUo_6Wbywjo6qXtcTKF6X1N5sOXmerANXICArDtoDNE62BJnaU7FNBx6DQc24f-j0aMrV7S4CfdK4-rahP4SalWh3X7hIeSqdiNcoT4bBXc_Nq6nmwL-ftwIG1JcFgpHxPeSG49-l94kYL6MjAMpDezziHhBzVuQo_jrCZJNDd0VfjrAVVVEk-QxfK9YLCcYTa2f4uqIPs7YYeykrg9sw",
        "accessToken": "eyJraWQiOiJZbXFzYnArNTJGNktrWDBHSzdITlRjTlVzY1g5TStSWWliWFJHRWtkeDgwPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkMzllODhkZi0yODMxLTQyZjktODc3Zi1iMjQ3Njg0ZTcwZTUiLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb21cL3VzLWVhc3QtMV9YdFN6VGJrdGMiLCJjbGllbnRfaWQiOiI0dGVpcDQxOWs1NWhqdXB0aXA4amJvYWxudiIsIm9yaWdpbl9qdGkiOiI3YjMyNjZkNy1mMmY4LTQ4MjMtOTIyYi1iNzE4OWViZWI5NWYiLCJldmVudF9pZCI6ImQ4NzBhMzAzLWQ0YzYtNDc5MC1iODM3LWIxNmQ5NzEzM2E0MSIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE3MjE4ODMxMDEsImV4cCI6MTcyMTg4NDMwMSwiaWF0IjoxNzIxODgzMTAxLCJqdGkiOiJiMzM1YTU3OC05YTc4LTQ3ZTktYTY2OC1iZGFlZDgzZmQ5ZTAiLCJ1c2VybmFtZSI6ImQzOWU4OGRmLTI4MzEtNDJmOS04NzdmLWIyNDc2ODRlNzBlNSJ9.RrP8B8AD42J7WpDtFhjBXHYu_TnhbT5SZmZsoR2Kl9rPl8eShLQFUuxobK6sxdgiunznbAL3e4hlEGxJIrVbVMxBJc9qURQ6Q6f7M8VsioMXfJNDDFGwYrWfBAXg-SgLyIn7-3krEyed7SdtRFO1LWsB5UVHrCGcAIuQWWYLYnGnvuSAANpdSYphDA3rlykcd21XAG6TpT8Wn5MQEzLzQ9cv_1mxZqvvM6hV0ebLUq31o7O6h9n5n6lGQ91TB7TLjXEM57ruvFJaPyGqVgn4KSE4yhRwUzSOpuvqh8C-F0_-rSx5jUy7f5mNzENXT-Y9QnzxTtsWup4CDeu5Q1OIuw",
        "refreshToken": "eyJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwiYWxnIjoiUlNBLU9BRVAifQ.IB1Yu1VwgZtXo7nz1-weAiyqQBf722bIEeWxM9gcbDyw1v1XEUpA5bdDvia9C0qFO4AiHg7gHWPXXSI7OVAHgaVpruvRhnj4W8elh_55qsXBwfe6Sm0ecQZHw8vim0XmMGeHSxPs5vhq0g-6bC9USsZF7yo3n15K-4tUeOxERJFwGobEQKck496U1HGcsjBTZO2BRu3Mw6C3VyOYtHaFpOGKLQIu0KH0BdbzssbYfjNji9Gvsak_Ngs0fB2IAdn3I6XU6n9MQrPRQU7fgu_F8lcKnScbK7SSjL077vlLZBbH27NClpxrNxFg0omWNKdW3VrzXVL_kXGBXFXvN8wvtQ.L0a2gopIQUkYudSE.mJsiKinKTRfvjaKRRwtWOQ5J8-9WhA1cJSyqDLc_5t9pjABthOaT4J6p4Fv1fcfNVhjqlb6Y3Rd3NIDbwKQi3yfsQdP0rL4eL9xWkt1_-6Y9vb_cu8aN6arUX5TqKlY68BD1X8TPGvj3raoY8WK5Bkxv8BnCXoUFsZ9X62kj2-p1Zla5qEvBUI2e2bGO6FgblRuQqXBaIhKBxn9C8pI2FiELqf995n-Uv855jnxRaSr-DtXcmDiqe0Sh-DDfusKvX7NTa3oriqtR5mZh1C5mRdkFwUqUOc708ll2etX6rGQrebU1MY2eR5lXS4TdmzAuDyGTKN4NaBR1JFZQL2TAmoa9hnvHLmIMtzwCGJcPRJx3IeA3OwtfoT-dFaRH0jTmxOtCA5Q1GUINBkUWrfbDs6GHeMw62vf6flqk6-cXDm3b1aI3f3S0dQ8IqVka7YCflLOq3Uqk3eI3OY0e3CcxerHGmH7xVSyMtgx0KmD8vRx7z0umjCENgV9Bn7yXVItf4qw1SzMwooqsu8Ud8CQNT-JhA3Ot09kVRz9q4fu9zsqysGymYec3nW_Vb00pgM0XVVs5AQCX7BfdPFSoc4KRD-AgHIlJW3EJDIjs5yN8xHJVjYsRIedUCEIC0wnHH6MQFVWfvwtNezJA_-zMnyt7FrKryAT4yC85Wl1K28w32ZTfF4bFGs0DWa1xzo5FHGtjHpLCJC6IG8uaxM89wLcKdJwa3maXWi8NyNokygFZkV6QtZ5xJhCyT67UbnvlwCAj7H4kMBEEpJA8J6gfQgxUJZtKugonKyrxJIHBPqYdqCmNzRcnTPwzBdek2Eou8TGsJpTQeZuGiFTIO__tZgULO5be3FsnvaLNzyU08EbY6bthvrmG9iXwpVYh4suFnptisxAs5gTN2MetqcxkfXohxJWeJ--nQfFwovTY2biiu-Zz062OZUByyyD_Uf_GbJokRkAB-mtlDxhXUBzMSEqwxs_a-jSxdYFxkbIauxaudzN73-i93cwsHg50zaBptnF5wxhfO7Vt06sGJ2AHDMBXGCd36jcp_3jBYIU4T3WRyih3fP4emQFdSSEiWM5aizbzVkadXXmjcVD6J_iQvCzu4Rj0J0gAey6aJz5kdtEfuqsv_6gtmscQ_voNYwCIgJC8cvwrphkAVcWAeknlkyAtAJkyNP5sVf4MZQ67Yv_7u4KrCmgF5mskpvlFqvFgyBFQ2E-2KlDF49SXSagwp51iKnW_ZIasmdJc9Cy2wfEYjbJyHZF-OoKCuqHZ1qmCbjQFN_X5_oT-lCByHVa0P_vCuQ-RExmsx8FyHATwGctwkqz9l-HiDUWpmKJlMSo.JPSnS_F9mZb7GnnQpGOTFw",
        "expiresIn": 1200,
        "authTime": 1721883101000,
        "status": "Logged in successfully"
    },
    "requestData": null
}
*/
type SignupResponse struct {
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
	StatusInfo    string `json:"statusInfo"`
	ResponseData  struct {
		UserName     string `json:"userName"`
		IDToken      string `json:"idToken"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresIn    int    `json:"expiresIn"`
		AuthTime     int    `json:"authTime"`
		Status       string `json:"status"`
	} `json:"responseData"`
	RequestData interface{} `json:"requestData"`
}

type VehicleList struct {
	BatteryID  string `json:"battery_id"`
	ChassisNum string `json:"chassis_num"`
}

/*
	{
	    "chargecycle": 9,
	    "current": 0,
	    "chargestatus": 0,
	    "DTE": 77,
	    "odometer": 279,
	    "soc": 69.3,
	    "ignitionstatus": 0,
	    "latitude": 28.582621,
	    "max_speed": 0,
	    "longitude": 77.200366,
	    "timestamp": "2024-07-19 17:27:54",
	    "VIN No": "MVU310823018"
	}
*/
type VehicleData struct {
	ChargeCycle    float64 `json:"chargecycle"`
	Current        float64 `json:"current"`
	ChargeStatus   float64 `json:"chargestatus"`
	DTE            float64 `json:"DTE"`
	Odometer       float64 `json:"odometer"`
	Soc            float64 `json:"soc"`
	IgnitionStatus int     `json:"ignitionstatus"`
	Latitude       float64 `json:"latitude"`
	MaxSpeed       float64 `json:"max_speed"`
	Longitude      float64 `json:"longitude"`
	Timestamp      string  `json:"timestamp"`
	VINNo          string  `json:"VIN No"`
}
type BikeLog struct {
	DeviceID            int
	DeviceName          string
	Location            entity.Location
	DeviceTotalDistance float64
	DeviceTime          string
	Type                string
}

func AddDeviceMoto() {
	data, shouldReturn := signup()
	if shouldReturn {
		return
	}
	if data.ResponseData.IDToken != "" {
		fmt.Println("Logged in successfully")
		url := viper.GetString("motoapi.urlapp") + "/getvehiclelist?app_code=" + viper.GetString("motoapi.appcode")
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", data.ResponseData.IDToken)

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

		var dataVehicle []VehicleList
		err = json.Unmarshal(body, &dataVehicle)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, d := range dataVehicle {
			url := viper.GetString("motoapi.urlapp") + "/vehiclelivedata?app_code=" + viper.GetString("motoapi.appcode") + "&chassis_num=" + d.ChassisNum
			method := "GET"

			client := &http.Client{}
			req, err := http.NewRequest(method, url, nil)

			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Add("Authorization", data.ResponseData.IDToken)

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

			var dataVehicleData VehicleData
			err = json.Unmarshal(body, &dataVehicleData)
			if err != nil {
				log.Println(err)
			}
			var iodDevice entity.IotBikeDB
			numberPart := ""
			for _, dnum := range d.ChassisNum {
				if unicode.IsNumber(dnum) {
					numberPart += string(dnum)
				}
			}
			iodDevice.DeviceId, _ = strconv.Atoi(numberPart)
			iodDevice.BatteryLevel = dataVehicleData.Soc
			iodDevice.Name = d.ChassisNum
			iodDevice.Location.Type = "Point"
			iodDevice.Location.Coordinates = []float64{dataVehicleData.Longitude, dataVehicleData.Latitude}
			iodDevice.Ignition = strconv.Itoa(dataVehicleData.IgnitionStatus)
			iodDevice.LastUpdate = dataVehicleData.Timestamp
			iodDevice.TotalDistanceFloat = dataVehicleData.Odometer
			iodDevice.TotalDistance = strconv.FormatFloat(dataVehicleData.Odometer, 'f', -1, 64)
			iodDevice.Speed = strconv.FormatFloat(dataVehicleData.MaxSpeed, 'f', -1, 64)
			iodDevice.Type = "moto"
			var updateFields bson.M
			conv, _ := bson.Marshal(iodDevice)
			bson.Unmarshal(conv, &updateFields)
			filter := bson.M{"deviceId": iodDevice.DeviceId}
			repo := iotbike.NewProfileRepository("iotBike")

			repo.UpdateOne(filter, bson.M{"$set": updateFields})
			bikeLog := BikeLog{
				DeviceID:            iodDevice.DeviceId,
				DeviceName:          iodDevice.Name,
				Location:            iodDevice.Location,
				DeviceTotalDistance: iodDevice.TotalDistanceFloat,
				DeviceTime:          dataVehicleData.Timestamp,
				Type:                iodDevice.Type,
			}
			repoDev := generic.NewRepository("bikeLog")
			repoDev.InsertOne(bikeLog)
		}
	}
}

func signup() (SignupResponse, bool) {
	url := viper.GetString("motoapi.url") + "/auth/signin"
	method := "POST"

	payload := strings.NewReader(`{
	"userName": "` + viper.GetString("motoapi.userName") + `",
	"userPoolId": "` + viper.GetString("motoapi.userPoolId") + `",
	"password": "` + viper.GetString("motoapi.password") + `"
	}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return SignupResponse{}, true
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return SignupResponse{}, true
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return SignupResponse{}, true
	}

	var data SignupResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return SignupResponse{}, true
	}
	return data, false
}

func ImmoblizeDevice(immoblize int, chasisNumber string) {
	data, shouldReturn := signup()
	if shouldReturn {
		return
	}
	if data.ResponseData.IDToken != "" {
		url := viper.GetString("motoapi.urlapp") + "/cloudControlCommand"
		method := "POST"
		payloadS := map[string]interface{}{
			"app_code":    viper.GetString("motoapi.appcode"),
			"chassis_num": chasisNumber,
			"command":     "IMMOBILIZE",
			"value":       immoblize, // No need for strconv.Itoa if immoblize is already an int
			"param_name":  "DEVICE_UPDATE",
		}

		payload, err := json.Marshal(payloadS)
		if err != nil {
			fmt.Println(err)
			return
		}
		client := &http.Client{}
		req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", data.ResponseData.IDToken)

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
}
