package persistence

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"strings"
)

func clearGraphDB(transaction neo4j.Transaction) error {
	_, err := transaction.Run(
		"MATCH (n) DETACH DELETE n",
		map[string]interface{}{})
	return err
}

func initGraphDB(transaction neo4j.Transaction) error {
	_, err := transaction.Run(
		"CREATE (S1:SKILL{name: \"C\"}), (S2:SKILL{name: \"C++\"}), (S3:SKILL{name: \"C#\"}), (S4:SKILL{name: \"GO\"}), (S5:SKILL{name: \"Java\"}), (S6:SKILL{name: \"JavaScript\"}), (S7:SKILL{name: \"TypeScript\"}), (S8:SKILL{name: \"PhP\"}), (S9:SKILL{name: \"Perl\"}), (S10:SKILL{name: \"Ruby\"}), (S11:SKILL{name: \"Python\"}), (S12:SKILL{name: \"Rust\"}), (S13:SKILL{name: \"Haskell\"}),  (S14:SKILL{name: \"Vue.js\"}), (S15:SKILL{name: \"React\"}), (S16:SKILL{name: \"Angular\"}), (S17:SKILL{name: \"Bootstrap\"}), (S18:SKILL{name: \"Bulma\"}), (S19:SKILL{name: \"Sass\"}), (S20:SKILL{name: \"HTML\"}), (S21:SKILL{name: \"CSS\"}), (S22:SKILL{name: \"Materialize\"}), (S23:SKILL{name: \"Tailwind\"}), (S24:SKILL{name: \"Babel\"}),  (S25:SKILL{name: \"nodeJS\"}), (S26:SKILL{name: \"Spring Boot\"}), (S27:SKILL{name: \"Express\"}), (S28:SKILL{name: \"Kafka\"}), (S29:SKILL{name: \"Nginx\"}), (S30:SKILL{name: \"RabbitMQ\"}), (S31:SKILL{name: \".NET\"}), (S32:SKILL{name: \"Django\"}), (S33:SKILL{name: \"Electron\"}), (S34:SKILL{name: \"Rails\"}), (S35:SKILL{name: \"Laravel\"}),  (S36:SKILL{name: \"Android\"}), (S37:SKILL{name: \"Flutter\"}), (S38:SKILL{name: \"Dart\"}), (S39:SKILL{name: \"Kotlin\"}), (S40:SKILL{name: \"ReactNative\"}), (S41:SKILL{name: \"iOS\"}), (S42:SKILL{name: \"Swift\"}),  (S43:SKILL{name: \"Tensorflow\"}), (S44:SKILL{name: \"Pytorch\"}), (S45:SKILL{name: \"Pandas\"}), (S46:SKILL{name: \"Seaborn\"}), (S47:SKILL{name: \"OpenCv\"}), (S48:SKILL{name: \"Scikit Learn\"}),  (S49:SKILL{name: \"Mongo DB\"}), (S50:SKILL{name: \"My SQL\"}), (S51:SKILL{name: \"PostgreSQL\"}), (S52:SKILL{name: \"Oracle\"}), (S53:SKILL{name: \"MariaDB\"}), (S54:SKILL{name: \"CouchDB\"}),  (S55:SKILL{name: \"AWS\"}), (S56:SKILL{name: \"Docker\"}), (S57:SKILL{name: \"Jenkins\"}), (S58:SKILL{name: \"GCP\"}), (S59:SKILL{name: \"Kubernetes\"}), (S60:SKILL{name: \"Bash\"}), (S61:SKILL{name: \"CircleCI\"}), (S62:SKILL{name: \"TravisCI\"}),  (S63:SKILL{name: \"Cypress\"}), (S64:SKILL{name: \"Selenium\"}), (S65:SKILL{name: \"Karma\"}), (S66:SKILL{name: \"Junit\"}), (S67:SKILL{name: \"Jest\"}), (S68:SKILL{name: \"Mocha\"}),  (S69:SKILL{name: \"Unity\"}),  (S70:SKILL{name: \"Unreal Engine\"}),   (rasti:USER{userID: \"62752bf27407f54ce1839cb9\", username: \"rasti\"}), (zarko:USER{userID: \"62752bf27407f54ce1839cb6\", username: \"zarkoo\"}), (tara:USER{userID: \"62752bf27407f54ce1839cb7\", username: \"Jelovceva\"}), (djordje:USER{userID: \"62752bf27407f54ce1839cb8\", username: \"djole\"}), (srdjan:USER{userID: \"62752bf27407f54ce1839cb3\", username: \"srdjan\"}), (marko:USER{userID: \"62752bf27407f54ce1839cb2\", username: \"marko99\"}), (nikola:USER{userID: \"62752bf27407f54ce1839cb4\", username: \"nikola93\"}), (svetozar:USER{userID: \"62752bf27407f54ce1839cb5\", username: \"svetozar\"}),  (rasti) -[:KNOWS]-> (S5), (rasti) <-[:KNOWS]- (S5), (rasti) -[:KNOWS]-> (S11), (rasti) <-[:KNOWS]- (S11), (rasti) -[:KNOWS]-> (S16), (rasti) <-[:KNOWS]- (S16),  (djordje) -[:KNOWS]-> (S5), (djordje) <-[:KNOWS]- (S5), (djordje) -[:KNOWS]-> (S4), (djordje) <-[:KNOWS]- (S4), (djordje) -[:KNOWS]-> (S56), (djordje) <-[:KNOWS]- (S56),  (tara) -[:KNOWS]-> (S5), (tara) <-[:KNOWS]- (S5), (tara) -[:KNOWS]-> (S16), (tara) <-[:KNOWS]- (S16), (tara) -[:KNOWS]-> (S18), (tara) <-[:KNOWS]- (S18), (tara) -[:KNOWS]-> (S7), (tara) <-[:KNOWS]- (S7),  (zarko) -[:KNOWS]-> (S43), (zarko) <-[:KNOWS]- (S43), (zarko) -[:KNOWS]-> (S56), (zarko) <-[:KNOWS]- (S56), (zarko) -[:KNOWS]-> (S61), (zarko) <-[:KNOWS]- (S61),    (J1:JOB{   Id:           \"62752bf27407f54ce1839cb2\",   Position:     \"Test Engineer\",   Seniority:    \"Junior\",   Description:  \"This Software Tester is responsible for testing our mobile and web applications by executing test cases and clearly documenting the results of testing. This position works collaboratively with other testers and developers.\",   CompanyName:  \"Joberty\",   UserId:       \"62752bf27407f54ce1839cb8\" }), (J2:JOB{   Id:           \"62752bf27407f54ce1839cb3\",   Position:     \"Software Engineer\",   Seniority:    \"Senior\",   Description:  \"Analyze, design, develop and produce various Mainframe applications. Prepare technical design documents and conduct design reviews. Modify and enhance existing applications. Create transmission file setups on mainframe.\",   CompanyName:  \"Naovis\",   UserId:       \"62752bf27407f54ce1839cb8\" }), (J3:JOB{   Id:           \"62752bf27407f54ce1839cb4\",   Position:     \"Software Engineer\",   Seniority:    \"Senior\",   Description:  \"A Computer Programmer, or Systems Programmer, writes code to help software applications operate more efficiently. Their duties include designing and updating software solutions, writing and updating source-code and managing various operating systems.\",   CompanyName:  \"Synechron\",   UserId:       \"62752bf27407f54ce1839cb9\" }), (J4:JOB{   Id:           \"62752bf27407f54ce1839cb5\",   Position:     \"Software Architect\",   Seniority:    \"Senior\",   Description:  \"A Computer Programmer, or Systems Programmer, writes code to help software applications operate more efficiently. Their duties include designing and updating software solutions, writing and updating source-code and managing various operating systems.\",   CompanyName:  \"Microsoft\",   UserId:       \"62752bf27407f54ce1839cb6\" }),   (J1) -[:NEED]-> (S66), (J1) <-[:NEED]- (S66), (J1) -[:NEED]-> (S64), (J1) <-[:NEED]- (S64),  (J2) -[:NEED]-> (S5), (J2) <-[:NEED]- (S5), (J2) -[:NEED]-> (S6), (J2) <-[:NEED]- (S6), (J2) -[:NEED]-> (S17), (J2) <-[:NEED]- (S17),  (J3) -[:NEED]-> (S5), (J3) <-[:NEED]- (S5), (J3) -[:NEED]-> (S30), (J3) <-[:NEED]- (S30), (J3) -[:NEED]-> (S26), (J3) <-[:NEED]- (S26),  (J4) -[:NEED]-> (S43), (J4) <-[:NEED]- (S43), (J4) -[:NEED]-> (S44), (J4) <-[:NEED]- (S44), (J4) -[:NEED]-> (S48), (J4) <-[:NEED]- (S48), (J4) -[:NEED]-> (S47), (J4) <-[:NEED]- (S47), (J4) -[:NEED]-> (S42), (J4) <-[:NEED]- (S42) ",
		map[string]interface{}{})
	return err
}

