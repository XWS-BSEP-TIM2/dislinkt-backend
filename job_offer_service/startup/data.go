package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jobOffers = []*domain.JobOffer{
	{
		Id:           getIdFromHex("62752bf27407f54ce1839cb2"),
		Position:     "Test Engineer",
		Seniority:    "Junior",
		Description:  "This Software Tester is responsible for testing our mobile and web applications by executing test cases and clearly documenting the results of testing. This position works collaboratively with other testers and developers.",
		Technologies: []string{"Selenium", "JUnit"},
		CompanyName:  "Joberty",
		UserId:       getObjectId("62752bf27407f54ce1839cb8"),
	},
	{
		Id:           getIdFromHex("62752bf27407f54ce1839cb3"),
		Position:     "Software Engineer",
		Seniority:    "Senior",
		Description:  "Analyze, design, develop and produce various Mainframe applications. Prepare technical design documents and conduct design reviews. Modify and enhance existing applications. Create transmission file setups on mainframe.",
		Technologies: []string{"Java", "Spring", "MVC"},
		CompanyName:  "Naovis",
		UserId:       getObjectId("62752bf27407f54ce1839cb8"),
	},
	{
		Id:           getIdFromHex("62752bf27407f54ce1839cb4"),
		Position:     "Software Engineer",
		Seniority:    "Senior",
		Description:  "A Computer Programmer, or Systems Programmer, writes code to help software applications operate more efficiently. Their duties include designing and updating software solutions, writing and updating source-code and managing various operating systems.",
		Technologies: []string{"Java", "Spring", "MVC"},
		CompanyName:  "Synechron",
		UserId:       getObjectId("62752bf27407f54ce1839cb9"),
	},
	{
		Id:           getIdFromHex("62752bf27407f54ce1839cb5"),
		Position:     "Software Architect",
		Seniority:    "Senior",
		Description:  "A Computer Programmer, or Systems Programmer, writes code to help software applications operate more efficiently. Their duties include designing and updating software solutions, writing and updating source-code and managing various operating systems.",
		Technologies: []string{"Java", "Spring", "MVC", "UML", "Microservices"},
		CompanyName:  "Microsoft ",
		UserId:       getObjectId("62752bf27407f54ce1839cb6"),
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}

func getIdFromHex(userID string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(userID)
	return id
}
