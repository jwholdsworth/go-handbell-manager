package main

type KeyMap struct {
	Controller1 string
	Controller2 string
	Button1     string
	Button2     string
	Button3     string
	Button4     string
}

type ButtonPress struct {
	First  bool
	Second bool
}

type Stroke int64
