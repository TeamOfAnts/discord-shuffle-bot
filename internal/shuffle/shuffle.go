package shuffle

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Shuffle(members []string, numberOfTeam int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	random.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})

	teams := make([][]string, 0)

	for _, member := range members {
		if len(teams) == 0 {
			teams = append(teams, []string{member})
			continue
		}

		if len(teams[len(teams)-1]) < numberOfTeam {
			teams[len(teams)-1] = append(teams[len(teams)-1], member)
		} else {
			teams = append(teams, []string{member})
		}
	}

	result := ""
	for i, team := range teams {
		result += fmt.Sprintf("%d조 : %s \n", i+1, strings.Join(team, ", "))
	}

	return result
}
