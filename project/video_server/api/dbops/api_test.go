package dbops

import (
	"testing"
	"strconv"
	"time"
	"fmt"
	"video_server/api/utils"
)

var (
	tempvid string
	tempsid string
)
// test用例编写步骤 ：第一步初始化，第二步运行tests ， 第三步清除数据
// init(dblogin, truncate tables)-> run tests -> clear data(truncate tables)

//清除表数据
func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

//测试的工作流（顺序） 先添加，然后查看，接着删除，最后再次查看
func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("avenssi", "123")
	if err != nil {
		t.Errorf("Error of AddUser: %v", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("avenssi")
	if pwd != "123" || err != nil {
		t.Errorf("Error of GetUser")
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("avenssi", "123")
	if err != nil {
		t.Errorf("Error of DeleteUser: %v", err)
	}
}

func testRegetUser(t *testing.T) {
	pwd, err := GetUserCredential("avenssi")
	if err != nil {
		t.Errorf("Error of RegetUser: %v", err)
	}
	if pwd != "" {
		t.Errorf("Deleting user test failed")
	}
}

//测试工作流
func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	vi, err:=AddNewVideo(1, "my-video")
	if err!=nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tempvid = vi.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err:=GetVideoInfo(tempvid)
	if err!=nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err:=DeleteVideoInfo(tempvid)
	if err!=nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	vi, err:= GetVideoInfo(tempvid)
	if err!=nil || vi !=nil{
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddComments", testAddComments)
	t.Run("ListComments", testListComments)
}

func testAddComments(t * testing.T) {
	vid:="12345"
	aid:=1
	content:= "I like this video"
	err:=AddNewComments(vid, aid, content)
	if err!=nil {
		t.Errorf("Error of AddComments: %v", err)
	}
}

func testListComments(t *testing.T) {
	vid:="12345"
	from:=1514764800
	to, _:=strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))
	res, err:=ListComments(vid, from, to)
	if err!=nil{
		t.Errorf("Error of ListComments: %v", err)
	}
	for i, ele:=range res {
		fmt.Printf("comment: %d, %+v \n", i, ele)
	}
}

func TestSessions(t *testing.T) {
	clearTables()
	t.Run("AddSession", testAddSession)
	t.Run("RetriveOneSession", testRetriveSession)
	clearTables()
}

func testAddSession(t * testing.T) {
	sid, err:=utils.NewUUID()
	if err!=nil{
		t.Errorf("Error of UUID, %v", err)
	}
	tempsid = sid
	ttl:=int64(129183174987124)
	err= InsertSession(sid, ttl, "skyone")
	if err!=nil{
		t.Errorf("Error of InsertSession: %v", err)
	}
}

func testRetriveSession(t *testing.T) {
	res, err:=RetrieveSession(tempsid)
	if err!=nil{
		t.Errorf("Error of RetriveSession: %v", err)
	}
	fmt.Printf("session: %+v", res)
}