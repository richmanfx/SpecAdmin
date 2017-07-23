package helpers

// Получить sessid из БД
func SessionIdExistInBD(sessidFromBrowser string) bool {

	var err error
	var sessidExist bool = false
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }


	db.Close()

	return sessidExist
}
