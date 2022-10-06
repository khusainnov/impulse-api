package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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
	housesUpr     = "houses_upr.json"
	planetsPower  = "planets_power.json"
	filename      = "Aspects.txt"
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

// DataWorkerWithoutTime
// @Tags 		working with data which haven't time (default time = 12:00)
// @Description handling data by the criteria
// @Accept  	io.Reader, sex string
// @Produce 	entity.ResponseUpr
// @Param   	dataBody entity.ResponseWithoutTime
// @Param   	checkData []entity.CheckVars
// @Param   	localAspects []entity.Aspects
// @Param   	responseBody entity.ResponseWithoutTime
// @Param   	mapElement map[string]int
// @Param   	mapCrest map[string]int
// @Param   	planetsName []string
// @Success 	{object} entity.ResponseUpr
// @Failure 	{object} error
// DataWorkerWithoutTime - change response param from entity.Summary on string (for easy sending data into chat)
func (ws *WesternHoroscope) DataWorkerWithoutTime(r io.Reader, sex string) (entity.ResponseWithoutTime, error) {
	var dataBody entity.ResponseWithoutTime

	//localCheckData := make([]entity.CheckVars, 1096)
	//checkData := make([]entity.CheckVars, 365)
	var checkData []entity.CheckVars
	localAspects := make([]entity.Aspects, 0, 1000)
	var responseBody entity.ResponseWithoutTime

	mapElement := map[string]int{fire: 0, ground: 0, air: 0, water: 0}
	mapCrest := map[string]int{cardinalCross: 0, fixedCross: 0, mutableCross: 0}
	planetsName := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn", "Uranus", "Neptune", "Pluto"}

	dBody, err := ioutil.ReadAll(r)
	if err != nil {
		return entity.ResponseWithoutTime{}, err
	}

	err = json.Unmarshal(dBody, &dataBody)
	if err != nil {
		return entity.ResponseWithoutTime{}, err
	}

	// planets power json unmarshalling
	localPlanetsPower, err := JsonPowerReader()
	if err != nil {
		return entity.ResponseWithoutTime{}, err
	}

	// getting response data for messages from local .txt file
	checkData, err = TxtDataWorker()

	// assignments elements and crests for every planet with zodiac sign
	// присваивание элеменотов и крестов каждой планете по знаку задиака
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

		// checking planets for their statement
		// проврека планет на их состояние
		if i > 0 && (((dataBody.Planets[i].FullDegree - sunDegree2) >= -4) && ((dataBody.Planets[i].FullDegree - sunDegree) <= 4)) {
			dataBody.Planets[i].Burred = burred
		} else {
			dataBody.Planets[i].Burred = intact
		}

		// counting elements and crests if planets contains in planetsName
		// элементы и кресты считаются если планета существует в planetsName
		for _, v := range planetsName {
			if v == dataBody.Planets[i].Name {
				mapElement[dataBody.Planets[i].Element]++
				mapCrest[dataBody.Planets[i].Crest]++
			}
		}

		responseBody.AllElems = mapElement
		responseBody.AllCrests = mapCrest

		// assigning planets power (only for 0, 1, 6)
		// присваивание силы планет (только для 0, 1, 6)
		for _, v := range localPlanetsPower.Params {
			if dataBody.Planets[i].Name == v.Planet {
				if dataBody.Planets[i].Sign == v.House {
					dataBody.Planets[i].Power = 6
				} else if dataBody.Planets[i].Sign == v.Exile {
					dataBody.Planets[i].Power = 1
				} else if dataBody.Planets[i].Sign == v.Fall {
					dataBody.Planets[i].Power = 0
				} else {
					dataBody.Planets[i].Power = nil
				}
			}
		}
	}

	tmp := 0
	tmpName := ""

	// look over mapElement and taking max value in mapElement body and writing it into responseBody
	// перебираем mapElement и выбираем переменную с максимальным значением и добавлем его в responseBody
	for i, v := range mapElement {
		if mapElement[i] > tmp {
			tmp = v
		} else {
			continue
		}
		if mapElement[i] == tmp && i != responseBody.PrevVal.FirstElem {
			responseBody.PrevVal.FirstElem = i
			tmpName = i
			//responseBody.PrevElem = fmt.Sprintf("%s: %d", i, v)
		}
	}
	responseBody.TestElems = fmt.Sprintf("%s\n", responseBody.PrevVal.FirstElem)

	for i, v := range mapElement {
		if tmp == v && tmpName != i {
			responseBody.PrevVal.SecondElem = i
			//responseBody.SndPrevE = fmt.Sprintf("%s: %d", i, v)
			responseBody.TestElems += fmt.Sprintf("%s\n", responseBody.PrevVal.SecondElem)
			break
		}
	}
	for i, v := range mapElement {
		if tmp == v && tmpName != i && responseBody.PrevVal.SecondElem != i {
			responseBody.PrevVal.ThirdElem = i
			responseBody.TestElems += fmt.Sprintf("%s\n", responseBody.PrevVal.ThirdElem)
		}

		if tmp == v && tmpName != i && responseBody.PrevVal.SecondElem != i && responseBody.PrevVal.ThirdElem != i {
			responseBody.PrevVal.FourthElem = i
			responseBody.TestElems += fmt.Sprintf("%s\n", responseBody.PrevVal.FourthElem)
		}
	}

	tmp = 0
	tmpName = ""

	// look over mapCrest and taking max value in mapCrest body and writing it into responseBody
	// перебираем mapCrest и выбираем переменную с максимальным значением и добавлем его в responseBody
	for i, v := range mapCrest {
		if mapCrest[i] > tmp {
			tmp = v
		} else {
			continue
		}
		if mapCrest[i] == tmp && i != responseBody.PrevCrest.FirstCrest {
			responseBody.PrevCrest.FirstCrest = i
			tmpName = i
			//responseBody.PrevCrest = fmt.Sprintf("%s: %d", i, v)
		}
	}
	responseBody.TestCrests = fmt.Sprintf("%s\n", responseBody.PrevCrest.FirstCrest)

	for i, v := range mapCrest {
		if tmp == v && tmpName != i {
			responseBody.PrevCrest.SecondCrest = i
			//responseBody.SndPrevC = fmt.Sprintf("%s: %d", i, v)
			responseBody.TestCrests += fmt.Sprintf("%s\n", responseBody.PrevCrest.SecondCrest)
		}
	}
	for i, v := range mapCrest {
		if tmp == v && tmpName != i && responseBody.PrevCrest.SecondCrest != i {
			responseBody.PrevCrest.ThirdCrest = i
			//responseBody.SndPrevC = fmt.Sprintf("%s: %d", i, v)
			responseBody.TestCrests += fmt.Sprintf("%s\n", responseBody.PrevCrest.ThirdCrest)
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
	// cycle for look over all dataAspects
	// цикл для перебора всех аспектов из тела dataAspects
	for i, _ := range dataAspects {

		// if Mars and Venus contains in dataAspects as (aspected or aspecting) planet then we add it into localAspects array
		// если Марс и Венера содератся в dataAspects как ожидаемая планета тогда мы добавляем данный аспект в массив localAspects
		if ((dataAspects[i].AspectedPlanet == "Mars" && dataAspects[i].AspectingPlanet == "Venus") || (dataAspects[i].AspectedPlanet == "Venus" && dataAspects[i].AspectingPlanet == "Mars")) && dataAspects[i].Type == "Conjunction" {
			localAspects = append(localAspects, dataAspects[i])
		}

		// if planet aspects contains type Square or Opposition then this aspect will be added into localAspects array
		// если аспекты содержат тип Квадрат или Оппозиция тогда мы добавляем данный элемент в массив localAspects
		if dataAspects[i].Type == "Square" || dataAspects[i].Type == "Opposition" {

			// skip if this item is already in array
			// пропускаем элемент если он уже содержится в массиве
			if len(localAspects) > 0 {
				if localAspects[0].AspectingPlanet == dataAspects[i].AspectingPlanet && localAspects[0].AspectedPlanet == dataAspects[i].AspectedPlanet {
					continue
				}
			}

			localAspects = append(localAspects, dataAspects[i])
		}

		//
		//
		if (((dataAspects[i].AspectingPlanet == "Mars" || dataAspects[i].AspectingPlanet == "Venus") &&
			(dataAspects[i].AspectedPlanet == "Uranus" && dataAspects[i].AspectedPlanet == "Neptune" && dataAspects[i].AspectedPlanet == "Pluto")) ||
			((dataAspects[i].AspectedPlanet == "Mars" || dataAspects[i].AspectedPlanet == "Venus") &&
				(dataAspects[i].AspectingPlanet == "Uranus" && dataAspects[i].AspectingPlanet == "Neptune" && dataAspects[i].AspectingPlanet == "Pluto"))) &&
			dataAspects[i].Type == "Conjunction" {
			localAspects = append(localAspects, dataAspects[i])
		}
		// Ending p.3
	}

	//dataBody.Aspects = localAspects
	dataBody.Aspects = localAspects

	// adding data into responseBody struct
	// добавляем элементы в выводящую структуру responseBody
	responseBody.Ascendant = dataBody.Ascendant
	responseBody.Midheaven = dataBody.Midheaven
	responseBody.Lilith = dataBody.Lilith

	// preparation data for response body
	// добавляем текст к respMsg если обработанные аспекты совпадают со значениями из файла
	for i := 0; i < len(dataBody.Aspects); i++ {
		switch dataBody.Aspects[i].Type {
		case "Conjunction":
			dataBody.Aspects[i].Type = 0
			break
		case "Sextile":
			dataBody.Aspects[i].Type = 60
			break
		case "Square":
			dataBody.Aspects[i].Type = 90
			break
		case "Trine":
			dataBody.Aspects[i].Type = 120
			break
		case "Opposition":
			dataBody.Aspects[i].Type = 180
			break
		default:
			dataBody.Aspects[i].Type = "null"
			break
		}

		_ = checkData
		for j := 0; j < len(checkData); j++ {
			//for _, v := range checkData {
			if dataBody.Aspects[i].AspectingPlanetID == checkData[j].CheckAspectingID &&
				dataBody.Aspects[i].AspectedPlanetID == checkData[j].CheckAspectedID {

				if dataBody.Aspects[i].Type == checkData[j].CheckType && dataBody.Aspects[i].Type != "null" {
					responseBody.Aspects = append(responseBody.Aspects, dataBody.Aspects[i])
					responseBody.RespMsg += fmt.Sprintf( /*%s\n*/ "%s\n\n" /*v.Soed,*/, checkData[j].Body)
				}
			}
		}
	}

	for _, v := range dataBody.Planets {
		responseBody.Planets = append(responseBody.Planets, v)
	}

	for _, v := range dataBody.Houses {
		responseBody.Houses = append(responseBody.Houses, v)
	}

	return responseBody, nil
}

// DataWorkerWithTime
// @Tags 		working with data which have time
// @Description handling data by the criteria
// @Accept  	io.Reader
// @Produce 	entity.ResponseUpr
// @Param   	dataBody entity.Summary
// @Param   	localDataBody entity.Summary
// @Param   	responseBody entity.ResponseUpr
// @Param   	checkData []entity.CheckVars
// @Param   	mapElement map[string]int
// @Param   	mapCrest map[string]int
// @Param   	planetsName []string
// @Success 	{object} entity.ResponseUpr
// @Failure 	{object} error
// DataWorkerWithTime - change response param from entity.ResponseUpr on string (for easy sending data into chat)
func (ws *WesternHoroscope) DataWorkerWithTime(r io.Reader) (entity.ResponseUpr, error) {
	var dataBody entity.Summary
	var localDataBody entity.Summary
	var responseBody entity.ResponseUpr
	checkData := make([]entity.CheckVars, 1096)

	mapElement := map[string]int{fire: 0, ground: 0, air: 0, water: 0}
	mapCrest := map[string]int{cardinalCross: 0, fixedCross: 0, mutableCross: 0}
	planetsName := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn", "Uranus", "Neptune", "Pluto"}

	dBody, err := ioutil.ReadAll(r)
	if err != nil {
		return entity.ResponseUpr{}, err
	}

	err = json.Unmarshal(dBody, &dataBody)
	if err != nil {
		return entity.ResponseUpr{}, err
	}

	localPlanetsPower, err := JsonPowerReader()

	// assignments elements and crests for every planet with zodiac sign
	// присваивание элеменотов и крестов каждой планете по знаку задиака
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

		// checking planets for their statement
		// проврека планет на их состояние
		if i > 0 && (((dataBody.Planets[i].FullDegree - sunDegree2) >= -4) && ((dataBody.Planets[i].FullDegree - sunDegree) <= 4)) {
			dataBody.Planets[i].Burred = burred
		} else {
			dataBody.Planets[i].Burred = intact
		}

		// counting elements and crests if planets contains in planetsName
		// элементы и кресты считаются если планета существует в planetsName
		for _, v := range planetsName {
			if v == dataBody.Planets[i].Name {
				mapElement[dataBody.Planets[i].Element]++
				mapCrest[dataBody.Planets[i].Crest]++
			}
		}

		responseBody.AllElems = mapElement
		responseBody.AllCrests = mapCrest

		// assigning planets power (only for 0, 1, 6)
		// присваивание силы планет (только для 0, 1, 6)
		for _, k := range localPlanetsPower.Params {
			if dataBody.Planets[i].Name == k.Planet {
				if dataBody.Planets[i].Sign == k.House {
					dataBody.Planets[i].Power = 6
					break
				} else if dataBody.Planets[i].Sign == k.Exile {
					dataBody.Planets[i].Power = 1
					break
				} else if dataBody.Planets[i].Sign == k.Fall {
					dataBody.Planets[i].Power = 0
					break
				} else {
					dataBody.Planets[i].Power = nil
				}
			}
		}
	}

	tmp := 0
	tmpName := ""

	// look over mapElement and taking max value in mapElement body and writing it into responseBody
	// перебираем mapElement и выбираем переменную с максимальным значением и добавлем его в responseBody
	for i, v := range mapElement {
		if mapElement[i] > tmp {
			tmp = v
		} else {
			continue
		}
		if mapElement[i] == tmp && i != responseBody.PrevVal.FirstElem {
			responseBody.PrevVal.FirstElem = i
			tmpName = i
			//responseBody.PrevElem = fmt.Sprintf("%s: %d", i, v)
		}
	}
	responseBody.TestElems = fmt.Sprintf("%s\n", responseBody.PrevVal.FirstElem)

	for i, v := range mapElement {
		if tmp == v && tmpName != i {
			responseBody.PrevVal.SecondElem = i
			responseBody.TestElems += fmt.Sprintf("%s\n", responseBody.PrevVal.SecondElem)
			break
		}
	}
	for i, v := range mapElement {
		if tmp == v && tmpName != i && responseBody.PrevVal.SecondElem != i {
			responseBody.PrevVal.ThirdElem = i
			responseBody.TestElems += fmt.Sprintf("%s\n", responseBody.PrevVal.ThirdElem)
		}

		if tmp == v && tmpName != i && responseBody.PrevVal.SecondElem != i && responseBody.PrevVal.ThirdElem != i {
			responseBody.PrevVal.FourthElem = i
			responseBody.TestElems += fmt.Sprintf("%s\n", responseBody.PrevVal.FourthElem)
		}
	}

	tmp = 0
	tmpName = ""

	// look over mapCrest and taking max value in mapCrest body and writing it into responseBody
	// перебираем mapCrest и выбираем переменную с максимальным значением и добавлем его в responseBody
	for i, v := range mapCrest {
		if mapCrest[i] > tmp {
			tmp = v
		} else {
			continue
		}
		if mapCrest[i] == tmp && i != responseBody.PrevCrest.FirstCrest {
			responseBody.PrevCrest.FirstCrest = i
			tmpName = i
		}
	}
	responseBody.TestCrests = fmt.Sprintf("%s\n", responseBody.PrevCrest.FirstCrest)

	for i, v := range mapCrest {
		if tmp == v && tmpName != i {
			responseBody.PrevCrest.SecondCrest = i
			responseBody.TestCrests += fmt.Sprintf("%s\n", responseBody.PrevCrest.SecondCrest)
		}
	}
	for i, v := range mapCrest {
		if tmp == v && tmpName != i && responseBody.PrevCrest.SecondCrest != i {
			responseBody.PrevCrest.ThirdCrest = i
			responseBody.TestCrests += fmt.Sprintf("%s\n", responseBody.PrevCrest.ThirdCrest)
		}
	}

	responseBody.Planets = dataBody.Planets

	// start of p.1
	// Planets with house number = 7 will be added into localDataBody
	// Планеты у которых номер дома = 7 добавляются в localDataBody
	for _, v := range dataBody.Houses {
		if v.House == 7 {
			localDataBody.Houses = append(localDataBody.Houses, v)
		}
	}

	localHousesUpr, err := jsonUprReader()
	if err != nil {
		return entity.ResponseUpr{}, err
	}

	// assigning UprPlanet values into responseBody
	// присваивание данных планеты-управителя в responseBody
	for _, v := range localHousesUpr.Hoe {
		if v.Sign == localDataBody.Houses[0].Sign {
			responseBody.House = localDataBody.Houses[0].House
			responseBody.Sign = localDataBody.Houses[0].Sign
			responseBody.Upr = v.Upr
		}
	}

	// appending into responseBody.Aspects aspect which contains UprPlanet
	// добавляем в слайс responseBody.Aspects аспекты в которых содержится Планета-управитель
	for _, v := range dataBody.Aspects {
		if v.AspectedPlanet == responseBody.Upr || v.AspectingPlanet == responseBody.Upr {
			responseBody.Aspects = append(responseBody.Aspects, v)
		}
	}

	// getting response data for messages from local .txt file
	checkData, err = TxtDataWorker()

	// if planets with house #7 exist in localDataBody we continue work with p.1
	// продолжение второй части ТЗ, пока делать не нужно
	/*if len(localDataBody.Planets) != 0 {
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
	}*/

	// adding power number where upr = 7 if UprPlanetName equals on of responseBody.Planet
	// присваем управителю 7 планеты значение силы если Имя планеты совпадает с планетой из списка responseBody.Planet
	for _, v := range responseBody.Planets {
		if v.Name == responseBody.Upr {
			responseBody.Power = v.Power
		}
	}

	// preparation data for response body
	// добавляем текст к respMsg если обработанные аспекты совпадают со значениями из файла
	for i := 0; i < len(responseBody.Aspects); i++ {
		switch responseBody.Aspects[i].Type {
		case "Conjunction":
			responseBody.Aspects[i].Type = 0
			break
		case "Sextile":
			responseBody.Aspects[i].Type = 60
			break
		case "Square":
			responseBody.Aspects[i].Type = 90
			break
		case "Trine":
			responseBody.Aspects[i].Type = 120
			break
		case "Opposition":
			responseBody.Aspects[i].Type = 180
			break
		default:
			responseBody.Aspects[i].Type = "null"
			break
		}

		for _, v := range checkData {
			if responseBody.Aspects[i].AspectingPlanetID == v.CheckAspectingID &&
				responseBody.Aspects[i].AspectedPlanetID == v.CheckAspectedID {

				if responseBody.Aspects[i].Type == v.CheckType && responseBody.Aspects[i].Type != "null" {
					responseBody.RespMsg += fmt.Sprintf( /*%s\n*/ "%s\n\n" /*v.Soed,*/, v.Body)
				}
			}
		}
	}

	return responseBody, nil
}

