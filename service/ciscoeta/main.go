package ciscoeta

import "github.com/zartbot/goflow/service/ciscoeta/idp"

func ExtendField(d map[string]interface{}) {
	idp.DecodeIDPField(d)
}
