package lottery

import (
	"math/rand"
	"pick-a-bro/internal/commons"
	"pick-a-bro/internal/data"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func InitMembersList() (*data.MembersList, error) {
	membersList, err := prepareLottery()
	if err != nil {
		return nil, err
	}

	return membersList, nil
}

// prepareLottery prepares the lottery by retrieving preferences, getting the members list,
// applying the chances rule, excluding winners if necessary, shuffling the members list,
// and setting the enhanced members list as the new members list.
// It returns the original members list and an error, if any.
func prepareLottery() (*data.MembersList, error) {
	preferences := commons.GetPreferences()
	membersList := data.GetMembersAndTiers()

	chancesRule := preferences.StringWithFallback(commons.ChancesRule, commons.GetTranslation(commons.ChancesRules[0]))

	enhancedMembersList := prepareMembersList(chancesRule, membersList.PatreonMembers)

	if preferences.BoolWithFallback(commons.ExcludeWinners, false) {
		enhancedMembersList = excludeWinners(enhancedMembersList)
	}
	shuffleMembers(enhancedMembersList)
	data.SetMembersList(enhancedMembersList)
	return membersList, nil
}

// prepareMembersList prepares the members list based on the chances rule and returns the updated list.
// It takes a chancesRule string and a membersList []data.PatreonMember as input parameters.
// The chancesRule determines how the members list will be prepared.
// The membersList is the list of Patreon members to be prepared.
// The function iterates over the membersList and duplicates each member based on the chances rule.
// If the chancesPerUser is 1, the function returns the original membersList.
// If the chancesPerUser is greater than 1, the function duplicates each member in the list by the chancesPerUser value.
// If the chances rule is based on the tier of each member, the function duplicates each member in the list based on the chances value associated with their tier.
// The function returns the updated membersList.
func prepareMembersList(chancesRule string, membersList []data.PatreonMember) []data.PatreonMember {
	switch chancesRule {
	case commons.GetTranslation(commons.ChancesRules[0]):
		chancesPerUser := commons.GetPreferences().IntWithFallback(commons.ChancesPerUser, 1)
		if chancesPerUser == 1 {
			return membersList
		}
		for _, d := range membersList {
			for i := 1; i < chancesPerUser; i++ {
				membersList = append(membersList, d)
			}
		}
	case commons.GetTranslation(commons.ChancesRules[1]):
		for _, d := range membersList {
			for i := 1; i < commons.GetPreferences().IntWithFallback("chances"+d.Tier, 1); i++ {
				membersList = append(membersList, d)
			}
		}
	}

	return membersList
}

// excludeWinners removes the winners from the enhancedMembersList.
// It takes a slice of PatreonMember structs as input and returns a new slice
// with the winners excluded.
func excludeWinners(enhancedMembersList []data.PatreonMember) []data.PatreonMember {
	winners := GetWinnersList()
	for _, winner := range winners {
		for i, member := range enhancedMembersList {
			if member.FullName == winner.FullName {
				enhancedMembersList = append(enhancedMembersList[:i], enhancedMembersList[i+1:]...)
			}
		}
	}
	return enhancedMembersList
}

func shuffleMembers(membersList []data.PatreonMember) {
	r.Shuffle(len(membersList), func(i, j int) {
		membersList[i], membersList[j] = membersList[j], membersList[i]
	})
}
