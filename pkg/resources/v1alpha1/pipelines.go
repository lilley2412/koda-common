package v1alpha1

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PipelineRun struct {
	ID             interface{}         `json:"id" bson:"_id,omitempty"`
	UUID           string              `json:"uuid" bson:"uuid,omitempty"`
	Name           string              `json:"name" bson:"name,omitempty"`
	Namespace      string              `json:"namespace"`
	Labels         map[string]string   `json:"labels,omitempty"`
	Annotations    map[string]string   `json:"annotations,omitempty"`
	CompletedAt    *time.Time          `json:"completedAt,omitempty"`
	StartedAt      *time.Time          `json:"startedAt,omitempty"`
	Duration       *time.Duration      `json:"duration,omitempty"`
	CreatedAt      time.Time           `json:"createdAt,omitempty"`
	Status         PipelineStatus      `json:"status,omitempty"`
	Tasks          map[string]*TaskRun `json:"tasks,omitempty"`
	TotalTasks     int                 `json:"totalTasks,omitempty"`
	RunningTasks   int                 `json:"runningTasks,omitempty"`
	PendingTasks   int                 `json:"pendingTasks,omitempty"`
	FailedTasks    int                 `json:"failedTasks,omitempty"`
	SucceededTasks int                 `json:"succeededTasks,omitempty"`
	CompleteTasks  int                 `json:"completeTasks,omitempty"`
}

type TaskRun struct {
	TaskName string         `json:"taskName"`
	PodName  string         `json:"podName"`
	Status   PipelineStatus `json:"status,omitempty"`
	// Parents  []*TaskRun     `json:"parents"`
	Parents []*TaskRun `json:"parents"`
}

type PipelineStatus int16

const (
	Undefined          PipelineStatus = iota
	PipelineRunPending                // 1
	NotStarted                        // 2
	Running                           // 3
	Success                           // 4
	Failed                            // 5
	Pending                           // 6
)

func (p PipelineStatus) String() string {
	switch p {
	case PipelineRunPending:
		return "PipelineRunPending"
	case NotStarted:
		return "NotStarted"
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failed:
		return "Failed"
	case Pending:
		return "Pending"
	}
	return "unknown"
}

func (p *PipelineRun) addTaskSpecs(tasks []v1beta1.PipelineTask) error {
	p.TotalTasks = len(tasks)

	// index all tasks
	for _, task := range tasks {
		p.Tasks[task.Name] = &TaskRun{
			TaskName: task.Name,
		}
	}

	// add relationships
	for _, task := range tasks {
		tr := p.Tasks[task.Name]
		for _, ra := range task.RunAfter {
			tr.Parents = append(tr.Parents, p.Tasks[ra])
		}
	}

	// if len(p.Tasks) == 0 {
	// 	if len(t.RunAfter) > 0 {
	// 		return fmt.Errorf("task %s: runAfter not allowed on first task", t.Name)
	// 	}
	// 	p.Tasks = append(p.Tasks, &TaskRun{
	// 		TaskSpecName: t.Name,
	// 	})
	// 	return nil
	// }

	// for _, ra := range t.RunAfter {
	// 	// find all parents
	// }

	return nil
}

