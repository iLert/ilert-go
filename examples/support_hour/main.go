package main

import (
	"log"

	"github.com/iLert/ilert-go/v3"
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
}
