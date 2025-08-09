package disk

import (
	"fmt"
	"strings"
)

func encodeList(list []string) ([]byte, error) {
	if len(list) == 0 {
		return nil, nil
	}
	return []byte("[\"" + strings.Join(list, "\",\"") + "\"]"), nil
}

func encodeMap(kv map[string]string) ([]byte, error) {
	if len(kv) == 0 {
		return nil, nil
	}

	var sb strings.Builder
	sb.WriteString("{")
	count := len(kv)
	for k, v := range kv {
		sb.WriteString(fmt.Sprintf("\"%s\":\"%s\"", k, v))
		count--
		if count > 0 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("}")
	return []byte(sb.String()), nil
}

func decodeList(data []byte, list *[]string) error {
	if len(data) == 0 {
		return nil
	}

	str := strings.Trim(string(data), "[]")
	if str == "" {
		return nil
	}

	items := strings.Split(str, "\",\"")
	for _, item := range items {
		*list = append(*list, strings.Trim(item, "\""))
	}
	return nil
}

func decodeMap(data []byte, kv *map[string]string) error {
	if len(data) == 0 {
		return nil
	}

	str := strings.Trim(string(data), "{}")
	if str == "" {
		return nil
	}

	*kv = make(map[string]string)

	items := strings.Split(str, "\",\"")
	for _, item := range items {
		parts := strings.SplitN(item, "\":\"", 2)
		if len(parts) != 2 {
			continue // skip malformed entries
		}
		key := strings.Trim(parts[0], "\"")
		value := strings.Trim(parts[1], "\"")
		(*kv)[key] = value
	}
	return nil
}
