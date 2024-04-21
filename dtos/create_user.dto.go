package dtos

type CreateUserInput struct {
	Provider   string
	Fullname   string
	Username   string
	Email      string
	Bio        string
	PictureURL string
}
