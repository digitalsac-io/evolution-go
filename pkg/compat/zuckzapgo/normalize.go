package zuckzapgo

// ResolveRecipient returns number if set, otherwise phone.
func ResolveRecipient(phone, number string) string {
	if number != "" {
		return number
	}
	return phone
}

// FirstNonEmpty returns the first non-empty string.
func FirstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
