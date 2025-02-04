package distfilesystem

import (
	"errors"
	"log"
	"strings"

	"github.com/enteprise/etl-central/app/mainapp/database/models"

	"github.com/enteprise/etl-central/app/mainapp/database"
	wrkerconfig "github.com/enteprise/etl-central/app/workers/config"

	"gorm.io/gorm"
)

/*
 1. Check the cache that files are up to date.
    1a. Cache at node level
    1b. Cache at file level
 2. Find all the files for this node.
 3. Batch download files

NOTE: Node ID is the node in the graph and not the worker ID.
*/
func DistributedStoragePipelineDownload(environmentID string, directoryRun string, nodeID string) error {

	// log.Println("folder", folder)
	if strings.Contains(directoryRun, "/pipeline/") == false {
		log.Println("Folder incorrect format - doesn't contain /pipeline/")
		return errors.New("Folder incorrect format - doesn't contain /pipeline/")
	}

	// if strings.Contains(folder, "_Platform") == false {
	// 	log.Println("Folder incorrect format - doesn't contain _Platform")
	// 	return errors.New("Folder incorrect format - doesn't contain _Platform")
	// }
	/* node cache is a higher level cache for the node */
	nodeCache := models.CodeNodeCache{}
	err := database.DBConn.Select("cache_valid").Where("node_id = ? and environment_id = ? and worker_id = ?", nodeID, environmentID, wrkerconfig.WorkerID).First(&nodeCache).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	/*
		If the cache is invalid, download only the changed files, otherwise do nothing
	*/
	if nodeCache.CacheValid != true {

		if wrkerconfig.Debug == "true" {
			log.Println("Node cache invalid, updating node:", nodeID)
		}

		FilesOutput := []*models.CodeFilesCacheOutput{}

		/* Retrieve all the files which are not in lower level cache in table code_node_cache */
		query := `
		select 
		cf.file_id, 
		cf.folder_id, 
		cf.file_name, 
		cs.checksum_md5,
		cs.file_store,
		cfolder.level
		from code_files cf, code_files_store cs, code_folders cfolder
		where 
		cf.folder_id = cfolder.folder_id and
		cf.file_id = cs.file_id and 
		cf.environment_id = cs.environment_id and 
		cf.node_id = ? and 
		cs.run_include = true and 
		cf.environment_id =?
		and NOT EXISTS
        (
        SELECT  file_id
        FROM code_files_cache cfc
        WHERE   
		cf.file_id = cfc.file_id and 
		cfc.worker_id = ? and 
		cf.node_id = cfc.node_id and 
		cf.environment_id = cfc.environment_id and 
		cf.environment_id =?
        )
		`

		err = database.DBConn.Raw(query, nodeID, environmentID, wrkerconfig.WorkerID, environmentID).Scan(&FilesOutput).Error
		if err != nil {
			log.Println("Download cached files from DB: ", err)
			return err
		}

		// distfilesystem.BatchFileWrite(FilesOutput, folderID, environmentID, folder)
		// log.Println("==== FS:", RunType, version)
		errSave := DistributedStorageFileSave(FilesOutput, directoryRun, nodeID, environmentID, "pipeline", "latest")
		if errSave != nil {
			log.Println("Error saved to cache: ", errSave)
			return errSave
		}

	}

	return nil
}