func jsonUprReader() (entity.HouseUpr, error) {
	var localHousesUpr entity.HouseUpr

	jsonFile, err := os.Open(housesUpr)
	if err != nil {
		logrus.Errorf("Cannot open the file: %s, due to error: %s", housesUpr, err.Error())
		return entity.HouseUpr{}, err
	}
	defer jsonFile.Close()

	byteData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logrus.Errorf("Cannot read the file: %s, due to error: %s", housesUpr, err.Error())
		return entity.HouseUpr{}, err
	}

	//byteData = bytes.TrimPrefix(byteData, []byte("\xef\xbb\xbf"))

	err = json.Unmarshal(byteData, &localHousesUpr)
	if err != nil {
		logrus.Errorf("Cannot unmaeshal the file: %s, due to error: %s", housesUpr, err.Error())
		return entity.HouseUpr{}, err
	}

	return localHousesUpr, nil
}

// JsonPowerReader
// @Tags 		json planets power reader
// @Description getting data from planets_power.json
// @Accept  	json
// @Produce 	entity.PlanetPower
// @Param   	localPlanetsPower entity.PlanetPower
// @Success 	{object} entity.PlanetPower
// @Failure 	{object} error
func JsonPowerReader() (entity.PlanetPower, error) {
	var localPlanetPower entity.PlanetPower

	jsonFile, err := os.Open(planetsPower)
	if err != nil {
		logrus.Errorf("Cannot open the file: %s, due to error: %s", housesUpr, err.Error())
		return entity.PlanetPower{}, err
	}
	defer jsonFile.Close()

	byteData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logrus.Errorf("Cannot read the file: %s, due to error: %s", housesUpr, err.Error())
		return entity.PlanetPower{}, err
	}

	//byteData = bytes.TrimPrefix(byteData, []byte("\xef\xbb\xbf")) - magic method

	err = json.Unmarshal(byteData, &localPlanetPower)
	if err != nil {
		logrus.Errorf("Cannot unmaeshal the file: %s, due to error: %s", housesUpr, err.Error())
		return entity.PlanetPower{}, err
	}

	return localPlanetPower, nil
}

