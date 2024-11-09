package types

import "time"

type InviteeStore interface {
	GetInviteeById(id string) (*Invitees, error)
	GetInviteeWithName(name string) (*Invitees, error)
	CreateInvitee(Invitees) error
}

type RsvpStore interface {
	GetRsvpByInviteeId(inviteeId string) (*Rsvp, error)
	CreateRsvp(Rsvp) error
	UpdateRsvpWithId(Rsvp) error
}

type Invitees struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type Register struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Rsvp struct {
	Id        string    `json:"id"`
	InviteeId string    `json:"inviteeId"`
	Rsvp      string    `json:"rsvp"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateRsvpPayload struct {
	InviteeId string `json:"inviteeId"`
	Rsvp      string `json:"rsvp"`
}
