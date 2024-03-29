package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DbSettings struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
	SSLMode  string
}

type ParserSettings struct {
	Link string
	Path string
	Num  int
}

type Settings struct {
	DB     *DbSettings
	Parser *ParserSettings
}

func InitConfig() error {
	pathToEnvFile, ok := os.LookupEnv("PATH_TO_ENV_FILE")
	if !ok {
		return fmt.Errorf("env PATH_TO_ENV_FILE is not set")
	}

	err := godotenv.Load(pathToEnvFile)
	if err != nil {
		return fmt.Errorf("error load env variables: %w", err)
	}

	return nil
}

func ReadConfig() (Settings, error) {
	dbSettings, err := readDbConfig()
	if err != nil {
		return Settings{}, err
	}

	parserSettings, err := readParserConfig()
	if err != nil {
		return Settings{}, err
	}

	return Settings{
		DB:     &dbSettings,
		Parser: &parserSettings,
	}, nil
}

func readDbConfig() (DbSettings, error) {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return DbSettings{}, fmt.Errorf("env DB_HOST is not set")
	}

	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		return DbSettings{}, fmt.Errorf("env DB_USER is not set")
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return DbSettings{}, fmt.Errorf("env DB_PASSWORD is not set")
	}

	name, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return DbSettings{}, fmt.Errorf("env DB_NAME is not set")
	}

	_port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return DbSettings{}, fmt.Errorf("env DB_PORT is not set")
	}
	port, err := strconv.Atoi(_port)
	if err != nil {
		return DbSettings{}, fmt.Errorf("convert port to int error: %w", err)
	}

	sslMode, ok := os.LookupEnv("DB_SSL_MODE")
	if !ok {
		return DbSettings{}, fmt.Errorf("env DB_SSL_MODE is not set")
	}

	return DbSettings{
		Host:     host,
		User:     user,
		Password: password,
		Name:     name,
		Port:     port,
		SSLMode:  sslMode,
	}, nil
}

func readParserConfig() (ParserSettings, error) {
	var (
		link string
		path string
	)

	link, ok := os.LookupEnv("PARSE_LINK")
	if !ok {
		path, ok = os.LookupEnv("PARSE_PATH")
		if !ok {
			return ParserSettings{}, fmt.Errorf("one of env PARSE_LINK and PARSE_PATH should be set")
		}
	}

	_num, ok := os.LookupEnv("NUMBER_OF_PRODUCTS")
	if !ok {
		return ParserSettings{}, fmt.Errorf("env NUMBER_OF_PRODUCTS is not set")
	}
	num, err := strconv.Atoi(_num)
	if err != nil {
		return ParserSettings{}, err
	}

	return ParserSettings{
		Link: link,
		Path: path,
		Num:  num,
	}, nil
}
