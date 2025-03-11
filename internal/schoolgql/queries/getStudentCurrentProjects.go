// current project

package queries

const (
	GetStudentCurrentProjects TOperationName = `getStudentCurrentProjects`

	getStudentCurrentProjectsQuery = `query getStudentCurrentProjects($userId: ID!) {
  student {
    getStudentCurrentProjects(userId: $userId) {
      ...StudentProjectItem
    }
  }
}

fragment StudentProjectItem on StudentItem {
  goalId
  name
  description
  experience
  dateTime
  finalPercentage
  laboriousness
  executionType
  goalStatus
  courseType
  displayedCourseStatus
  amountAnswers
  amountMembers
  amountJoinedMembers
  amountReviewedAnswers
  amountCodeReviewMembers
  amountCurrentCodeReviewMembers
  groupName
  localCourseId
}`
)

type VarsGetStudentCurrentProjects struct {
	UserID string `json:"userId"`
}

type ResponseGetStudentCurrentProjects struct {
	BaseResponse
	Data struct {
		Student struct {
			GetStudentCurrentProjects []struct {
				AmountAnswers                  any     `json:"amountAnswers"`
				AmountCodeReviewMembers        any     `json:"amountCodeReviewMembers"`
				AmountCurrentCodeReviewMembers any     `json:"amountCurrentCodeReviewMembers"`
				AmountJoinedMembers            any     `json:"amountJoinedMembers"`
				AmountMembers                  any     `json:"amountMembers"`
				AmountReviewedAnswers          any     `json:"amountReviewedAnswers"`
				CourseType                     *string `json:"courseType"`
				DateTime                       any     `json:"dateTime"`
				Description                    *string `json:"description"`
				DisplayedCourseStatus          *string `json:"displayedCourseStatus"`
				ExecutionType                  any     `json:"executionType"`
				Experience                     int     `json:"experience"`
				FinalPercentage                any     `json:"finalPercentage"`
				GoalID                         *string `json:"goalId"`
				GoalStatus                     *string `json:"goalStatus"`
				GroupName                      *string `json:"groupName"`
				Laboriousness                  *int    `json:"laboriousness"`
				LocalCourseID                  *int    `json:"localCourseId"`
				Name                           string  `json:"name"`
			} `json:"getStudentCurrentProjects"`
		} `json:"student"`
	} `json:"data"`
}
