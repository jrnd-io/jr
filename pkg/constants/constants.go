package constants

import "time"

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