func checkIfJobOfferExist(jobId string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (existing_job:JOB) WHERE existing_job.Id = $jobId RETURN existing_job.Id",
		map[string]interface{}{"jobId": jobId})

	if result != nil && result.Next() && result.Record().Values[0] == jobId {
		return true
	}
	return false
}

func checkIfUserExist(userID string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (existing_user:USER) WHERE existing_user.userID=$userID RETURN existing_user.userID",
		map[string]interface{}{"userID": userID})

	if result != nil && result.Next() && result.Record().Values[0] == userID {
		return true
	}
	return false
}

func checkIfSkillExist(skillName string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (existing_skill:SKILL) WHERE existing_skill.name=$skillName RETURN existing_skill.name",
		map[string]interface{}{"skillName": skillName})

	if result != nil && result.Next() && result.Record().Values[0] == skillName {
		return true
	}
	return false
}

func checkIfSkillIsPresentInJobOffer(jobId, skillName string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (j:JOB) -[n:NEED]-> (s:SKILL) WHERE j.Id=$jobId AND s.name=$skillName RETURN s.name ",
		map[string]interface{}{"jobId": jobId, "skillName": skillName})

	if result != nil && result.Next() && result.Record().Values[0] == skillName {
		return true
	}
	return false
}

func checkIfSkillIsPresentInUser(userID, skillName string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (u:USER) -[k:KNOWS]-> (s:SKILL) WHERE u.userID=$uID AND s.name=$skillName RETURN s.name",
		map[string]interface{}{"uID": userID, "skillName": skillName})

	if result != nil && result.Next() && result.Record().Values[0] == skillName {
		return true
	}
	return false
}

