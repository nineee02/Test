package constant

const (
	RequestIdKey   = "request_id"
	AcceptLanguage = "accept_language"
)

const (
	Phase1 = iota + 1 // 120
	Phase2            // 60
	Phase3            // 180
)

const (
	StatusCompleted   string = "COMPLETED"
	StatusInCompleted string = "INCOMPLETED"
	StatusInProgress  string = "INPROGRESS"
	StatusTerminated  string = "TERMINATED"
	StatusDraft       string = "DRAFT"

	LanguageTypeEN string = "EN"
	LanguageTypeTH string = "TH"
)
