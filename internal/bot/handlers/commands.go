package handlers

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = map[string]func(*discordgo.Session, *discordgo.MessageCreate){
	"ping":  Ping,
	"oi":    Greetings,
	"ola":   Greetings,
	"eae":   Greetings,
	"hey":   Greetings,
	"help":  Help,
	"ajuda": Help,
	"addwallet":  AddWalletCommand,
	"carteiras": ListWalletsCommand,
	"excluir": DeleteWalletCommand,
	"resumo": ListPositionsCommand,
	"pools": ListDetailedPositionsCommand,
}