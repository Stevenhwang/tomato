package utils

type MD struct {
	Total       string `json:"total" validate:"required"`
	Free        string `json:"free" validate:"required"`
	UsedPercent string `json:"usedpercent" validate:"required"`
}

type LD struct {
	Load1  string `json:"load1" validate:"required"`
	Load5  string `json:"load5" validate:"required"`
	Load15 string `json:"load15" validate:"required"`
}

type Info struct {
	Mem  MD `json:"mem" validate:"required,dive,required"`
	Disk MD `json:"disk" validate:"required,dive,required"`
	Load LD `json:"load" validate:"required,dive,required"`
}

type Register struct {
	Name string `json:"name" validate:"required"`
	Info Info   `json:"info" validate:"required,dive,required"`
}

type CMD struct {
	Agents []string `json:"agents" validate:"required,gt=0,dive,required"`
	CMD    string   `json:"cmd" validate:"required"`
}

type Response struct {
	Agent string `json:"agent"`
	Msg   string `json:"msg"`
}
