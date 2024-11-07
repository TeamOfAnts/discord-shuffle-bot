package teams

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"
)

var dbName = "./team.json"

type Team struct {
	Id      int      `json:"id"`
	Members []string `json:"members"`
	Size    int      `json:"size"`
}

func AddMember(names []string) ([]string, error) {
	teams, err := GetTeams()
	if err != nil {
		return nil, err
	}

	Id := teams[0].Id + 1

	members := teams[0].Members
	for _, name := range names {
		if !slices.Contains(members, name) {
			members = append(members, name)
		}
	}

	newTeam := Team{Id: Id, Members: members}
	teams = append([]Team{newTeam}, teams...)

	if err := save(teams); err != nil {
		return nil, err
	}

	return newTeam.Members, nil
}

func GetTeams() ([]Team, error) {
	file, err := os.Open(dbName)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []Team
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("JSON 디코딩 실패:", err)
	}

	return data, nil
}

func UpdateSize(size int) (int, error) {
	teams, err := GetTeams()
	if err != nil {
		return 0, err
	}

	newTeam := Team{Id: teams[0].Id + 1, Members: teams[0].Members, Size: size}
	teams = append([]Team{newTeam}, teams...)

	if err := save(teams); err != nil {
		return 0, err
	}

	return size, nil
}

func save(teams []Team) error {
	newFile, err := os.Create(dbName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	encoder := json.NewEncoder(newFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(teams); err != nil {
		fmt.Println("JSON 인코딩 실패:", err)
	}
	return nil
}

func Shuffle() string {
	teams, err := GetTeams()
	if err != nil {
		return "멤버가 없습니다."
	}

	members := teams[0].Members

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	shuffledTeams := make([][]string, 0)

	for _, member := range members {
		if len(shuffledTeams) == 0 {
			shuffledTeams = append(shuffledTeams, []string{member})
			continue
		}

		if len(shuffledTeams[len(shuffledTeams)-1]) < teams[0].Size {
			shuffledTeams[len(shuffledTeams)-1] = append(shuffledTeams[len(shuffledTeams)-1], member)
		} else {
			shuffledTeams = append(shuffledTeams, []string{member})
		}
	}

	result := ""
	for i, team := range shuffledTeams {
		result += fmt.Sprintf("%d조 : %s \n", i+1, strings.Join(team, ", "))
	}

	return result
}
