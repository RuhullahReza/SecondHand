package db

import (
	"log"
	"github.com/RuhullahReza/SecondHand/util/config"

	"github.com/cloudinary/cloudinary-go"
)

func NewCloudinaryConnection() *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(config.EnvCloudName(), config.EnvCloudAPIKey(), config.EnvCloudAPISecret())
	if err != nil {
		log.Fatalln(err)
	}

	return cld
}




