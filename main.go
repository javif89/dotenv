package dotenv

import (
	"bufio"
	"os"
	"strings"
)

type EnvValue struct {
	Key        string
	Value      string
	IsComment  bool
	HasComment bool
	Comment    string
}

type EnvFile struct {
	Values []EnvValue
	Path   string
}

func (e *EnvFile) Keys() []string {
	keys := []string{}
	for _, v := range e.Values {
		if v.IsComment {
			continue
		}
		keys = append(keys, v.Key)
	}
	return keys
}

func (e *EnvFile) Add(key string, value string) {
	key = standardizeKey(key)

	envVal := newValue("", "")

	// Check if the value has a comment attached
	if strings.Contains(value, "#") {
		split := strings.Split(value, "#")
		value = split[0]
		envVal.Comment = cleanString(split[1])
		envVal.HasComment = true
	}

	value = cleanString(value)

	if e.Has(key) {
		for i, v := range e.Values {
			if v.Key == key {
				e.Values[i].Value = value
			}
		}
		return
	}

	envVal.Key = key
	envVal.Value = value

	e.Values = append(e.Values, envVal)
}

func (e *EnvFile) Remove(key string) {
	key = standardizeKey(key)
	for i, v := range e.Values {
		if v.Key == key {
			e.Values = append(e.Values[:i], e.Values[i+1:]...)
		}
	}
}

func (e *EnvFile) Has(key string) bool {
	for _, v := range e.Values {
		if v.Key == key {
			return true
		}
	}
	return false
}

func (e *EnvFile) Get(key string) string {
	key = standardizeKey(key)
	for _, v := range e.Values {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

func (e *EnvFile) Set(key string, value string) {
	e.Add(key, value)
}

func (e *EnvFile) AddComment(comment string) {
	if !strings.HasPrefix(comment, "#") {
		comment = "# " + comment
	}
	val := newValue("", comment)
	val.IsComment = true
	e.Values = append(e.Values, val)
}

func (e *EnvFile) DiffKeys(comparable *EnvFile) []string {
	diff := []string{}
	// TODO: Inneficient but we'll do this for now
	for _, k := range comparable.Keys() {
		if !e.Has(k) {
			diff = append(diff, k)
		}
	}

	return diff
}

func (e *EnvFile) Merge(comparable *EnvFile, overwrite bool) {
	var keys []string

	switch {
	case overwrite:
		keys = comparable.Keys()
	default:
		keys = e.DiffKeys(comparable)
	}

	for _, k := range keys {
		e.Set(k, comparable.Get(k))
	}
}

func New(path string) *EnvFile {
	return &EnvFile{Path: path}
}

func LoadOrCreate(path string) *EnvFile {
	if !fileExists(path) {
		return New(path)
	}

	return Load(path)
}

func Load(path string) *EnvFile {
	file, e := os.Open(path)
	if e != nil {
		panic(e)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	env := EnvFile{Path: path}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		kv := strings.Split(line, "=")

		// Check if it's a comment
		if len(kv) == 1 && strings.HasPrefix(kv[0], "#") {
			env.AddComment(kv[0])
			continue
		}

		// Not sure how this would happen but let's do it just in case
		if len(kv) < 2 {
			continue
		}

		// We should have a key and a value
		key := kv[0]
		value := kv[1]

		env.Add(key, value)
	}

	return &env
}

func (e *EnvFile) Save() {
	// Check if file exists if not create it
	if _, err := os.Stat(e.Path); os.IsNotExist(err) {
		file, err := os.Create(e.Path)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	file, err := os.OpenFile(e.Path, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, v := range e.Values {
		if v.IsComment {
			file.WriteString(v.Value + "\n")
		} else {
			file.WriteString(v.Key + "=" + formatValueForPrint(v.Value, v.Comment) + "\n")
		}
	}
}

func newValue(key string, value string) EnvValue {
	return EnvValue{
		Key:        key,
		Value:      value,
		IsComment:  false,
		HasComment: false,
		Comment:    "",
	}
}
