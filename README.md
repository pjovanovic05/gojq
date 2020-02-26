# gojq

This is a module for easier manipulation of JSON data in Go.
As a typed language, Go requires one to know the types of
structures and fields in JSON data before deserializing them. This can be too cumbersome, for instance, if you need to access or change a deeply nested value in a JSON object unknown at compile time.

This package was inspired by the [GoJSONQ](https://github.com/thedevsaddam/gojsonq)
which is way more capable when it comes to ways of accessing
data, but it could not modify the underlying JSON object.
So, the key feature of gojq is that it supports setting
values at properties at any depth in the JSON.

## Installation

```
go get -u github.com/pjovanovic05/gojq
```

## Example use

```golang
src := `{
    "foo": {
        "bar": {
            "baz": [
                
            ]
        }
    }
}`

jq := gojq.FromBytes([]byte(src))
jq.Select()

```