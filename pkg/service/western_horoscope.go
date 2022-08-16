package service

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"impulse-api/internal/entity"
	"impulse-api/pkg/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const (
	clientSecret = "a99e7d506d3701c5c04de3db1913eeee"
)

type WesternHoroscope struct {
	repo repository.ZodiacAPI
}

type tokenClaims struct {
	jwt.StandardClaims
	GrandType    string `json:"grand_type"`
	ClientId     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func NewWesternHoroscope(repo repository.ZodiacAPI) *WesternHoroscope {
	return &WesternHoroscope{repo: repo}
}

func (ws *WesternHoroscope) DataWorker(r io.Reader) (entity.Summary, error) {
	var dataBody entity.Summary
	//var dataBody.Planets []entity.Planets

	dBody, err := ioutil.ReadAll(r)
	if err != nil {
		return entity.Summary{}, err
	}

	err = json.Unmarshal(dBody, &dataBody)
	if err != nil {
		return entity.Summary{}, err
	}

	for i, _ := range dataBody.Planets {
		switch dataBody.Planets[i].Sign {
		case "Aries":
			dataBody.Planets[i].Element = "Огонь"
			break
		case "Taurus":
			dataBody.Planets[i].Element = "Земля"
			break
		case "Gemini":
			dataBody.Planets[i].Element = "Воздух"
			break
		case "Cancer":
			dataBody.Planets[i].Element = "Вода"
			break
		case "Leo":
			dataBody.Planets[i].Element = "Огонь"
			break
		case "Virgo":
			dataBody.Planets[i].Element = "Земля"
			break
		case "Libra":
			dataBody.Planets[i].Element = "Воздух"
			break
		case "Scorpio":
			dataBody.Planets[i].Element = "Вода"
			break
		case "Sagittarius":
			dataBody.Planets[i].Element = "Огонь"
			break
		case "Capricorn":
			dataBody.Planets[i].Element = "Земля"
			break
		case "Aquarius":
			dataBody.Planets[i].Element = "Воздух"
			break
		case "Pisces":
			dataBody.Planets[i].Element = "Вода"
			break
		default:
			logrus.Infoln("Данного знака зодиака не найдено")
		}
	}

	return dataBody, nil
}

func (ws *WesternHoroscope) GenerateToken(clientID int, clientSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		GrandType:    "client_credentials",
		ClientId:     clientID,
		ClientSecret: clientSecret,
	})

	return token.SigningString()
}