func NewPipelineRun(uns *unstructured.Unstructured) (*PipelineRun, error) {
	// defer instrument.Duration(instrument.Track("v1alpha1.NewPipelineRun"))

	pr := &PipelineRun{
		Name:        uns.GetName(),
		Namespace:   uns.GetNamespace(),
		Labels:      uns.GetLabels(),
		Annotations: uns.GetAnnotations(),
		CreatedAt:   uns.GetCreationTimestamp().Time,
		UUID:        string(uns.GetUID()),
		Status:      NotStarted,
		Tasks:       make(map[string]*TaskRun),
		ID:          string(uns.GetUID()),
	}

	uns.UnstructuredContent()

	var spec map[string]interface{}
	specInt, ok := uns.Object["spec"]
	if ok {
		spec = specInt.(map[string]interface{})
		if status, ok := spec["status"]; ok {
			if strings.EqualFold(fmt.Sprintf("%s", status), "PipelineRunPending") {
				pr.Status = PipelineRunPending
			}
		}
	}

	var pStatus v1beta1.PipelineRunStatus
	stInt, ok := uns.Object["status"]
	if ok {
		d, err := json.Marshal(stInt)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(d, &pStatus)
		if err != nil {
			return nil, err
		}
	} else {
		return pr, nil
	}

	err := pr.addTaskSpecs(pStatus.PipelineSpec.Tasks)
	if err != nil {
		return nil, err
	}

	// if ok {
	// 	if st, ok = stInt.(map[string]interface{}); !ok {
	// 		return pr, nil
	// 	}
	// } else {
	// 	return pr, nil
	// }
	if pStatus.CompletionTime != nil {
		pr.CompletedAt = &pStatus.CompletionTime.Time
	}

	if pStatus.StartTime != nil {
		pr.StartedAt = &pStatus.StartTime.Time
	}

	pr.setDuration()

	// if r, ok := st["completionTime"]; ok {
	// 	t, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", r))
	// 	if err == nil {
	// 		pr.CompletedAt = &t
	// 	}
	// }

	// if r, ok := st["startTime"]; ok {
	// 	t, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", r))
	// 	if err == nil {
	// 		pr.StartedAt = &t
	// 	}
	// }

	// var conditions []interface{}
	// if c, ok := st["conditions"]; ok {
	// 	if conditions, ok = c.([]interface{}); ok {
	for _, condition := range pStatus.Conditions {
		// if condition, ok := conditionInt.(map[string]interface{}); ok {
		// reason := fmt.Sprintf("%s", condition["reason"])
		// status := fmt.Sprintf("%s", condition["status"])
		// cType := fmt.Sprintf("%s", condition["type"])
		if strings.EqualFold(string(condition.Type), "Succeeded") {
			if strings.EqualFold(condition.Reason, "Running") {
				pr.Status = Running
			} else if strings.EqualFold(condition.Reason, "Succeeded") {
				if strings.EqualFold(string(condition.Status), "True") {
					pr.Status = Success
				} else {
					pr.Status = Failed
				}
			} else if strings.EqualFold(condition.Reason, "Failed") {
				pr.Status = Failed
			}
		}
	}
	// 	}
	// }
	// }

	// prs := v1beta1.PipelineRunStatus{}
	// prs := make(map[string]*v1beta1.PipelineRunTaskRunStatus)
	// if truns, ok := st["taskRuns"]; ok {
	// 	d, err := json.Marshal(truns)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	err = json.Unmarshal(d, &prs)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	for _, taskRunStatus := range pStatus.TaskRuns {
		if taskRunStatus.Status != nil && len(taskRunStatus.Status.Conditions) > 0 {
			cond := taskRunStatus.Status.Conditions[0]

			pr.Tasks[taskRunStatus.PipelineTaskName].PodName = taskRunStatus.Status.PodName

			if strings.EqualFold(cond.Reason, "running") {
				pr.Tasks[taskRunStatus.PipelineTaskName].Status = Running
				// pr.Tasks = append(pr.Tasks, &TaskRun{
				// 	Status: Running,
				// })
				pr.RunningTasks++
				continue
			}

			if strings.EqualFold(cond.Reason, "pending") {
				// pr.Tasks = append(pr.Tasks, &TaskRun{
				// 	Status: Pending,
				// })
				pr.Tasks[taskRunStatus.PipelineTaskName].Status = Pending
				pr.PendingTasks++
				continue
			}

			if strings.EqualFold(cond.Reason, "Succeeded") && strings.EqualFold(string(cond.Status), "True") {
				// pr.Tasks = append(pr.Tasks, &TaskRun{
				// 	Status: Success,
				// })
				pr.Tasks[taskRunStatus.PipelineTaskName].Status = Success
				pr.SucceededTasks++
				continue
			}

			if strings.EqualFold(cond.Reason, "Succeeded") && strings.EqualFold(string(cond.Status), "False") {
				// pr.Tasks = append(pr.Tasks, &TaskRun{
				// 	Status: Failed,
				// })
				pr.Tasks[taskRunStatus.PipelineTaskName].Status = Failed
				pr.FailedTasks++
				continue
			}

			if strings.EqualFold(cond.Reason, "Failed") && strings.EqualFold(string(cond.Status), "False") && strings.EqualFold(string(cond.Type), "Succeeded") {
				pr.Tasks[taskRunStatus.PipelineTaskName].Status = Failed
				pr.FailedTasks++
				continue
			}
		} else {
			// pr.Tasks = append(pr.Tasks, &TaskRun{
			// 	Status: NotStarted,
			// })
			pr.Tasks[taskRunStatus.PipelineTaskName].Status = NotStarted
		}
		// for _, c := range status.Status.Conditions {
		// 	pr.Tasks = append(pr.Tasks, &TaskRun{
		// 		Status: ,
		// 	})
		// }
	}

	pr.CompleteTasks = pr.FailedTasks + pr.SucceededTasks
	// }
	// taskRuns := make(map[string]interface{})
	// if truns, ok := st["taskRuns"]; ok {
	// 	if taskRuns, ok = truns.(map[string]interface{}); ok {
	// 		for _, t := range taskRuns {
	// 			if task, ok := t.()
	// 		}
	// 	}
	// }

	return pr, nil
}

