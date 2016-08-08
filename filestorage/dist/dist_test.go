package dist

import "testing"

func TestConvertPath(t *testing.T) {
  var tests = []struct {
    path string
    distSteps int
    distRange int
    expected string
  }{
    {"foo/bar/baz.qux", 3, 1000, "foo/bar/905/358/068/baz.qux"},
    {"FooBar", 2, 10, "2/3/FooBar"},
    {"foo_bar/baz_qux", 5, 2, "foo_bar/0/1/1/0/1/baz_qux"},
    {"foo/bar/baz", 0, 0, "foo/bar/baz"},
    {"/foo/bar/baz/qux", 5, 1000, "/foo/bar/baz/000/002/800/005/064/qux"},
    {"./foo/bar", 2, 100, "./foo/91/78/bar"},
    {"fooBar/bazQux/", 1, 1000, "fooBar/229/bazQux/"},
    {"", 2, 10000, "0000/0000"},
  }
  for i, c := range tests {
    actual := ConvertPath(c.path, c.distSteps, c.distRange)
    if actual != c.expected {
      t.Errorf("ConvertPath(%v, %v, %v) == %v, expected %v in case %v", c.path, c.distSteps, c.distRange, actual, c.expected, i)
    }
  }
}
