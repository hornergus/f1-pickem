package main

import (
	"blackmichael/f1-pickem/pkg/util"
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	UserId string `json:"userId"`
}

type Response struct {
	Leagues []LeagueResponse `json:"leagues"`
}

type LeagueResponse struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	NumOfMembers int    `json:"num_of_members"`
	Season       string `json:"season"`
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := Response{
		Leagues: []LeagueResponse{
			{
				Id:           "1",
				Name:         "Fast Boiz",
				NumOfMembers: 11,
				Season:       "2022",
			},
		},
	}

	respStr, err := json.Marshal(resp)
	if err != nil {
		return util.MessageResponse(500, "failed to render response"), errors.New("unable to serialize response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(respStr),
		Headers:    util.CorsHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
