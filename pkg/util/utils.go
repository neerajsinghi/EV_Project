package utils

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/generic"
	iotbike "bikeRental/pkg/repo/iot_bike"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	commonGo "github.com/Trestx-technology/trestx-common-go-lib"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
)

var client *twilio.RestClient
var clientFirebase *messaging.Client
var svc *ses.SES

func init() {
	commonGo.LoadConfig()
	accountSid := viper.GetString("twilio.sid")
	authToken := viper.GetString("twilio.token")
	pathTojsonFile := viper.GetString("firebase.jsonpath")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	opt := option.WithCredentialsFile(pathTojsonFile)
	app, _ := firebase.NewApp(context.Background(), nil, opt)
	clientFirebase, _ = app.Messaging(context.Background())
	svc, _ = createSeSSession()

}
func SendNotification(title, body, token string) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: token,
	}
	if clientFirebase != nil {
		log.Println(clientFirebase.Send(context.Background(), message))
	}
}
func ContainsString(list []string, val string) bool {
	for _, value := range list {
		if value == val {
			return true
		}
	}
	return false
}
func Containsint(list []int, val int) bool {
	for _, value := range list {
		if value == val {
			return true
		}
	}
	return false
}
func EmailLoginOTP(email, name, verificationCode, typ string) (string, error) {
	subject := "Your Verification Code"
	htmlBody := "Mr/Mrs" + name + "\nYour verification Code is<h1>" + verificationCode + "</h1>"
	if typ == "Signup" {
		subject = "Thank you registering with us"
		htmlBody = "Hi " + name + ",<br><br>" + "Thank you for registering with us. "
	}

	return sendEmail(email, subject, htmlBody, htmlBody)
}

func SendVerificationCode(email, name, verificationCode string) (string, error) {
	url := createUrl(verificationCode, "verifyemail")
	subject := viper.GetString("email.subject")
	htmlBody := viper.GetString("email.initial") + name + viper.GetString("email.mid") + " href=" + url + ">Verify Email Now</a>" + viper.GetString("email.end")
	textBody := viper.GetString("email.initial") + name + viper.GetString("email.mid") + " href=" + url + ">Verify Email Now</a>" + viper.GetString("email.end")
	return sendEmail(email, subject, htmlBody, textBody)
}

func SendVarificationSMS(phone, verificationCode string) {
	params := &verify.CreateVerificationParams{}
	params.SetTo(phone)
	params.SetChannel("sms")
	serviceSID := viper.GetString("twilio.service")
	resp, err := client.VerifyV2.CreateVerification(serviceSID, params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
}
func VerifyOTP(phone, code string) (bool, error) {
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(code)
	serviceSID := viper.GetString("twilio.service")
	resp, err := client.VerifyV2.CreateVerificationCheck(serviceSID, params)
	if err != nil {
		fmt.Println("Error verifying phone number: " + err.Error())
		return false, err
	}
	if resp.Valid != nil && *resp.Valid {

		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
		return true, nil
	}
	return false, errors.New("Invalid OTP")
}
func sendEmail(email, subject, htmlBody, textBody string) (string, error) {

	from := viper.GetString("email.from")
	to := email
	input := &ses.SendEmailInput{
		Source: &from,
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(htmlBody),
				},
				Text: &ses.Content{
					Data: aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Destination: &ses.Destination{
			ToAddresses: []*string{&to},
		},
	}
	_, err := svc.SendEmail(input)
	if err != nil {
		commonGo.ECLog3("send email verification failed", err, logrus.Fields{"email": email, "htmlBody": htmlBody})
		return "", err
	}
	return "Sent Successfully", nil
}

func createUrl(verificationcode, path string) string {
	cart := viper.GetString("website.url")
	website := cart
	if strings.Contains(cart, "https") {
		cartSplit := strings.Split(cart, "/")
		website = cartSplit[2]
	}
	u := &url.URL{
		Scheme: "https",
		Host:   website,
		Path:   path + "/" + verificationcode,
	}
	return u.String()
}

func createSeSSession() (*ses.SES, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(viper.GetString("aws.region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws.aws_access_key_id"),
			viper.GetString("aws.aws_secret_access_key"), "")},
	)
	if err != nil {
		commonGo.ECLog2("creating ses session", err)
		return nil, err
	}
	svc := ses.New(sess)
	return svc, nil
}

func CreatePreSignedDownloadUrl(url string) string {
	s := strings.Split(url, "?")
	if len(s) > 0 {
		o := strings.Split(s[0], "/")
		if len(o) > 3 {
			fileName := o[4]
			path := o[3]
			downUrl, _ := commonGo.PreSignedDownloadUrlAWS(fileName, path)
			return downUrl
		}
	}
	return ""
}

func GetDataFromPullAPI() {
	username := viper.GetString("pullapi.username")
	password := viper.GetString("pullapi.password")
	url := fmt.Sprintf("https://pullapi-s1.track360.co.in/api/v1/auth/pull_api?username=%s&password=%s", username, password)
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
	var data response
	err = json.Unmarshal(body, &data)

	if err != nil {
		fmt.Println(err)

	}
	var long, lat float64
	repo := iotbike.NewProfileRepository("iotBike")
	for i, d := range data.Data {
		if data.Data[i].TotalDistance == "" {
			continue
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

	}
}

type response struct {
	Data []entity.IotBikeDB `json:"data"`
}

func CheckError(err error, w http.ResponseWriter) bool {
	if err != nil {
		commonGo.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return true
	}
	return false
}

type AddLog struct {
	UserId      string
	UserRole    string
	ApiName     string
	Err         string
	Data        string
	Request     string
	RequestType string
	Body        string
}

func SendOutput(err error, w http.ResponseWriter, r *http.Request, data, body interface{}, apiName string) {

	token := r.Header.Get("Authorization")
	if token != "" && (r.Method != "GET" || body != nil) {
		tokens := strings.Split(token, " ")
		if len(tokens) > 1 {
			token = tokens[1]
		}
		da, lErr := commonGo.DecodeToken(token)
		if lErr == nil {

			dat, _ := json.Marshal(data)
			logData := AddLog{
				Request:     r.URL.String(),
				RequestType: r.Method,
				ApiName:     apiName,

				Data: string(dat),
			}
			if err != nil {
				logData.Err = err.Error()
			}
			if body != nil {
				dat, _ := json.Marshal(body)

				logData.Body = string(dat)
			}
			if da["userid"] != nil {
				logData.UserId = da["userid"].(string)
			}
			if da["email"] != nil {
				logData.UserRole = da["email"].(string)
			}
			generic.NewRepository("logs").InsertOne(logData)
		}
	}
	if err != nil {
		commonGo.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something went wrong"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

func SetOutput(w http.ResponseWriter) {
	startTime := time.Now()
	commonGo.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
