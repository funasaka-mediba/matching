package main

import (
	"fmt"
)

// Clinic 配属先クリニック
type Clinic struct {
	ID          int
	Name        string
	DesiredRank []int   // 希望順位(ユーザーID)
	tmpMatch    []*User // 仮マッチリスト
	Limit       int     // 受け入れ人数限界
}

// User ユーザー
type User struct {
	ID          int
	Name        string
	DesiredRank map[int]*Clinic // 希望順位
}

func main() {
	a := Clinic{1, "a", []int{3, 7}, []*User{}, 2}
	b := Clinic{2, "b", []int{7, 8, 5, 1, 2, 3, 4, 6}, []*User{}, 2}
	c := Clinic{3, "c", []int{2, 5, 8, 1, 3, 4, 7}, []*User{}, 2}
	d := Clinic{4, "d", []int{2, 5, 1, 3, 6, 4, 7}, []*User{}, 2}

	satou := User{1, "satou", map[int]*Clinic{1: &b}}
	suzuki := User{2, "suzuki", map[int]*Clinic{1: &b, 2: &a}}
	takahashi := User{3, "takahashi", map[int]*Clinic{1: &b, 2: &a}}
	tanaka := User{4, "tanaka", map[int]*Clinic{1: &a, 2: &b, 3: &c, 4: &d}}
	watanabe := User{5, "watanabe", map[int]*Clinic{1: &b, 2: &a, 3: &d, 4: &c}}
	yamamoto := User{6, "yamamoto", map[int]*Clinic{1: &b, 2: &c, 3: &a, 4: &d}}
	kobayashi := User{7, "kobayashi", map[int]*Clinic{1: &b, 2: &a, 3: &d, 4: &c}}
	abe := User{8, "abe", map[int]*Clinic{1: &d, 2: &b, 3: &a, 4: &c}}

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

	var unMatchUsers []*User

	for _, user := range users {
		unMatchUser := CreateMatch(user)
		if unMatchUser != nil {
			unMatchUsers = append(unMatchUsers, unMatchUser)
		}
	}
	for _, v := range a.tmpMatch {
		fmt.Printf("a.tmpMatch: %+v\n", v.Name)
	}
	for _, v := range b.tmpMatch {
		fmt.Printf("b.tmpMatch: %+v\n", v.Name)
	}
	for _, v := range c.tmpMatch {
		fmt.Printf("c.tmpMatch: %+v\n", v.Name)
	}
	for _, v := range d.tmpMatch {
		fmt.Printf("d.tmpMatch: %+v\n", v.Name)
	}
	for i, v := range unMatchUsers {
		fmt.Printf("unMatchUser%v: %v\n", i, v.Name)
	}

	AttemptCreateMatch(unMatchUsers)

	for _, v := range a.tmpMatch {
		fmt.Printf("a.tmpMatch: %+v\n", v.Name)
	}
	for _, v := range b.tmpMatch {
		fmt.Printf("b.tmpMatch: %+v\n", v.Name)
	}
	for _, v := range c.tmpMatch {
		fmt.Printf("c.tmpMatch: %+v\n", v.Name)
	}
	for _, v := range d.tmpMatch {
		fmt.Printf("d.tmpMatch: %+v\n", v.Name)
	}
	for i, v := range unMatchUsers {
		fmt.Printf("unMatchUser%v: %v\n", i, v.Name)
	}
}

// AttemptUnMatchUserMatch アンマッチユーザーのマッチング
func AttemptUnMatchUserMatch(unMatchUsers []*User) []*User {
	var unMatchUsers2 []*User
	cnt := 0
	// unMatchUsersとunMatchUsers2が同じにならない限り無限ループ
	for i := 0; ; i++ {
		if i%2 == 0 {
			for _, u := range unMatchUsers {
				unMatchUser := CreateMatch(u)
				if unMatchUser != nil {
					unMatchUsers2 = append(unMatchUsers2, unMatchUser)
				}
			}
			if unMatchUsers2 == nil {
				return nil
			}
			/* TODO
			ここがダメっぽい。unMatchUsersに含まれるユーザーが持ってる
			DesiredRankの中身に入っているtmpMatchの内容はCreateMatchごとに変わっていくから、
			IDか何かをみて比較しなくちゃダメ。
			そもそも、ユーザー自体に、Clinicの構造体が含まれるmapを持たせるのが設計的に良くないかも。
			ClinicのIDだけ持たせてなんとか設計を変えるか。
			*/
			cnt = 0
			for i := 0; i < len(unMatchUsers); i++ {
				if unMatchUsers[i].ID != unMatchUsers2[i].ID {
					break
				}
				cnt++
			}
			if cnt == len(unMatchUsers) && cnt == len(unMatchUsers2) {
				return unMatchUsers
			}
			if i == 0 {
				continue
			}
			unMatchUsers = nil
		}
		if i%2 == 1 {
			for _, u := range unMatchUsers2 {
				unMatchUser := CreateMatch(u)
				if unMatchUser != nil {
					unMatchUsers = append(unMatchUsers, unMatchUser)
				}
			}
			if unMatchUsers == nil {
				return nil
			}
			cnt = 0
			for i := 0; i < len(unMatchUsers); i++ {
				if unMatchUsers[i].ID != unMatchUsers2[i].ID {
					break
				}
				cnt++
			}
			if cnt == len(unMatchUsers) && cnt == len(unMatchUsers2) {
				return unMatchUsers
			}
			unMatchUsers2 = nil
		}
	}
}

