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
	UUID           string            `json:"uuid"`
	Name           string            `json:"name"`
	Namespace      string            `json:"namespace"`
	Labels         map[string]string `json:"labels,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	CompletedAt    *time.Time        `json:"completedAt,omitempty"`
	StartedAt      *time.Time        `json:"startedAt,omitempty"`
	CreatedAt      time.Time         `json:"createdAt,omitempty"`
	Status         PipelineStatus    `json:"status,omitempty"`
	Tasks          []*TaskRun        `json:"tasks,omitempty"`
	TotalTasks     int
	RunningTasks   int
	PendingTasks   int
	FailedTasks    int
	SucceededTasks int
}

type TaskRun struct {
	Status PipelineStatus `json:"status,omitempty"`
}

type PipelineStatus int16

const (
	Undefined PipelineStatus = iota
	PipelineRunPending
	NotStarted
	Running
	Success
	Failed
	Pending
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

	var st map[string]interface{}
	stInt, ok := uns.Object["status"]
	if ok {
		if st, ok = stInt.(map[string]interface{}); !ok {
			return pr, nil
		}
	} else {
		return pr, nil
	}

	if r, ok := st["completionTime"]; ok {
		t, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", r))
		if err == nil {
			pr.CompletedAt = &t
		}
	}

	if r, ok := st["startTime"]; ok {
		t, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", r))
		if err == nil {
			pr.StartedAt = &t
		}
	}

	var conditions []interface{}
	if c, ok := st["conditions"]; ok {
		if conditions, ok = c.([]interface{}); ok {
			for _, conditionInt := range conditions {
				if condition, ok := conditionInt.(map[string]interface{}); ok {
					reason := fmt.Sprintf("%s", condition["reason"])
					status := fmt.Sprintf("%s", condition["status"])
					cType := fmt.Sprintf("%s", condition["type"])
					if strings.EqualFold(cType, "Succeeded") {
						if strings.EqualFold(reason, "Running") {
							pr.Status = Running
						} else if strings.EqualFold(reason, "Succeeded") {
							if strings.EqualFold(status, "True") {
								pr.Status = Success
							} else {
								pr.Status = Failed
							}
						} else if strings.EqualFold(reason, "Failed") {
							pr.Status = Failed
						}
					}
				}
			}
		}
	}

	// prs := v1beta1.PipelineRunStatus{}
	prs := make(map[string]*v1beta1.PipelineRunTaskRunStatus)
	if truns, ok := st["taskRuns"]; ok {
		d, err := json.Marshal(truns)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(d, &prs)
		if err != nil {
			return nil, err
		}
		for _, status := range prs {
			pr.TotalTasks++
			if status.Status != nil && len(status.Status.Conditions) > 0 {
				cond := status.Status.Conditions[0]

				if strings.EqualFold(cond.Reason, "running") {
					pr.Tasks = append(pr.Tasks, &TaskRun{
						Status: Running,
					})
					pr.RunningTasks++
					continue
				}

				if strings.EqualFold(cond.Reason, "pending") {
					pr.Tasks = append(pr.Tasks, &TaskRun{
						Status: Pending,
					})
					pr.PendingTasks++
					continue
				}

				if strings.EqualFold(cond.Reason, "Succeeded") && strings.EqualFold(string(cond.Status), "True") {
					pr.Tasks = append(pr.Tasks, &TaskRun{
						Status: Success,
					})
					pr.SucceededTasks++
					continue
				}

				if strings.EqualFold(cond.Reason, "Succeeded") && strings.EqualFold(string(cond.Status), "False") {
					pr.Tasks = append(pr.Tasks, &TaskRun{
						Status: Failed,
					})
					pr.FailedTasks++
					continue
				}
			} else {
				pr.Tasks = append(pr.Tasks, &TaskRun{
					Status: NotStarted,
				})
			}
			// for _, c := range status.Status.Conditions {
			// 	pr.Tasks = append(pr.Tasks, &TaskRun{
			// 		Status: ,
			// 	})
			// }
		}
	}
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
	return p.TotalTasks != other.TotalTasks || p.RunningTasks != other.RunningTasks
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
