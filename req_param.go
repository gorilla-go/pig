package pig

type ReqParams map[string]*ReqParamV

type ReqParamV struct {
	value []string
}

func NewReqParamV(value []string) *ReqParamV {
	return &ReqParamV{
		value: value,
	}
}
