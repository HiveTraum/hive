package middlewares

//
//func AuthenticationMiddleware(loginController infrastructure.LoginControllerInterface) func(methods map[string][]string) http.HandlerFunc {
//	return func(methods map[string][]string) func(next http.HandlerFunc) http.HandlerFunc {
//		return func(next http.HandlerFunc) http.HandlerFunc {
//			return func(w http.ResponseWriter, r *http.Request) {
//
//				req := functools.Request{Request: r}
//
//				status, user, ctx := loginController.Login(r.Context(), req.GetAuthorizationHeader())
//				if status == enums.Ok && user != nil {
//					r = r.WithContext(ctx)
//				}
//				next(w, r)
//			}
//		}
//	}
//}
