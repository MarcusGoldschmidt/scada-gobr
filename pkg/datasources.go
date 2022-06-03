package pkg

import (
	"context"
	"github.com/MarcusGoldschmidt/scadagobr/pkg/datasources"
)

func LoadDataSourcesRuntimeManager(ctx context.Context, s *Scadagobr) ([]datasources.DataSourceRuntimeManager, error) {
	ds, err := s.dataSourcePersistence.GetDadaSources(ctx)

	if err != nil {
		return nil, err
	}

	var dsrm []datasources.DataSourceRuntimeManager

	for _, d := range ds {

		manager, err := DataSourceToRuntimeManager(s, d)
		if err != nil {
			return nil, err
		}

		dsrm = append(dsrm, manager)
	}

	return dsrm, nil
}
