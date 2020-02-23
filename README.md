# gojq

This is a module for easier manipulation of JSON data in Go.
As a typed language, Go requires one to know the types of
structures and fields in JSON data before deserializing them,
which can be too cumbersome.

This utility was inspired by the [GoJSONQ](https://github.com/thedevsaddam/gojsonq)
which is way more capable when it comes to ways of accessing
data, but it could not modify the underlying JSON object.
So, the key feature of gojq is that it supports setting
values at properties at any depth in the JSON.

TODO example usage...
