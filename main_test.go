package main

import (
	"reflect"
	"testing"
)

func TestCreateMatch(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "success",
			args: args{
				user: &User{
					ID:   3,
					Name: "takahashi",
					DesiredRank: map[int]*Clinic{1: {
						ID:          3,
						Name:        "c",
						DesiredRank: []int{2, 5, 8, 1, 3, 4, 7},
						tmpMatch: []*User{
							{
								ID:          1,
								Name:        "satou",
								DesiredRank: map[int]*Clinic{},
							},
							{
								ID:          4,
								Name:        "tanaka",
								DesiredRank: map[int]*Clinic{},
							},
						},
						Limit: 2,
					}},
				},
			},
			want: &User{
				ID:          4,
				Name:        "tanaka",
				DesiredRank: map[int]*Clinic{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateMatch(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsUserID(t *testing.T) {
	type args struct {
		desiredRank []int
		ID          int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				desiredRank: []int{7, 8, 6, 1, 2, 3, 4, 6},
				ID:          2,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsUserID(tt.args.desiredRank, tt.args.ID); got != tt.want {
				t.Errorf("ContainsUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClinic_InsertTmpMatch(t *testing.T) {
	type fields struct {
		ID          int
		Name        string
		DesiredRank []int
		tmpMatch    []*User
		Limit       int
	}
	type args struct {
		u *User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				ID:          1,
				Name:        "a",
				DesiredRank: []int{3, 7},
				tmpMatch:    []*User{},
				Limit:       2,
			},
			args: args{
				u: &User{
					ID:          3,
					Name:        "takahashi",
					DesiredRank: map[int]*Clinic{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinic := &Clinic{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				DesiredRank: tt.fields.DesiredRank,
				tmpMatch:    tt.fields.tmpMatch,
				Limit:       tt.fields.Limit,
			}
			clinic.InsertTmpMatch(tt.args.u)
		})
	}
}

func TestClinic_UpdateTmpMatch(t *testing.T) {
	type fields struct {
		ID          int
		Name        string
		DesiredRank []int
		tmpMatch    []*User
		Limit       int
	}
	type args struct {
		u           *User
		unMatchUser *User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				ID:          3,
				Name:        "c",
				DesiredRank: []int{2, 5, 8, 1, 3, 4, 7},
				tmpMatch: []*User{
					{
						ID:          1,
						Name:        "satou",
						DesiredRank: map[int]*Clinic{},
					},
					{
						ID:          4,
						Name:        "tanaka",
						DesiredRank: map[int]*Clinic{},
					},
				},
				Limit: 2,
			},
			args: args{
				u: &User{
					ID:          3,
					Name:        "takahashi",
					DesiredRank: map[int]*Clinic{},
				},
				unMatchUser: &User{
					ID:          4,
					Name:        "tanaka",
					DesiredRank: map[int]*Clinic{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clinic := &Clinic{
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				DesiredRank: tt.fields.DesiredRank,
				tmpMatch:    tt.fields.tmpMatch,
				Limit:       tt.fields.Limit,
			}
			clinic.UpdateTmpMatch(tt.args.u, tt.args.unMatchUser)
		})
	}
}

func TestFindUnMatchUser(t *testing.T) {
	type args struct {
		clinic Clinic
		u      *User
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "success",
			args: args{
				clinic: Clinic{
					ID:          3,
					Name:        "c",
					DesiredRank: []int{2, 5, 8, 1, 3, 4, 7},
					tmpMatch: []*User{
						{
							ID:          1,
							Name:        "satou",
							DesiredRank: map[int]*Clinic{},
						},
						{
							ID:          4,
							Name:        "tanaka",
							DesiredRank: map[int]*Clinic{},
						},
					},
					Limit: 2,
				},
				u: &User{
					ID:          3,
					Name:        "takahashi",
					DesiredRank: map[int]*Clinic{},
				},
			},
			want: &User{
				ID:          4,
				Name:        "tanaka",
				DesiredRank: map[int]*Clinic{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindUnMatchUser(&tt.args.clinic, tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUnMatchUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
