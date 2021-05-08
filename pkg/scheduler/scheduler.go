package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/haroldlomo15/scheduler/pkg/db"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const timeFormat = "2006-01-02T15:04:05-07:00"

type Appointment struct {
	Id        int    `json:"Id,omitempty"`
	UserId        int    `json:"user_id,omitempty"`
	TrainerId int    `json:"trainer_id"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
}

func GetAvailableAppointment(write http.ResponseWriter, req *http.Request) {
	// Gets a list of available appointment times for a trainer between two dates
	trainerId, startsAt, endsAt, err := validateAppointmentReq(req)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
		return
	}

	appointments, err := queryDBForAppointments()
	if err != nil {
		log.Error(err)
		http.Error(write, err.Error(), http.StatusInternalServerError)
	}

	availableAppointments, err := getTrainerAvailableDateTime(appointments, trainerId, startsAt, endsAt)
	if err != nil {
		http.Error(write, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(write).Encode(availableAppointments)
}


func GetScheduledAppointment(write http.ResponseWriter, req *http.Request) {
	// Gets a list of scheduled appointments for a trainer
	trainerId, ok := req.URL.Query()["trainer_id"]
	if !ok {
		http.Error(write, "trainer_id is required", http.StatusBadRequest)
		return
	}
	trainerIdVal, err := strconv.Atoi(trainerId[0])
	if err != nil {
		http.Error(write, "trainer_id must be an integer", http.StatusBadRequest)
		return
	}

	appointments, err := queryDBForAppointments()
	if err != nil {
		log.Error(err)
		http.Error(write, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(write).Encode(getTrainerAppointments(appointments, trainerIdVal))
}

func PostAppointment(write http.ResponseWriter, req *http.Request) {
	var appointmentBody Appointment
	if req.Method != "POST" {
		http.Error(write, "POST method only", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&appointmentBody)
	if err != nil {
		log.Errorf("error decoding req: %v", err)
		http.Error(write, "Bad Request", http.StatusBadRequest)
		return
	}

	appointments, err := queryDBForAppointments()
	if err != nil {
		log.Error(err)
		http.Error(write, err.Error(), http.StatusInternalServerError)
	}

	appointments = append(appointments, appointmentBody)
	dataBytes, err := json.Marshal(appointments)
	if err != nil {
		log.Error(err)
		http.Error(write, err.Error(), http.StatusInternalServerError)
	}

	err = db.WriteToDbFile(dataBytes)
	if err != nil {
		log.Error(err)
		http.Error(write, err.Error(), http.StatusInternalServerError)
	}
	write.Write([]byte("posted successfully"))
}

func getTrainerAvailableDateTime(appointments []Appointment, trainerId int, startsAt, endsAt string) (*[]string, error) {

	trainerAppointments := getTrainerAppointments(appointments, trainerId)
	trainerAppointmentsMap := make(map[int64]int64) //
	for _, val := range trainerAppointments {
		startTime, err := time.Parse(timeFormat, val.StartsAt)
		if err != nil {
			log.Error(err)
			continue
		}
		endTime, err := time.Parse(timeFormat, val.EndsAt)
		if err != nil {
			log.Error(err)
			continue
		}
		trainerAppointmentsMap[startTime.Unix()] = endTime.Unix()
	}

	start, end, err := parseStartsAtEndsAtTime(startsAt, endsAt)
	if err != nil {
		return nil, err
	}

	result := getAvailableResult(*start, *end, trainerAppointmentsMap)
	return &result, nil
}

func getAvailableResult(start, end time.Time, trainerAppointmentsMap map[int64]int64) []string {
	// Loop through startTime until endTime
	// Removing trainer scheduled appointment time
	// Removing non business hours M-F 8am - 5pm
	var result []string
	availableTime := start

	for !availableTime.After(end) {
		if _, ok := trainerAppointmentsMap[availableTime.Unix()]; ok {
			// Trainer has an appointment so skip appointment date
			availableTime = availableTime.Add(30 * time.Minute)
			continue
		}

		if int(availableTime.Weekday()) == 0 || int(availableTime.Weekday()) == 6 {
			// Time is sunday or saturday we skip appointment date
			availableTime = availableTime.Add(30 * time.Minute)
			continue
		}

		hourTime := availableTime.Hour()
		if hourTime < 8 || hourTime > 16 {
			// hour time is not between 8 - 5 so we skip appointment date
			availableTime = availableTime.Add(30 * time.Minute)
			continue
		}

		if len(result) > 999 {
			// Break loop after a 1000 appointmwnt
			log.Info("1000 available appointment result reached")
			break
		}

		result = append(result, availableTime.Format(timeFormat))
		availableTime = availableTime.Add(30 * time.Minute)
	}
	return result
}

func parseStartsAtEndsAtTime(startsAt, endsAt string) (*time.Time, *time.Time, error) {
	startTime, err := time.Parse(timeFormat, startsAt)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing startsAt time - %v", err)
	}

	endTime, err := time.Parse(timeFormat, endsAt)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing endsAt time - %v", err)
	}

	err = validateStartsAtEndsAtTime(startTime, endTime)
	if err != nil {
		return nil, nil, err
	}

	return &startTime, &endTime, nil
}

func validateStartsAtEndsAtTime(startsAt, endsAt time.Time) error {
	if startsAt.Minute() != 0 && startsAt.Minute() != 30 {
		return fmt.Errorf("endsAt time minutes must be :00 or :30 ")
	}

	if endsAt.Minute() != 0 && endsAt.Minute() != 30 {
		return fmt.Errorf("endsAt time minutes must be :00 or :30 ")
	}
	return nil
}

func validateAppointmentReq(req *http.Request) (int, string, string, error) {
	trainerId, ok := req.URL.Query()["trainer_id"]
	if !ok {
		return 0, "", "", fmt.Errorf("trainer_id is required")
	}
	startsAt, ok := req.URL.Query()["starts_at"]
	if !ok {
		return 0, "", "", fmt.Errorf("starts_at is required")
	}
	endsAt, ok := req.URL.Query()["ends_at"]
	if !ok {
		return 0, "", "", fmt.Errorf("ends_at is required")
	}
	trainerIdVal, err := strconv.Atoi(trainerId[0])
	if err != nil {
		return 0, "", "", fmt.Errorf("trainer_id must be an integer - %v", err)
	}
	return trainerIdVal, startsAt[0], endsAt[0], nil
}


func getTrainerAppointments(appointments []Appointment, trainerId int) []Appointment {
	var trainerAppointments []Appointment
	for _, val := range appointments {
		if val.TrainerId == trainerId {
			trainerAppointments = append(trainerAppointments, val)
		}
	}
	return trainerAppointments
}

func queryDBForAppointments() ([]Appointment, error) {
	var appointments []Appointment

	data, err := db.GetAppointments()
	if err != nil {
		return appointments, err
	}

	err = json.Unmarshal(data, &appointments)
	if err != nil {
		return appointments, fmt.Errorf("error unmarshalling appointment - %v", err)
	}
	return appointments, nil
}
