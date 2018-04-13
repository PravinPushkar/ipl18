package handler

import "github.wdf.sap.corp/I334816/ipl18/backend/auth"

var tokenManager = auth.NewTokenManager(auth.SignMethodSHA512)
