package kjson

import "strconv"

func JsonFlatten(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	decodeStruct("", m, result)
	return result

}
func decodeStruct(prefix string, m map[string]interface{}, r map[string]interface{}) {
	for k, v := range m {
		var key string
		if prefix == "" {
			key = k
		} else {
			key = prefix + "." + k
		}
		switch vv := v.(type) {
		case bool:
			r[key] = vv
		case string:
			r[key] = vv
		case int:
			r[key] = vv
		case float64:
			r[key] = vv
		case []interface{}:
			for i, u := range vv {
				uu, valid := u.(map[string]interface{})

				if valid {
					decodeStruct(key+"."+strconv.Itoa(i), uu, r)
				} else {
					key := prefix + "." + k + "." + strconv.Itoa(i)
					r[key] = u
				}
			}
		case map[string]interface{}:
			decodeStruct(key, vv, r)
		default:
			r[key] = vv
		}
	}
}
