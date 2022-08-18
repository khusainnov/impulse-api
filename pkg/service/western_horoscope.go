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
	burred        = "true"
	intact        = "false"
	clientSecret  = "a99e7d506d3701c5c04de3db1913eeee"
	fire          = "Огонь"
	ground        = "Земля"
	air           = "Воздух"
	water         = "Вода"
	cardinalCross = "Кардинальный крест"
	fixedCross    = "Фиксированный крест"
	mutableCross  = "Мутабельный крест"
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

func (ws *WesternHoroscope) DataWorkerWithoutTime(r io.Reader, sex string) (entity.Summary, error) {
	var dataBody entity.Summary

	localAspects := make([]entity.Aspects, 0, 1000)

	dBody, err := ioutil.ReadAll(r)
	if err != nil {
		return entity.Summary{}, err
	}

	err = json.Unmarshal(dBody, &dataBody)
	if err != nil {
		return entity.Summary{}, err
	}

	// Assignments elements and crests for every planet with zodiac sign
	for i, _ := range dataBody.Planets {
		sunDegree := dataBody.Planets[0].FullDegree
		sunDegree2 := dataBody.Planets[0].FullDegree
		switch dataBody.Planets[i].Sign {
		case "Aries":
			dataBody.Planets[i].Element = fire
			dataBody.Planets[i].Crest = cardinalCross
			break
		case "Taurus":
			dataBody.Planets[i].Element = ground
			dataBody.Planets[i].Crest = fixedCross
			break
		case "Gemini":
			dataBody.Planets[i].Element = air
			dataBody.Planets[i].Crest = mutableCross
			break
		case "Cancer":
			dataBody.Planets[i].Element = water
			dataBody.Planets[i].Crest = cardinalCross
			break
		case "Leo":
			dataBody.Planets[i].Element = fire
			dataBody.Planets[i].Crest = fixedCross
			break
		case "Virgo":
			dataBody.Planets[i].Element = ground
			dataBody.Planets[i].Crest = mutableCross
			break
		case "Libra":
			dataBody.Planets[i].Element = air
			dataBody.Planets[i].Crest = cardinalCross
			break
		case "Scorpio":
			dataBody.Planets[i].Element = water
			dataBody.Planets[i].Crest = fixedCross
			break
		case "Sagittarius":
			dataBody.Planets[i].Element = fire
			dataBody.Planets[i].Crest = mutableCross
			break
		case "Capricorn":
			dataBody.Planets[i].Element = ground
			dataBody.Planets[i].Crest = cardinalCross
			break
		case "Aquarius":
			dataBody.Planets[i].Element = air
			dataBody.Planets[i].Crest = fixedCross
			break
		case "Pisces":
			dataBody.Planets[i].Element = water
			dataBody.Planets[i].Crest = mutableCross
			break
		default:
			logrus.Infoln("Данного знака зодиака не найдено")
		}

		// Checking planets for their statement
		if i > 0 && (((dataBody.Planets[i].FullDegree - sunDegree2) >= -4) && ((dataBody.Planets[i].FullDegree - sunDegree) <= 4)) {
			dataBody.Planets[i].Burred = burred
		} else {
			dataBody.Planets[i].Burred = intact
		}
	}

	dataAspects := dataBody.Aspects

	// for male
	if sex == "male" {
		for i, _ := range dataAspects {
			if (dataAspects[i].AspectingPlanet == "Sun" && (dataAspects[i].AspectedPlanet == "Venus" || dataAspects[i].AspectedPlanet == "Moon")) ||
				((dataAspects[i].AspectingPlanet == "Venus" || dataAspects[i].AspectingPlanet == "Moon") && dataAspects[i].AspectedPlanet == "Sun") {
				localAspects = append(localAspects, dataAspects[i])
			}
		}
	}

	// starting p.3
	for i, _ := range dataAspects {
		if ((dataAspects[i].AspectedPlanet == "Mars" && dataAspects[i].AspectingPlanet == "Venus") || (dataAspects[i].AspectedPlanet == "Venus" && dataAspects[i].AspectingPlanet == "Mars")) && dataAspects[i].Type == "Conjunction" {
			localAspects = append(localAspects, dataAspects[i])
		}

		if dataAspects[i].Type == "Square" || dataAspects[i].Type == "Opposition" {
			localAspects = append(localAspects, dataAspects[i])
		}

		if (((dataAspects[i].AspectingPlanet == "Mars" || dataAspects[i].AspectingPlanet == "Venus") &&
			(dataAspects[i].AspectedPlanet == "Uranus" && dataAspects[i].AspectedPlanet == "Neptune" && dataAspects[i].AspectedPlanet == "Pluto")) ||
			((dataAspects[i].AspectedPlanet == "Mars" || dataAspects[i].AspectedPlanet == "Venus") &&
				(dataAspects[i].AspectingPlanet == "Uranus" && dataAspects[i].AspectingPlanet == "Neptune" && dataAspects[i].AspectingPlanet == "Pluto"))) &&
			dataAspects[i].Type == "Conjunction" {
			localAspects = append(localAspects, dataAspects[i])
		}
		// Ending p.3
	}

	dataBody.Aspects = localAspects

	return dataBody, nil
}

