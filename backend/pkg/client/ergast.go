package client

import (
	"blackmichael/f1-pickem/pkg/domain"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type RaceDataClient interface {
	GetRaceResults(ctx context.Context, season, raceNumber string) (*domain.RaceResults, error)
	GetRaces(ctx context.Context, season string) (domain.Races, error)
}

type ergastClient struct {
	baseUrl string
}

func NewErgastClient(baseUrl string) RaceDataClient {
	return &ergastClient{
		baseUrl: baseUrl,
	}
}

// generated by https://mholt.github.io/json-to-go/
type raceResultsResponse struct {
	MRData struct {
		Xmlns     string `json:"xmlns"`
		Series    string `json:"series"`
		Limit     string `json:"limit"`
		Offset    string `json:"offset"`
		Total     string `json:"total"`
		RaceTable struct {
			Season string `json:"season"`
			Round  string `json:"round"`
			Races  []struct {
				Season   string `json:"season"`
				Round    string `json:"round"`
				URL      string `json:"url"`
				RaceName string `json:"raceName"`
				Circuit  struct {
					CircuitID   string `json:"circuitId"`
					URL         string `json:"url"`
					CircuitName string `json:"circuitName"`
					Location    struct {
						Lat      string `json:"lat"`
						Long     string `json:"long"`
						Locality string `json:"locality"`
						Country  string `json:"country"`
					} `json:"Location"`
				} `json:"Circuit"`
				Date    string `json:"date"`
				Time    string `json:"time"`
				Results []struct {
					Number       string `json:"number"`
					Position     string `json:"position"`
					PositionText string `json:"positionText"`
					Points       string `json:"points"`
					Driver       struct {
						DriverID        string `json:"driverId"`
						PermanentNumber string `json:"permanentNumber"`
						Code            string `json:"code"`
						URL             string `json:"url"`
						GivenName       string `json:"givenName"`
						FamilyName      string `json:"familyName"`
						DateOfBirth     string `json:"dateOfBirth"`
						Nationality     string `json:"nationality"`
					} `json:"Driver"`
					Constructor struct {
						ConstructorID string `json:"constructorId"`
						URL           string `json:"url"`
						Name          string `json:"name"`
						Nationality   string `json:"nationality"`
					} `json:"Constructor"`
					Grid   string `json:"grid"`
					Laps   string `json:"laps"`
					Status string `json:"status"`
					Time   struct {
						Millis string `json:"millis"`
						Time   string `json:"time"`
					} `json:"Time"`
					FastestLap struct {
						Rank string `json:"rank"`
						Lap  string `json:"lap"`
						Time struct {
							Time string `json:"time"`
						} `json:"Time"`
						AverageSpeed struct {
							Units string `json:"units"`
							Speed string `json:"speed"`
						} `json:"AverageSpeed"`
					} `json:"FastestLap"`
				} `json:"Results"`
			} `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

func (api ergastClient) GetRaceResults(ctx context.Context, season, raceNumber string) (*domain.RaceResults, error) {
	log.Printf("fetching race results for season:%s, race_number:%s\n", season, raceNumber)

	client := http.DefaultClient
	url := fmt.Sprintf("%s/api/f1/%s/%s/results.json", api.baseUrl, season, raceNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Printf("ERROR: failed to fetch race results, code: %d\n", resp.StatusCode)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: unable to read response body (%s)\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var raceResults raceResultsResponse
	if err := json.Unmarshal(body, &raceResults); err != nil {
		log.Printf("ERROR: unable to parse response body (%s)\n", err.Error())
		return nil, err
	}

	// if the results aren't available then this is probably empty
	if len(raceResults.MRData.RaceTable.Races) == 0 {
		return nil, nil
	}

	if len(raceResults.MRData.RaceTable.Races) != 1 {
		log.Printf("ERROR: unexpected number of races, results: %#v\n", raceResults)
		return nil, errors.New("unexpected number of races")
	}

	race := raceResults.MRData.RaceTable.Races[0]
	// if the results aren't available then this might be empty instead
	if len(race.Results) == 0 {
		return nil, nil
	}

	if len(race.Results) != 20 {
		log.Printf("ERROR: unexpected number of race results, results: %#v\n", raceResults)
		return nil, errors.New("unexpected number of race results")
	}

	results := make([]string, 20, 20)
	for _, result := range race.Results {
		position, err := strconv.Atoi(result.Position)
		if err != nil {
			log.Printf("ERROR: failed to parse race result position, results: %#v\n", raceResults)
			return nil, errors.Wrap(err, "failed to parse position")
		}

		if position < 1 || position > 20 {
			log.Printf("ERROR: invalid position found, results: %#v\n", raceResults)
			return nil, errors.New("invalid position found")
		}

		name := fmt.Sprintf("%s %s", result.Driver.GivenName, result.Driver.FamilyName)
		results[position-1] = name
	}

	return &domain.RaceResults{
		Season:     season,
		RaceNumber: raceNumber,
		RaceDate:   race.Date,
		Results:    results,
	}, nil
}

type racesResponse struct {
	MRData struct {
		RaceTable struct {
			Races []struct {
				Circuit struct {
					Location struct {
						Country  string `json:"country"`
						Lat      string `json:"lat"`
						Locality string `json:"locality"`
						Long     string `json:"long"`
					} `json:"Location"`
					CircuitID   string `json:"circuitId"`
					CircuitName string `json:"circuitName"`
					URL         string `json:"url"`
				} `json:"Circuit"`
				FirstPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"FirstPractice"`
				Qualifying struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"Qualifying"`
				SecondPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"SecondPractice"`
				Sprint struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"Sprint"`
				ThirdPractice struct {
					Date string `json:"date"`
					Time string `json:"time"`
				} `json:"ThirdPractice"`
				Date     string `json:"date"`
				RaceName string `json:"raceName"`
				Round    string `json:"round"`
				Season   string `json:"season"`
				Time     string `json:"time"`
				URL      string `json:"url"`
			} `json:"Races"`
			Season string `json:"season"`
		} `json:"RaceTable"`
		Limit  string `json:"limit"`
		Offset string `json:"offset"`
		Series string `json:"series"`
		Total  string `json:"total"`
		URL    string `json:"url"`
		Xmlns  string `json:"xmlns"`
	} `json:"MRData"`
}

func (api ergastClient) GetRaces(ctx context.Context, season string) (domain.Races, error) {
	log.Printf("fetching races for season:%s\n", season)

	client := http.DefaultClient
	url := fmt.Sprintf("%s/api/f1/%s.json", api.baseUrl, season)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Printf("ERROR: failed to fetch races, code: %d\n", resp.StatusCode)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: unable to read response body (%s)\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var _racesResponse racesResponse
	if err := json.Unmarshal(body, &_racesResponse); err != nil {
		log.Printf("ERROR: unable to parse response body (%s)\n", err.Error())
		return nil, err
	}

	numOfRaces := len(_racesResponse.MRData.RaceTable.Races)
	if numOfRaces == 0 {
		return nil, nil
	}

	// TODO handle this better
	if _racesResponse.MRData.Total >= _racesResponse.MRData.Limit {
		log.Printf("unhandled paginated race schedule, season: %s\n", season)
		return nil, errors.New("unhandled paginated result")
	}

	races := make(domain.Races, numOfRaces, numOfRaces)
	for i, race := range _racesResponse.MRData.RaceTable.Races {
		_, err := strconv.Atoi(race.Round)
		if err != nil {
			log.Printf("ERROR: failed to parse race round, race: %#v\n", race)
			return nil, errors.Wrap(err, "failed to parse race round")
		}

		races[i] = &domain.Race{
			RaceName:   race.RaceName,
			RaceNumber: race.Round,
			Season:     race.Season,
			RaceDate:   race.Date,
			RaceId:     domain.GetRaceId(race.Season, race.Round),
		}
		datetime := fmt.Sprintf("%sT%s", race.Date, race.Time)
		races[i].StartTime, err = time.Parse(time.RFC3339, datetime)
		if err != nil {
			log.Printf("ERROR: failed to parse race start time: %v\n", race)
			return nil, errors.Wrap(err, "failed to parse race start time")
		}
	}

	return races, nil
}
