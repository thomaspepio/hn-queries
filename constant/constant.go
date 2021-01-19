package constant

const (
	// DateFormat is the date format in which we expect HN dates
	DateFormat = "2006-01-02 15:04:05"

	// Tab : tabulation constant
	Tab = "	"

	// Dash : dash constant
	Dash = "-"

	// DashRune : dash constant as a rune
	DashRune = '-'

	// SpaceRune: space constant as a rune
	SpaceRune = ' '

	// ColonRune : colon constant as a rune
	ColonRune = ':'

	// DateAsString : example of a HN date as string
	DateAsString = "2015-08-01 00:03:43"

	// URLAsString : example of a HN url as string
	URLAsString = "http%3A%2F%2Ftechacute.com%2F10-essentials-every-desk-needs%2F"

	// CorrectLine : example of a HN line to parse
	CorrectLine = DateAsString + Tab + URLAsString
)