func createNewSkill(skillName string, transaction neo4j.Transaction) bool {
	_, err := transaction.Run(
		"CREATE (:SKILL{name:$skillName})",
		map[string]interface{}{"skillName": skillName})

	if err != nil {
		return false
	}
	return true
}

func getUserJobOffersIds(userID string, transaction neo4j.Transaction) ([]string, error) {
	result, err := transaction.Run(
		"MATCH (j:JOB) where j.UserId=$uID RETURN j.Id ",
		map[string]interface{}{"uID": userID})

	var jobOfferIds []string

	if result == nil {
		return jobOfferIds, err
	}
	for result.Next() {
		jobOfferIds = append(jobOfferIds, result.Record().Values[0].(string))
	}
	return jobOfferIds, nil
}

func searchJobOffersIds(position string, transaction neo4j.Transaction) ([]string, error) {
	result, err := transaction.Run(
		"MATCH (j:JOB) RETURN j.Id, j.Position ",
		map[string]interface{}{})

	var jobOfferIds []string

	if result == nil {
		return jobOfferIds, err
	}
	for result.Next() {
		itPos := result.Record().Values[1].(string)
		if strings.Contains(itPos, position) {
			jobOfferIds = append(jobOfferIds, result.Record().Values[0].(string))
		}
	}
	return jobOfferIds, nil
}

func getAllJobOffersIds(transaction neo4j.Transaction) ([]string, error) {
	result, err := transaction.Run(
		"MATCH (j:JOB) RETURN j.Id ",
		map[string]interface{}{})

	var jobOfferIds []string

	if result == nil {
		return jobOfferIds, err
	}
	for result.Next() {
		jobOfferIds = append(jobOfferIds, result.Record().Values[0].(string))
	}
	return jobOfferIds, nil
}

func getJobOfferSkills(jobId string, transaction neo4j.Transaction) ([]string, error) {
	result, err := transaction.Run(
		"MATCH (j:JOB) -[n:NEED]-> (s:SKILL) WHERE j.Id=$jobId return s.name ",
		map[string]interface{}{"jobId": jobId})

	var skills []string

	if err != nil || result == nil {
		return nil, err
	}
	for result.Next() {
		skills = append(skills, result.Record().Values[0].(string))
	}
	return skills, nil
}

