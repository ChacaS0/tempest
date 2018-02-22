# Test related stuff

## Naming hint
``<package_name>_test.go`` or ``<file_name>_test.go`` ?

## Sample of a test function

```go
import "testing"

func Test(t *testing.T) {
    var tests = []struct {
        s, want string
    }{
        {"Hello", "HelloHello"},
        {"Good Bye!", "Good Bye!Good Bye!"},
        {"日本", "日本日本"},
        {"", ""},
    }

    for _, c := range tests {
        got := Double(c.s)
        if got != c.want {
            t.Errorf("Double(%s) == %s, want %s", c.s, got, c.want)
        }
    }
}
```