package translate

var ValidMap = map[string]map[string]string{
	"valid_password": {"en": "Must be at least 8 in length and contain at least one letter and one number",
		"zh": "长度至少为8，至少含有一个字母和一个数字"},
	"valid_iplist": {"en": "invalid ip format", "zh": "不符合ip格式"},
}

func GetValidMsg(valildKey, language string) string {
	val, ok := ValidMap[valildKey]
	if !ok {
		return ""
	}
	msg, ok := val[language]
	if !ok {
		return val["en"]
	}
	return msg
}