func getJobOfferData(jobId string, transaction neo4j.Transaction) ([]string, error) {
	result, err := transaction.Run(
		"MATCH (j:JOB) WHERE j.Id=$jobId RETURN j.CompanyName, j.Description, j.Id, j.Position, j.Seniority, j.UserId ",
		map[string]interface{}{"jobId": jobId})

	var data []string

	if result == nil {
		return nil, err
	}
	for result.Next() {
		data = append(data, result.Record().Values[0].(string))
		data = append(data, result.Record().Values[1].(string))
		data = append(data, result.Record().Values[2].(string))
		data = append(data, result.Record().Values[3].(string))
		data = append(data, result.Record().Values[4].(string))
		data = append(data, result.Record().Values[5].(string))
	}
	if len(data) != 6 {
		return nil, nil
	}
	return data, nil
}

func getJobOffer(jobId string, transaction neo4j.Transaction) (*domain.JobOffer, error) {
	jobOfferData, err1 := getJobOfferData(jobId, transaction)
	if err1 != nil || jobOfferData == nil {
		return nil, err1
	}
	skills, err2 := getJobOfferSkills(jobId, transaction)
	if err2 != nil || skills == nil {
		return nil, err2
	}
	//j.CompanyName, j.Description, j.Id, j.Position, j.Seniority, j.UserId
	jobOffer := &domain.JobOffer{CompanyName: jobOfferData[0], Description: jobOfferData[1], Id: jobOfferData[2], Position: jobOfferData[3], Seniority: jobOfferData[4], UserId: jobOfferData[5]}
	jobOffer.Technologies = skills
	return jobOffer, nil
}

func createNewJobOffer(jobOffer *domain.JobOffer, transaction neo4j.Transaction) error {
	_, err := transaction.Run(
		"CREATE (:JOB{Id:$id,CompanyName:$companyName, Description:$description, Position:$position, Seniority:$seniority, UserId:$userId}) ",
		map[string]interface{}{"id": jobOffer.Id, "companyName": jobOffer.CompanyName, "description": jobOffer.Description, "position": jobOffer.Position, "seniority": jobOffer.Seniority, "userId": jobOffer.UserId})
	return err
}

func addSkillToJobOffer(jobId, skillName string, transaction neo4j.Transaction) (bool, error) {
	result, err := transaction.Run(
		"MATCH (j:JOB) WHERE j.Id=$jobId "+
			"MATCH (s:SKILL) WHERE s.name=$skillName "+
			"CREATE (j) -[:NEED]-> (s) "+
			"CREATE (s) -[:NEED]-> (j) "+
			"RETURN j.Id, s.name ",
		map[string]interface{}{"jobId": jobId, "skillName": skillName})
	if err != nil {
		return false, err
	}

	if result != nil && result.Next() && result.Record().Values[0].(string) == jobId {
		return true, nil
	}
	return false, nil
}

func updateJobOfferData(jobOffer *domain.JobOffer, transaction neo4j.Transaction) (bool, error) {
	result, err := transaction.Run(
		"MATCH (j:JOB) WHERE j.Id=$id SET j.CompanyName=$companyName, j.Description=$description, j.Position=$position, j.Seniority=$seniority, j.UserId=$userId RETURN j.Id ",
		map[string]interface{}{"id": jobOffer.Id, "companyName": jobOffer.CompanyName, "description": jobOffer.Description, "position": jobOffer.Position, "seniority": jobOffer.Seniority, "userId": jobOffer.UserId})

	if err != nil {
		return false, err
	}

	if result != nil && result.Next() && result.Record().Values[0].(string) == jobOffer.Id {
		return true, nil
	}

	return false, nil
}

func removeSkillFromJobOffer(jobId, skillName string, transaction neo4j.Transaction) (bool, error) {
	result, err := transaction.Run(
		"MATCH (j:JOB)-[n:NEED]-(s:SKILL) WHERE j.Id=$jobId AND s.name=$skillName DELETE n RETURN s.name ",
		map[string]interface{}{"jobId": jobId, "skillName": skillName})

	if err != nil {
		return false, err
	}

	if result != nil && result.Next() && result.Record().Values[0].(string) == skillName {
		return true, nil
	}

	return false, nil
}