func (p *PipelineRun) GetStartedAtString() string {
	if p.StartedAt == nil {
		return ""
	}
	return p.StartedAt.Format(time.RFC3339)
}

func (p *PipelineRun) GetCompletedAtString() string {
	if p.CompletedAt == nil {
		return ""
	}
	return p.CompletedAt.Format(time.RFC3339)
}

func (p *PipelineRun) String() string {
	return fmt.Sprintf("%s | %s | %s | %s | %s", p.Name, p.CreatedAt, p.GetStartedAtString(), p.GetCompletedAtString(), p.Status)
}

func (p *PipelineRun) NotEqual(other *PipelineRun) bool {
	return p.StartedAtNotEqual(other) || p.CompletedAtNotEqual(other) || p.StartedAtNotEqual(other) || p.TaskStatusNotEqual(other)
}

func (p *PipelineRun) StatusNotEqual(other *PipelineRun) bool {
	return p.Status != other.Status
}

func (p *PipelineRun) CompletedAtNotEqual(other *PipelineRun) bool {
	return p.GetCompletedAtString() != other.GetCompletedAtString()
}

func (p *PipelineRun) StartedAtNotEqual(other *PipelineRun) bool {
	return p.GetStartedAtString() != other.GetStartedAtString()
}

func (p *PipelineRun) TaskStatusNotEqual(other *PipelineRun) bool {
	return p.TotalTasks != other.TotalTasks || p.RunningTasks != other.RunningTasks || p.SucceededTasks != other.SucceededTasks || p.PendingTasks != other.PendingTasks
}

func (p *PipelineRun) setDuration() {
	if p.StartedAt == nil {
		return
	}

	var t time.Duration

	if p.CompletedAt == nil {
		t = time.Since(*p.StartedAt)
	} else {
		t = p.CompletedAt.Sub(*p.StartedAt)
	}
	p.Duration = &t
}

// func (p *PipelineRun) SetStatus(status map[string]interface{}) {

// }

