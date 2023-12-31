package identity

import (
	"strconv"

	"github.com/pharmaeasy/tokenizer/internal/models/db"

	"github.com/pharmaeasy/tokenizer/internal/models/encryption"
)

// TODO:: Add Gateway, Core Oms Service with Level
func IdentifierMap() map[string]int {
	IdentifierMapper := map[string]int{
		"IRON": 2,
		//"PRODUCT_OMS":                     2,
	}

	return IdentifierMapper
}

// AuthorizeLevelForEncryption checks the level of the identifier for an encryption request
func AuthorizeLevelForEncryption(requestData *encryption.EncryptRequest) bool {
	for i := 0; i < len(requestData.EncryptRequestData); i++ {
		if !authorizeIdentifierByLevel(requestData.Identifier, requestData.EncryptRequestData[i].Level) {
			return false
		}
	}

	return true
}

func authorizeIdentifierByLevel(identifier string, level string) bool {
	IdentifierMap := IdentifierMap()
	src := IdentifierMap[identifier]
	i, err := strconv.Atoi(level)
	if err != nil {
		return false
	}
	if src > i {
		return false
	}

	return true
}

// AuthenticateRequest checks for a valid identifier
func AuthenticateRequest(accessToken string) bool {

	IdentifierMap := IdentifierMap()
	for key, _ := range IdentifierMap {
		if key == accessToken {
			return true
		}
	}

	return false
}

func getAccessLevelByIdentifier(identifier string) *int {
	if level, ok := IdentifierMap()[identifier]; ok {
		return &level
	}

	return nil
}

// AuthorizeTokenAccess authorizes token access using identifer and corresponding level
func AuthorizeTokenAccess(tokenData *map[string]db.TokenData, identifier string) bool {
	levelOfIdentifier := getAccessLevelByIdentifier(identifier)

	for _, token := range *tokenData {
		tokenLevel, _ := strconv.Atoi(token.Level)
		if tokenLevel < *levelOfIdentifier {
			return false
		}
	}

	return true
}
