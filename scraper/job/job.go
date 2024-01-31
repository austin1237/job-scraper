package job

type Job struct {
	Title   string
	Company string
	Link    string
}

func DeduplicatedLinks(jobs []Job) []Job {
	seen := make(map[string]bool)
	deduplicated := []Job{}

	for _, possibleJob := range jobs {
		if !seen[possibleJob.Link] {
			seen[possibleJob.Link] = true
			deduplicated = append(deduplicated, possibleJob)
		}
	}
	return deduplicated
}
