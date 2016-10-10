package api

import (
	"fmt"
	"net/http"

	"github.com/zhonglin6666/universal/util"
)

func getMasterTest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("getMasterTest ......")

	util.Response(nil, http.StatusOK, w)
}
