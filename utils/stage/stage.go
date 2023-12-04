package stage

import "strings"

type StageType string

const (
	Local StageType = "LOCAL"
	Int   StageType = "INT"
	Dev   StageType = "DEV"
	Stag  StageType = "STAG"
	PROD  StageType = "PROD"
)

func (s *StageType) UpCase() {
	if s != nil {
		*s = StageType(strings.ToUpper(string(*s)))
	}
}

func (s StageType) ToString() string {
	return string(s)
}

func (s StageType) IsValid() bool {
	switch s {
	case Dev, Stag, PROD, Int, Local:
		return true
	default:
		return false
	}
}

func (s StageType) IsProd() bool {
	return s == PROD
}

func (s StageType) IsStaging() bool {
	return s == Stag
}

func (s StageType) IsDev() bool {
	return s == Dev
}

func (s StageType) IsLocal() bool {
	return s == Local
}

func (s StageType) IsInt() bool {
	return s == Int
}
