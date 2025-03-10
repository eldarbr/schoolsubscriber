package queries

// goalId -> taskId
const (
	GetProjectInfoByStudent TOperationName = `getProjectInfoByStudent`

	getProjectInfoByStudentQuery = `query getProjectInfoByStudent($goalId: ID!, $studentId: UUID!) {
  school21 {
    getModuleById(goalId: $goalId, studentId: $studentId) {
      ...ProjectInfo
    }
    getModuleCoverInformation(goalId: $goalId, studentId: $studentId) {
      ...ModuleCoverInfo
    }
    getP2PChecksInfo(goalId: $goalId, studentId: $studentId) {
      ...P2PInfo
    }
    getGoalRetryInfo(goalId: $goalId, studentId: $studentId) {
      ...StudentGoalRetryInfo
    }
  }
}

fragment TimelineItemChildren on ProjectTimelineItem {
  type
  elementType
  status
  start
  end
  order
}

fragment ProjectReviewsInfo on ProjectReviewsInfo {
  reviewByStudentCount
  relevantReviewByStudentsCount
  reviewByInspectionStaffCount
  relevantReviewByInspectionStaffCount
}

fragment TimelineItem on ProjectTimelineItem {
  type
  status
  start
  end
  children {
    ...TimelineItemChildren
  }
}

fragment CurrentInternshipTaskInfo on StudentTask {
  id
  taskId
  task {
    id
    assignmentType
    taskSolutionType
    studentTaskAdditionalAttributes {
      cookiesCount
      maxCodeReviewCount
      codeReviewCost
      ciCdMode
    }
    checkTypes
    taskSolutionType
  }
  lastAnswer {
    id
  }
  teamSettings {
    ...teamSettingsInfo
  }
}

fragment RetrySettings on ModuleAttemptsSettings {
  maxModuleAttempts
  isUnlimitedAttempts
}

fragment teamSettingsInfo on TeamSettings {
  teamCreateOption
  minAmountMember
  maxAmountMember
  enableSurrenderTeam
}

fragment StudentGoalRetryInfo on StudentGoalRetryInfo {
  totalRetryValue
  usedRetryCount
  unlimitedAttempts
}

fragment P2PInfo on P2PChecksInfo {
  cookiesCount
  periodOfVerification
  projectReviewsInfo {
    ...ProjectReviewsInfo
  }
}

fragment ModuleCoverInfo on ModuleCoverInformation {
  isOwnStudentTimeline
  softSkills {
    softSkillId
    softSkillName
    totalPower
    maxPower
    currentUserPower
    achievedUserPower
    teamRole
  }
  timeline {
    ...TimelineItem
  }
}

fragment ProjectInfo on StudentModule {
  id
  moduleTitle
  finalPercentage
  finalPoint
  goalExecutionType
  displayedGoalStatus
  accessBeforeStartProgress
  resultModuleCompletion
  finishedExecutionDateByScheduler
  durationFromStageSubjectGroupPlan
  currentAttemptNumber
  isDeadlineFree
  isRetryAvailable
  localCourseId
  courseBaseParameters {
    isGradedCourse
  }
  teamSettings {
    ...teamSettingsInfo
  }
  studyModule {
    id
    idea
    duration
    goalPoint
    retrySettings {
      ...RetrySettings
    }
    levels {
      id
      goalElements {
        id
        tasks {
          id
          taskId
        }
      }
    }
  }
  currentTask {
    ...CurrentInternshipTaskInfo
  }
}`
)

type VarsGetProjectInfoByStudent struct {
	GoalID    string `json:"goalId"`
	StudentID string `json:"studentId"`
}

