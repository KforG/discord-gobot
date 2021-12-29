package functions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KforG/discord-price-bot/logging"
	"github.com/KforG/discord-price-bot/util"
	"github.com/bwmarrin/discordgo"
)

func DetermineCommand(dg *discordgo.Session, message *discordgo.MessageCreate) {

	if message.Author.Bot { //If the message originates from a bot there's no need to respond
		return
	}
	logging.Infof("Message: %s\n", message.Content)
	message.Content = strings.ToLower(message.Content)

	if message.Content == "" { //If the message only contain an image, the content will be an empty string, which will cause a runtime error
		return
	}

	argument := strings.SplitN(message.Content, " ", len(message.Content))

	switch { //Remember to update -help command when a new command is added
	case argument[0] == "-help":
		respondHelp(dg, message)

	case argument[0] == "-halving":
		respondHalving(dg, message)

	case argument[0] == "-whyhalving":
		respondWhyHalving(dg, message)

	case argument[0] == "-ocm":
		respondOCM(dg, message)

	case argument[0] == "-verthash-ocm":
		respondVerthashOCM(dg, message)

	case argument[0] == "-vhm":
		respondVHM(dg, message)

	case argument[0] == "-block":
		respondBlockHeight(dg, message)

	case argument[0] == "-estearnings":
		// If no argument is given with the command, a runtime error will occur
		//We need to check that there is something after the command
		if len(argument) == 1 {
			logging.Errorf("Command for -estearningsdetected, no hashrate argument supplied.. Sending error message to user.")
			response := "No hashrate argument supplied \n?estearnings [hashrate in kH/s] is the correct syntax.\nExample: ``?estearnings 600``"
			_, err := dg.ChannelMessageSendReply(message.ChannelID, response, message.Reference())
			if err != nil {
				logging.Errorf("error responding to estimated earnings error \n", err)
			}
		} else {
			respondEstimatedEarnings(dg, message, argument[1])
		}

	default:
		return
	}
}

func respondHelp(dg *discordgo.Session, message *discordgo.MessageCreate) {
	help := "• ``-help`` - Displays the help command you're currently seeing."
	halving := "\n\n• ``-halving`` - Bot responds with the remaining blocks and a time estimate before the next halving."
	whyhalving := "\n\n• ``-whyhalving`` - Bot responds with a message detailing why halvings takes place."
	ocm := "\n\n• ``-ocm`` - Sends a link to the newest release of the Vertcoin One-Click-Miner."
	vertocm := "\n\n• ``-verthash-ocm`` - Sends a link to the newest release of the Verthash One-Click-Miner."
	vhm := "\n\n• ``-vhm`` - Sends a link to the newest release of Verthashminer."
	block := "\n\n• ``-block`` - Bot sends a message with the current block height."
	estearnings := "\n\n• ``-estearnings [hashrate in kH/s]`` - Bot will calculate the estimated number of VTC a miner will recieve per day with the supplied hashrate."
	_, err := dg.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
		Title:       "Commands",
		Description: help + halving + whyhalving + ocm + vertocm + vhm + block + estearnings,
		Color:       296535,
	})
	if err != nil {
		logging.Errorf("error responding to help command\n", err)
	}
}

func respondHalving(dg *discordgo.Session, message *discordgo.MessageCreate) {
	halvingInterval := 840000

	//get current blockheight
	blockHeight, _, err := util.GetNetworkStats()
	if err != nil {
		logging.Errorf("Error fetching network stats")
		response := "Unable to get current block height"
		_, err := dg.ChannelMessageSend(message.ChannelID, response)
		if err != nil {
			logging.Errorf("error responding to halving command\n", err)
		}
		return
	}

	halvings := blockHeight / int64(halvingInterval)

	if halvings >= 64 { //Vertcoin and Bitcoin will undergo 64 halvings before zero emission
		response := "Vertcoin has completed it's last halving and is no longer emitting coins"
		_, err := dg.ChannelMessageSend(message.ChannelID, response)
		if err != nil {
			logging.Errorf("error responding to halving command\n", err)
		}
	} else {
		halvings++ //We can't divide by zero and it's easier to do this with the rest of the code anyway

		nextHalvingBlockHeight := int64(halvingInterval) * halvings
		blocksRemain := nextHalvingBlockHeight - blockHeight

		response := fmt.Sprintf("%d blocks remaining before halving!\nEstimated time left before halving: %s", blocksRemain, util.TimeBeforeHalving(blocksRemain))

		_, err := dg.ChannelMessageSend(message.ChannelID, response)
		if err != nil {
			logging.Errorf("error responding to halving command\n", err)
		}
	}
}

