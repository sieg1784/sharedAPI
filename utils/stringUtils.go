package utils

// [+886931068729]
func TrimBrackets(str string) string {
	if '[' == str[0] && ']' == str[len(str)-1] {
		return str[1 : len(str)-1]
	} else {
		return str
	}
}

func ExtractCountryCode(str string) string {
	return str[0:4]
}

func ExtractPhoneNum(str string) string {
	return str[4:]
}

func NormalizePhoneNum(phoneNum string) string {
	if len(phoneNum) == 10 && phoneNum[0] == '0' {
		return phoneNum[1:]
	} else {
		return phoneNum
	}
}

func GenderInt2String(input int64) string {
	genderMap := make(map[int64]string)
	genderMap[1] = "male"
	genderMap[2] = "female"
	genderMap[3] = "noneOfAbove"
	genderMap[9] = "notAvailable"
	return genderMap[input]
}

func GenderString2Int(input string) int64 {
	genderMap := make(map[string]int64)
	genderMap["male"] = 1
	genderMap["female"] = 2
	genderMap["noneOfAbove"] = 3
	genderMap["notAvailable"] = 9
	return genderMap[input]
}
