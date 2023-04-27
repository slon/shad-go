// Package console provides wrapper to javascript console.
//
// Some functions support format specifiers. The following specifiers
// are supported:
//
//	%s – Formats the value as a string
//	%d – Formats the value as an integer
//	%i – Same as %d
//	%f – Formats the value as a floating point value.
//	%o – Formats the value as an expandable DOM element (as in the Elements panel).
//	%O – Formats the value as an expandable JavaScript object.
//	%c – Formats the output string according to CSS styles you provide.
//
// This package does not provide functions for the aliases
// console.debug and console.info – use Log instead.
//
// For a more detailed explanation of the APIs, see for example
// Google's documentation at
// https://developers.google.com/chrome-developer-tools/docs/console-api.
//
// Portions of this documentation are modifications based on work
// created and shared by Google and used according to terms described
// in the Creative Commons 3.0 Attribution License.
package console

import (
	"bytes"
	"syscall/js"
)

var c = js.Global().Get("console")

// Assert writes msg to the console if b is false.
func Assert(b bool, msg interface{}) {
	c.Call("assert", b, msg)
}

// Clear clears the console.
func Clear() {
	c.Call("clear")
}

// Count writes the number of times that it has been invoked at the
// same line and with the same label.
func Count(label string) {
	c.Call("count", label)
}

// Dir prints a JavaScript representation of the specified object. If
// the object being logged is an HTML element, then the properties of
// its DOM representation are displayed.
func Dir(obj interface{}) {
	c.Call("dir", obj)
}

// DirXML Prints an XML representation of the specified object. For
// HTML elements, calling this method is equivalent to calling Log.
func DirXML(obj interface{}) {
	c.Call("dirxml", obj)
}

// Error is like Log but also includes a stack trace.
func Error(objs ...interface{}) {
	c.Call("error", objs...)
}

// Group starts a new logging group with an optional title.
//
// All console output that occurs after calling this function appears
// in the same visual group. Groups can be nested.
//
// The title will be generated according to the rules of the Log
// function.
func Group(objs ...interface{}) {
	c.Call("group", objs...)
}

// GroupCollapsed is like Group, except that the newly created
// group starts collapsed instead of open.
func GroupCollapsed(objs ...interface{}) {
	c.Call("groupCollapsed", objs...)
}

// GroupEnd closes the currently active logging group.
func GroupEnd() {
	c.Call("groupEnd")
}

// Log displays a message in the console. You pass one or more objects
// to this method, each of which are evaluated and concatenated into a
// space-delimited string. The first parameter you pass to Log may
// contain format specifiers.
func Log(objs ...interface{}) {
	c.Call("log", objs...)
}

// Profile starts a new CPU profile with an optional label.
func Profile(label interface{}) {
	c.Call("profile", label)
}

// ProfileEnd stops the currently running CPU profile, if any.
func ProfileEnd() {
	c.Call("profileEnd")
}

// Time starts a new timer with an associated label. Calling TimeEnd
// with the same label will stop the timer and print the elapsed time.
func Time(label interface{}) {
	c.Call("time", label)
}

// TimeEnd ends a timer that was started with Time and prints the
// elapsed time.
func TimeEnd(label interface{}) {
	c.Call("timeEnd", label)
}

// Timestamp adds an event to the timeline during a recording session.
func Timestamp(label interface{}) {
	c.Call("timeStamp", label)
}

// Trace prints a stack trace, starting from the point where Trace was
// called.
func Trace() {
	c.Call("trace")
}

// Warn is like Log but displays a different icon alongside the
// message.
func Warn(objs ...interface{}) {
	c.Call("warn", objs...)
}

// Writer implements an io.Writer on top of the JavaScript console.
// Writes will be buffered until a newline is encountered, which will
// cause flushing up to the newline.
type Writer struct {
	buf *bytes.Buffer
}

func (w *Writer) Write(buf []byte) (n int, err error) {
	if len(buf) == 0 {
		return 0, nil
	}

	for i := len(buf); i >= 0; i-- {
		if buf[i] == '\n' {
			w.buf.Write(buf[:i])
			Log(w.buf.String())
			w.buf.Reset()
			w.buf.Write(buf[i+1:])
			break
		}
	}

	return len(buf), nil
}

func (w *Writer) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

// Flush will flush the current line to the console.
func (w *Writer) Flush() {
	w.Write([]byte{'\n'})
}

func New() *Writer {
	return &Writer{buf: new(bytes.Buffer)}
}
