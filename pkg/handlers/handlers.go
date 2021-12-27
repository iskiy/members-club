package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Member struct {
	Email 	string `json:"email"`
	Name	string `json:"name"`
	Date 	string `json:"date,omitempty"`
}

var membersEmails = make(map[string]bool)
var members []Member

func Home(c *gin.Context){
}

func GetMembers(c *gin.Context){
	membersJSON, err := getMembersJSON()
	if err != nil{
		sendJSONError(c, http.StatusInternalServerError, err, "database error")
	}
	log.Println(string(membersJSON))
	c.Data(http.StatusOK, "application/json",membersJSON)
}

func AddMember(c *gin.Context){
	member := validateRequestMemberBody(c)
	if member == nil{
		return
	}
	members = append(members, *member)
	membersEmails[member.Email] = true
	GetMembers(c)
}

func validateRequestMemberBody(c *gin.Context) *Member{
	member := Member{}
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		sendJSONError(c, http.StatusInternalServerError, err, "read body error")
		return nil
	}

	if err = json.Unmarshal(body, &member); err != nil{
		sendJSONError(c, http.StatusInternalServerError, err, "parse body error")
		return nil
	}

	if err = checkMemberEmail(member.Email); err != nil{
		sendJSONError(c, http.StatusBadRequest, err, err.Error())
		return nil
	}
	member.Date = time.Now().Format("01.02.2006")

	return &member
}

func getMembersJSON() ([]byte, error) {
	res, err := json.Marshal(members)
	if err != nil{
		return []byte{}, fmt.Errorf("marshal members error")
	}
	return res, nil
}

func checkMemberEmail(email string) error{
	if _, ok := membersEmails[email]; ok{
		return fmt.Errorf(fmt.Sprintf("members already contains member with email: %s", email ))
	}
	return nil
}

func sendJSONError(c *gin.Context, status int, err error, message string){
	c.Data(status, "application/json", []byte(fmt.Sprintf("{ \"error\": \"%s\" }", message)))
	log.Println(err)
	return
}