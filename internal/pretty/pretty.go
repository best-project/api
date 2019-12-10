package pretty

import (
	"fmt"
	"github.com/go-playground/validator"
)

type Kind int

const (
	Course Kind = iota
	Courses
	CourseResult
	CourseResults
	User
	Users
	Task
	Tasks
)

func (k Kind) String() string {
	switch k {
	case Course:
		return "Course"
	case Courses:
		return "Courses"
	case CourseResult:
		return "Result"
	case CourseResults:
		return "Results"
	case User:
		return "User"
	case Users:
		return "Users"
	case Task:
		return "Task"
	case Tasks:
		return "Tasks"
	default:
		return ""
	}
}

func NewCreateMessage(k Kind) string {
	return fmt.Sprintf("%s created successfully", k.String())
}

func NewRemoveMessage(k Kind) string {
	return fmt.Sprintf("%s removed successfully", k.String())
}

func NewUpdateMessage(k Kind) string {
	return fmt.Sprintf("%s updated successfully", k.String())
}

const intErr = "Internal error:"

func NewDecodeError(k Kind) string {
	return fmt.Sprintf("%s cannot decode %s", intErr, k)
}

func NewNotFoundError(k Kind) string {
	return fmt.Sprintf("%s %s not found", intErr, k)
}

func NewForbiddenError(k Kind) string {
	return fmt.Sprintf("Operation on %s forbidden", k)
}

func NewBadRequest() string {
	return "Bad request"
}

func NewAlreadyExistError(k Kind) string {
	return fmt.Sprintf("%s %s already exist", intErr, k)
}

func NewErrorSave(k Kind) string {
	return fmt.Sprintf("%s cannot save %s", intErr, k)
}

func NewErrorRemove(k Kind) string {
	return fmt.Sprintf("%s cannot remove %s", intErr, k)
}

func NewErrorValidate(k Kind, errs validator.ValidationErrors) string {
	msg := ""
	for _, err := range errs {
		msg = fmt.Sprintf("%s '%v' failed on the '%v' tag", msg, err.Field(), err.Tag())
	}
	return fmt.Sprintf("%s validation error:%s", k, msg)
}

func NewErrorUpdate(k Kind) string {
	return fmt.Sprintf("%s cannot update %s", intErr, k)
}

func NewErrorList(k Kind) string {
	return fmt.Sprintf("%s cannot list %s", intErr, k)
}

func NewErrorGet(k Kind) string {
	return fmt.Sprintf("%s cannot get %s", intErr, k)
}

func NewErrorConvert(k Kind) string {
	return fmt.Sprintf("%s cannot convert %s", intErr, k)
}

func NewInternalError() string {
	return intErr
}
