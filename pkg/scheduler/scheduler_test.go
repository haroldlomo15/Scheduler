package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testData = `[
  {"id":1,"trainer_id":1,"starts_at":"2020-01-24T09:00:00-08:00","ends_at":"2020-01-24T09:30:00-08:00"},
  {"id":2,"trainer_id":1,"starts_at":"2020-01-24T10:00:00-08:00","ends_at":"2020-01-24T10:30:00-08:00"},
  {"id":3,"trainer_id":1,"starts_at":"2020-01-25T10:00:00-08:00","ends_at":"2020-01-25T10:30:00-08:00"},
  {"id":4,"trainer_id":2,"starts_at":"2020-01-24T09:00:00-08:00","ends_at":"2020-01-24T09:30:00-08:00"},
  {"id":5,"trainer_id":2,"starts_at":"2020-01-26T10:00:00-08:00","ends_at":"2020-01-26T10:30:00-08:00"},
  {"id":6,"trainer_id":3,"starts_at":"2020-01-26T12:00:00-08:00","ends_at":"2020-01-26T12:30:00-08:00"}
]
`

func TestGetTrainerAppointments(t *testing.T) {
	type testDef struct {
		name      string
		trainerId int
		expected []Appointment
	}

	tests := []testDef{
		{name: "trainer 2 scheduled appointments", trainerId: 2, expected: []Appointment{
			{
				Id:        4,
				TrainerId: 2,
				StartsAt:  "2020-01-24T09:00:00-08:00",
				EndsAt:    "2020-01-24T09:30:00-08:00",
			},
			{
				Id:        5,
				TrainerId: 2,
				StartsAt:  "2020-01-26T10:00:00-08:00",
				EndsAt:    "2020-01-26T10:30:00-08:00",
			},
		}},
		{name: "trainer 3 scheduled appointments", trainerId: 3, expected: []Appointment{
			{
				Id:        6,
				TrainerId: 3,
				StartsAt:  "2020-01-26T12:00:00-08:00",
				EndsAt:    "2020-01-26T12:30:00-08:00",
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appointment, err := getAppointmetTestData()
			assert.NoError(t, err, "error getting appointment test data")
			result := getTrainerAppointments(appointment, tt.trainerId)
			assert.Equal(t, tt.expected, result)
			fmt.Println(result)
		})
	}

}

func TestGetTrainerAvailableDateTime(t *testing.T) {

	type testDef struct {
		name      string
		trainerId int
		startsAt  string
		endsAt    string
		expected  []string
	}

	tests := []testDef{
		{name: "no appointments", trainerId: 1, startsAt: "2020-01-06T08:00:00-08:00", endsAt: "2020-01-06T08:00:00-08:00", expected: []string(nil)},
		{name: "one appointment", trainerId: 1, startsAt: "2020-01-06T08:00:00-08:00", endsAt: "2020-01-06T08:30:00-08:00", expected: []string{"2020-01-06T08:00:00-08:00"}},
		{name: "two appointment", trainerId: 1, startsAt: "2020-01-06T08:00:00-08:00", endsAt: "2020-01-06T09:00:00-08:00", expected: []string{"2020-01-06T08:00:00-08:00", "2020-01-06T08:30:00-08:00"}},
		{name: "appointment already scheduled for 2020-01-24T09:00:00-08:00", trainerId: 1, startsAt: "2020-01-24T09:00:00-08:00", endsAt: "2020-01-24T10:00:00-08:00", expected: []string{"2020-01-24T09:30:00-08:00"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appointment, err := getAppointmetTestData()
			assert.NoError(t, err, "error getting appointment test data")

			result, err := getTrainerAvailableDateTime(appointment, tt.trainerId, tt.startsAt, tt.endsAt)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, *result)

		})
	}
}

func getAppointmetTestData() ([]Appointment, error) {
	var appointment []Appointment
	err := json.Unmarshal([]byte(testData), &appointment)
	if err != nil {
		return appointment, err
	}
	return appointment, nil
}
