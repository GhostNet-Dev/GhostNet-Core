package workload

type IWorkload interface {
	PrepareRun()
	Run()
	CheckRunning() bool
}
