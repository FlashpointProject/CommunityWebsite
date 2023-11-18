package service

import "github.com/FlashpointProject/CommunityWebsite/constants"

func dberr(err error) error {
	return constants.DatabaseError{Err: err}
}
