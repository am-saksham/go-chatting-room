package main

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
}

type Room struct {
	ID         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Slug       string `db:"slug" json:"slug"`
	IsPrivate  bool   `db:"is_private" json:"is_private"`
	InviteToken string `db:"invite_token" json:"invite_token"`
	CreatedBy  int    `db:"created_by" json:"created_by"`
}