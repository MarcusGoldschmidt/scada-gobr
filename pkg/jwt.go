package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/auth"
	"net/http"
	"strings"
)

func (s *Scadagobr) jwtMiddleware(function RequestHandlerFunction) RequestHandlerFunction {
	return func(scada *Scadagobr, w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		claims, err := scada.JwtHandler.ValidateJwt(reqToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.ContextClaimsKey, claims)

		r = r.WithContext(ctx)

		function(scada, w, r)
	}
}
