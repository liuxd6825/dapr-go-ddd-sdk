package ddd_errors

func IsMongoNoDocumentsInResult(err error) bool {
	if err.Error() == "no documents in result" {
		return true
	}
	return false
}
