package ccs

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)


func createEmbed(Project *Project) (*discordgo.MessageEmbed) {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Status",
			Value:  Project.State,
			Inline: false,
		},
		{
			Name:   "Target",
			Value:  fmt.Sprintf("%0.2f VTC", Project.TargetAmount),
			Inline: true,
		},
		{
			Name:   "Raised",
			Value:  fmt.Sprintf("%0.2f VTC", Project.RaisedAmount),
			Inline: true,
		},
	}

	if Project.State == "COMPLETED" {
		fields[2].Value =  fmt.Sprintf("%0.2f VTC", Project.TargetAmount)
	}

	// Makes the number an int if it's the same value as the float
	if float64(int(Project.TargetAmount)) == Project.TargetAmount {
		fields[1].Value = fmt.Sprintf("%d VTC", int(Project.TargetAmount))
		if Project.State == "COMPLETED" {
			fields[2].Value = fmt.Sprintf("%d VTC", int(Project.TargetAmount))
		}
	}

	if Project.State == "FUNDING-REQUIRED" {
		fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Donation address",
				Value:  fmt.Sprintf("``%s``", Project.Address),
				Inline: false,
			})
	}

	return &discordgo.MessageEmbed{
		Title:       Project.Title,
		//URL:		 ,
		Description: "",
		Timestamp:   convertCCSDate(Project.Date),
		Color:       colorOfEmbed(Project.State),
		Fields:      fields,
		Footer:		 &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Author: %s", Project.Author),},
	}
}

type EmbedColor struct {
	Ideas 				int
	FundingRequired 	int
	WorkInProgress  	int
	Completed 			int
}

func colorOfEmbed(state string ) int {
	embedColor := EmbedColor{
		Ideas: 15844367,
		FundingRequired: 15158332,
		WorkInProgress: 206694,
		Completed: 296535,
	}

	switch {
	case state == "COMPLETED":
		return embedColor.Completed
	
	case state == "FUNDING-REQUIRED":
		return embedColor.FundingRequired
	
	case state == "WORK-IN-PROGRESS":
		return embedColor.WorkInProgress
	
	case state == "IDEA":
		return embedColor.Ideas
	}
	return 0
}

// Discord uses the ISO8601 standart for dates, we need to convert to that.
func convertCCSDate(date string) string {
	dateSplit := strings.SplitN(date, " ", len(date))
	for i := 0; i < len(dateSplit); i++ {
		dateSplit[i] = strings.TrimSuffix(dateSplit[i], ",")
	}
	switch {
	case dateSplit[0] == "January":
		dateSplit[0] = "01"

	case dateSplit[0] == "February":
		dateSplit[0] = "02"

	case dateSplit[0] == "March":
		dateSplit[0] = "03"

	case dateSplit[0] == "April":
		dateSplit[0] = "04"

	case dateSplit[0] == "May":
		dateSplit[0] = "05"

	case dateSplit[0] == "June":
		dateSplit[0] = "06"

	case dateSplit[0] == "July":
		dateSplit[0] = "07"

	case dateSplit[0] == "August":
		dateSplit[0] = "08"

	case dateSplit[0] == "September":
		dateSplit[0] = "09"

	case dateSplit[0] == "October":
		dateSplit[0] = "10"

	case dateSplit[0] == "November":
		dateSplit[0] = "11"

	case dateSplit[0] == "December":
		dateSplit[0] = "12"

	default:
	}

	return fmt.Sprintf("%s-%s-%s", dateSplit[2], dateSplit[0], dateSplit[1])
}