package misc

// MessageSplit splits message to chunks each limited by chunkLimit
func MessageSplit(message string, chunkLimit int) []string {
	return messageMakeChunks(messageSplitSeparator(message), chunkLimit)
}

// messageMakeChunks makes messages chunks limited chunkLimit from strings
func messageMakeChunks(strs []string, chunkLimit int) []string {

	var (
		r   string
		res []string
	)

	for i, s := range strs {

		// If current string has len more then available
		if len(s) > chunkLimit {

			// Calc len of rest space in `r`
			l := chunkLimit - len(r)

			r += s[0:l]

			// Add to result first part of current string
			res = messageAppend(res, r)

			// Prepare rest of strings
			rs := []string{}
			if i+1 < len(strs) {
				// If not last string
				rs = strs[i+1:]
			}

			// Join current result with the result of truncate with the rest of strings
			return append(res, messageMakeChunks(append(
				[]string{s[l:]},
				rs...), chunkLimit)...)
		}

		if len(r)+len(s) > chunkLimit {
			res = messageAppend(res, r)
			r = s
			continue
		}

		r += s
	}

	return messageAppend(res, r)
}

// messageSplitSeparator splits long message to substrings using spaces, tabs and EOLs as separator
func messageSplitSeparator(str string) []string {

	strs := []string{}
	r := ""

	for _, s := range str {
		if s == ' ' || s == '\t' || s == '\n' {
			strs = append(strs, r+string(s))
			r = ""
		} else {
			r += string(s)
		}
	}

	return append(strs, r)
}

// messageTrim trims all trailing separators
func messageTrim(str string) string {

	if len(str) == 0 {
		return ""
	}

	s := str[len(str)-1]

	if s == ' ' || s == '\t' || s == '\n' {
		return messageTrim(str[:len(str)-1])
	}

	return str
}

// messageAppend appends trimmed `s` to `strs` if `s` len is positive
func messageAppend(strs []string, s string) []string {

	r := messageTrim(s)

	if len(r) > 0 {
		return append(strs, r)
	}

	return strs
}
