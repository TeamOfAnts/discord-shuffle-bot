package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/TeamOfAnts/discord-shuffle-bot/internal/shuffle"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var members = []string{}
var teamSize = 0

func Run() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("세션 생성 오류:", err)
	}

	discord.AddHandler(messageCreate)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
func messageCreate(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.HasPrefix(message.Content, "!ping"):
		discord.ChannelMessageSend(message.ChannelID, "pong🏓")
	case strings.HasPrefix(message.Content, "!members"), strings.HasPrefix(message.Content, "!멤버"):
		if len(members) == 0 {
			discord.ChannelMessageSend(message.ChannelID, "멤버가 없습니다.")
			return
		}
		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("현재 멤버: %s", strings.Join(members, ", ")))
	case strings.HasPrefix(message.Content, "!add"), strings.HasPrefix(message.Content, "!추가"):
		replacer := strings.NewReplacer(
			"!add", "",
			"!추가", "",
		)
		name := replacer.Replace(message.Content)
		members = append(members, name)
		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s 멤버 추가 완료", name))
	case strings.HasPrefix(message.Content, "!팀크기"), strings.HasPrefix(message.Content, "!teamSize"):
		replacer := strings.NewReplacer(
			"!teamSize", "",
			"!팀크기", "",
		)
		size := replacer.Replace(message.Content)
		if len(size) == 0 {
			discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("현재 팀 사이즈는 %d명 입니다.", teamSize))
			return
		}
		s, err := strconv.Atoi(size)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "숫자를 입력해주세요.")
			return
		}
		teamSize = s
		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("팀 사이즈가 %d명으로 변경되었습니다.", teamSize))
	case strings.HasPrefix(message.Content, "!shuffle"):
		teams := shuffle.Shuffle(members, teamSize)
		discord.ChannelMessageSend(message.ChannelID, teams)
	case strings.HasPrefix(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "명령어 목록\n- 멤버 추가\n  - !add [이름]\n  - !추가 [이름]\n- 멤버 목록\n  - !members\n  - !멤버\n- 팀당 인원수 확인\n  - !teamSize\n  - !팀크기\n- 팀당 인원수 변경\n  - !teamSize [숫자]\n  - !팀크기 [숫자]\n- 팀 나누기\n  - !shuffle\n- 도움말\n  - !help")
	}
}
