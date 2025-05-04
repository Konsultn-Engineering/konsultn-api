package crud

type QueryParams struct {
	Page   int               `form:"page"`
	Limit  int               `form:"limit"`
	Sort   string            `form:"sort"`
	Order  string            `form:"order"`
	Filter map[string]string `form:"filter"`
	Search string            `form:"q"`
}

type JoinClause struct {
	Table    string // e.g., "team_members"
	On       string // e.g., "team_members.team_id = teams.id"
	JoinType string // "JOIN", "LEFT JOIN", etc.
}

type WhereClause struct {
	Query string // e.g., "team_members.user_id = ?"
	Args  []any  // e.g., []any{userID}
}

type AdvancedQuery struct {
	QueryParams
	Joins   []JoinClause
	Wheres  []WhereClause
	Preload []string
}