type ResponseGetProjectInfoByStudent struct {
	BaseResponse
	Data struct {
		School21 struct {
			// GetGoalRetryInfo struct {
			// 	TotalRetryValue   int  `json:"totalRetryValue"`
			// 	UnlimitedAttempts bool `json:"unlimitedAttempts"`
			// 	UsedRetryCount    int  `json:"usedRetryCount"`
			// } `json:"getGoalRetryInfo"`
			GetModuleById struct {
				AccessBeforeStartProgress bool `json:"accessBeforeStartProgress"`
				CourseBaseParameters      struct {
					IsGradedCourse bool `json:"isGradedCourse"`
				} `json:"courseBaseParameters"`
				CurrentAttemptNumber int `json:"currentAttemptNumber"`
				CurrentTask          struct {
					ID         string `json:"id"`
					LastAnswer struct {
						ID string `json:"id"`
					} `json:"lastAnswer"`
					Task struct {
						AssignmentType                  string   `json:"assignmentType"`
						CheckTypes                      []string `json:"checkTypes"`
						ID                              string   `json:"id"`
						StudentTaskAdditionalAttributes struct {
							CiCdMode           string `json:"ciCdMode"`
							CodeReviewCost     int    `json:"codeReviewCost"`
							CookiesCount       int    `json:"cookiesCount"`
							MaxCodeReviewCount int    `json:"maxCodeReviewCount"`
						} `json:"studentTaskAdditionalAttributes"`
						TaskSolutionType string `json:"taskSolutionType"`
					} `json:"task"`
					TaskID       string `json:"taskId"`
					TeamSettings any    `json:"teamSettings"`
				} `json:"currentTask"`
				DisplayedGoalStatus               string `json:"displayedGoalStatus"`
				DurationFromStageSubjectGroupPlan any    `json:"durationFromStageSubjectGroupPlan"`
				FinalPercentage                   any    `json:"finalPercentage"`
				FinalPoint                        any    `json:"finalPoint"`
				FinishedExecutionDateByScheduler  any    `json:"finishedExecutionDateByScheduler"`
				GoalExecutionType                 string `json:"goalExecutionType"`
				ID                                string `json:"id"`
				IsDeadlineFree                    bool   `json:"isDeadlineFree"`
				IsRetryAvailable                  bool   `json:"isRetryAvailable"`
				LocalCourseId                     string `json:"localCourseId"`
				ModuleTitle                       string `json:"moduleTitle"`
				ResultModuleCompletion            any    `json:"resultModuleCompletion"`
				StudyModule                       struct {
					Duration  int    `json:"duration"`
					GoalPoint int    `json:"goalPoint"`
					ID        string `json:"id"`
					Idea      string `json:"idea"`
					Levels    []struct {
						GoalElements []struct {
							ID    string `json:"id"`
							Tasks []struct {
								ID     string `json:"id"`
								TaskID string `json:"taskId"`
							} `json:"tasks"`
						} `json:"goalElements"`
						ID string `json:"id"`
					} `json:"levels"`
					RetrySettings struct {
						IsUnlimitedAttempts bool `json:"isUnlimitedAttempts"`
						MaxModuleAttempts   int  `json:"maxModuleAttempts"`
					} `json:"retrySettings"`
				} `json:"studyModule"`
				TeamSettings any `json:"teamSettings"`
			} `json:"getModuleById"`
			// GetModuleCoverInformation struct {
			// 	IsOwnStudentTimeline bool `json:"isOwnStudentTimeline"`
			// 	SoftSkills           []struct {
			// 		AchievedUserPower any    `json:"achievedUserPower"`
			// 		CurrentUserPower  int    `json:"currentUserPower"`
			// 		MaxPower          int    `json:"maxPower"`
			// 		SoftSkillId       int    `json:"softSkillId"`
			// 		SoftSkillName     string `json:"softSkillName"`
			// 		TeamRole          any    `json:"teamRole"`
			// 		TotalPower        int    `json:"totalPower"`
			// 	} `json:"softSkills"`
			// 	Timeline []struct {
			// 		Children []struct {
			// 			ElementType string `json:"elementType"`
			// 			End         any    `json:"end"`
			// 			Order       int    `json:"order"`
			// 			Start       any    `json:"start"`
			// 			Status      string `json:"status"`
			// 			Type        string `json:"type"`
			// 		} `json:"children"`
			// 		End    any    `json:"end"`
			// 		Start  any    `json:"start"`
			// 		Status string `json:"status"`
			// 		Type   string `json:"type"`
			// 	} `json:"timeline"`
			// } `json:"getModuleCoverInformation"`
			// GetP2PChecksInfo struct {
			// 	CookiesCount         int `json:"cookiesCount"`
			// 	PeriodOfVerification int `json:"periodOfVerification"`
			// 	ProjectReviewsInfo   struct {
			// 		RelevantReviewByInspectionStaffCount int `json:"relevantReviewByInspectionStaffCount"`
			// 		RelevantReviewByStudentsCount        int `json:"relevantReviewByStudentsCount"`
			// 		ReviewByInspectionStaffCount         int `json:"reviewByInspectionStaffCount"`
			// 		ReviewByStudentCount                 int `json:"reviewByStudentCount"`
			// 	} `json:"projectReviewsInfo"`
			// } `json:"getP2PChecksInfo"`
		} `json:"school21"`
	} `json:"data"`
}
