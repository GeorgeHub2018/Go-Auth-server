package main

type (
	//UserData user data
	UserData struct {
		Name   string
		Age    string
		AuthID string
	}

	//ListData list data
	ListData struct {
		Users map[string]UserData
	}

	//FullUserData full user data
	FullUserData struct {
		User UserData
		List ListData
	}

	//ErrorData error data
	ErrorData struct {
		ErrorText string
	}
)
