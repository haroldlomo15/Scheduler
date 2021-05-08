# scheduler

Appointment Scheduler to post a schedule appointment and to retrieve a trainer's appointments from an `appointments.json` file
which serves as a db. The file can be found under `./pkg/db/appointments.json` .


### Setup
- Using Goland or any IDE that supports Go
- Run the project by running the main.go file under `cmd/scheduler/main.go`

Host url: http://localhost:7005/v1

### Endpoints
- [Health](#health)
- [GetScheduled](#getscheduled)
- [GetAppointments](#getappointments)
- [PostAppointment](#postappointment)

<br>

#### Health
Returns `ok` if service can accept requests

<br>

#### GetScheduled
Returns all scheduled trainer's appointments

`GET` /getscheduled?trainer_id={ trainer_id }

Sample Response
```
[
  {
    "Id": 1,
    "trainer_id": 1,
    "starts_at": "2020-01-24T09:00:00-08:00",
    "ends_at": "2020-01-24T09:30:00-08:00"
  },
  {
    "Id": 2,
    "trainer_id": 1,
    "starts_at": "2020-01-24T10:00:00-08:00",
    "ends_at": "2020-01-24T10:30:00-08:00"
  }
]
```
<br>

#### GetAppointments

Returns a list of available appointment times for a trainer between two dates

`GET` /getappointments?trainer_id=1&starts_at={ start_date_time }&ends_at={ end_date_time }

Sample Request & Response

`/getappointments?trainer_id=1&starts_at=2020-01-06T08:00:00-08:00&ends_at=2021-01-26T10:30:00-08:00`

```
[
  "2020-01-06T08:00:00-08:00",
  "2020-01-06T08:30:00-08:00",
  "2020-01-06T09:00:00-08:00",
  "2020-01-06T09:30:00-08:00",
  "2020-01-06T10:00:00-08:00",
  "2020-01-06T10:30:00-08:00",
  "2020-01-06T11:00:00-08:00",
  "2020-01-06T11:30:00-08:00",
  "2020-01-06T12:00:00-08:00",
  "2020-01-06T12:30:00-08:00",
  "2020-01-06T13:00:00-08:00",
  "2020-01-06T13:30:00-08:00",
  "2020-01-06T14:00:00-08:00",
  "2020-01-06T14:30:00-08:00",
]
```
<br>

#### PostAppointment
Post a schedule appointment

`POST` /postappointment

Sample Request Body
```
{
    "trainer_id": 9,
    "user_id": 2,
    "starts_at": "2019-01-25T09:00:00-08:00",
    "ends_at": "2019-01-25T09:30:00-08:00"
}
```
