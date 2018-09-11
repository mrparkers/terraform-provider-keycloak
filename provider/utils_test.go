package provider

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"math/rand"
)

func randomBool() bool {
	return rand.Intn(2) == 0
}

func randomStringInSlice(slice []string) string {
	return slice[acctest.RandIntRange(0, len(slice)-1)]
}
