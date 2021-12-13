package server

import "gopkg.in/guregu/null.v4"

type handleRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type handleGetCollectionByIdResponse struct {
	Id    int64          `json:"id"`
	Name  string         `json:"name"`
	Words []responseWord `json:"words"`
}

type responseWord struct {
	English string `json:"english"`
	French  string `json:"french"`
}

type collectionShortInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type handleChangeUserInfoRequest struct {
	FirstName   null.String `json:"first_name"`
	LastName    null.String `json:"last_name"`
	DateOfBirth null.Time   `json:"date_of_birth"`
}
