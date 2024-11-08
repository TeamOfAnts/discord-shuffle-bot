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

var teamDbName = "./team.json"
var shuffleDbName = "./shuffle-version.json"

type Team struct {
	Id      int      `json:"id"`
	Members []string `json:"members"`
	Size    int      `json:"size"`
}

type ShuffleVersion struct {
	Teams                []ShuffledTeam `json:"teams"`
	ShuffleAvailableDate string         `json:"shuffleAvailableDate"`
	ShuffledDate         string         `json:"shuffledDate"`
	ValidDate            string         `json:"validDate"`
}

type ShuffledTeam struct {
	Id      int      `json:"id"`
	Members []string `json:"members"`
}

func (s *ShuffleVersion) GetStringTeams() string {
	str := ""

	for _, team := range s.Teams {
		str += fmt.Sprintf("%d조 : %s \n", team.Id, strings.Join(team.Members, ", "))
	}

	return str
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

	newTeam := Team{Id: Id, Members: members, Size: teams[0].Size}
	teams = append([]Team{newTeam}, teams...)

	if err := saveTeams(teams); err != nil {
		return nil, err
	}

	return newTeam.Members, nil
}

func GetTeams() ([]Team, error) {
	file, err := os.Open(teamDbName)

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

	if err := saveTeams(teams); err != nil {
		return 0, err
	}

	return size, nil
}

func saveTeams(teams []Team) error {
	newFile, err := os.Create(teamDbName)
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

func Shuffle() (string, error) {
	teams, err := GetTeams()
	if err != nil {
		return "", err
	}

	shuffleVersions, err := GetShuffleVersions()
	if err != nil {
		return "", err
	}

	if len(shuffleVersions) != 0 && shuffleVersions[0].ShuffleAvailableDate > time.Now().Format("2006-01-02") {
		return "", fmt.Errorf("팀 셔플은 **`%s`** 이후 가능합니다.\n현재 조원\n %s", shuffleVersions[0].ShuffleAvailableDate, shuffleVersions[0].GetStringTeams())
	}

	members := teams[0].Members

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	shuffledTeams := make([][]string, 0)

	newShuffleVersion := ShuffleVersion{
		Teams:                []ShuffledTeam{},
		ShuffleAvailableDate: time.Now().Add(time.Hour * 24 * 10).Format("2006-01-02"),
		ShuffledDate:         time.Now().Format("2006-01-02"),
		ValidDate:            time.Now().Add(time.Hour * 24 * 14).Format("2006-01-02"),
	}

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
		newShuffleVersion.Teams = append(newShuffleVersion.Teams, ShuffledTeam{Id: i + 1, Members: team})
		result += fmt.Sprintf("%d조 : %s \n", i+1, strings.Join(team, ", "))
	}

	result += fmt.Sprintf("\n팀 셔플 완료\n\n**`%s`** 까지 스터디를 진행해주세요!", newShuffleVersion.ValidDate)

	saveShuffleVersion([]ShuffleVersion{newShuffleVersion})

	return result, nil
}

func GetShuffleVersions() ([]ShuffleVersion, error) {

	file, err := os.Open(shuffleDbName)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []ShuffleVersion
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("JSON 디코딩 실패:", err)
	}

	return data, nil
}
func saveShuffleVersion(shuffleVersions []ShuffleVersion) error {
	newFile, err := os.Create(shuffleDbName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	encoder := json.NewEncoder(newFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(shuffleVersions); err != nil {
		fmt.Println("JSON 인코딩 실패:", err)
	}
	return nil
}

func GetShuffledTeams() (string, error) {
	shuffleVersions, err := GetShuffleVersions()
	if err != nil {
		return "", err
	}

	if len(shuffleVersions) == 0 {
		return "팀 셔플 기록이 없습니다.", nil
	}

	return fmt.Sprintf("현재 조원\n %s \n **`%s`** 까지 스터디 진행 부탁드립니다.", shuffleVersions[0].GetStringTeams(), shuffleVersions[0].ValidDate), nil
}
