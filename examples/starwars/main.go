package main

import (
	"github.com/lukewilliamboswell/graphql"
	"log"
	"os"
)

func main() {

	episodeEnum := graphql.EnumType{
		Name:        "Episode",
		Description: "One of the films in the Star Wars Trilogy",
		Values: []graphql.EnumValue{
			graphql.EnumValue{Name: "NEWHOPE", Description: "Released in 1977."},
			graphql.EnumValue{Name: "EMPIRE", Description: "Released in 1980."},
			graphql.EnumValue{Name: "JEDI", Description: "Released in 1983."},
		},
	}

	humanObject := graphql.ObjectType{
		Name:        "Human",
		Description: "Represents a human.",
		Children: []graphql.ObjectChild{
			graphql.ScalarType{
				Type: graphql.GWL_PRIMITIVE_STRING,
				Name: "id",
			},
			graphql.ScalarType{
				Type: graphql.GWL_PRIMITIVE_STRING,
				Name: "name",
			},
		},
	}

	log.Println(episodeEnum, episodeEnum.IsValid())
	log.Println(humanObject)

	episodeEnum.MarshalGraphQL(os.Stdout)
	humanObject.MarshalGraphQL(os.Stdout)

}
