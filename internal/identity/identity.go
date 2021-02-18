package identity

import "strconv"

func IdentifierMap() map[string]int {
	IdentifierMapper := map[string]int{
		"IRON":  1,
		"ALLOY": 2,
		"OMS":   3,
	}

	return IdentifierMapper
}

// AuthorizeTokenAccessForEncryption checks the level of the identifier
func AuthorizeTokenAccessForEncryption(identifier string, level string) bool {
	IdentifierMap := IdentifierMap()
	src := IdentifierMap[identifier]
	i, _ := strconv.Atoi(level)
	if src < i {
		return false
	}

	return true
}

// AuthorizeRequest checks for a valid identifier
func AuthorizeRequest(accessToken string) bool {

	IdentifierMap := IdentifierMap()
	for key, _ := range IdentifierMap {
		if key == accessToken {
			return true
		}
	}

	return false
}