// status
/*
{
   "completionTime":"2022-06-02T12:53:59Z",
   "conditions":[
      {
         "lastTransitionTime":"2022-06-02T12:53:59Z",
         "message":"Tasks Completed: 1 (Failed: 0, Cancelled 0), Skipped: 0",
         "reason":"Succeeded",
         "status":"True",
         "type":"Succeeded"
      }
   ],
   "pipelineSpec":{
      "params":[
         {
            "name":"MESSAGE",
            "type":"string"
         }
      ],
      "tasks":[
         {
            "name":"echo-message",
            "params":[
               {
                  "name":"MESSAGE",
                  "value":"$(params.MESSAGE)"
               }
            ],
            "taskSpec":{
               "metadata":{

               },
               "params":[
                  {
                     "name":"MESSAGE",
                     "type":"string"
                  }
               ],
               "spec":null,
               "steps":[
                  {
                     "image":"ubuntu",
                     "name":"echo",
                     "resources":{

                     },
                     "script":"#!/usr/bin/env bash\necho \"$(params.MESSAGE)\"            \n"
                  }
               ]
            }
         }
      ]
   },
   "startTime":"2022-06-02T12:53:54Z",
   "taskRuns":{
      "echo-6gsmt-echo-message":{
         "pipelineTaskName":"echo-message",
         "status":{
            "completionTime":"2022-06-02T12:53:59Z",
            "conditions":[
               {
                  "lastTransitionTime":"2022-06-02T12:53:59Z",
                  "message":"All Steps have completed executing",
                  "reason":"Succeeded",
                  "status":"True",
                  "type":"Succeeded"
               }
            ],
            "podName":"echo-6gsmt-echo-message-pod",
            "startTime":"2022-06-02T12:53:54Z",
            "steps":[
               {
                  "container":"step-echo",
                  "imageID":"docker.io/library/ubuntu@sha256:26c68657ccce2cb0a31b330cb0be2b5e108d467f641c62e13ab40cbec258c68d",
                  "name":"echo",
                  "terminated":{
                     "containerID":"containerd://4f31b45a40b6b376c45f3006549929ee6340d404c4b25f8ddfc1195fc6fd7992",
                     "exitCode":0,
                     "finishedAt":"2022-06-02T12:53:59Z",
                     "reason":"Completed",
                     "startedAt":"2022-06-02T12:53:59Z"
                  }
               }
            ],
            "taskSpec":{
               "params":[
                  {
                     "name":"MESSAGE",
                     "type":"string"
                  }
               ],
               "steps":[
                  {
                     "image":"ubuntu",
                     "name":"echo",
                     "resources":{

                     },
                     "script":"#!/usr/bin/env bash\necho \"$(params.MESSAGE)\"            \n"
                  }
               ]
            }
         }
      }
   }
}
*/

/* status transitions
[map[lastTransitionTime:2022-06-02T17:13:55Z message:PipelineRun "echo-t7km8" is pending reason:PipelineRunPending status:Unknown type:Succeeded]]

[map[lastTransitionTime:2022-06-02T17:07:05Z message:Tasks Completed: 0 (Failed: 0, Cancelled 0), Incomplete: 1, Skipped: 0 reason:Running status:Unknown type:Succeeded]]
[map[lastTransitionTime:2022-06-02T17:07:05Z message:Tasks Completed: 0 (Failed: 0, Cancelled 0), Incomplete: 1, Skipped: 0 reason:Running status:Unknown type:Succeeded]]
[map[lastTransitionTime:2022-06-02T17:07:05Z message:Tasks Completed: 0 (Failed: 0, Cancelled 0), Incomplete: 1, Skipped: 0 reason:Running status:Unknown type:Succeeded]]
[map[lastTransitionTime:2022-06-02T17:07:05Z message:Tasks Completed: 0 (Failed: 0, Cancelled 0), Incomplete: 1, Skipped: 0 reason:Running status:Unknown type:Succeeded]]
[map[lastTransitionTime:2022-06-02T17:07:05Z message:Tasks Completed: 0 (Failed: 0, Cancelled 0), Incomplete: 1, Skipped: 0 reason:Running status:Unknown type:Succeeded]]
[map[lastTransitionTime:2022-06-02T17:07:05Z message:Tasks Completed: 0 (Failed: 0, Cancelled 0), Incomplete: 1, Skipped: 0 reason:Running status:Unknown type:Succeeded]]
[map[lastTransitionTime:2022-06-02T17:07:05Z message:Tasks Completed: 0 (Failed: 0, Cancelled 0), Incomplete: 1, Skipped: 0 reason:Running status:Unknown type:Succeeded]]
[map[lastTransitionTime:2022-06-02T17:07:11Z message:Tasks Completed: 1 (Failed: 0, Cancelled 0), Skipped: 0 reason:Succeeded status:True type:Succeeded]]
*/
