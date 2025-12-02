package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.isokol.dev/internal/models"
	"snippetbox.isokol.dev/internal/validator"
)

type (
	// Snippet creation form data.
	snippetCreateForm struct {
		// Extend from validator for form validation.
		validator.Validator `form:"-"`

		// Snippet title in form data.
		Title string `form:"title"`
		// Snippet content in form data.
		Content string `form:"content"`
		// Snippet expiration in form data.
		Expires int `form:"expires"`
	}
	// User signup form.
	userSignupForm struct {
		// Extend from validator for form validation.
		validator.Validator `form:"-"`

		// User name in form data.
		Name string `form:"name"`
		// User email in form data.
		Email string `form:"email"`
		// User plaintext password in form data.
		Password string `form:"password"`
	}
	userLoginForm struct {
		// Extend from validator for form validation.
		validator.Validator `form:"-"`

		// Email in form data.
		Email string `form:"email"`
		// Password in form data.
		Password string `form:"password"`
	}
)

const (
	// Minimal logical ID for entities.
	minID = 1
	// Title length limit.
	titleLengthLimit = 100
	// Expiration option - 1 day.
	expiresInDay = 1
	// Expiration option - 1 week.
	expiresInWeek = 7
	// Expiration option - 1 year.
	expiresInYear = 365
	// Blank field validation error text.
	validationErrorBlank = "This field cannot be blank"
	// Email format validation error text.
	validationEmailInvalid = "This field must be a valid email address"
	// Home template file name.
	homeTemplateName = "home.tmpl.html"
	// Snippet view template file name.
	viewTemplateName = "view.tmpl.html"
	// Snippet creation template file name.
	createTemplateName = "create.tmpl.html"
	// User signup template file name.
	signupTemplateName = "signup.tmpl.html"
	// Login template file name.
	loginTemplateName = "login.tmpl.html"
	// Form field title.
	fieldTitle = "title"
	// Form field content.
	fieldContent = "content"
	// Form field expires.
	fieldExpires = "expires"
	// Form field name.
	fieldName = "name"
	// Form field email.
	fieldEmail = "email"
	// Form field password.
	fieldPassword = "password"
	// Password minimal length limit.
	passwordMinLength = 8
)

// Handler for home page.
func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	snippets, err := app.repositories.Snippet.Latest(request.Context())
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	data := app.newTemplateData(request)
	data.Snippets = snippets

	app.renderTemplate(writer, request, http.StatusOK, homeTemplateName, data)
}

// Handler for snippet view page.
func (app *application) snippetView(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.PathValue("id"))

	if err != nil || id < minID {
		http.NotFound(writer, request)

		return
	}

	snippet, err := app.repositories.Snippet.Get(request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(writer, request)
		} else {
			app.serverError(writer, request, err)
		}

		return
	}

	data := app.newTemplateData(request)
	data.Snippet = &snippet

	app.renderTemplate(writer, request, http.StatusOK, viewTemplateName, data)
}

// Handler for snippet create page.
func (app *application) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	data := app.newTemplateData(request)
	data.Form = snippetCreateForm{
		Expires: expiresInYear,
	}

	app.renderTemplate(writer, request, http.StatusOK, createTemplateName, data)
}

// Handler for snippet creation request.
func (app *application) snippetCreatePost(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var form snippetCreateForm

	err := app.decodePostForm(request, &form)
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)

		return
	}

	form.validate()

	if !form.Valid() {
		data := app.newTemplateData(request)
		data.Form = form
		app.renderTemplate(writer, request, http.StatusUnprocessableEntity, createTemplateName, data)

		return
	}

	id, err := app.repositories.Snippet.Insert(request.Context(), form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	app.sessionManager.Put(request.Context(), sessionFlashField, "Snippet successfully created!")
	http.Redirect(writer, request, fmt.Sprintf(snippetViewRoute+"/%d", id), http.StatusSeeOther)
}

// Validate snippet creation form.
func (form *snippetCreateForm) validate() {
	validator.CheckField(
		&form.Validator,
		validator.CreateNotBlankValidator(),
		form.Title,
		fieldTitle,
		validationErrorBlank,
	)

	validator.CheckField(
		&form.Validator,
		validator.CreateMaxCharsValidator(titleLengthLimit),
		form.Title,
		fieldTitle,
		fmt.Sprintf("This field cannot be more than %d characters long", titleLengthLimit),
	)

	validator.CheckField(
		&form.Validator, validator.CreateNotBlankValidator(), form.Content, fieldContent, validationErrorBlank)

	validator.CheckField(
		&form.Validator,
		validator.CreatePermittedValueValidator(expiresInDay, expiresInWeek, expiresInYear),
		form.Expires,
		fieldExpires,
		fmt.Sprintf(
			"This field must be either %d, %d or %d",
			expiresInDay,
			expiresInWeek,
			expiresInYear,
		),
	)
}

