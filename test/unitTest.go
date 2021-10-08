package test

import (
	"bookAPI/apiservice/utils"
	"fmt"
)

func main() {
	phoneFromIdToken := "+886931068729"
	fmt.Println(phoneFromIdToken)

	countryCodeFromFirebase := utils.ExtractCountryCode(phoneFromIdToken)
	fmt.Println(countryCodeFromFirebase)

	phoneNumFromIdToken := utils.ExtractPhoneNum(phoneFromIdToken)
	fmt.Println(phoneNumFromIdToken)
}
