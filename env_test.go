package dotenv

import (
	"os"
	"reflect"
	"testing"
)

func TestDiffKeys(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")
	e.Add("test2", "test2")

	comparable := EnvFile{}
	comparable.Add("test", "test")
	comparable.Add("test2", "test2")
	comparable.Add("test3", "test3")

	keys := e.DiffKeys(&comparable)
	if len(keys) != 1 {
		t.Errorf("Expected length of 1, got %d", len(keys))
	}

	if keys[0] != "TEST3" {
		t.Errorf("Expected key to be TEST3, got %s", keys[0])
	}
}

func TestMergeNoOverwrite(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")
	e.Add("test2", "test2")

	comparable := EnvFile{}
	comparable.Add("test", "test")
	comparable.Add("test2", "test2")
	comparable.Add("test3", "test3")

	e.Merge(&comparable, false)
	if len(e.Values) != 3 {
		t.Errorf("Expected length of 3, got %d", len(e.Values))
	}

	if e.Values[2].Key != "TEST3" {
		t.Errorf("Expected key to be TEST3, got %s", e.Values[2].Key)
	}
}

func TestMergeOverwrite(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")
	e.Add("test2", "test2")

	comparable := EnvFile{}
	comparable.Add("test", "overwritten")
	comparable.Add("test2", "overwritten")
	comparable.Add("test3", "overwritten")

	e.Merge(&comparable, true)
	if len(e.Values) != 3 {
		t.Errorf("Expected length of 3, got %d", len(e.Values))
	}

	for _, v := range e.Values {
		if v.Value != "overwritten" {
			t.Errorf("Expected value to be overwritten, got %s", v.Value)
		}
	}
}

func TestAdd(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")
	if len(e.Values) != 1 {
		t.Errorf("Expected length of 1, got %d", len(e.Values))
	}

	if !e.Has("TEST") {
		t.Errorf("Expected to have key TEST")
	}
}

func TestLineWithComment(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test # comment")
	if len(e.Values) != 1 {
		t.Errorf("Expected length of 2, got %d", len(e.Values))
	}

	if !e.Values[0].HasComment {
		t.Errorf("Expected HasComment to be true, got %t", e.Values[1].IsComment)
	}

	if e.Values[0].Comment != "comment" {
		t.Errorf("Expected comment to be comment, got %s", e.Values[1].Comment)
	}
}

func TestHas(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")
	if !e.Has("TEST") {
		t.Errorf("Expected to have key TEST")
	}
}

func TestSet(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")
	e.Set("test", "test2")
	if e.Values[0].Value != "test2" {
		t.Errorf("Expected value to be test2, got %s", e.Values[0].Value)
	}
}

func TestRemove(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")
	e.Remove("test")
	if len(e.Values) != 0 {
		t.Errorf("Expected length of 0, got %d", len(e.Values))
	}
}

func TestGet(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")
	if e.Get("test") != "test" {
		t.Errorf("Expected value to be test, got %s", e.Get("test"))
	}
}

func TestAddComment(t *testing.T) {
	e := EnvFile{}
	e.AddComment("test")
	if len(e.Values) != 1 {
		t.Errorf("Expected length of 1, got %d", len(e.Values))
	}

	if !e.Values[0].IsComment {
		t.Errorf("Expected IsComment to be true, got %t", e.Values[0].IsComment)
	}
}

func TestKeys(t *testing.T) {
	e := EnvFile{}
	e.Add("test", "test")

	keys := e.Keys()
	if len(keys) != 1 {
		t.Errorf("Expected length of 1, got %d", len(keys))
	}

	if keys[0] != "TEST" {
		t.Errorf("Expected key to be TEST, got %s", keys[0])
	}

	// Test that comments are skipped
	e.AddComment("test")
	keys = e.Keys()
	if len(keys) != 1 {
		t.Errorf("Expected length of 1, got %d", len(keys))
	}

	if keys[0] != "TEST" {
		t.Errorf("Expected key to be TEST, got %s", keys[0])
	}
}

