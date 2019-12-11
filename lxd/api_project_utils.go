package main

import (
	"github.com/lxc/lxd/lxd/db"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
	"github.com/pkg/errors"
)

func doProjectUpdateContainer(d *Daemon, project *api.Project, req api.ProjectPut, args db.InstanceArgs) error {

	c := containerLXCInstantiate(d.State(), args)

	if project.Config["limits.cpu"] != req.Config["limits.cpu"] {
		logger.Infof("UPDATING CPU LIMIT FOR PROJECT TO %s", req.Config["limits.cpu"])
		c.expandedConfig["limits.cpu"] = req.Config["limits.cpu"]
	}

	return c.Update(db.InstanceArgs{
		Architecture: c.Architecture(),
		Config:       c.LocalConfig(),
		Description:  c.Description(),
		Devices:      c.LocalDevices(),
		Ephemeral:    c.IsEphemeral(),
		Profiles:     c.Profiles(),
		Project:      c.Project(),
		Type:         c.Type(),
		Snapshot:     c.IsSnapshot(),
	}, true)
}

func getProjectContainersInfo(cluster *db.Cluster, project string) ([]db.InstanceArgs, error) {
	// grab all containers given th project

	var names []string
	var err error
	err = cluster.Transaction(func (tx *db.ClusterTx) error {
		var err error
		names, err = tx.ContainerNames(project)
		return err
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query containers with project '%s'", project)
	}

	containers := []db.InstanceArgs{}
	err = cluster.Transaction(func(tx *db.ClusterTx) error {
		for _, ctName := range names {
				container, err := tx.InstanceGet(project, ctName)
				if err != nil {
					return err
				}
				containers = append(containers, db.ContainerToArgs(container))
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to fetch containers")
	}

	return containers, nil
}
