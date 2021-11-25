package graph

import (
	"go-webrtc/infra/firebase"
	"net/http"
)

const (
	authKey      = "Authorization"
	debugAuthKey = "X-User-Id"
)

type Authenticate func(base http.Handler) http.Handler

func NewAuthenticate(
	contextProvider ContextProvider,
	fireClient firebase.Client) Authenticate {
	return func(base http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if uid := r.Header.Get(debugAuthKey); uid != "" {
				base.ServeHTTP(w, r.WithContext(contextProvider.WithAuthUID(ctx, firebase.UID(uid))))
			} else {
				token := r.Header.Get(authKey)
				if len(token) <= 7 {
					http.Error(w, "unauthorized", 401)
					return
				}

				uid, err := fireClient.VerifyToken(ctx, token[7:])
				if err != nil {
					http.Error(w, "unauthorized", 401)
					return
				}

				base.ServeHTTP(w, r.WithContext(contextProvider.WithAuthUID(ctx, uid)))
			}
		})
	}
}