/*func withoutTimeMan(r io.Reader) (entity.Summary, error) {
	var dataBody entity.Summary

	localAspects := make([]entity.Aspects, 0, 1000)

	dBody, err := ioutil.ReadAll(r)
	if err != nil {
		return entity.Summary{}, err
	}

	err = json.Unmarshal(dBody, &dataBody)
	if err != nil {
		return entity.Summary{}, err
	}

	// Assignments elements and crests for every planet with zodiac sign
	for i, _ := range dataBody.Planets {
		sunDegree := dataBody.Planets[0].FullDegree
		sunDegree2 := dataBody.Planets[0].FullDegree
		switch dataBody.Planets[i].Sign {
		case "Aries":
			dataBody.Planets[i].Element = fire
			dataBody.Planets[i].Crest = cardinalCross
			break
		case "Taurus":
			dataBody.Planets[i].Element = ground
			dataBody.Planets[i].Crest = fixedCross
			break
		case "Gemini":
			dataBody.Planets[i].Element = air
			dataBody.Planets[i].Crest = mutableCross
			break
		case "Cancer":
			dataBody.Planets[i].Element = water
			dataBody.Planets[i].Crest = cardinalCross
			break
		case "Leo":
			dataBody.Planets[i].Element = fire
			dataBody.Planets[i].Crest = fixedCross
			break
		case "Virgo":
			dataBody.Planets[i].Element = ground
			dataBody.Planets[i].Crest = mutableCross
			break
		case "Libra":
			dataBody.Planets[i].Element = air
			dataBody.Planets[i].Crest = cardinalCross
			break
		case "Scorpio":
			dataBody.Planets[i].Element = water
			dataBody.Planets[i].Crest = fixedCross
			break
		case "Sagittarius":
			dataBody.Planets[i].Element = fire
			dataBody.Planets[i].Crest = mutableCross
			break
		case "Capricorn":
			dataBody.Planets[i].Element = ground
			dataBody.Planets[i].Crest = cardinalCross
			break
		case "Aquarius":
			dataBody.Planets[i].Element = air
			dataBody.Planets[i].Crest = fixedCross
			break
		case "Pisces":
			dataBody.Planets[i].Element = water
			dataBody.Planets[i].Crest = mutableCross
			break
		default:
			logrus.Infoln("Данного знака зодиака не найдено")
			break
		}

		// Checking planets for their statement
		if i > 0 && (((dataBody.Planets[i].FullDegree - sunDegree2) >= -4) && ((dataBody.Planets[i].FullDegree - sunDegree) <= 4)) {
			dataBody.Planets[i].Burred = burred
		} else {
			dataBody.Planets[i].Burred = intact
		}
	}

	dataAspects := dataBody.Aspects

	for i, _ := range dataAspects {
		if (dataAspects[i].AspectingPlanet == "Sun" && (dataAspects[i].AspectedPlanet == "Venus" || dataAspects[i].AspectedPlanet == "Moon")) ||
			((dataAspects[i].AspectingPlanet == "Venus" || dataAspects[i].AspectingPlanet == "Moon") && dataAspects[i].AspectedPlanet == "Sun") {
			localAspects = append(localAspects, dataAspects[i])
		}
	}

	// starting p.3
	for i, _ := range dataAspects {
		if ((dataAspects[i].AspectedPlanet == "Mars" && dataAspects[i].AspectingPlanet == "Venus") || (dataAspects[i].AspectedPlanet == "Venus" && dataAspects[i].AspectingPlanet == "Mars")) && dataAspects[i].Type == "Conjunction" {
			localAspects = append(localAspects, dataAspects[i])
		}

		if dataAspects[i].Type == "Square" || dataAspects[i].Type == "Opposition" {
			localAspects = append(localAspects, dataAspects[i])
		}
	}

	return entity.Summary{}, nil
}*/

func (ws *WesternHoroscope) DataWorkerWithTime(r io.Reader) (entity.Summary, error) {
	var dataBody entity.Summary
	var localDataBody entity.Summary

	dBody, err := ioutil.ReadAll(r)
	if err != nil {
		return entity.Summary{}, err
	}

	err = json.Unmarshal(dBody, &dataBody)
	if err != nil {
		return entity.Summary{}, err
	}

	// start of p.1
	// Planets with house number 7 will add into localDataBody
	for _, v := range dataBody.Planets {
		if v.House == 7 {
			localDataBody.Planets = append(localDataBody.Planets, v)
		}
	}

	dataBodyAspects := dataBody.Aspects

	// If planets with house #7 exist in localDataBody we continue work with p.1
	if len(localDataBody.Planets) != 0 {
		for _, v := range localDataBody.Planets {
			for i, _ := range dataBodyAspects {
				if (dataBodyAspects[i].AspectingPlanet == v.Name) || (dataBodyAspects[i].AspectedPlanet == v.Name) {
					localDataBody.Aspects = append(localDataBody.Aspects, dataBodyAspects[i])
				}
			}
		}

		// end of p.1
	} else {
		// if planets with house #7 do not exist
		// start of p.2

		// Retro planet exist or not
		for i, _ := range dataBody.Planets {
			if dataBody.Planets[i].IsRetro == "true" {
				localDataBody.Planets = append(localDataBody.Planets, dataBody.Planets[i])
			}
		}

		// tense aspect between Sun and Moon
		dbAsp := dataBody.Aspects
		for i, _ := range dbAsp {
			if (dbAsp[i].Type == "Square" || dbAsp[i].Type == "Opposition") &&
				((dbAsp[i].AspectingPlanet == "Moon" && dbAsp[i].AspectedPlanet == "Sun") ||
					(dbAsp[i].AspectingPlanet == "Sun" && dbAsp[i].AspectedPlanet == "Moon")) {
				localDataBody.Aspects = append(localDataBody.Aspects, dataBody.Aspects[i])
			}
		}
		// end of p.2
	}

	dataBody = localDataBody

	return dataBody, nil
}

func (ws *WesternHoroscope) GenerateToken(clientID int) (string, error) {
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