// txtReader
// @Tags 		txt reader
// @Description getting data from Aspects.txt and reading it into []string and removing dots and parentheses
// @Accept  	txt
// @Produce 	[]string
// @Param   	lines []string
// @Success 	{object} []string
// @Failure 	{object} error
func txtReader(fname string) ([]string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := 0; i < len(lines); i++ {
		switch i {
		default:
			lines[i] = strings.Replace(lines[i], ".", "", -1)
			lines[i] = strings.Replace(lines[i], "[", "", -1)
			lines[i] = strings.Replace(lines[i], "]", "", -1)
			i += 2
			break
		}
	}

	return lines, scanner.Err()
}

// TxtDataWorker
// @Tags 		data worker
// @Description splitting aspects from txt file and fills []entity.CheckVars struct array
// @Accept  	txt
// @Produce 	[]entity.CheckVars
// @Param   	procBody []entity.CheckVars
// @Success 	{object} []entity.CheckVars
// @Failure 	{object} error
func TxtDataWorker() ([]entity.CheckVars, error) {
	procData := make([]entity.CheckVars, 1096)

	localReadData, err := txtReader(filename)
	if err != nil {
		return []entity.CheckVars{}, err
	}

	var counter = 0

	for i := 0; i < len(localReadData); i += 3 {
		var aspType int
		var aspectedID int

		aspectingID, _ := strconv.Atoi(localReadData[i][:2])
		aspType, _ = strconv.Atoi(localReadData[i][2:5])
		aspectedID, _ = strconv.Atoi(localReadData[i][5:len(localReadData[i])])

		procData[counter].CheckAspectingID = aspectingID
		procData[counter].CheckType = aspType
		procData[counter].CheckAspectedID = aspectedID

		counter++
	}
	counter = 0
	for i := 1; i < len(localReadData); i += 2 {
		procData[counter].Soed = localReadData[i]
		i++
		procData[counter].Body = localReadData[i]
		counter++
	}

	return procData, nil
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