func updateSkillsForJobOffer(jobOffer *domain.JobOffer, transaction neo4j.Transaction) (bool, error) {
	skills, errS := getJobOfferSkills(jobOffer.Id, transaction)
	if errS != nil || skills == nil {
		return false, errS
	}

	//delete old skills
	currentSkillPresent := false
	for _, currentSkill := range skills {
		currentSkillPresent = false
		for _, newSkill := range jobOffer.Technologies {
			if currentSkill == newSkill {
				currentSkillPresent = true
				break
			}
		}

		if !currentSkillPresent {
			//obrisi vezu imzejdu jobOffera i skila
			isRemoved, err := removeSkillFromJobOffer(jobOffer.Id, currentSkill, transaction)
			if err != nil || isRemoved == false {
				return false, err
			}
			fmt.Println("Removed skill " + currentSkill + " from jobOffer")
		}
	}

	// add new skills to job offer
	newSkillNotPresent := true
	for _, newSkill := range jobOffer.Technologies {
		newSkillNotPresent = true
		for _, oldSkill := range skills {
			if newSkill == oldSkill {
				newSkillNotPresent = false
				break
			}
		}

		if newSkillNotPresent {
			if !checkIfSkillExist(newSkill, transaction) {
				createNewSkill(newSkill, transaction)
			}
			isAdded, err := addSkillToJobOffer(jobOffer.Id, newSkill, transaction)
			if err != nil || isAdded == false {
				return false, err
			}
			fmt.Println("Added skill " + newSkill + " to jobOffer")
		}
	}

	return true, nil
}

func createNewUser(userID string, transaction neo4j.Transaction) bool {
	_, err := transaction.Run(
		"CREATE (:USER{userID:$uID})",
		map[string]interface{}{"uID": userID})

	if err != nil {
		return false
	}
	return true
}

func updateSkillsForUser(userID string, newSkills []string, transaction neo4j.Transaction) (bool, error) {
	skills, errS := getUserSkills(userID, transaction)
	if errS != nil || skills == nil {
		return false, errS
	}

	//delete old skills
	currentSkillPresent := false
	for _, currentSkill := range skills {
		currentSkillPresent = false
		for _, newSkill := range newSkills {
			if currentSkill == newSkill {
				currentSkillPresent = true
				break
			}
		}

		if !currentSkillPresent {
			//obrisi vezu imzejdu jobOffera i skila
			isRemoved, err := removeSkillFromUser(userID, currentSkill, transaction)
			if err != nil || isRemoved == false {
				return false, err
			}
			fmt.Println("Removed skill " + currentSkill + " from user")
		}
	}

	// add new skills to user
	newSkillNotPresent := true
	for _, newSkill := range newSkills {
		newSkillNotPresent = true
		for _, oldSkill := range skills {
			if newSkill == oldSkill {
				newSkillNotPresent = false
				break
			}
		}

		if newSkillNotPresent {
			if !checkIfSkillExist(newSkill, transaction) {
				createNewSkill(newSkill, transaction)
			}
			isAdded, err := addSkillToUser(userID, newSkill, transaction)
			if err != nil || isAdded == false {
				return false, err
			}
			fmt.Println("Added skill " + newSkill + " to user")
		}
	}

	return true, nil
}

func getUserSkills(userID string, transaction neo4j.Transaction) ([]string, error) {
	result, err := transaction.Run(
		"MATCH (u:USER) -[k:KNOWS]-> (s:SKILL) WHERE u.userID=$uID return s.name ",
		map[string]interface{}{"uID": userID})

	var skills []string

	if err != nil || result == nil {
		return nil, err
	}
	for result.Next() {
		skills = append(skills, result.Record().Values[0].(string))
	}
	return skills, nil
}

func removeSkillFromUser(userID, skillName string, transaction neo4j.Transaction) (bool, error) {
	result, err := transaction.Run(
		"MATCH (u:USER)-[k:KNOWS]-(s:SKILL) WHERE u.userID=$uID AND s.name=$skillName DELETE k RETURN s.name ",
		map[string]interface{}{"uID": userID, "skillName": skillName})

	if err != nil {
		return false, err
	}

	if result != nil && result.Next() && result.Record().Values[0].(string) == skillName {
		return true, nil
	}

	return false, nil
}

func addSkillToUser(userID, skillName string, transaction neo4j.Transaction) (bool, error) {
	result, err := transaction.Run(
		"MATCH (u:USER) WHERE u.userID=$uID "+
			"MATCH (s:SKILL) WHERE s.name=$skillName "+
			"CREATE (u) -[:KNOWS]-> (s) "+
			"CREATE (s) -[:KNOWS]-> (u) "+
			"RETURN u.userID, s.name ",
		map[string]interface{}{"uID": userID, "skillName": skillName})
	if err != nil {
		return false, err
	}

	if result != nil && result.Next() && result.Record().Values[0].(string) == userID {
		return true, nil
	}
	return false, nil
}
