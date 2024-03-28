package internals

type FormValidationErr struct {
	ElementID string
	Message   string
}

type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type DynamicNavigateForm struct {
	Page string `form:"page" binding:"required"`
}
