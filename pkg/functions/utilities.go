// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package functions

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/ugol/jr/pkg/ctx"
)

// Counter creates a counter named c, starting from start and incrementing by step
func Counter(c string, start, step int) int {
	ctx.JrContext.CtxCountersLock.Lock()
	defer ctx.JrContext.CtxCountersLock.Unlock()
	val, exists := ctx.JrContext.CtxCounters[c]
	if exists {
		ctx.JrContext.CtxCounters[c] = val + step
		return ctx.JrContext.CtxCounters[c]
	} else {
		ctx.JrContext.CtxCounters[c] = start
		return start
	}
}

// Image generates a random Image url of given width, height and type
func Image(width int, height int) string {
	imageType := []string{"abstract", "animals", "business", "cats", "city", "fashion", "food", "nature", "nightlife", "people", "sport", "technics", "transport"}
	return ImageOf(
		width,
		height,
		imageType[Random.Intn(len(imageType))],
	)
}

// ImageOf generates a random Image url of given width, height and type
func ImageOf(width int, height int, imageType string) string {
	return fmt.Sprintf("https://loremflickr.com/%d/%d/%s", width, height, imageType)
}

// RandomBool returns a random boolean
func RandomBool() string {
	b := Random.Intn(2)
	if b == 0 {
		return "false"
	} else {
		return "true"
	}
}

// UniqueId returns a random uuid
func UniqueId() string {
	return uuid.New().String()
}

// YesOrNo returns a random yes or no
func YesOrNo() string {
	b := Random.Intn(2)
	if b == 0 {
		return "no"
	} else {
		return "yes"
	}
}

// Contains checks if the str string is present in an array of string []string
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Inject is used to inject a different value with a given probability, typically used to generate a bad value
func Inject(probability float64, injected, original any) any {
	if Random.Float64() < probability {
		return injected
	}
	return original
}

// FromCsv gets the label value from csv file
func FromCsv(c string) string {
	ctx.JrContext.CtxCSVLock.Lock()
	defer ctx.JrContext.CtxCSVLock.Unlock()

	if len(ctx.JrContext.CtxCSV) > 0 {
		return ctx.JrContext.CtxCSV[(ctx.JrContext.CurrentIterationLoopIndex-1)%len(ctx.JrContext.CtxCSV)][c]
	} else {
		return ""
	}

}
