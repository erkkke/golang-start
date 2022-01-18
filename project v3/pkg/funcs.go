package pkg

import (
	"context"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/models"
	"net/http"
)

func IsUserAdmin(ctx context.Context, w http.ResponseWriter) bool {
	if userInfo := ctx.Value(CtxKeyUser).(*models.AuthorizedUserInfo); userInfo.Role != models.Admin {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "err: Insufficient rights to access data")
		return false
	}

	return true
}

