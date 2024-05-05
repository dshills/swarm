package task

type List interface {
	Push(...string)
	Pop() string
	Peak(idx int) string
	Len() int
}

func New() List {
	return &_taskList{}
}

type _taskList struct {
	tasks []string
}

func (tl *_taskList) Push(tasks ...string) {
	for _, task := range tasks {
		if task != "" {
			tl.tasks = append(tl.tasks, task)
		}
	}
}

func (tl *_taskList) Pop() string {
	if len(tl.tasks) == 0 {
		return ""
	}
	task := ""
	task, tl.tasks = tl.tasks[len(tl.tasks)-1], tl.tasks[:len(tl.tasks)-1]
	return task
}

func (tl *_taskList) Peak(idx int) string {
	if idx >= len(tl.tasks) {
		return ""
	}
	return tl.tasks[idx]
}

func (tl *_taskList) Len() int {
	return len(tl.tasks)
}
