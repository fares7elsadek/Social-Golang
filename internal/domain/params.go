package domain

type UpdateUserParams struct {
    Username *string
    Email    *string
}

type UpdatePostParams struct {
    Title *string
    Content    *string
}