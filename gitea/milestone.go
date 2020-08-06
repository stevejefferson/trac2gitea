package gitea

import "log"

func (accessor *Accessor) AddMilestone(name string, content string, closed bool, deadlineTS int64, closedTS int64) {
	_, err := accessor.db.Exec(`
		INSERT INTO
			milestone(repo_id,name,content,is_closed,deadline_unix,closed_date_unix)
			SELECT $1,$2,$3,$4,$5,$6 WHERE
				NOT EXISTS (SELECT * FROM milestone WHERE repo_id = $1 AND name = $2)`,
		accessor.repoID, name, content, closed, deadlineTS, closedTS)
	if err != nil {
		log.Fatal(err)
	}
}
