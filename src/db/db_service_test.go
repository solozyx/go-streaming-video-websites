package db

import (
	"testing"
	"fmt"
	"strconv"
	"time"
)

/*
init(dblogin,truncate tables) 每次跑测试表都是空的 初始化状态
-->
run tests
-->
clear data (truncate tables)
避免test和其他包test相互干扰 清理数据-->测试-->清理数据
*/

// 在整个test case没办法传参 用一个全局变量
var tempvid string

/*
test初始化
*/
func TestMain(m *testing.M){
	clearTables()
	// 跑所有的test
	m.Run()
	clearTables()
}

func clearTables(){
	dbConn.Exec("TRUNCATE users")
	dbConn.Exec("TRUNCATE video_info")
	dbConn.Exec("TRUNCATE comments")
	dbConn.Exec("TRUNCATE sessions")
}

/*
1条测试路径 规整 子test 调用顺序
*/
func TestUserWorkFlow(t *testing.T){
	t.Run("Add",testAddUser)
	t.Run("Get",testGetUser)
	t.Run("Delete",testDeleteUser)
	t.Run("ReGet",testReGetUser)
}

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddCommnets", testAddComments)
	t.Run("ListComments", testListComments)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("solozyx","123")
	if err != nil{
		t.Errorf("mysql 插入 user 错误 %s",err)
	}
}

func testGetUser(t *testing.T) {
	pwd,err := GetUserCredential("solozyx")
	if pwd != "123" || err != nil {
		t.Errorf("mysql 查询 user 错误 %s",err)
	}
	fmt.Println("pwd = ",pwd)
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("solozyx","123")
	if err != nil {
		t.Errorf("mysql 删除 user 错误 %s",err)
	}
}

func testReGetUser(t *testing.T) {
	pwd,err := GetUserCredential("solozyx")
	if err != nil {
		t.Errorf("mysql user 错误 %s",err)
	}
	if pwd != ""{
		t.Errorf("mysql 删除 user 错误 %s",err)
	}
}

func testAddVideoInfo(t *testing.T) {
	vi, err := AddNewVideo(1, "my-video")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tempvid = vi.Id
	fmt.Println("tempvid = ",tempvid)
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	vi, err := GetVideoInfo(tempvid)
	if err != nil || vi != nil{
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "I like this video"
	err := AddNewComments(vid, aid, content)
	if err != nil {
		t.Errorf("Error of AddComments: %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "12345"
	// 2周前的一个随机时间
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/(1000*1000*1000), 10))
	res, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("Error of ListComments: %v", err)
	}
	for i, ele := range res {
		// 在test case中不应该 fmt.Printf 打印结果
		// 需要使用另外的机制 判断 是否与预期结果一致
		// 这里简化处理 直接打印结果
		fmt.Printf("comment: %d, %v \n", i, ele)
	}
}

/*
C:\goenv\githubProj\go-streaming-video-websites\src\db>go test -v
=== RUN   TestUserWorkFlow
=== RUN   TestUserWorkFlow/Add
=== RUN   TestUserWorkFlow/Get
pwd =  123
=== RUN   TestUserWorkFlow/Delete
=== RUN   TestUserWorkFlow/ReGet
--- PASS: TestUserWorkFlow (0.05s)
    --- PASS: TestUserWorkFlow/Add (0.02s)
    --- PASS: TestUserWorkFlow/Get (0.00s)
    --- PASS: TestUserWorkFlow/Delete (0.02s)
    --- PASS: TestUserWorkFlow/ReGet (0.00s)
=== RUN   TestVideoWorkFlow
=== RUN   TestVideoWorkFlow/PrepareUser
=== RUN   TestVideoWorkFlow/AddVideo
tempvid =  51af7f37-6836-4761-bb11-10a90ef71a06
=== RUN   TestVideoWorkFlow/GetVideo
=== RUN   TestVideoWorkFlow/DelVideo
=== RUN   TestVideoWorkFlow/RegetVideo
--- PASS: TestVideoWorkFlow (0.19s)
    --- PASS: TestVideoWorkFlow/PrepareUser (0.02s)
    --- PASS: TestVideoWorkFlow/AddVideo (0.04s)
    --- PASS: TestVideoWorkFlow/GetVideo (0.00s)
    --- PASS: TestVideoWorkFlow/DelVideo (0.02s)
    --- PASS: TestVideoWorkFlow/RegetVideo (0.00s)
=== RUN   TestComments
=== RUN   TestComments/AddUser
=== RUN   TestComments/AddCommnets
=== RUN   TestComments/ListComments
comment: 0, &{36c4d188-af70-4378-be5f-53b8ee5be456 12345 solozyx I like this video}
--- PASS: TestComments (0.15s)
    --- PASS: TestComments/AddUser (0.02s)
    --- PASS: TestComments/AddCommnets (0.02s)

*/