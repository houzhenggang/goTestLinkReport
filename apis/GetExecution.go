package apis

import (
	"encoding/json"
	"fmt"
	// "github.com/astaxie/beego"
	"github.com/alphazero/Go-Redis"
	"github.com/royburns/goTestLinkReport/models"
	// "strconv"
	// "strings"
)

func (this *ApiController) GetLastExecution() {

	expiration := int64(60 * 60)
	spec := redis.DefaultSpec().Db(0)
	client, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		fmt.Println("Failed to create the client: ", err)
		return
	}

	fmt.Println("In ApiController.GetLastExecution()...")
	// var executions []*V_testlink_testexecution_tree
	executions := []models.V_testlink_testexecution_tree{}
	testplan := this.Input().Get("testplan")
	fmt.Println(testplan)
	if testplan != "" {
		key := testplan
		value, err := client.Get(key)
		if err != nil {
			fmt.Println("Failed to get value: ", err)
			return
		}

		if value == nil {
			fmt.Println("The data is not exists. We will query it from mysql and then store them in redis.")

			// Get Executions
			maxTPNum := 10000
			testexecution_where, err := models.GetAllExecutionsWhere(testplan)
			tpCount := len(testexecution_where)
			if err != nil || tpCount > maxTPNum {
				testplan = ""
			}

			// executions = testexecution_where
			for _, v := range testexecution_where {

				v.LastTimeRun = fmt.Sprintf("%04d-%02d-%02d", v.Date.Year(), v.Date.Month(), v.Date.Day())
				switch v.Status {
				case "b":
					v.Status = "blocked"
					break
				case "p":
					v.Status = "passed"
					break
				case "f":
					v.Status = "failed"
					break
				case "":
					v.Status = "not run"
					break
				default:
					break
				}
				executions = append(executions, v)
			}
			value, err := json.Marshal(executions)
			if err != nil {
				fmt.Println("Failed to marshal: ", err)
			}
			client.Set(key, value)
			client.Expire(key, expiration)
		} else {
			fmt.Println("The plan is exists. We will unmarshal them.")
			// var temp planinfo
			err := json.Unmarshal(value, &executions)
			if err != nil {
				fmt.Println("Failed to unmarshal: ", err)
				return
			}
		}
	}

	this.Data["json"] = &executions

	// this.TplNames = "report.tpl"
	this.ServeJson()
}

func (this *ApiController) GetSprintExecution() {

	fmt.Println("In ApiController.GetSprintExecution()...")

	expiration := int64(60 * 60)
	spec := redis.DefaultSpec().Db(0)
	client, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		fmt.Println("Failed to create the client: ", err)
		return
	}

	// var executions []*V_testlink_testexecution_tree
	executions := []models.V_toad_sprint_report{}
	sp_id := this.Input().Get("sp_id")
	sp_product := this.Input().Get("sp_product")
	sp_version := this.Input().Get("sp_version")

	fmt.Printf("%s,%s,%s\n", sp_id, sp_product, sp_version)

	if sp_id != "" {
		key := fmt.Sprintf("%s-%s-%s", sp_id, sp_product, sp_version)
		value, err := client.Get(key)
		if err != nil {
			fmt.Println("Failed to get value: ", err)
			return
		}

		if value == nil {
			fmt.Println("The data is not exists. We will query it from mysql and then store them in redis.")

			// Get Executions
			maxTPNum := 10000
			res, err := models.GetSprintExecutionsWhere(sp_id, sp_product, sp_version)
			tpCount := len(res)
			if err != nil || tpCount > maxTPNum {
				sp_id = ""
			}

			// executions = res
			for _, v := range res {

				v.LastTimeRun = fmt.Sprintf("%04d-%02d-%02d", v.Date.Year(), v.Date.Month(), v.Date.Day())
				// fmt.Println(v.Status)
				switch v.Status {
				case "b":
					v.Status = "blocked"
					break
				case "p":
					v.Status = "passed"
					break
				case "f":
					v.Status = "failed"
					break
				case "":
					v.Status = "not run"
					break
				default:
					break
				}
				executions = append(executions, v)
			}
			value, err := json.Marshal(executions)
			if err != nil {
				fmt.Println("Failed to marshal: ", err)
			}
			client.Set(key, value)
			client.Expire(key, expiration)
		} else {
			fmt.Println("The data is exists. We will unmarshal them.")
			// var temp planinfo
			err := json.Unmarshal(value, &executions)
			if err != nil {
				fmt.Println("Failed to unmarshal: ", err)
				return
			}
		}
	}

	fmt.Println(len(executions))
	this.Data["json"] = &executions

	// this.TplNames = "report.tpl"
	this.ServeJson()
}
