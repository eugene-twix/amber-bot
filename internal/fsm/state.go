// internal/fsm/state.go
package fsm

type State string

const (
	StateNone State = ""

	// NewTeam flow
	StateNewTeamName       State = "new_team:name"
	StateNewTeamAddMembers State = "new_team:add_members" // после создания - добавить участников?
	StateNewTeamMoreMember State = "new_team:more_member" // ещё участника?
	StateNewTeamAddResult  State = "new_team:add_result"  // записать результат?

	// AddMember flow
	StateAddMemberTeam State = "add_member:team"
	StateAddMemberName State = "add_member:name"

	// NewTournament flow
	StateNewTournamentName     State = "new_tournament:name"
	StateNewTournamentDate     State = "new_tournament:date"
	StateNewTournamentLocation State = "new_tournament:location"

	// Result flow
	StateResultTournament State = "result:tournament"
	StateResultTeam       State = "result:team"
	StateResultPlace      State = "result:place"

	// Grant flow
	StateGrantUser State = "grant:user"
	StateGrantRole State = "grant:role"
)

type Data map[string]any

func (d Data) GetInt64(key string) int64 {
	if v, ok := d[key].(float64); ok {
		return int64(v)
	}
	if v, ok := d[key].(int64); ok {
		return v
	}
	return 0
}

func (d Data) GetString(key string) string {
	if v, ok := d[key].(string); ok {
		return v
	}
	return ""
}
