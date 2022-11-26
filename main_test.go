package main

import (
	"reflect"
	"testing"
)

func TestCreateMatch(t *testing.T) {
	type args struct {
		users   []*User
		clinics []*Clinic
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "success_one_user_one_clinic_case1",
			args: args{
				users: []*User{
					{1, "satou", []int{1}},
				},
				clinics: []*Clinic{
					{1, "a", []int{2}, []int{}, 2},
				},
			},
			want: []int{1},
		},
		{
			name: "success_one_user_one_clinic_case2",
			args: args{
				users: []*User{
					{1, "satou", []int{1}},
				},
				clinics: []*Clinic{
					{1, "a", []int{1}, []int{}, 2},
				},
			},
			want: []int{},
		},
		{
			name: "success_one_user_one_clinic_case3",
			args: args{
				users: []*User{
					{1, "satou", []int{1}},
				},
				clinics: []*Clinic{
					{1, "a", []int{1, 2, 3}, []int{2, 3}, 2},
				},
			},
			want: []int{3},
		},
		{
			name: "success_two_user_one_clinic_case1",
			args: args{
				users: []*User{
					{1, "satou", []int{1}},
					{2, "suzuki", []int{1}},
				},
				clinics: []*Clinic{
					{1, "a", []int{1}, []int{}, 2},
				},
			},
			want: []int{2},
		},
		{
			name: "success_two_user_one_clinic_case2",
			args: args{
				users: []*User{
					{1, "satou", []int{1}},
					{2, "suzuki", []int{1}},
				},
				clinics: []*Clinic{
					{1, "a", []int{1, 2}, []int{}, 2},
				},
			},
			want: []int{},
		},
		{
			name: "success_three_user_one_clinic_case1",
			args: args{
				users: []*User{
					{1, "satou", []int{1}},
					{2, "suzuki", []int{1}},
					{3, "takahashi", []int{1}},
				},
				clinics: []*Clinic{
					{1, "a", []int{1, 2}, []int{}, 2},
				},
			},
			want: []int{3},
		},
		{
			name: "success_three_user_one_clinic_case2",
			args: args{
				users: []*User{
					{1, "satou", []int{1}},
					{2, "suzuki", []int{1}},
					{3, "takahashi", []int{1}},
				},
				clinics: []*Clinic{
					{1, "a", []int{2, 3, 1}, []int{}, 2},
				},
			},
			want: []int{1},
		},
		{
			name: "success_three_user_three_clinic_case1",
			args: args{
				users: []*User{
					{1, "satou", []int{1, 2, 3}},
					{2, "suzuki", []int{2, 1, 3}},
					{3, "takahashi", []int{3, 2, 1}},
				},
				clinics: []*Clinic{
					{1, "a", []int{2, 3, 1}, []int{}, 2},
					{2, "b", []int{2, 3, 1}, []int{}, 2},
					{3, "c", []int{2, 3, 1}, []int{}, 2},
				},
			},
			want: []int{},
		},
		{
			name: "success_eight_user_four_clinic_case1",
			args: args{
				users: []*User{
					{1, "satou", []int{2}},
					{2, "suzuki", []int{2, 1}},
					{3, "takahashi", []int{2, 1}},
					{4, "tanaka", []int{1, 2, 3, 4}},
					{5, "watanabe", []int{2, 1, 4, 3}},
					{6, "yamamoto", []int{2, 3, 1, 4}},
					{7, "kobayashi", []int{2, 1, 4, 3}},
					{8, "abe", []int{4, 2, 1, 3}},
				},
				clinics: []*Clinic{
					{1, "a", []int{3, 7}, []int{}, 2},
					{2, "b", []int{7, 8, 5, 1, 2, 3, 4, 6}, []int{}, 2},
					{3, "c", []int{2, 5, 8, 1, 3, 4, 7}, []int{}, 2},
					{4, "d", []int{2, 5, 1, 3, 6, 4, 7}, []int{}, 2},
				},
			},
			want: []int{2, 1, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateMatch(tt.args.users, tt.args.clinics); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttemptMatch(t *testing.T) {
	type args struct {
		user   *User
		clinic *Clinic
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "success_no_userID_in_desired_clinic",
			args: args{
				user:   &User{1, "satou", []int{1}},
				clinic: &Clinic{1, "a", []int{2}, []int{}, 2},
			},
			want: 1,
		},
		{
			name: "success_desired_clinic_tmp_match_empty",
			args: args{
				user:   &User{1, "satou", []int{1}},
				clinic: &Clinic{1, "a", []int{1, 2}, []int{}, 2},
			},
			want: 0,
		},
		{
			name: "success_desired_clinic_tmp_match_one_vacancy",
			args: args{
				user:   &User{1, "satou", []int{1}},
				clinic: &Clinic{1, "a", []int{1, 2}, []int{2}, 2},
			},
			want: 0,
		},
		{
			name: "success_desired_clinic_tmp_match_no_vacancy_but_match",
			args: args{
				user:   &User{1, "satou", []int{1}},
				clinic: &Clinic{1, "a", []int{1, 2, 3}, []int{2, 3}, 2},
			},
			want: 3,
		},
		{
			name: "success_desired_clinic_tmp_match_no_vacancy_and_unmatch",
			args: args{
				user:   &User{1, "satou", []int{1}},
				clinic: &Clinic{1, "a", []int{2, 3, 1}, []int{2, 3}, 2},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AttemptMatch(tt.args.user, tt.args.clinic); got != tt.want {
				t.Errorf("AttemptMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindUnMatchUser(t *testing.T) {
	type args struct {
		c *Clinic
		u *User
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "success",
			args: args{
				c: &Clinic{
					1, "a", []int{2, 3, 1}, []int{2, 3}, 2,
				},
				u: &User{
					1, "satou", []int{1},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindUnMatchUser(tt.args.c, tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUnMatchUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
