package queries

const (
	GetLocalCourseGoals TOperationName = `getLocalCourseGoals`

	getLocalCourseGoalsQuery = `query getLocalCourseGoals($localCourseId: ID!) {
  course {
    getLocalCourseGoals(localCourseId: $localCourseId) {
      localCourseId
      globalCourseId
      courseName
      courseType
      localCourseGoals {
        ...LocalCourse
      }
    }
  }
}

fragment RetrySettings on ModuleAttemptsSettings {
  maxModuleAttempts
  isUnlimitedAttempts
}

fragment LocalCourse on LocalCourseGoalInformation {
  localCourseGoalId
  goalId
  goalName
  description
  projectHours
  signUpDate
  beginDate
  deadlineDate
  checkDate
  isContentAvailable
  executionType
  finalPoint
  finalPercentage
  status
  periodSettings
  retriesUsed
  statusUpdateDate
  retrySettings {
    ...RetrySettings
  }
}`
)

type VarsGetLocalCourseGoals struct {
	LocalCourseID int `json:"localCourseId"`
}

type ResponseGetLocalCourseGoals struct {
	BaseResponse
	Data struct {
		Course struct {
			GetLocalCourseGoals struct {
				CourseName       string `json:"courseName"`
				CourseType       string `json:"courseType"`
				GlobalCourseId   string `json:"globalCourseId"`
				LocalCourseId    string `json:"localCourseId"`
				LocalCourseGoals []struct {
					BeginDate          any    `json:"beginDate"`
					CheckDate          any    `json:"checkDate"`
					DeadlineDate       any    `json:"deadlineDate"`
					Description        string `json:"description"`
					ExecutionType      string `json:"executionType"`
					FinalPercentage    int    `json:"finalPercentage"`
					FinalPoint         int    `json:"finalPoint"`
					GoalId             string `json:"goalId"`
					GoalName           string `json:"goalName"`
					IsContentAvailable bool   `json:"isContentAvailable"`
					PeriodSettings     string `json:"periodSettings"`
					ProjectHours       int    `json:"projectHours"`
					RetriesUsed        int    `json:"retriesUsed"`
					RetrySettings      struct {
						IsUnlimitedAttempts bool `json:"isUnlimitedAttempts"`
						MaxModuleAttempts   int  `json:"maxModuleAttempts"`
					} `json:"retrySettings"`
					SignUpDate       any    `json:"signUpDate"`
					Status           string `json:"status"`
					StatusUpdateDate string `json:"statusUpdateDate"`
				} `json:"localCourseGoals"`
			} `json:"getLocalCourseGoals"`
		} `json:"course"`
	} `json:"data"`
}
