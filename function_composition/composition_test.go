package function_composition

import (
	"github.com/mikejlong60/golangz/arrays"
	"strings"
	"testing"
)

func TestCompose(t *testing.T) {
	normalizeUserDn := func(userDN string) string {

		lowercase := func(s string) string {
			return strings.ToLower(s)
		}

		trimSpaces := func(s string) string {
			return strings.TrimSpace(s)
		}

		g := Compose(lowercase, trimSpaces)

		removeEmptyStrings := func(s string) bool {
			if s == "" {
				return false
			} else {
				return true
			}
		}

		normalized := arrays.Filter(arrays.Map(strings.Split(strings.ReplaceAll(userDN, "/", ","), ","), g), removeEmptyStrings)

		var result string
		if len(normalized) == 0 {
			return ""
		}
		if strings.HasPrefix(normalized[0], "cn") {
			result = strings.Join(normalized, ",")
		} else {
			f := arrays.Appender[string]
			reversed := arrays.FoldRight(normalized, []string{}, f) //Reverse the list
			result = strings.Join(reversed, ",")
		}
		return result

	}

	expectedUsers := []string{
		"cn=test tester02,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester03,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"c=us,o=u.s. government,ou=chimera,ou=dae,ou=people,dn=test tester04",
		"cn=test tester05,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester06,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester07,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester08,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester09,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"",
		"c=us,o=u.s. government,ou=chimera,ou=dae,ou=people,dn=test tester10",
		"",
		"",
	}

	testUsers := []string{
		",cn=test tester02, ou=pEOple,,,, ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester03,ou=people,ou=dae,ou=chimera,o=u.s. governmEnt,c=us",
		"DN=test tester04/  OU=people/ou=dae/ou=chimera,o=u.s. government,c=US",
		"cn=test tester05,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester06,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester07,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester08,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"cn=test tester09,ou=people,ou=dae,ou=chimera,o=u.s. government,c=us",
		"",
		"DN=test tester10,  OU=people,ou=dae,ou=chimera,o=u.s. government,c=US",
		",,/,,/,,",
		",,,,,",
	}

	var actualUsers []string
	for _, rawUserDn := range testUsers {
		actualUsers = append(actualUsers, normalizeUserDn(rawUserDn))
	}
	equality := func(l, r string) bool {
		if l == r {
			return true
		} else {
			return false
		}
	}

	if !arrays.SetEquality(actualUsers, expectedUsers, equality) {
		t.Errorf("actual:%v, expected:%v", actualUsers, expectedUsers)
	}
}
