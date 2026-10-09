package main

import (
	"log"

	"github.com/iLert/ilert-go/v4"
)

func main() {
	var apiToken = "your API token"
	client := ilert.NewClient(ilert.WithAPIToken(apiToken))

	createSupportHourInput := &ilert.CreateSupportHourInput{
		SupportHour: &ilert.SupportHour{
			Name:     "example",
			Timezone: "Europe/Berlin",
			SupportDays: &ilert.SupportDays{
				MONDAY: &ilert.SupportDay{
					Start: "09:00",
					End:   "18:00",
				},
				WEDNESDAY: &ilert.SupportDay{
					Start: "09:00",
					End:   "18:00",
				},
				FRIDAY: &ilert.SupportDay{
					Start: "09:00",
					End:   "18:00",
				},
			},
		},
	}

	result, err := client.CreateSupportHour(createSupportHourInput)
	if err != nil {
		log.Println(result)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Support hour:\n\n %+v\n", result.SupportHour)

	// coverage that one window per day cannot hold, a split shift on Monday and the weekend from
	// Friday 17:00 to Monday 08:00, is sent as windows, without support days next to them
	supportWindows := []ilert.SupportWindow{
		{
			From: &ilert.TimeOfWeek{DayOfWeek: ilert.DayOfWeek.Monday, Time: "09:00"},
			To:   &ilert.TimeOfWeek{DayOfWeek: ilert.DayOfWeek.Monday, Time: "12:00"},
		},
		{
			From: &ilert.TimeOfWeek{DayOfWeek: ilert.DayOfWeek.Monday, Time: "13:00"},
			To:   &ilert.TimeOfWeek{DayOfWeek: ilert.DayOfWeek.Monday, Time: "17:00"},
		},
		{
			From: &ilert.TimeOfWeek{DayOfWeek: ilert.DayOfWeek.Friday, Time: "17:00"},
			To:   &ilert.TimeOfWeek{DayOfWeek: ilert.DayOfWeek.Monday, Time: "08:00"},
		},
	}
	windowsResult, err := client.CreateSupportHour(&ilert.CreateSupportHourInput{
		SupportHour: &ilert.SupportHour{
			Name:           "example with windows",
			Timezone:       "Europe/Berlin",
			SupportWindows: &supportWindows,
		},
	})
	if err != nil {
		log.Println(windowsResult)
		log.Fatalln("ERROR:", err)
	}
	log.Printf("Support hour windows:\n\n %+v\n", *windowsResult.SupportHour.SupportWindows)
}
