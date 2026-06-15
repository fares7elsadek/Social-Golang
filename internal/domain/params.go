package domain

type UpdateUserParams struct {
    Username *string
    Email    *string
    Password *string
}

type UpdatePostParams struct {
    Title *string
    Content    *string
}

