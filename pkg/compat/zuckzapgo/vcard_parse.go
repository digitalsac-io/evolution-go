package zuckzapgo

import (
	"strings"
)

// ParseSimpleVCard extrai FN, primeira TEL e ORG de um VCARD em texto (compatível com MessageContact.Vcard do OpenAPI).
func ParseSimpleVCard(raw string) (fullName, phone, org string) {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		upper := strings.ToUpper(line)
		if strings.HasPrefix(upper, "FN:") {
			fullName = strings.TrimSpace(line[3:])
			continue
		}
		if strings.HasPrefix(upper, "ORG:") {
			org = strings.TrimSpace(line[4:])
			continue
		}
		if strings.HasPrefix(upper, "TEL") {
			if idx := strings.LastIndex(line, ":"); idx >= 0 && idx < len(line)-1 {
				candidate := strings.TrimSpace(line[idx+1:])
				if candidate != "" && phone == "" {
					phone = candidate
				}
			}
		}
	}
	return fullName, phone, org
}
