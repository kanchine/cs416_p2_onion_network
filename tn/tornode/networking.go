package tornode

import (
	"math/rand"
	"strconv"
	"time"
)

// GetNewUnusedPort generate a random local ip port
func getNewUnusedPort() string {
	// TODO: you should check whether this random port is really unused
	// assuming IPs are all localhost
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	// generate a random port between [min, max]
	min := 3000
	max := 60000
	port := random.Intn(max-min) + min
	return strconv.Itoa(port)
}
