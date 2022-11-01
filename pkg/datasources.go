package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
)

func LoadDataSourcesRuntimeManager(ctx context.Context, s *Scadagobr) ([]datasources.DataSourceRuntimeManager, error) {
	ds, err := s.dataSourcePersistence.GetDataSources(ctx)

	if err != nil {
		return nil, err
	}

	dsrm := make([]datasources.DataSourceRuntimeManager, len(ds))

	for i, d := range ds {

		manager, err := DataSourceToRuntimeManager(s, d)
		if err != nil {
			return nil, err
		}

		dsrm[i] = manager
	}

	return dsrm, nil
}
