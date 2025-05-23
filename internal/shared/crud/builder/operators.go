package builder

type Operator string

const (
	EQ      Operator = "="
	NEQ     Operator = "!="
	LT      Operator = "<"
	GT      Operator = ">"
	GTE     Operator = ">="
	LTE     Operator = "<="
	IN      Operator = "IN"
	NIN     Operator = "NOT IN"
	BETWEEN Operator = "BETWEEN"
	NULL    Operator = "NULL"
	NOTNULL Operator = "NOT NULL"
)

type JoinType string

const (
	JoinInner JoinType = "JOIN"
	JoinLeft  JoinType = "LEFT JOIN"
	JoinRight JoinType = "RIGHT JOIN"
	JoinCross JoinType = "CROSS JOIN"
)

func (jt JoinType) String() string {
	return string(jt)
}
