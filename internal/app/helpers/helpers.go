package helpers

import (
	"fmt"
	"zkKYC-backend/internal/app/storage"
)

func CreateDid(user storage.User) string {

	return fmt.Sprintf("did:ethr:%s", user.EthAddress)

}
