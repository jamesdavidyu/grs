package types

import "time"

type InviteeStore interface {
	GetInviteeById(id string) (*Invitees, error)
	GetInviteeWithName(name string) (*Invitees, error)
	CreateInvitee(Invitees) error
	GetRoleById(id string) (*Invitees, error)
}

type RsvpStore interface {
	GetRsvpByInviteeId(inviteeId string) (*Rsvp, error)
	CreateRsvp(Rsvp) error
	UpdateRsvpWithId(Rsvp) error
}

type GuestsStore interface {
	GetGuestsByInviteeId(inviteeId string) (*Guests, error)
	CreateGuests(Guests) error
	UpdateGuestsWithId(Guests) error
}

type DashboardStore interface {
	GetRsvps() ([]Dashboard, error)
}

type NewInviteeStore interface {
	CreateNewInvitee(NewInvitee) error
}

type Invitees struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
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

type Guests struct {
	Id        string    `json:"id"`
	InviteeId string    `json:"inviteeId"`
	Guests    string    `json:"guests"`
	CreatedAt time.Time `json:"createdAt"`
}

type GuestsPayload struct {
	InviteeId string `json:"inviteeId"`
	Guests    string `json:"guests"`
}

type Dashboard struct {
	Id     string `json:"id"`
	Rsvp   string `json:"rsvp"`
	Name   string `json:"name"`
	Guests string `json:"guests"`
}

type NewInvitee struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Rsvp      string    `json:"rsvp"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewInviteePayload struct {
	Name string `json:"name"`
	Rsvp string `json:"rsvp"`
}
