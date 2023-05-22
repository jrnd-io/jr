/*
Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/*
JR is a CLI program that helps you to create quality random data for your applications.

JR is very straightforward to use. To list existing templates:
> jr list

Templates are in the directory $HOME/.jr/templates.
You can override with the --templatePath command flag Templates with parsing issues are showed in red, Templates with no parsing issues are showed in green

To use for example one of the predefined templates, net-device:

> jr template run net-device

Using -n option you can create more data in each pass. This example creates 3 net-device objects at once:

> jr template run net-device -n 3

Using --frequency option you can repeat the creation every f milliseconds. This example creates 2 net-device every second, for ever:

> jr template run net-device -n 2 -f 1s

Using --duration option you can time bound the entire object creation. This example creates 2 net-device every 100ms for 1 minute:

> jr template run net-device -n 2 -f 100ms -d 1m
*/
package main

import (
	"github.com/ugol/jr/pkg/cmd"
)

func main() {
	cmd.Execute()
}
