package main

import (
	"fmt"
	"reflect"
)

// Clinic 配属先クリニック
type Clinic struct {
	ID          int
	Name        string
	DesiredRank []int // 希望順位(ユーザーID)
	TmpMatch    []int // 仮マッチ
	Limit       int   // 受け入れ人数限界
}

// User ユーザー
type User struct {
	ID          int
	Name        string
	DesiredRank []int // 希望順位
}

func main() {
	a := Clinic{1, "a", []int{3, 7}, []int{}, 2}
	b := Clinic{2, "b", []int{7, 8, 5, 1, 2, 3, 4, 6}, []int{}, 2}
	c := Clinic{3, "c", []int{2, 5, 8, 1, 3, 4, 7}, []int{}, 2}
	d := Clinic{4, "d", []int{2, 5, 1, 3, 6, 4, 7}, []int{}, 2}

	clinics := []*Clinic{
		&a,
		&b,
		&c,
		&d,
	}

	satou := User{1, "satou", []int{2}}
	suzuki := User{2, "suzuki", []int{2, 1}}
	takahashi := User{3, "takahashi", []int{2, 1}}
	tanaka := User{4, "tanaka", []int{1, 2, 3, 4}}
	watanabe := User{5, "watanabe", []int{2, 1, 4, 3}}
	yamamoto := User{6, "yamamoto", []int{2, 3, 1, 4}}
	kobayashi := User{7, "kobayashi", []int{2, 1, 4, 3}}
	abe := User{8, "abe", []int{4, 2, 1, 3}}

	users := []*User{
		&satou,
		&suzuki,
		&takahashi,
		&tanaka,
		&watanabe,
		&yamamoto,
		&kobayashi,
		&abe,
	}

	fmt.Println("1 try")
	unMatchUserIDs := CreateMatch(users, clinics)
	for _, c := range clinics {
		fmt.Printf("clinic name: %v, tmpMatch: %v\n", c.Name, c.TmpMatch)
	}
	fmt.Printf("unMatchUserIDs: %+v\n", unMatchUserIDs)
	var unMatchUsers []*User
	for _, ID := range unMatchUserIDs {
		for _, u := range users {
			if ID == u.ID {
				unMatchUsers = append(unMatchUsers, u)
			}
		}
	}

	unMatchUserIDsDic := [][]int{}
	for i := 0; ; i++ {
		fmt.Printf("%v try\n", i+2)
		unMatchUserIDs, unMatchUsers = RetryMatch(users, unMatchUsers, clinics)
		unMatchUserIDsDic = append(unMatchUserIDsDic, unMatchUserIDs)
		if i > 0 {
			if reflect.DeepEqual(unMatchUserIDsDic[i-1], unMatchUserIDsDic[i]) {
				break
			}
		}
	}

	fmt.Println("もう仮マッチもアンマッチも変動ないので、処理終了")
}

func RetryMatch(users, unMatchUsers []*User, clinics []*Clinic) ([]int, []*User) {
	unMatchUserIDs := CreateMatch(unMatchUsers, clinics)
	for _, c := range clinics {
		fmt.Printf("clinic name: %v, tmpMatch: %v\n", c.Name, c.TmpMatch)
	}
	fmt.Printf("unMatchUserIDs: %+v\n", unMatchUserIDs)

	var unMatchUsers2 []*User
	for _, ID := range unMatchUserIDs {
		for _, u := range users {
			if ID == u.ID {
				unMatchUsers2 = append(unMatchUsers2, u)
			}
		}
	}
	return unMatchUserIDs, unMatchUsers2
}

// CreateMatch ユーザーを希望するクリニックとマッチさせる
func CreateMatch(users []*User, clinics []*Clinic) []int {
	unMatchUserID := 0
	unMatchUserIDs := []int{}
	for _, user := range users {
	loop:
		for i := 0; i < len(user.DesiredRank); i++ {
			for _, clinic := range clinics {
				if user.DesiredRank[i] == clinic.ID {
					unMatchUserID = AttemptMatch(user, clinic)
					switch unMatchUserID {
					case 0:
						// 空きがあって仮マッチできた場合 -> 別のuserへ
						break loop
					case user.ID:
						// 仮マッチできなかった場合 -> 別のクリニックへ
						continue
					default:
						// 空きがなかったが、仮マッチできた場合 -> 別のuserへ
						unMatchUserIDs = append(unMatchUserIDs, unMatchUserID)
						break loop
					}
				}
			}
		}
		// 全ての希望クリニックを確認しても仮マッチできなかった場合
		if unMatchUserID == user.ID {
			unMatchUserIDs = append(unMatchUserIDs, user.ID)
		}
	}
	return unMatchUserIDs
}

// AttemptMatch 一人のUserと一つのClinic間で仮マッチを試みる
func AttemptMatch(user *User, clinic *Clinic) int {
	unMatchUserID := 0

	if !ContainsUserID(clinic.DesiredRank, user.ID) {
		// 希望しているclinicの希望順位にuserのIDがない場合
		// このクリニックとはアンマッチ
		return user.ID
	}
	// この行に到達するということは、userの希望しているclinicの希望順位にuserのIDが含まれているのは確定

	// まずuserの希望しているclinicのTmpMatchに空きがあるかチェック
	if len(clinic.TmpMatch) < clinic.Limit {
		// 空きがある場合
		clinic.TmpMatch = append(clinic.TmpMatch, user.ID) // 仮マッチ
		return 0
	}
	// この行に到達した時点で、userの希望しているclinicのTmpMatchに空きがない
	unMatchUserID = FindUnMatchUser(clinic, user) // 最下位のユーザーのIDが決定

	if unMatchUserID != user.ID {
		// 仮マッチリストの最下位ユーザーのIDを新規ユーザーIDで更新
		clinic.UpdateTmpMatch(user.ID, unMatchUserID)
		return unMatchUserID
	}
	return unMatchUserID
}

// ContainsUserID クリニックの希望順位リストにIDが含まれているか判定
func ContainsUserID(desiredRank []int, ID int) bool {
	for _, v := range desiredRank {
		if v == ID {
			return true
		}
	}
	return false
}

// UpdateTmpMatch 仮マッチリストを新規ユーザーIDで更新
func (clinic *Clinic) UpdateTmpMatch(userID, unMatchUserID int) {
	for i, v := range clinic.TmpMatch {
		if v == unMatchUserID {
			clinic.TmpMatch[i] = userID
		}
	}
}

// FindUnMatchUser TmpMatch内のユーザーと新規ユーザーの中で最下位のユーザーのIDを判定する
func FindUnMatchUser(c *Clinic, u *User) int {
	userIDs := []int{u.ID}
	userIDs = append(userIDs, c.TmpMatch...)

	desiredRank := c.DesiredRank

	// userIDsのIDの中で、desiredRankの中で一番右にあるIDを特定する
	worstID := 0
	for i := 0; i < len(desiredRank); i++ {
		for j := 0; j < len(userIDs); j++ {
			if desiredRank[i] == userIDs[j] {
				worstID = userIDs[j]
			}
		}
	}
	return worstID
}
