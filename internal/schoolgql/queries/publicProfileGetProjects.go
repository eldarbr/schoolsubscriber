package queries

const (
	PublicProfileGetProjects TOperationName = `publicProfileGetProjects`

	publicProfileGetProjectsQuery TQuery = `query publicProfileGetProjects($studentId: UUID!, $stageGroupId: ID!) {
  school21 {
    getStudentProjectsForPublicProfileByStageGroup(
      studentId: $studentId
      stageGroupId: $stageGroupId
    ) {
      groupName
      name
      experience
      finalPercentage
      goalId
      goalStatus
      amountAnswers
      amountReviewedAnswers
      executionType
      localCourseId
      courseType
      displayedCourseStatus
      __typename
    }
    __typename
  }
}
`
)

type VarsPublicProfileGetProjects struct {
	StageGroupID string `json:"stageGroupId"`
	StudentID    string `json:"studentId"`
}

type ResponsePublicProfileGetProjects struct {
	BaseResponse
	Data struct {
		School21 struct {
			GetStudentProjectsForPublicProfileByStageGroup []ProjectPlatform `json:"getStudentProjectsForPublicProfileByStageGroup"`
		} `json:"school21"`
	} `json:"data"`
}

type ProjectPlatform struct {
	GroupName             *string  `json:"groupName"`
	Name                  *string  `json:"name"`
	Experience            *int     `json:"experience"`
	FinalPercentage       *float64 `json:"finalPercentage"`
	GoalID                *string  `json:"goalId"`
	GoalStatus            *string  `json:"goalStatus"`
	AmountAnswers         *int     `json:"amountAnswers"`
	AmountReviewedAnswers *int     `json:"amountReviewedAnswers"`
	ExecutionType         *string  `json:"executionType"`
	LocalCourseID         *int     `json:"localCourseId"`
	CourseType            *string  `json:"courseType"`
	DisplayedCourseStatus *string  `json:"displayedCourseStatus"`
}

func (project *ProjectPlatform) GetIsBeingReviewed() bool {
	if project != nil && project.GoalStatus != nil && *project.GoalStatus == "" {
		return true
	}

	return false
}
