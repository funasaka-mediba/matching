package main

import "fmt"

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

	unMatchUsers := []*User{}

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

	// TODO: アンマッチユーザーで再度仮マッチを試行する
	secondUnMatchUser := []*User{}
	for _, v := range unMatchUsers {
		unMatchUser := CreateMatch(v)
		if unMatchUser != nil {
			secondUnMatchUser = append(secondUnMatchUser, unMatchUser)
		}
	}
	for i, v := range secondUnMatchUser {
		fmt.Printf("secondUnMatchUser%v: %v\n", i, v.Name)
	}

	// TODO: アンマッチユーザーを何回やり直してあげればいいのか？？？
}

// CreateMatch ユーザーを希望するクリニックとマッチさせる
func CreateMatch(user *User) *User {
	var unMatchUser *User
	for i := 1; i <= len(user.DesiredRank); i++ {
		desiredClinic := user.DesiredRank[i]
		fmt.Printf("\nuser: %+v\n", user.Name)
		fmt.Printf("desiredClinic: %+v\n", desiredClinic.Name)
		// desiredClinicのDesiredRankのvalueの中に登録しようとしているユーザーのIDが含まれていなければ、そのユーザーはそのクリニックとはアンマッチなのでスキップ
		if !ContainsUserID(desiredClinic.DesiredRank, user.ID) {
			fmt.Println("残念、クリニックはあなたいらないって")
			fmt.Printf("user.Name: %+v, clinic.Name: %+v, unmatch\n", user.Name, desiredClinic.Name)
			continue
		}
		fmt.Println("クリニックの希望リストにユーザーが含まれてるよ")
		if len(desiredClinic.tmpMatch) < desiredClinic.Limit {
			// ユーザーの希望しているクリニックの仮マッチリストに空きがあるので仮マッチ
			fmt.Printf("len(desiredClinic.tmpMatch): %+v\n", len(desiredClinic.tmpMatch))
			fmt.Printf("user.Name: %+v, clinic.Name: %+v, tmpMatch-1\n", user.Name, desiredClinic.Name)
			desiredClinic.tmpMatch = append(desiredClinic.tmpMatch, user)
			return nil
		}
		// この時点でユーザーの希望クリニックの仮マッチリストの要素が埋まっている
		fmt.Println("仮マッチ埋まってる")
		// 最下位のユーザーを判定する
		unMatchUser = FindUnMatchUser(desiredClinic, user)
		fmt.Printf("UnMatchUser.Name: %+v, clinic.Name: %+v\n", unMatchUser.Name, desiredClinic.Name)
		if unMatchUser != user {
			// ユーザーは仮マッチできる
			fmt.Printf("user.Name: %+v, clinic.Name: %+v, tmpMatch-2\n", user.Name, desiredClinic.Name)
			desiredClinic.UpdateTmpMatch(user, unMatchUser)
			return unMatchUser
		}
		// 希望したクリニックと仮マッチできなかったため次のループへ
	}
	// 全ての希望クリニックを確認しても仮マッチできなかった場合
	return unMatchUser
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
	fmt.Printf("FindUnMatchUser:: clinic: %v, u: %v\n", clinic.Name, u.Name)
	for _, v := range clinic.tmpMatch {
		userIDs = append(userIDs, v.ID)
	}
	fmt.Printf("userIDs: %+v\n", userIDs)
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
	fmt.Printf("worstUser.Name: %+v\n", worstUser.Name)
	return worstUser
}
