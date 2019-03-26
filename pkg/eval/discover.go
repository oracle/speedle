package eval

import (
	log "github.com/sirupsen/logrus"
	"gitlab-odx.oracledx.com/wcai/speedle/api/ads"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/errors"
	"gitlab-odx.oracledx.com/wcai/speedle/pkg/store"
)

func (p *PolicyEvalImpl) Discover(ctx ads.RequestContext) (bool, ads.Reason, error) {
	if d, ok := p.Store.(store.DiscoverRequestManager); ok {
		err := d.SaveDiscoverRequest(&ctx)
		if err != nil {
			log.Warn("error in saving discover request, ", err)
		}
		return true, ads.DISCOVER_MODE, err
	}
	return true, ads.DISCOVER_MODE, errors.Errorf(errors.DiscoverError, "unsupported store type of discovery function:%s", p.Store.Type())
}
