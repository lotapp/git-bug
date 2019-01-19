package launchpad

import (
	"fmt"
	"time"

	"github.com/MichaelMure/git-bug/bridge/core"
	"github.com/MichaelMure/git-bug/bug"
	"github.com/MichaelMure/git-bug/cache"
	"github.com/MichaelMure/git-bug/identity"
	"github.com/pkg/errors"
)

type launchpadImporter struct {
	conf core.Configuration
}

func (li *launchpadImporter) Init(conf core.Configuration) error {
	li.conf = conf
	return nil
}

const keyLaunchpadID = "launchpad-id"
const keyLaunchpadLogin = "launchpad-login"

func (li *launchpadImporter) makePerson(repo *cache.RepoCache, owner LPPerson) (*identity.Identity, error) {
	// Look first in the cache
	i, err := repo.ResolveIdentityImmutableMetadata(keyLaunchpadLogin, owner.Login)
	if err == nil {
		return i, nil
	}
	if _, ok := err.(identity.ErrMultipleMatch); ok {
		return nil, err
	}

	return repo.NewIdentityRaw(
		owner.Name,
		"",
		owner.Login,
		"",
		map[string]string{
			keyLaunchpadLogin: owner.Login,
		},
	)
}

func (li *launchpadImporter) ImportAll(repo *cache.RepoCache) error {
	lpAPI := new(launchpadAPI)

	err := lpAPI.Init()
	if err != nil {
		return err
	}

	lpBugs, err := lpAPI.SearchTasks(li.conf["project"])
	if err != nil {
		return err
	}

	for _, lpBug := range lpBugs {
		lpBugID := fmt.Sprintf("%d", lpBug.ID)
		_, err := repo.ResolveBugCreateMetadata(keyLaunchpadID, lpBugID)
		if err != nil && err != bug.ErrBugNotExist {
			return err
		}

		owner, err := li.makePerson(repo, lpBug.Owner)
		if err != nil {
			return err
		}

		if err == bug.ErrBugNotExist {
			createdAt, _ := time.Parse(time.RFC3339, lpBug.CreatedAt)
			_, err := repo.NewBugRaw(
				owner,
				createdAt.Unix(),
				lpBug.Title,
				lpBug.Description,
				nil,
				map[string]string{
					keyLaunchpadID: lpBugID,
				},
			)
			if err != nil {
				return errors.Wrapf(err, "failed to add bug id #%s", lpBugID)
			}
		} else {
			/* TODO: Update bug */
			fmt.Println("TODO: Update bug")
		}

	}
	return nil
}

func (li *launchpadImporter) Import(repo *cache.RepoCache, id string) error {
	fmt.Println("IMPORT")
	return nil
}
