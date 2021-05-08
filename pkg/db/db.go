package db

import (
	"fmt"
	"io/ioutil"
	"os"
)

const dbFilePath = "./pkg/db/appointments.json"

func New() error {
	_, err := os.OpenFile(dbFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("could not open db file - %v", err)
	}
	return nil

}

func WriteToDbFile(data []byte) error {
	err := ioutil.WriteFile(dbFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing to dbfile - %v", err)
	}
	return nil
}

func GetAppointments() ([]byte, error) {
	bytes, err := ioutil.ReadFile(dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read db file - %v", err)
	}

	return bytes, nil
}