func respondWhyHalving(dg *discordgo.Session, message *discordgo.MessageCreate) {
	response1 := `
	Vertcoin like many cryptocurrencies are designed as a deflationary currency. Vertcoin accomplishes this by reducing the blockreward
by half every 840,000 blocks (approximately every 4 years), after a number of halvings the emission of new coins will be zero. 
	`
	response2 := "\nLearn more here: https://alvie.github.io/vertcoin-block-height-live/"
	response := response1 + response2
	_, err := dg.ChannelMessageSend(message.ChannelID, response)
	if err != nil {
		logging.Errorf("error responding to whyhalving command\n", err)
	}
}

func respondOCM(dg *discordgo.Session, message *discordgo.MessageCreate) {
	response := "https://github.com/vertcoin-project/one-click-miner-vnext/releases"
	_, err := dg.ChannelMessageSend(message.ChannelID, response)
	if err != nil {
		logging.Errorf("error responding to OCM command\n", err)
	}
}

func respondVerthashOCM(dg *discordgo.Session, message *discordgo.MessageCreate) {
	response := "https://github.com/vertiond/verthash-one-click-miner/releases"
	_, err := dg.ChannelMessageSend(message.ChannelID, response)
	if err != nil {
		logging.Errorf("error responding to Verthash-OCM command\n", err)
	}
}

func respondVHM(dg *discordgo.Session, message *discordgo.MessageCreate) {
	response := "https://github.com/CryptoGraphics/VerthashMiner/releases"
	_, err := dg.ChannelMessageSend(message.ChannelID, response)
	if err != nil {
		logging.Errorf("error responding to VHM command\n", err)
	}
}

func respondBlockHeight(dg *discordgo.Session, message *discordgo.MessageCreate) {
	height, _, err := util.GetNetworkStats()
	if err != nil {
		logging.Errorf("Unable to respond to blockheight command, unable to fetch current block height.\n")
	} else {
		response := fmt.Sprintf("Current blockheight: %d", height)
		_, err := dg.ChannelMessageSend(message.ChannelID, response)
		if err != nil {
			logging.Errorf("error responding to blockheight command, discord  \n", err)
		}
	}
}

func respondEstimatedEarnings(dg *discordgo.Session, message *discordgo.MessageCreate, arg string) {
	//First of all we need to check if the hashrate argument given in the command is valid
	//If the argument can't be converted to a float value, we need to return an error message
	hashrate, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		logging.Errorf("Error parsing argument in estimated earningscommand. argument (%s) may not be valid %s\n", arg, err)
		response := "Invalid hashrate argument \n?estearnings [hashrate in kH/s] is the correct syntax.\nExample: ``?estearnings 600``"
		_, err = dg.ChannelMessageSendReply(message.ChannelID, response, message.Reference())
		if err != nil {
			logging.Errorf("error responding to estimated earnings error \n", err)
			return
		}
		return
	}

	//Check current revenue
	blockHeight, diff, err := util.GetNetworkStats()
	if err != nil {
		logging.Errorf("Unable to fetch network stats, returning error message to ?estearnings command")
		response := "Something went wrong, unable to complete command request"
		_, err = dg.ChannelMessageSendReply(message.ChannelID, response, message.Reference())
		if err != nil {
			logging.Errorf("error responding to estimated earnings error \n", err)
			return
		}
		return
	}

	netHash := util.GetNethash(diff)
	blockReward, _ := util.GetBlockReward(int(blockHeight))

	// blockreward * blocks per day * minerhashrate / networkhashrate
	revenue := blockReward * (86400 / 150) * ((hashrate * 1000) / float64(netHash))

	response := fmt.Sprintf("A hashrate of %0.2f kH/s will give an estimated of %0.3f VTC per day", hashrate, revenue)
	_, err = dg.ChannelMessageSendReply(message.ChannelID, response, message.Reference())
	if err != nil {
		logging.Errorf("error responding to estimated earnings command, discord error \n", err)
	}
}