// CreateMatch ユーザーを希望するクリニックとマッチさせる
func CreateMatch(user *User) *User {
	var unMatchUser *User
	for i := 1; i <= len(user.DesiredRank); i++ {
		desiredClinic := user.DesiredRank[i]
		// desiredClinicのDesiredRankのvalueの中に登録しようとしているユーザーのIDが含まれていなければ、そのユーザーはそのクリニックとはアンマッチなのでスキップ
		if !ContainsUserID(desiredClinic.DesiredRank, user.ID) {
			continue
		}

		if len(desiredClinic.tmpMatch) < desiredClinic.Limit {
			// ユーザーの希望しているクリニックの仮マッチリストに空きがあるので仮マッチ
			desiredClinic.tmpMatch = append(desiredClinic.tmpMatch, user)
			return nil
		}
		// この時点でユーザーの希望クリニックの仮マッチリストの要素が埋まっている

		// 最下位のユーザーを判定する
		unMatchUser = FindUnMatchUser(desiredClinic, user)

		if unMatchUser != user {
			// ユーザーは仮マッチできる
			desiredClinic.UpdateTmpMatch(user, unMatchUser)
			return unMatchUser
		}
		// 希望したクリニックと仮マッチできなかったため次のループへ
	}
	// 全ての希望クリニックを確認しても仮マッチできなかった場合
	if unMatchUser == nil {
		unMatchUser = user
	}
	return unMatchUser
}

// AttemptCreateMatch ユーザーを希望するクリニックとマッチさせる
func AttemptCreateMatch(users []*User) {
	var unMatchUsers []*User
	for _, user := range users {
		var unMatchUser *User
		for i := 1; i <= len(user.DesiredRank); i++ {
			desiredClinic := user.DesiredRank[i]
			// desiredClinicのDesiredRankのvalueの中に登録しようとしているユーザーのIDが含まれていなければ、そのユーザーはそのクリニックとはアンマッチなのでスキップ
			if !ContainsUserID(desiredClinic.DesiredRank, user.ID) {
				continue
			}

			if len(desiredClinic.tmpMatch) < desiredClinic.Limit {
				// ユーザーの希望しているクリニックの仮マッチリストに空きがあるので仮マッチ
				desiredClinic.tmpMatch = append(desiredClinic.tmpMatch, user)
				break
			}
			// この時点でユーザーの希望クリニックの仮マッチリストの要素が埋まっている

			// 最下位のユーザーを判定する
			unMatchUser = FindUnMatchUser(desiredClinic, user)

			if unMatchUser != user {
				// ユーザーは仮マッチできる
				desiredClinic.UpdateTmpMatch(user, unMatchUser)
				break
			}
			// 希望したクリニックと仮マッチできなかったため次のループへ
		}
		// 全ての希望クリニックを確認しても仮マッチできなかった場合
		if unMatchUser != nil {
			unMatchUsers = append(unMatchUsers, unMatchUser)
		}
	}

	if len(unMatchUsers) != 0 {
		AttemptCreateMatch(unMatchUsers)
	}
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

// InsertTmpMatch ユーザーを仮マッチリストに登録
func (clinic *Clinic) InsertTmpMatch(u *User) {
	clinic.tmpMatch = append(clinic.tmpMatch, u)
}

// UpdateTmpMatch 仮マッチリストを新規ユーザーで更新
func (clinic *Clinic) UpdateTmpMatch(u, unMatchUser *User) {
	for i, v := range clinic.tmpMatch {
		if v == unMatchUser {
			clinic.tmpMatch[i] = u
		}
	}
}

// FindUnMatchUser tmpMatch内のユーザーと新規ユーザーの中で最下位のユーザーを判定する
func FindUnMatchUser(clinic *Clinic, u *User) *User {
	userIDs := []int{u.ID}
	for _, v := range clinic.tmpMatch {
		userIDs = append(userIDs, v.ID)
	}
	desiredRank := clinic.DesiredRank

	// userIDsのIDの中で、desiredRankの中で一番右にあるIDを特定する
	worstID := 0
	for i := 0; i < len(desiredRank); i++ {
		for j := 0; j < len(userIDs); j++ {
			if desiredRank[i] == userIDs[j] {
				worstID = userIDs[j]
			}
		}
	}

	// 最下位のIDを持つユーザーを特定する
	worstUser := &User{}
	if worstID == u.ID {
		worstUser = u
	} else {
		for _, v := range clinic.tmpMatch {
			if v.ID == worstID {
				worstUser = v
			}
		}
	}
	return worstUser
}
