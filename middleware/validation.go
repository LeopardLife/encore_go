package middleware

import (
	"encore.app/global"
	"encore.dev/beta/auth"
	"encore.dev/middleware"
)

//encore:middleware global target=all
func ValidationMiddlewareALL(req middleware.Request, next middleware.Next) middleware.Response {

	// If the payload has a Validate method, use it to validate the request.
	payload := req.Data()

	data := auth.Data().(*global.DataAuth)

	print("testeeee	")
	print(payload)
	print(data.Username)
	// if validator, ok := payload.(interface{ Validate() error }); ok {
	// 	if err := validator.Validate(); err != nil {
	// 		// If the validation fails, return an InvalidArgument error.
	// 		err = errs.WrapCode(err, errs.InvalidArgument, "validation failed")
	// 		return middleware.Response{Err: err}
	// 	}
	// }
	return next(req)
}
