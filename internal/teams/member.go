package teams

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

var dbName = "./team.json"

type Team struct {
	Id      int      `json:"id"`
	Members []string `json:"members"`
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

	newFile, err := os.Create(dbName)
	if err != nil {
		return nil, err
	}
	defer newFile.Close()

	encoder := json.NewEncoder(newFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(teams); err != nil {
		fmt.Println("JSON 인코딩 실패:", err)
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
