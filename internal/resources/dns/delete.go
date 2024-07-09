package dns

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Delete(ctx clustercontext.ClusterContext) {
	records, err := Get(ctx)
	if err != nil {
		return
	}

	for _, record := range records {
		recordId := record.Name + " | " + record.Type + " | " + record.Value

		logger.LogResourceEvent(logger.DNSRecord, logger.Delete, recordId, logger.Initialized)

		err := ctx.DNSClient.DeleteRecord(ctx.Context, record.ID)

		if err != nil {
			logger.LogResourceEvent(logger.DNSRecord, logger.Delete, recordId, logger.Failure, err)
		} else {
			logger.LogResourceEvent(logger.DNSRecord, logger.Delete, recordId, logger.Success)
		}
	}

}
