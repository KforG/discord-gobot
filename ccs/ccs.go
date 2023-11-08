package ccs

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/KforG/discord-price-bot/logging"
	"github.com/KforG/discord-price-bot/util"

	"github.com/bwmarrin/discordgo"
)

const channelID string = "959174405005668422"

func UpdateCCSChannel(dg *discordgo.Session) {
	for {
		allProposals := AllCCSProposals{}
		err := getAllProposals(&allProposals)
		if err != nil {
			time.Sleep(301 * time.Second)
			continue
		}
		logging.Infof("reading CCS JSON db\n")
		messageIDs, err := readJSONdb()
		if err != nil {
			logging.Warnf("This error can potentially be ignored if experienced on first startup: %s", err)
		}

		IDCounter := 0

		// Loop over completed proposals
		for i := 0; i < len(allProposals.Proposal); i++ {
			if allProposals.Proposal[i].State == "COMPLETED" {
				if len(messageIDs) == 0 || len(messageIDs) == IDCounter {
					logging.Infof("Not enough messages found, posting %s\n", allProposals.Proposal[i].Title)
					newMessage := createCCSMessage(allProposals.Proposal[i], dg)
					messageIDs = append(messageIDs, newMessage)
					IDCounter++
					continue
				} else {
					logging.Infof("Editing existing message for %s with messageID: %s\n", allProposals.Proposal[i].Title, messageIDs[IDCounter].ID)
					err = editCCSMessage(messageIDs[IDCounter].ID, allProposals.Proposal[i], dg)
					IDCounter++
					if err != nil {
						logging.Warnf("Could not edit %s\n", messageIDs[IDCounter].ID)
					}
				}
			}
		}

		// Loop over WIP proposals
		for i := 0; i < len(allProposals.Proposal); i++ {
			if allProposals.Proposal[i].State == "WORK-IN-PROGRESS" {
				if len(messageIDs) == 0 || len(messageIDs) == IDCounter {
					logging.Infof("Not enough messages found, posting %s\n", allProposals.Proposal[i].Title)
					newMessage := createCCSMessage(allProposals.Proposal[i], dg)
					messageIDs = append(messageIDs, newMessage)
					IDCounter++
					continue
				} else {
					logging.Infof("Editing existing message for %s with messageID: %s\n", allProposals.Proposal[i].Title, messageIDs[IDCounter].ID)
					err = editCCSMessage(messageIDs[IDCounter].ID, allProposals.Proposal[i], dg)
					IDCounter++
					if err != nil {
						logging.Warnf("Could not edit %s\n", messageIDs[IDCounter].ID)
					}
				}
			}
		}

		// Loop over FR proposals
		for i := 0; i < len(allProposals.Proposal); i++ {
			if allProposals.Proposal[i].State == "FUNDING-REQUIRED" {
				if len(messageIDs) == 0 || len(messageIDs) == IDCounter {
					logging.Infof("Not enough messages found, posting %s\n", allProposals.Proposal[i].Title)
					newMessage := createCCSMessage(allProposals.Proposal[i], dg)
					messageIDs = append(messageIDs, newMessage)
					IDCounter++
					continue
				} else {
					logging.Infof("Editing existing message for %s with messageID: %s\n", allProposals.Proposal[i].Title, messageIDs[IDCounter].ID)
					err = editCCSMessage(messageIDs[IDCounter].ID, allProposals.Proposal[i], dg)
					IDCounter++
					if err != nil {
						logging.Warnf("Could not edit %s\n", messageIDs[IDCounter].ID)
					}
				}
			}
		}

		// Loop over IDEAS proposals
		for i := 0; i < len(allProposals.Proposal); i++ {
			if allProposals.Proposal[i].State == "IDEA" {
				if len(messageIDs) == 0 || len(messageIDs) == IDCounter {
					logging.Infof("Not enough messages found, posting %s\n", allProposals.Proposal[i].Title)
					newMessage := createCCSMessage(allProposals.Proposal[i], dg)
					messageIDs = append(messageIDs, newMessage)
					IDCounter++
					continue
				} else {
					logging.Infof("Editing existing message for %s with messageID: %s\n", allProposals.Proposal[i].Title, messageIDs[IDCounter].ID)
					err = editCCSMessage(messageIDs[IDCounter].ID, allProposals.Proposal[i], dg)
					IDCounter++
					if err != nil {
						logging.Warnf("Could not edit %s\n", messageIDs[IDCounter].ID)
					}
				}
			}
		}

		logging.Infof("writing to JSON db\n")
		err = writeJSONdb(messageIDs)
		if err != nil {
			logging.Errorf("error writing to JSON db: %s", err)
		}
		logging.Infof("Sleep for 30 minutes before next update..")
		time.Sleep(30 * time.Minute)
	}
}

type AllCCSProposals struct {
	Proposal []Project `json:"data"`
}

type Project struct {
	Contributions    int     `json:"contributions"`
	PercentageFunded int     `json:"percentage_funded"`
	RaisedAmount     float64 `json:"raised_amount"`
	TargetAmount     float64 `json:"target_amount"`
	Address          string  `json:"address"`
	Author           string  `json:"author"`
	Date             string  `json:"date"`
	Title            string  `json:"title"`
	State            string  `json:"state"`
}

type Message struct {
	ID string `json:"id"`
}

// https://ccs.vertcoin.io/index.php/projects
// This function fetches all merged projects from ccs.vertcoin.io
func getAllProposals(jsonPayload *AllCCSProposals) error {
	err := util.GetJson("https://ccs.vertcoin.io/index.php/projects", jsonPayload)
	if err != nil {
		logging.Errorf("Error fetching projects from ccs.vertcoin.io\n", err)
		return err
	}
	return nil
}

func createCCSMessage(Project Project, dg *discordgo.Session) Message {
	message, err := dg.ChannelMessageSendEmbed(channelID, createEmbed(&Project))
	if err != nil {
		logging.Errorf("error posting new CCS funding request\n", err)
		return Message{}
	}

	return Message{
		ID: message.ID,
	}
}

func editCCSMessage(messageID string, Project Project, dg *discordgo.Session) error {
	_, err := dg.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      messageID,
		Channel: channelID,
		Embed:   createEmbed(&Project),
	})
	if err != nil {
		logging.Errorf("error editing CCS message ID: %s\n", messageID, err)
		return err
	}
	return nil
}

func readJSONdb() ([]Message, error) {
	messageIDs := []Message{}

	file, err := ioutil.ReadFile("ccsmessage.json")
	if err != nil {
		logging.Errorf("error reading ccsmessage file: %s\n", err)
		return messageIDs, err
	}

	err = json.Unmarshal(file, &messageIDs)

	if err != nil {
		logging.Errorf("error unmarshaling file: %s\n", err)
		return messageIDs, err
	}
	return messageIDs, nil
}

func writeJSONdb(messageIDs []Message) error {
	file, err := json.MarshalIndent(messageIDs, "", " ")
	if err != nil {
		logging.Errorf("error serializing struct to JSON\n", err)
		return err
	}

	err = ioutil.WriteFile("ccsmessage.json", file, 0644)
	if err != nil {
		logging.Errorf("error writing JSON to file\n", err)
		return err
	}

	return nil
}