// Handler for user signup page.
func (app *application) userSignup(writer http.ResponseWriter, request *http.Request) {
	data := app.newTemplateData(request)
	data.Form = userSignupForm{}
	app.renderTemplate(writer, request, http.StatusOK, signupTemplateName, data)
}

// Handler for user signup request.
func (app *application) userSignupPost(writer http.ResponseWriter, request *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(request, &form)
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)

		return
	}

	form.validate()

	if !form.Valid() {
		data := app.newTemplateData(request)
		data.Form = form
		app.renderTemplate(writer, request, http.StatusUnprocessableEntity, signupTemplateName, data)

		return
	}

	err = app.repositories.User.Insert(request.Context(), form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError(fieldEmail, "Email address is already in use")

			data := app.newTemplateData(request)
			data.Form = form
			app.renderTemplate(writer, request, http.StatusUnprocessableEntity, signupTemplateName, data)
		} else {
			app.serverError(writer, request, err)
		}

		return
	}

	app.sessionManager.Put(request.Context(), sessionFlashField, "Your signup was successful. Please log in.")

	http.Redirect(writer, request, userLoginRoute, http.StatusSeeOther)
}

// User signup form validation.
func (form *userSignupForm) validate() {
	validator.CheckField(
		&form.Validator,
		validator.CreateNotBlankValidator(),
		form.Name,
		fieldName,
		validationErrorBlank,
	)
	validator.CheckField(
		&form.Validator,
		validator.CreateNotBlankValidator(),
		form.Email,
		fieldEmail,
		validationErrorBlank,
	)
	validator.CheckField(
		&form.Validator,
		validator.CreateMatchesRegexValidator(validator.EmailRX),
		form.Email,
		fieldEmail,
		validationEmailInvalid,
	)
	validator.CheckField(
		&form.Validator,
		validator.CreateNotBlankValidator(),
		form.Password,
		fieldPassword,
		validationErrorBlank,
	)
	validator.CheckField(
		&form.Validator,
		validator.CreateMinCharsValidator(passwordMinLength),
		form.Password,
		fieldPassword,
		fmt.Sprintf("This field must be at least %d characters long", passwordMinLength),
	)
}

// Handler for user login page.
func (app *application) userLogin(writer http.ResponseWriter, request *http.Request) {
	data := app.newTemplateData(request)
	data.Form = userLoginForm{}
	app.renderTemplate(writer, request, http.StatusOK, loginTemplateName, data)
}

// Handler for user login request.
func (app *application) userLoginPost(writer http.ResponseWriter, request *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(request, &form)
	if err != nil {
		app.clientError(writer, http.StatusBadRequest)

		return
	}

	form.validate()

	if !form.Valid() {
		data := app.newTemplateData(request)
		data.Form = form
		app.renderTemplate(writer, request, http.StatusUnprocessableEntity, loginTemplateName, data)

		return
	}

	id, err := app.repositories.User.Authenticate(request.Context(), form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(request)
			data.Form = form
			app.renderTemplate(writer, request, http.StatusUnprocessableEntity, loginTemplateName, data)
		} else {
			app.serverError(writer, request, err)
		}

		return
	}

	err = app.sessionManager.RenewToken(request.Context())
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	app.sessionManager.Put(request.Context(), sessionAuthenticatedUserField, id)

	http.Redirect(writer, request, snippetCreateRoute, http.StatusSeeOther)
}

// User login form validation.
func (form *userLoginForm) validate() {
	validator.CheckField(
		&form.Validator,
		validator.CreateNotBlankValidator(),
		form.Email,
		fieldEmail,
		validationErrorBlank,
	)
	validator.CheckField(
		&form.Validator,
		validator.CreateMatchesRegexValidator(validator.EmailRX),
		form.Email,
		fieldEmail,
		validationEmailInvalid,
	)
	validator.CheckField(
		&form.Validator,
		validator.CreateNotBlankValidator(),
		form.Password,
		fieldPassword,
		validationErrorBlank,
	)
}

// Handler for user logout request.
func (app *application) userLogoutPost(writer http.ResponseWriter, request *http.Request) {
	err := app.sessionManager.RenewToken(request.Context())
	if err != nil {
		app.serverError(writer, request, err)

		return
	}

	app.sessionManager.Remove(request.Context(), sessionAuthenticatedUserField)

	app.sessionManager.Put(request.Context(), sessionFlashField, "You've been logged out successfully!")

	http.Redirect(writer, request, homeRoute, http.StatusSeeOther)
}

// Handler for health check.
func (app *application) healthCheck(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte("OK"))
	if err != nil {
		app.serverError(writer, request, err)
	}
}
