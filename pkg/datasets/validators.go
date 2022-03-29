package datasets

func IsProviderValid(provider string) bool {
	if provider == "ssm" || provider == "none" {
		return true
	}

	return false
}

func IsTypeValid(dataSetType string) bool {
	return dataSetType == "sql"
}
