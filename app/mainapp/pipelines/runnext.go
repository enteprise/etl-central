package pipelines

import (
	dpconfig "github.com/enteprise/etl-central/app/mainapp/config"

	"github.com/enteprise/etl-central/app/mainapp/database/models"
	"github.com/enteprise/etl-central/app/mainapp/logging"
	"github.com/enteprise/etl-central/app/mainapp/messageq"
)

func RunNextPipeline() {

	_, err := messageq.NATSencoded.QueueSubscribe("pipeline-run-next", "runnext", func(subj, reply string, msg models.WorkerTaskSend) {

		RunNext(msg)

	})

	if err != nil {
		if dpconfig.Debug == "true" {
			logging.PrintSecretsRedact(err)
		}

	}

}
