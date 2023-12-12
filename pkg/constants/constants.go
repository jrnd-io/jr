//Copyright © 2022 Ugo Landini <ugo.landini@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package constants

import "time"

var JRhome string

const NUM = 1
const LOCALE = "us"
const FREQUENCY = -1
const YEAR = time.Hour * 24 * 7 * 365
const DEFAULT_HOMEDIR = "$HOME/.jr"
const DEFAULT_KEY = "null"
const DEFAULT_OUTPUT = "stdout"
const DEFAULT_OUTPUT_TEMPLATE = "{{.V}}\n"
const DEFAULT_OUTPUT_KCAT_TEMPLATE = "{{.K}},{{.V}}\n"
const KAFKA_CONFIG = "./kafka/config.properties"
const DEFAULT_PARTITIONS = 6
const DEFAULT_REPLICA = 3
const DEFAULT_PRELOAD_SIZE = 0
const DEFAULT_ENV_PREFIX = "JR"
const DEFAULT_EMITTER_NAME = "emitter"
const DEFAULT_VALUE_TEMPLATE = "user"
const DEFAULT_TOPIC = "test"
const DEFAULT_HTTP_PORT = 7482
