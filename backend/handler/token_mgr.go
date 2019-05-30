package handler

import "github.com/PravinPushkar/ipl18/backend/auth"

var tokenManager = auth.NewTokenManager(auth.SignMethodSHA512)