func TestNew(t *testing.T) {
	e := New("randompath")
	if len(e.Values) != 0 {
		t.Errorf("Expected length of 0, got %d", len(e.Values))
	}

	if e.Path != "randompath" {
		t.Errorf("Expected path to be randompath, got %s", e.Path)
	}

	if reflect.TypeOf(e).Kind() != reflect.Ptr {
		t.Errorf("Expected type to be pointer, got %s", reflect.TypeOf(e).Kind())
	}

	if reflect.TypeOf(e).String() != "*dotenv.EnvFile" {
		t.Errorf("Expected type to be *dotenv.EnvFile, got %s", reflect.TypeOf(e).String())
	}
}

func setupFile(lines []string) *os.File {
	f, _ := os.CreateTemp("", ".env")
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f
}

func TestLoad(t *testing.T) {
	f := setupFile([]string{"TEST=test"})
	defer f.Close()
	defer os.Remove(f.Name())

	e := Load(f.Name())

	if len(e.Values) != 1 {
		t.Errorf("Expected length of 1, got %d", len(e.Values))
	}

	if e.Values[0].Value != "test" {
		t.Errorf("Expected value to be test, got %s", e.Values[0].Value)
	}

	if e.Values[0].Key != "TEST" {
		t.Errorf("Expected key to be TEST, got %s", e.Values[0].Key)
	}

	if e.Path != f.Name() {
		t.Errorf("Expected path to be %s, got %s", f.Name(), e.Path)
	}
}

func TestLoadPreserveComments(t *testing.T) {
	f := setupFile([]string{"# comment", "TEST=test"})
	defer f.Close()
	defer os.Remove(f.Name())

	e := Load(f.Name())

	if len(e.Values) != 2 {
		t.Errorf("Expected length of 1, got %d", len(e.Values))
	}

	if !e.Values[0].IsComment {
		t.Errorf("Expected comment. Not found")
	}

	if e.Values[0].Value != "# comment" {
		t.Errorf("Comment not read correctly")
	}

	if e.Values[1].Key != "TEST" {
		t.Errorf("Expected key to be TEST, got %s", e.Values[0].Key)
	}

	if e.Values[1].Value != "test" {
		t.Errorf("Expected key to be TEST, got %s", e.Values[0].Key)
	}
}

func TestSaveExistingFile(t *testing.T) {
	f := setupFile([]string{"TEST=test"})
	defer f.Close()
	defer os.Remove(f.Name())

	e := Load(f.Name())
	e.Add("testing", "test")
	e.Save()

	e = Load(f.Name())
	if len(e.Values) != 2 {
		t.Errorf("Expected length of 2, got %d", len(e.Values))
	}

	if e.Values[0].Key != "TEST" {
		t.Errorf("Expected key to be TEST, got %s", e.Values[0].Key)
	}

	if e.Values[1].Key != "TESTING" {
		t.Errorf("Expected key to be TEST, got %s", e.Values[0].Key)
	}
}

func TestSaveWithComment(t *testing.T) {
	f := setupFile([]string{"TEST=test"})
	defer f.Close()
	defer os.Remove(f.Name())

	e := Load(f.Name())
	e.AddComment("comment")

	e.Save()

	e = Load(f.Name())
	if len(e.Values) != 2 {
		t.Errorf("Expected length of 2, got %d", len(e.Values))
		t.FailNow()
	}

	if e.Values[0].Key != "TEST" {
		t.Errorf("Expected key to be TEST, got %s", e.Values[0].Key)
	}

	if !e.Values[1].IsComment {
		t.Errorf("Expected comment. Not found")
	}

	if e.Values[1].Value != "# comment" {
		t.Errorf("Comment not read correctly")
	}
}

func TestSaveNewFile(t *testing.T) {
	e := LoadOrCreate("./.env.testcase")
	e.Add("test", "test")
	e.Save()
	f, err := os.Open("./.env.testcase")

	if err != nil {
		t.Errorf("File was not created")
	}

	defer f.Close()
	defer os.Remove(f.Name())

	e = Load(f.Name())
	if len(e.Values) != 1 {
		t.Errorf("Expected length of 2, got %d", len(e.Values))
	}

	if e.Values[0].Key != "TEST" {
		t.Errorf("Expected key to be TEST, got %s", e.Values[0].Key)
	}

	if e.Values[0].Value != "test" {
		t.Errorf("Expected key to be TEST, got %s", e.Values[0].Key)
	}
}
